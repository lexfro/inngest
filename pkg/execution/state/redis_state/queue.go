package redis_state

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"golang.org/x/exp/rand"
	"golang.org/x/sync/semaphore"

	"github.com/go-redis/redis/v8"
	json "github.com/goccy/go-json"
	"github.com/google/uuid"
	osqueue "github.com/inngest/inngest/pkg/execution/queue"
	"github.com/oklog/ulid/v2"
	"github.com/uber-go/tally"
	"gonum.org/v1/gonum/stat/sampleuv"
)

const (
	PartitionSelectionMax   int64 = 20
	PartitionPeekMax        int64 = PartitionSelectionMax * 3
	PartitionLeaseDuration        = 2 * time.Second
	QueueSelectionMax       int64 = 50
	QueuePeekMax            int64 = 1000
	QueuePeekDefault        int64 = QueueSelectionMax * 3
	QueueLeaseDuration            = 10 * time.Second
	SequentialLeaseDuration       = 10 * time.Second
	SequentialLeaseMax            = 20 * time.Second

	PriorityMax     uint = 0
	PriorityDefault uint = 5
	PriorityMin     uint = 9
)

var (
	ErrQueueItemExists               = fmt.Errorf("queue item already exists")
	ErrQueueItemNotFound             = fmt.Errorf("queue item not found")
	ErrQueueItemAlreadyLeased        = fmt.Errorf("queue item already leased")
	ErrQueueItemLeaseMismatch        = fmt.Errorf("item lease does not match")
	ErrQueueItemNotLeased            = fmt.Errorf("queue item is not leased")
	ErrQueuePeekMaxExceedsLimits     = fmt.Errorf("peek exceeded the maximum limit of %d", QueuePeekMax)
	ErrPriorityTooLow                = fmt.Errorf("priority is too low")
	ErrPriorityTooHigh               = fmt.Errorf("priority is too high")
	ErrWeightedSampleRead            = fmt.Errorf("error reading from weighted sample")
	ErrPartitionNotFound             = fmt.Errorf("partition not found")
	ErrPartitionAlreadyLeased        = fmt.Errorf("partition already leased")
	ErrPartitionPeekMaxExceedsLimits = fmt.Errorf("peek exceeded the maximum limit of %d", PartitionPeekMax)
	ErrSequentialAlreadyLeased       = fmt.Errorf("sequential scanner already leased")
	ErrSequentialLeaseExceedsLimits  = fmt.Errorf("sequential lease duration exceeds the maximum of %d seconds", int(SequentialLeaseMax.Seconds()))
)

var (
	source rand.Source
	rnd    *rand.Rand
)

func init() {
	// For weighted shuffles generate a new rand.
	source = rand.NewSource(uint64(time.Now().UnixNano()))
	rnd = rand.New(source)
}

// PriorityFinder returns the priority for a given workflow.
type PriorityFinder func(ctx context.Context, workflowID uuid.UUID) uint

type QueueOpt func(q *queue)

func WithName(name string) func(q *queue) {
	return func(q *queue) {
		q.name = name
	}
}

func WithMetricsScope(scope tally.Scope) func(q *queue) {
	return func(q *queue) {
		q.metrics = scope
	}
}

func WithPriorityFinder(pf PriorityFinder) func(q *queue) {
	return func(q *queue) {
		q.pf = pf
	}
}

func WithQueueKeyGenerator(kg QueueKeyGenerator) func(q *queue) {
	return func(q *queue) {
		q.kg = kg
	}
}

func NewQueue(r *redis.Client, opts ...QueueOpt) *queue {
	q := &queue{
		r: r,
		pf: func(ctx context.Context, workflowID uuid.UUID) uint {
			return PriorityDefault
		},
		kg:         defaultQueueKey,
		numWorkers: defaultNumWorkers,
		metrics:    tally.NewTestScope("queue", map[string]string{}),
		wg:         &sync.WaitGroup{},
	}

	for _, opt := range opts {
		opt(q)
	}

	q.r = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	q.sem = &trackingSemaphore{Weighted: semaphore.NewWeighted(int64(q.numWorkers))}

	return q
}

type queue struct {
	// name is the identifiable name for this worker, for logging.
	name string
	// redis stores the redis connection to use.
	r  *redis.Client
	pf PriorityFinder
	kg QueueKeyGenerator
	// metrics allows reporting of metrics
	metrics tally.Scope

	// quit is a channel that any method can send on to trigger termination
	// of the Run loop.
	quit chan bool
	// seqLeaseID stores the lease ID if this queue is the sequential processor.
	// all runners attempt to claim this lease automatically.
	seqLeaseID *ulid.ULID
	// wg stores a waitgroup for all in-progress jobs
	wg *sync.WaitGroup
	// numWorkers stores the number of workers available to concurrently process jobs.
	numWorkers int32
	// workers is a buffered channel which allows scanners to send queue items
	// to workers to be processed
	workers chan QueueItem
	// sem stores a semaphore controlling the number of jobs currently
	// being processed.  This lets us check whether there's capacity in the queue
	// prior to leasing items.
	sem *trackingSemaphore
}

// QueueItem represents an individually queued work scheduled for some time in the
// future.
type QueueItem struct {
	// ID represents a unique identifier for the queue item.  The ULID
	// stores when the QueueItem was created.
	ID ulid.ULID `json:"id"`
	// At represents the current time that this QueueItem needs to be executed at,
	// as a millisecond epoch.  This is the millisecond-level granularity score of
	// the item.  Note that the score in Redis is second-level, ie this field / 1000.
	//
	// This is necessary for rescoring partitions and checking latencies.
	At          int64     `json:"at"`
	WorkflowID  uuid.UUID `json:"workflowID"`
	WorkspaceID uuid.UUID `json:"workspaceID"`
	Attempt     int       `json:"attempt"`
	MaxAttempts int       `json:"maxAttempts"`
	// LeaseID is a ULID which embeds a timestamp denoting when the lease expires.
	LeaseID *ulid.ULID `json:"leaseID"`
	// Data represents the enqueued data, eg. the edge to process or the pause
	// to resume.
	Data osqueue.Item `json:"data"`
}

func (q QueueItem) MarshalBinary() ([]byte, error) {
	return json.Marshal(q)
}

// QueuePartition represents an individual queue for a workflow.  It stores the
// time of the earliest job within the workflow.
type QueuePartition struct {
	WorkflowID uuid.UUID `json:"wid"`

	Priority uint `json:"p"`

	// At refers to the earliest QueueItem time within this partition, in
	// seconds as a unix epoch.
	//
	// This is updated when taking a lease.
	At int64 `json:"at"`

	// Last represents the time that this QueueItem was last leased, as a unix
	// epoch.
	//
	// This lets us inspect the time a QueuePartition was last leased, and figure
	// out whether we should garbage collect the partition index.
	Last int64 `json:"last"`

	// LeaseID represents a lease on this partition.  If the LeaseID is not nil,
	// this partition can be claimed by a shared-nothing worker to work on the
	// queue items within this partition.
	//
	// A lease is shortly held (eg seconds).  It should last long enough for
	// workers to claim QueueItems only.
	LeaseID *ulid.ULID `json:"leaseID"`
}

func (q QueuePartition) MarshalBinary() ([]byte, error) {
	return json.Marshal(q)
}

// EnqueueItem enqueues a QueueItem.  It creates a QueuePartition for the workspace
// if a partition does not exist.
//
// The QueueItem's ID can be a zero UUID;  if the ID is a zero value a new ID
// will be created for the queue item.
func (q queue) EnqueueItem(ctx context.Context, i QueueItem, at time.Time) (QueueItem, error) {
	if i.ID.Compare(ulid.ULID{}) == 0 {
		i.ID = ulid.MustNew(ulid.Now(), rnd)
	}

	priority := PriorityMin
	if q.pf != nil {
		priority = q.pf(ctx, i.WorkflowID)
	}

	if priority > PriorityMin {
		return i, ErrPriorityTooLow
	}
	if priority < PriorityMax {
		return i, ErrPriorityTooHigh
	}

	// Add the At timestamp.
	i.At = at.UnixMilli()

	qp := QueuePartition{WorkflowID: i.WorkflowID, Priority: priority, At: at.Unix()}
	keys := []string{
		q.kg.QueueItem(),                          // Queue item
		q.kg.QueueIndex(i.WorkflowID.String()),    // Queue sorted set
		q.kg.PartitionItem(),                      // Partition item, map
		q.kg.PartitionMeta(i.WorkflowID.String()), // Partition item
		q.kg.PartitionIndex(),                     // Global partition queue
	}
	status, err := scripts["queue/enqueue"].Eval(
		ctx,
		q.r,
		keys,

		i,
		i.ID.String(),
		at.Unix(),
		i.WorkflowID.String(),
		qp,
	).Int64()
	if err != nil {
		return i, fmt.Errorf("error enqueueing item: %w", err)
	}
	switch status {
	case 0:
		return i, nil
	case 1:
		return i, ErrQueueItemExists
	default:
		return i, fmt.Errorf("unknown response enqueueing item: %d", status)
	}
}

// Peek takes n items from a queue, up until QueuePeekMax.
//
// If limit is -1, this will return the first unleased item - representing the next available item in the
// queue.
func (q queue) Peek(ctx context.Context, workflowID uuid.UUID, until time.Time, limit int64) ([]*QueueItem, error) {
	// Check whether limit is -1, peeking next available time
	isPeekNext := limit == -1

	if limit > QueuePeekMax {
		// Lua's max unpack() length is 8000; don't allow users to peek more than
		// 1k at a time regardless.
		return nil, ErrQueuePeekMaxExceedsLimits
	}
	if limit <= 0 {
		limit = QueuePeekMax
	}

	items, err := scripts["queue/peek"].Eval(
		ctx,
		q.r,
		[]string{
			q.kg.QueueIndex(workflowID.String()),
			q.kg.QueueItem(),
		},
		until.Unix(),
		limit,
	).StringSlice()
	if err != nil {
		return nil, fmt.Errorf("error peeking queue items: %w", err)
	}

	// Create a slice up to items in length.  We're going to remove any items that are
	// leased here, so we may end up returning less than the total length.
	result := make([]*QueueItem, len(items))
	n := 0
	now := time.Now()

	for _, str := range items {
		qi := &QueueItem{}
		if err := json.Unmarshal([]byte(str), qi); err != nil {
			return nil, fmt.Errorf("error unmarshalling peeked queue item: %w", err)
		}
		if qi.LeaseID != nil && now.Before(ulid.Time(qi.LeaseID.Time())) {
			// Leased item, don't return.
			continue
		}
		result[n] = qi
		n++

		if isPeekNext {
			return []*QueueItem{qi}, nil
		}
	}

	return result[0:n], nil
}

// Lease temporarily dequeues an item from the queue by obtaining a lease, preventing
// other workers from working on this queue item at the same time.
//
// Obtaining a lease updates the vesting time for the queue item until now() +
// lease duration. This returns the newly acquired lease ID on success.
func (q queue) Lease(ctx context.Context, workflowID uuid.UUID, itemID ulid.ULID, duration time.Duration) (*ulid.ULID, error) {
	// TODO: Add custom throttling here.
	leaseID, err := ulid.New(ulid.Timestamp(time.Now().Add(duration).UTC()), rnd)
	if err != nil {
		return nil, fmt.Errorf("error generating id: %w", err)
	}

	keys := []string{
		q.kg.QueueItem(),
		q.kg.QueueIndex(workflowID.String()),
		q.kg.PartitionMeta(workflowID.String()),
	}
	status, err := scripts["queue/lease"].Eval(
		ctx,
		q.r,
		keys,
		itemID.String(),
		leaseID.String(),
		time.Now().UnixMilli(),
	).Int64()
	if err != nil {
		return nil, fmt.Errorf("error leasing pause: %w", err)
	}
	switch status {
	case 0:
		return &leaseID, nil
	case 1:
		return nil, ErrQueueItemNotFound
	case 2:
		return nil, ErrQueueItemAlreadyLeased
	default:
		return nil, fmt.Errorf("unknown response enqueueing item: %d", status)
	}
}

// ExtendLease extens the lease for a given queue item, given the queue item is currently
// leased with the given ID.  This returns a new lease ID if the lease is successfully ended.
//
// The existing lease ID must be passed in so that we can guarantee that the worker
// renewing the lease still owns the lease.
//
// Renewing a lease updates the vesting time for the queue item until now() +
// lease duration. This returns the newly acquired lease ID on success.
func (q queue) ExtendLease(ctx context.Context, i QueueItem, leaseID ulid.ULID, duration time.Duration) (*ulid.ULID, error) {
	newLeaseID, err := ulid.New(ulid.Timestamp(time.Now().Add(duration).UTC()), rnd)
	if err != nil {
		return nil, fmt.Errorf("error generating id: %w", err)
	}

	keys := []string{
		q.kg.QueueItem(),
		q.kg.QueueIndex(i.WorkflowID.String()),
	}
	status, err := scripts["queue/extendLease"].Eval(
		ctx,
		q.r,
		keys,
		i.ID.String(),
		leaseID.String(),
		newLeaseID.String(),
	).Int64()
	if err != nil {
		return nil, fmt.Errorf("error extending lease: %w", err)
	}
	switch status {
	case 0:
		return &newLeaseID, nil
	case 1:
		return nil, ErrQueueItemNotFound
	case 2:
		return nil, ErrQueueItemNotLeased
	case 3:
		return nil, ErrQueueItemLeaseMismatch
	default:
		return nil, fmt.Errorf("unknown response enqueueing item: %d", status)
	}
}

// Dequeue removes an item from the queue entirely.
func (q queue) Dequeue(ctx context.Context, i QueueItem) error {
	keys := []string{
		q.kg.QueueItem(),
		q.kg.QueueIndex(i.WorkflowID.String()),
		q.kg.PartitionMeta(i.WorkflowID.String()),
	}
	status, err := scripts["queue/dequeue"].Eval(
		ctx,
		q.r,
		keys,
		i.ID.String(),
	).Int64()
	if err != nil {
		return fmt.Errorf("error dequeueing item: %w", err)
	}
	switch status {
	case 0:
		return nil
	case 1:
		return ErrQueueItemNotFound
	default:
		return fmt.Errorf("unknown response dequeueing item: %d", status)
	}
}

// Requeue requeues an item in the future.
func (q queue) Requeue(ctx context.Context, i QueueItem, at time.Time) error {
	priority := PriorityMin
	if q.pf != nil {
		priority = q.pf(ctx, i.WorkflowID)
	}

	if priority > PriorityMin {
		return ErrPriorityTooLow
	}

	// Unset any lease ID as this is requeued.
	i.LeaseID = nil
	// Update the At timestamp.
	i.At = at.UnixMilli()

	qp := QueuePartition{WorkflowID: i.WorkflowID, Priority: priority, At: at.Unix()}
	keys := []string{
		q.kg.QueueItem(),
		q.kg.QueueIndex(i.WorkflowID.String()),
		q.kg.PartitionMeta(i.WorkflowID.String()),
		q.kg.PartitionIndex(),
	}
	status, err := scripts["queue/requeue"].Eval(
		ctx,
		q.r,
		keys,

		i,
		i.ID.String(),
		at.Unix(),
		qp.WorkflowID.String(),
		qp,
	).Int64()
	if err != nil {
		return fmt.Errorf("error requeueing item: %w", err)
	}
	switch status {
	case 0:
		return nil
	default:
		return fmt.Errorf("unknown response enqueueing item: %d", status)
	}
}

// PartitionLease leases a parititon for a given workflow ID.  It returns the new lease ID.
func (q queue) PartitionLease(ctx context.Context, wid uuid.UUID, duration time.Duration) (*ulid.ULID, error) {
	// XXX: Check for function throttling prior to leasing;  if it's throttled we can requeue
	// the pointer and back off.  A question here is enqueuing new items onto the partition
	// will reset the pointer update, leading to thrash.
	now := time.Now()
	leaseExpires := now.Add(duration).UTC().Truncate(time.Millisecond)
	leaseID, err := ulid.New(ulid.Timestamp(leaseExpires), rnd)
	if err != nil {
		return nil, fmt.Errorf("error generating id: %w", err)
	}

	keys := []string{
		q.kg.PartitionItem(),
		q.kg.PartitionIndex(),
	}
	status, err := scripts["queue/partitionLease"].Eval(
		ctx,
		q.r,
		keys,

		wid.String(),
		leaseID.String(),
		now.UnixMilli(),
		leaseExpires.Unix(),
	).Int64()
	if err != nil {
		return nil, fmt.Errorf("error leasing partition: %w", err)
	}
	switch status {
	case 0:
		return &leaseID, nil
	case 1:
		return nil, ErrPartitionNotFound
	case 2:
		return nil, ErrPartitionAlreadyLeased
	default:
		return nil, fmt.Errorf("unknown response enqueueing item: %d", status)
	}
}

// PartitionPeek returns up to PartitionSelectionMax partition items from the queue. This
// returns the indexes of partitions.
//
// If sequential is set to true this returns partitions in order from earliest to latest
// available lease times. Otherwise, this shuffles all partitions and picks partitions
// randomly, with higher priority partitions more likely to be selected.  This reduces
// lease contention amongst multiple shared-nothing workers.
func (q queue) PartitionPeek(ctx context.Context, sequential bool, until time.Time, limit int64) ([]*QueuePartition, error) {
	if limit > PartitionPeekMax {
		return nil, ErrPartitionPeekMaxExceedsLimits
	}
	if limit <= 0 {
		limit = PartitionPeekMax
	}

	encoded, err := scripts["queue/partitionPeek"].Eval(
		ctx,
		q.r,
		[]string{
			q.kg.PartitionIndex(),
			q.kg.PartitionItem(),
		},
		until.Unix(),
		limit,
	).StringSlice()
	if err != nil {
		return nil, fmt.Errorf("error peeking partition items: %w", err)
	}

	weights := []float64{}
	items := make([]*QueuePartition, len(encoded))
	for n, i := range encoded {
		item := &QueuePartition{}
		if err = json.Unmarshal([]byte(i), item); err != nil {
			return nil, fmt.Errorf("error reading partition item: %w", err)
		}
		items[n] = item
		weights = append(weights, float64(10-item.Priority))
	}

	// Some scanners run sequentially, ensuring we always work on the functions with
	// the oldest run at times in order, no matter the priority.
	if sequential {
		n := int(math.Min(float64(len(items)), float64(PartitionSelectionMax)))
		return items[0:n], nil
	}

	// We want to weighted shuffle the resulting array random.  This means that many
	// shared nothing scanners can query for outstanding partitions and receive a
	// randomized order favouring higher-priority queue items.  This reduces the chances
	// of contention when leasing.
	w := sampleuv.NewWeighted(weights, rnd)
	result := make([]*QueuePartition, len(items))
	for n := range result {
		idx, ok := w.Take()
		if !ok {
			return nil, ErrWeightedSampleRead
		}
		result[n] = items[idx]
	}

	return result, nil
}

// PartitionRequeue requeues a parition with a new score, ensuring that the partition will be
// read at (or very close to) the given time.
//
// This is used after peeking and passing all queue items onto workers; we then take the next
// unleased available time for the queue item and requeue the partition.
func (q queue) PartitionRequeue(ctx context.Context, workflowID uuid.UUID, at time.Time) error {
	keys := []string{
		q.kg.PartitionItem(),
		q.kg.PartitionIndex(),
	}
	status, err := scripts["queue/partitionRequeue"].Eval(
		ctx,
		q.r,
		keys,

		workflowID.String(),
		at.Unix(),
	).Int64()
	if err != nil {
		return fmt.Errorf("error requeueing partition: %w", err)
	}
	switch status {
	case 0:
		return nil
	case 1:
		return ErrPartitionNotFound
	default:
		return fmt.Errorf("unknown response requeueing item: %d", status)
	}
}

// PartitionDequeue removes a partition pointer from the queue.  This is used when peeking and
// receiving zero items to run.
func (q queue) PartitionDequeue(ctx context.Context, workflowID uuid.UUID, at time.Time) error {
	panic("unimplemented")
}

// PartitionReprioritize reprioritizes a workflow's QueueItems within the queue.
func (q queue) PartitionReprioritize(ctx context.Context, workflowID uuid.UUID, priority uint) error {
	if priority > PriorityMin {
		return ErrPriorityTooLow
	}
	if priority < PriorityMax {
		return ErrPriorityTooHigh
	}

	keys := []string{q.kg.PartitionItem()}
	status, err := scripts["queue/partitionReprioritize"].Eval(
		ctx,
		q.r,
		keys,
		workflowID.String(),
		priority,
	).Int64()
	if err != nil {
		return fmt.Errorf("error enqueueing item: %w", err)
	}
	switch status {
	case 0:
		return nil
	case 1:
		return ErrPartitionNotFound
	default:
		return fmt.Errorf("unknown response reprioritizing partition: %d", status)
	}
}

// LeaseSequential allows a worker to lease sequential processing.  If leased, this allows a worker
// to peek partitions sequentially.  Leasing this key works similar to leasing partitions or queue items:
//
// - If the key isn't leased, a new lease is accepted.
// - If the lease is expired, a new lease is accepted.
// - If the key is leased, you must pass in the existing lease ID to renew the lease.  Mismatches do not
//   grant a lease.
//
// This returns the new lease ID on success.
func (q queue) LeaseSequential(ctx context.Context, duration time.Duration, existingLeaseID ...*ulid.ULID) (*ulid.ULID, error) {
	if duration > SequentialLeaseMax {
		return nil, ErrSequentialLeaseExceedsLimits
	}

	now := time.Now()
	newLeaseID, err := ulid.New(ulid.Timestamp(now.Add(duration)), rnd)
	if err != nil {
		return nil, err
	}

	var existing string
	if len(existingLeaseID) > 0 && existingLeaseID[0] != nil {
		existing = existingLeaseID[0].String()
	}

	status, err := scripts["queue/sequentialLease"].Eval(
		ctx,
		q.r,
		[]string{q.kg.Sequential()},
		now.UnixMilli(),
		newLeaseID.String(),
		existing,
	).Int64()
	if err != nil {
		return nil, fmt.Errorf("error claiming sequential lease: %w", err)
	}
	switch status {
	case 0:
		return &newLeaseID, nil
	case 1:
		return nil, ErrSequentialAlreadyLeased
	default:
		return nil, fmt.Errorf("unknown response claiming sequential lease: %d", status)
	}
}
