// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package sqlc

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	ulid "github.com/oklog/ulid/v2"
)

const deleteApp = `-- name: DeleteApp :exec
UPDATE apps SET deleted_at = NOW() WHERE id = ?
`

func (q *Queries) DeleteApp(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteApp, id)
	return err
}

const deleteFunctionsByAppID = `-- name: DeleteFunctionsByAppID :exec
DELETE FROM functions WHERE app_id = ?
`

func (q *Queries) DeleteFunctionsByAppID(ctx context.Context, appID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteFunctionsByAppID, appID)
	return err
}

const deleteFunctionsByIDs = `-- name: DeleteFunctionsByIDs :exec
DELETE FROM functions WHERE id IN (/*SLICE:ids*/?)
`

func (q *Queries) DeleteFunctionsByIDs(ctx context.Context, ids []uuid.UUID) error {
	query := deleteFunctionsByIDs
	var queryParams []interface{}
	if len(ids) > 0 {
		for _, v := range ids {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:ids*/?", strings.Repeat(",?", len(ids))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:ids*/?", "NULL", 1)
	}
	_, err := q.db.ExecContext(ctx, query, queryParams...)
	return err
}

const getAllApps = `-- name: GetAllApps :many
SELECT id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url FROM apps
`

func (q *Queries) GetAllApps(ctx context.Context) ([]*App, error) {
	rows, err := q.db.QueryContext(ctx, getAllApps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*App
	for rows.Next() {
		var i App
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.SdkLanguage,
			&i.SdkVersion,
			&i.Framework,
			&i.Metadata,
			&i.Status,
			&i.Error,
			&i.Checksum,
			&i.CreatedAt,
			&i.DeletedAt,
			&i.Url,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getApp = `-- name: GetApp :one
SELECT id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url FROM apps WHERE id = ?
`

func (q *Queries) GetApp(ctx context.Context, id uuid.UUID) (*App, error) {
	row := q.db.QueryRowContext(ctx, getApp, id)
	var i App
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SdkLanguage,
		&i.SdkVersion,
		&i.Framework,
		&i.Metadata,
		&i.Status,
		&i.Error,
		&i.Checksum,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.Url,
	)
	return &i, err
}

const getAppByChecksum = `-- name: GetAppByChecksum :one
SELECT id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url FROM apps WHERE checksum = ? LIMIT 1
`

func (q *Queries) GetAppByChecksum(ctx context.Context, checksum string) (*App, error) {
	row := q.db.QueryRowContext(ctx, getAppByChecksum, checksum)
	var i App
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SdkLanguage,
		&i.SdkVersion,
		&i.Framework,
		&i.Metadata,
		&i.Status,
		&i.Error,
		&i.Checksum,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.Url,
	)
	return &i, err
}

const getAppByID = `-- name: GetAppByID :one
SELECT id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url FROM apps WHERE id = ? LIMIT 1
`

func (q *Queries) GetAppByID(ctx context.Context, id uuid.UUID) (*App, error) {
	row := q.db.QueryRowContext(ctx, getAppByID, id)
	var i App
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SdkLanguage,
		&i.SdkVersion,
		&i.Framework,
		&i.Metadata,
		&i.Status,
		&i.Error,
		&i.Checksum,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.Url,
	)
	return &i, err
}

const getAppByURL = `-- name: GetAppByURL :one
SELECT id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url FROM apps WHERE url = ? LIMIT 1
`

func (q *Queries) GetAppByURL(ctx context.Context, url string) (*App, error) {
	row := q.db.QueryRowContext(ctx, getAppByURL, url)
	var i App
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SdkLanguage,
		&i.SdkVersion,
		&i.Framework,
		&i.Metadata,
		&i.Status,
		&i.Error,
		&i.Checksum,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.Url,
	)
	return &i, err
}

const getAppFunctions = `-- name: GetAppFunctions :many
SELECT id, app_id, name, slug, config, created_at FROM functions WHERE app_id = ?
`

func (q *Queries) GetAppFunctions(ctx context.Context, appID uuid.UUID) ([]*Function, error) {
	rows, err := q.db.QueryContext(ctx, getAppFunctions, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Function
	for rows.Next() {
		var i Function
		if err := rows.Scan(
			&i.ID,
			&i.AppID,
			&i.Name,
			&i.Slug,
			&i.Config,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAppFunctionsBySlug = `-- name: GetAppFunctionsBySlug :many
SELECT functions.id, functions.app_id, functions.name, functions.slug, functions.config, functions.created_at FROM functions JOIN apps ON apps.id = functions.app_id WHERE apps.name = ?
`

func (q *Queries) GetAppFunctionsBySlug(ctx context.Context, name string) ([]*Function, error) {
	rows, err := q.db.QueryContext(ctx, getAppFunctionsBySlug, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Function
	for rows.Next() {
		var i Function
		if err := rows.Scan(
			&i.ID,
			&i.AppID,
			&i.Name,
			&i.Slug,
			&i.Config,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getApps = `-- name: GetApps :many
SELECT id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url FROM apps WHERE deleted_at IS NULL
`

func (q *Queries) GetApps(ctx context.Context) ([]*App, error) {
	rows, err := q.db.QueryContext(ctx, getApps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*App
	for rows.Next() {
		var i App
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.SdkLanguage,
			&i.SdkVersion,
			&i.Framework,
			&i.Metadata,
			&i.Status,
			&i.Error,
			&i.Checksum,
			&i.CreatedAt,
			&i.DeletedAt,
			&i.Url,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventBatchByRunID = `-- name: GetEventBatchByRunID :one
SELECT id, account_id, workspace_id, app_id, workflow_id, run_id, started_at, executed_at, event_ids FROM event_batches WHERE run_id = ?
`

func (q *Queries) GetEventBatchByRunID(ctx context.Context, runID ulid.ULID) (*EventBatch, error) {
	row := q.db.QueryRowContext(ctx, getEventBatchByRunID, runID)
	var i EventBatch
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.WorkspaceID,
		&i.AppID,
		&i.WorkflowID,
		&i.RunID,
		&i.StartedAt,
		&i.ExecutedAt,
		&i.EventIds,
	)
	return &i, err
}

const getEventByInternalID = `-- name: GetEventByInternalID :one
SELECT internal_id, account_id, workspace_id, source, source_id, received_at, event_id, event_name, event_data, event_user, event_v, event_ts FROM events WHERE internal_id = ?
`

func (q *Queries) GetEventByInternalID(ctx context.Context, internalID ulid.ULID) (*Event, error) {
	row := q.db.QueryRowContext(ctx, getEventByInternalID, internalID)
	var i Event
	err := row.Scan(
		&i.InternalID,
		&i.AccountID,
		&i.WorkspaceID,
		&i.Source,
		&i.SourceID,
		&i.ReceivedAt,
		&i.EventID,
		&i.EventName,
		&i.EventData,
		&i.EventUser,
		&i.EventV,
		&i.EventTs,
	)
	return &i, err
}

const getEventsByInternalIDs = `-- name: GetEventsByInternalIDs :many
SELECT internal_id, account_id, workspace_id, source, source_id, received_at, event_id, event_name, event_data, event_user, event_v, event_ts FROM events WHERE internal_id IN (/*SLICE:ids*/?)
`

func (q *Queries) GetEventsByInternalIDs(ctx context.Context, ids []ulid.ULID) ([]*Event, error) {
	query := getEventsByInternalIDs
	var queryParams []interface{}
	if len(ids) > 0 {
		for _, v := range ids {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:ids*/?", strings.Repeat(",?", len(ids))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:ids*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.InternalID,
			&i.AccountID,
			&i.WorkspaceID,
			&i.Source,
			&i.SourceID,
			&i.ReceivedAt,
			&i.EventID,
			&i.EventName,
			&i.EventData,
			&i.EventUser,
			&i.EventV,
			&i.EventTs,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventsTimebound = `-- name: GetEventsTimebound :many
SELECT DISTINCT e.internal_id, e.account_id, e.workspace_id, e.source, e.source_id, e.received_at, e.event_id, e.event_name, e.event_data, e.event_user, e.event_v, e.event_ts
FROM events AS e
LEFT OUTER JOIN function_runs AS r ON r.event_id = e.internal_id
WHERE
	e.received_at > ?
	AND e.received_at <= ?
	AND (
		-- Include internal events that triggered a run (e.g. an onFailure
		-- handler)
		r.run_id IS NOT NULL

		-- Optionally include internal events that did not trigger a run. It'd
		-- be better to use a boolean param instead of a string param but sqlc
		-- keeps making @include_internal a string.
		OR CASE WHEN e.event_name LIKE 'inngest/%' THEN 'true' ELSE 'false' END = ?
	)
ORDER BY e.received_at DESC
LIMIT ?
`

type GetEventsTimeboundParams struct {
	After           time.Time
	Before          time.Time
	IncludeInternal string
	Limit           int64
}

func (q *Queries) GetEventsTimebound(ctx context.Context, arg GetEventsTimeboundParams) ([]*Event, error) {
	rows, err := q.db.QueryContext(ctx, getEventsTimebound,
		arg.After,
		arg.Before,
		arg.IncludeInternal,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.InternalID,
			&i.AccountID,
			&i.WorkspaceID,
			&i.Source,
			&i.SourceID,
			&i.ReceivedAt,
			&i.EventID,
			&i.EventName,
			&i.EventData,
			&i.EventUser,
			&i.EventV,
			&i.EventTs,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFunctionByID = `-- name: GetFunctionByID :one
SELECT id, app_id, name, slug, config, created_at FROM functions WHERE id = ?
`

func (q *Queries) GetFunctionByID(ctx context.Context, id uuid.UUID) (*Function, error) {
	row := q.db.QueryRowContext(ctx, getFunctionByID, id)
	var i Function
	err := row.Scan(
		&i.ID,
		&i.AppID,
		&i.Name,
		&i.Slug,
		&i.Config,
		&i.CreatedAt,
	)
	return &i, err
}

const getFunctionBySlug = `-- name: GetFunctionBySlug :one
SELECT id, app_id, name, slug, config, created_at FROM functions WHERE slug = ?
`

func (q *Queries) GetFunctionBySlug(ctx context.Context, slug string) (*Function, error) {
	row := q.db.QueryRowContext(ctx, getFunctionBySlug, slug)
	var i Function
	err := row.Scan(
		&i.ID,
		&i.AppID,
		&i.Name,
		&i.Slug,
		&i.Config,
		&i.CreatedAt,
	)
	return &i, err
}

const getFunctionRun = `-- name: GetFunctionRun :one
SELECT function_runs.run_id, function_runs.run_started_at, function_runs.function_id, function_runs.function_version, function_runs.trigger_type, function_runs.event_id, function_runs.batch_id, function_runs.original_run_id, function_runs.cron, function_finishes.run_id, function_finishes.status, function_finishes.output, function_finishes.completed_step_count, function_finishes.created_at
  FROM function_runs
  LEFT JOIN function_finishes ON function_finishes.run_id = function_runs.run_id
  WHERE function_runs.run_id = ?1
`

type GetFunctionRunRow struct {
	FunctionRun    FunctionRun
	FunctionFinish FunctionFinish
}

func (q *Queries) GetFunctionRun(ctx context.Context, runID ulid.ULID) (*GetFunctionRunRow, error) {
	row := q.db.QueryRowContext(ctx, getFunctionRun, runID)
	var i GetFunctionRunRow
	err := row.Scan(
		&i.FunctionRun.RunID,
		&i.FunctionRun.RunStartedAt,
		&i.FunctionRun.FunctionID,
		&i.FunctionRun.FunctionVersion,
		&i.FunctionRun.TriggerType,
		&i.FunctionRun.EventID,
		&i.FunctionRun.BatchID,
		&i.FunctionRun.OriginalRunID,
		&i.FunctionRun.Cron,
		&i.FunctionFinish.RunID,
		&i.FunctionFinish.Status,
		&i.FunctionFinish.Output,
		&i.FunctionFinish.CompletedStepCount,
		&i.FunctionFinish.CreatedAt,
	)
	return &i, err
}

const getFunctionRunFinishesByRunIDs = `-- name: GetFunctionRunFinishesByRunIDs :many
SELECT run_id, status, output, completed_step_count, created_at FROM function_finishes WHERE run_id IN (/*SLICE:run_ids*/?)
`

func (q *Queries) GetFunctionRunFinishesByRunIDs(ctx context.Context, runIds []ulid.ULID) ([]*FunctionFinish, error) {
	query := getFunctionRunFinishesByRunIDs
	var queryParams []interface{}
	if len(runIds) > 0 {
		for _, v := range runIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:run_ids*/?", strings.Repeat(",?", len(runIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:run_ids*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*FunctionFinish
	for rows.Next() {
		var i FunctionFinish
		if err := rows.Scan(
			&i.RunID,
			&i.Status,
			&i.Output,
			&i.CompletedStepCount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFunctionRunHistory = `-- name: GetFunctionRunHistory :many
SELECT id, created_at, run_started_at, function_id, function_version, run_id, event_id, batch_id, group_id, idempotency_key, type, attempt, latency_ms, step_name, step_id, url, cancel_request, sleep, wait_for_event, wait_result, invoke_function, invoke_function_result, result FROM history WHERE run_id = ? ORDER BY created_at ASC
`

func (q *Queries) GetFunctionRunHistory(ctx context.Context, runID ulid.ULID) ([]*History, error) {
	rows, err := q.db.QueryContext(ctx, getFunctionRunHistory, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*History
	for rows.Next() {
		var i History
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.RunStartedAt,
			&i.FunctionID,
			&i.FunctionVersion,
			&i.RunID,
			&i.EventID,
			&i.BatchID,
			&i.GroupID,
			&i.IdempotencyKey,
			&i.Type,
			&i.Attempt,
			&i.LatencyMs,
			&i.StepName,
			&i.StepID,
			&i.Url,
			&i.CancelRequest,
			&i.Sleep,
			&i.WaitForEvent,
			&i.WaitResult,
			&i.InvokeFunction,
			&i.InvokeFunctionResult,
			&i.Result,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFunctionRunsFromEvents = `-- name: GetFunctionRunsFromEvents :many
SELECT function_runs.run_id, function_runs.run_started_at, function_runs.function_id, function_runs.function_version, function_runs.trigger_type, function_runs.event_id, function_runs.batch_id, function_runs.original_run_id, function_runs.cron, function_finishes.run_id, function_finishes.status, function_finishes.output, function_finishes.completed_step_count, function_finishes.created_at FROM function_runs
LEFT JOIN function_finishes ON function_finishes.run_id = function_runs.run_id
WHERE function_runs.event_id IN (/*SLICE:event_ids*/?)
`

type GetFunctionRunsFromEventsRow struct {
	FunctionRun    FunctionRun
	FunctionFinish FunctionFinish
}

func (q *Queries) GetFunctionRunsFromEvents(ctx context.Context, eventIds []ulid.ULID) ([]*GetFunctionRunsFromEventsRow, error) {
	query := getFunctionRunsFromEvents
	var queryParams []interface{}
	if len(eventIds) > 0 {
		for _, v := range eventIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:event_ids*/?", strings.Repeat(",?", len(eventIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:event_ids*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetFunctionRunsFromEventsRow
	for rows.Next() {
		var i GetFunctionRunsFromEventsRow
		if err := rows.Scan(
			&i.FunctionRun.RunID,
			&i.FunctionRun.RunStartedAt,
			&i.FunctionRun.FunctionID,
			&i.FunctionRun.FunctionVersion,
			&i.FunctionRun.TriggerType,
			&i.FunctionRun.EventID,
			&i.FunctionRun.BatchID,
			&i.FunctionRun.OriginalRunID,
			&i.FunctionRun.Cron,
			&i.FunctionFinish.RunID,
			&i.FunctionFinish.Status,
			&i.FunctionFinish.Output,
			&i.FunctionFinish.CompletedStepCount,
			&i.FunctionFinish.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFunctionRunsTimebound = `-- name: GetFunctionRunsTimebound :many
SELECT function_runs.run_id, function_runs.run_started_at, function_runs.function_id, function_runs.function_version, function_runs.trigger_type, function_runs.event_id, function_runs.batch_id, function_runs.original_run_id, function_runs.cron, function_finishes.run_id, function_finishes.status, function_finishes.output, function_finishes.completed_step_count, function_finishes.created_at FROM function_runs
LEFT JOIN function_finishes ON function_finishes.run_id = function_runs.run_id
WHERE function_runs.run_started_at > ? AND function_runs.run_started_at <= ?
ORDER BY function_runs.run_started_at DESC
LIMIT ?
`

type GetFunctionRunsTimeboundParams struct {
	After  time.Time
	Before time.Time
	Limit  int64
}

type GetFunctionRunsTimeboundRow struct {
	FunctionRun    FunctionRun
	FunctionFinish FunctionFinish
}

func (q *Queries) GetFunctionRunsTimebound(ctx context.Context, arg GetFunctionRunsTimeboundParams) ([]*GetFunctionRunsTimeboundRow, error) {
	rows, err := q.db.QueryContext(ctx, getFunctionRunsTimebound, arg.After, arg.Before, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetFunctionRunsTimeboundRow
	for rows.Next() {
		var i GetFunctionRunsTimeboundRow
		if err := rows.Scan(
			&i.FunctionRun.RunID,
			&i.FunctionRun.RunStartedAt,
			&i.FunctionRun.FunctionID,
			&i.FunctionRun.FunctionVersion,
			&i.FunctionRun.TriggerType,
			&i.FunctionRun.EventID,
			&i.FunctionRun.BatchID,
			&i.FunctionRun.OriginalRunID,
			&i.FunctionRun.Cron,
			&i.FunctionFinish.RunID,
			&i.FunctionFinish.Status,
			&i.FunctionFinish.Output,
			&i.FunctionFinish.CompletedStepCount,
			&i.FunctionFinish.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFunctions = `-- name: GetFunctions :many
SELECT id, app_id, name, slug, config, created_at FROM functions
`

func (q *Queries) GetFunctions(ctx context.Context) ([]*Function, error) {
	rows, err := q.db.QueryContext(ctx, getFunctions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Function
	for rows.Next() {
		var i Function
		if err := rows.Scan(
			&i.ID,
			&i.AppID,
			&i.Name,
			&i.Slug,
			&i.Config,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const hardDeleteApp = `-- name: HardDeleteApp :exec
DELETE FROM apps WHERE id = ?
`

func (q *Queries) HardDeleteApp(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, hardDeleteApp, id)
	return err
}

const insertApp = `-- name: InsertApp :one
INSERT INTO apps
	(id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, url) VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url
`

type InsertAppParams struct {
	ID          uuid.UUID
	Name        string
	SdkLanguage string
	SdkVersion  string
	Framework   sql.NullString
	Metadata    string
	Status      string
	Error       sql.NullString
	Checksum    string
	Url         string
}

func (q *Queries) InsertApp(ctx context.Context, arg InsertAppParams) (*App, error) {
	row := q.db.QueryRowContext(ctx, insertApp,
		arg.ID,
		arg.Name,
		arg.SdkLanguage,
		arg.SdkVersion,
		arg.Framework,
		arg.Metadata,
		arg.Status,
		arg.Error,
		arg.Checksum,
		arg.Url,
	)
	var i App
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SdkLanguage,
		&i.SdkVersion,
		&i.Framework,
		&i.Metadata,
		&i.Status,
		&i.Error,
		&i.Checksum,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.Url,
	)
	return &i, err
}

const insertEvent = `-- name: InsertEvent :exec

INSERT INTO events
	(internal_id, received_at, event_id, event_name, event_data, event_user, event_v, event_ts) VALUES
	(?, ?, ?, ?, ?, ?, ?, ?)
`

type InsertEventParams struct {
	InternalID ulid.ULID
	ReceivedAt time.Time
	EventID    string
	EventName  string
	EventData  string
	EventUser  string
	EventV     sql.NullString
	EventTs    time.Time
}

// Events
func (q *Queries) InsertEvent(ctx context.Context, arg InsertEventParams) error {
	_, err := q.db.ExecContext(ctx, insertEvent,
		arg.InternalID,
		arg.ReceivedAt,
		arg.EventID,
		arg.EventName,
		arg.EventData,
		arg.EventUser,
		arg.EventV,
		arg.EventTs,
	)
	return err
}

const insertEventBatch = `-- name: InsertEventBatch :exec
INSERT INTO event_batches
	(id, account_id, workspace_id, app_id, workflow_id, run_id, started_at, executed_at, event_ids) VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type InsertEventBatchParams struct {
	ID          ulid.ULID
	AccountID   uuid.UUID
	WorkspaceID uuid.UUID
	AppID       uuid.UUID
	WorkflowID  uuid.UUID
	RunID       ulid.ULID
	StartedAt   time.Time
	ExecutedAt  time.Time
	EventIds    []byte
}

func (q *Queries) InsertEventBatch(ctx context.Context, arg InsertEventBatchParams) error {
	_, err := q.db.ExecContext(ctx, insertEventBatch,
		arg.ID,
		arg.AccountID,
		arg.WorkspaceID,
		arg.AppID,
		arg.WorkflowID,
		arg.RunID,
		arg.StartedAt,
		arg.ExecutedAt,
		arg.EventIds,
	)
	return err
}

const insertFunction = `-- name: InsertFunction :one


INSERT INTO functions
	(id, app_id, name, slug, config, created_at) VALUES
	(?, ?, ?, ?, ?, ?) RETURNING id, app_id, name, slug, config, created_at
`

type InsertFunctionParams struct {
	ID        uuid.UUID
	AppID     uuid.UUID
	Name      string
	Slug      string
	Config    string
	CreatedAt time.Time
}

// functions
//
// note - this is very basic right now.
func (q *Queries) InsertFunction(ctx context.Context, arg InsertFunctionParams) (*Function, error) {
	row := q.db.QueryRowContext(ctx, insertFunction,
		arg.ID,
		arg.AppID,
		arg.Name,
		arg.Slug,
		arg.Config,
		arg.CreatedAt,
	)
	var i Function
	err := row.Scan(
		&i.ID,
		&i.AppID,
		&i.Name,
		&i.Slug,
		&i.Config,
		&i.CreatedAt,
	)
	return &i, err
}

const insertFunctionFinish = `-- name: InsertFunctionFinish :exec
INSERT INTO function_finishes
	(run_id, status, output, completed_step_count, created_at) VALUES
	(?, ?, ?, ?, ?)
`

type InsertFunctionFinishParams struct {
	RunID              ulid.ULID
	Status             sql.NullString
	Output             sql.NullString
	CompletedStepCount sql.NullInt64
	CreatedAt          sql.NullTime
}

func (q *Queries) InsertFunctionFinish(ctx context.Context, arg InsertFunctionFinishParams) error {
	_, err := q.db.ExecContext(ctx, insertFunctionFinish,
		arg.RunID,
		arg.Status,
		arg.Output,
		arg.CompletedStepCount,
		arg.CreatedAt,
	)
	return err
}

const insertFunctionRun = `-- name: InsertFunctionRun :exec

INSERT INTO function_runs
	(run_id, run_started_at, function_id, function_version, trigger_type, event_id, batch_id, original_run_id, cron) VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type InsertFunctionRunParams struct {
	RunID           ulid.ULID
	RunStartedAt    time.Time
	FunctionID      uuid.UUID
	FunctionVersion int64
	TriggerType     string
	EventID         ulid.ULID
	BatchID         ulid.ULID
	OriginalRunID   ulid.ULID
	Cron            sql.NullString
}

// function runs
func (q *Queries) InsertFunctionRun(ctx context.Context, arg InsertFunctionRunParams) error {
	_, err := q.db.ExecContext(ctx, insertFunctionRun,
		arg.RunID,
		arg.RunStartedAt,
		arg.FunctionID,
		arg.FunctionVersion,
		arg.TriggerType,
		arg.EventID,
		arg.BatchID,
		arg.OriginalRunID,
		arg.Cron,
	)
	return err
}

const insertHistory = `-- name: InsertHistory :exec

INSERT INTO history
	(id, created_at, run_started_at, function_id, function_version, run_id, event_id, batch_id, group_id, idempotency_key, type, attempt, latency_ms, step_name, step_id, url, cancel_request, sleep, wait_for_event, wait_result, invoke_function, invoke_function_result, result) VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type InsertHistoryParams struct {
	ID                   ulid.ULID
	CreatedAt            time.Time
	RunStartedAt         time.Time
	FunctionID           uuid.UUID
	FunctionVersion      int64
	RunID                ulid.ULID
	EventID              ulid.ULID
	BatchID              ulid.ULID
	GroupID              sql.NullString
	IdempotencyKey       string
	Type                 string
	Attempt              int64
	LatencyMs            sql.NullInt64
	StepName             sql.NullString
	StepID               sql.NullString
	Url                  sql.NullString
	CancelRequest        sql.NullString
	Sleep                sql.NullString
	WaitForEvent         sql.NullString
	WaitResult           sql.NullString
	InvokeFunction       sql.NullString
	InvokeFunctionResult sql.NullString
	Result               sql.NullString
}

// History
func (q *Queries) InsertHistory(ctx context.Context, arg InsertHistoryParams) error {
	_, err := q.db.ExecContext(ctx, insertHistory,
		arg.ID,
		arg.CreatedAt,
		arg.RunStartedAt,
		arg.FunctionID,
		arg.FunctionVersion,
		arg.RunID,
		arg.EventID,
		arg.BatchID,
		arg.GroupID,
		arg.IdempotencyKey,
		arg.Type,
		arg.Attempt,
		arg.LatencyMs,
		arg.StepName,
		arg.StepID,
		arg.Url,
		arg.CancelRequest,
		arg.Sleep,
		arg.WaitForEvent,
		arg.WaitResult,
		arg.InvokeFunction,
		arg.InvokeFunctionResult,
		arg.Result,
	)
	return err
}

const updateAppError = `-- name: UpdateAppError :one
UPDATE apps SET error = ? WHERE id = ? RETURNING id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url
`

type UpdateAppErrorParams struct {
	Error sql.NullString
	ID    uuid.UUID
}

func (q *Queries) UpdateAppError(ctx context.Context, arg UpdateAppErrorParams) (*App, error) {
	row := q.db.QueryRowContext(ctx, updateAppError, arg.Error, arg.ID)
	var i App
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SdkLanguage,
		&i.SdkVersion,
		&i.Framework,
		&i.Metadata,
		&i.Status,
		&i.Error,
		&i.Checksum,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.Url,
	)
	return &i, err
}

const updateAppURL = `-- name: UpdateAppURL :one
UPDATE apps SET url = ? WHERE id = ? RETURNING id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url
`

type UpdateAppURLParams struct {
	Url string
	ID  uuid.UUID
}

func (q *Queries) UpdateAppURL(ctx context.Context, arg UpdateAppURLParams) (*App, error) {
	row := q.db.QueryRowContext(ctx, updateAppURL, arg.Url, arg.ID)
	var i App
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SdkLanguage,
		&i.SdkVersion,
		&i.Framework,
		&i.Metadata,
		&i.Status,
		&i.Error,
		&i.Checksum,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.Url,
	)
	return &i, err
}

const updateFunctionConfig = `-- name: UpdateFunctionConfig :one
UPDATE functions SET config = ? WHERE id = ? RETURNING id, app_id, name, slug, config, created_at
`

type UpdateFunctionConfigParams struct {
	Config string
	ID     uuid.UUID
}

func (q *Queries) UpdateFunctionConfig(ctx context.Context, arg UpdateFunctionConfigParams) (*Function, error) {
	row := q.db.QueryRowContext(ctx, updateFunctionConfig, arg.Config, arg.ID)
	var i Function
	err := row.Scan(
		&i.ID,
		&i.AppID,
		&i.Name,
		&i.Slug,
		&i.Config,
		&i.CreatedAt,
	)
	return &i, err
}

const workspaceEvents = `-- name: WorkspaceEvents :many
SELECT internal_id, account_id, workspace_id, source, source_id, received_at, event_id, event_name, event_data, event_user, event_v, event_ts FROM events WHERE internal_id < ? AND received_at <= ? AND received_at >= ? ORDER BY internal_id DESC LIMIT ?
`

type WorkspaceEventsParams struct {
	Cursor ulid.ULID
	Before time.Time
	After  time.Time
	Limit  int64
}

func (q *Queries) WorkspaceEvents(ctx context.Context, arg WorkspaceEventsParams) ([]*Event, error) {
	rows, err := q.db.QueryContext(ctx, workspaceEvents,
		arg.Cursor,
		arg.Before,
		arg.After,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.InternalID,
			&i.AccountID,
			&i.WorkspaceID,
			&i.Source,
			&i.SourceID,
			&i.ReceivedAt,
			&i.EventID,
			&i.EventName,
			&i.EventData,
			&i.EventUser,
			&i.EventV,
			&i.EventTs,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const workspaceNamedEvents = `-- name: WorkspaceNamedEvents :many
SELECT internal_id, account_id, workspace_id, source, source_id, received_at, event_id, event_name, event_data, event_user, event_v, event_ts FROM events WHERE internal_id < ? AND received_at <= ? AND received_at >= ? AND event_name = ? ORDER BY internal_id DESC LIMIT ?
`

type WorkspaceNamedEventsParams struct {
	Cursor ulid.ULID
	Before time.Time
	After  time.Time
	Name   string
	Limit  int64
}

func (q *Queries) WorkspaceNamedEvents(ctx context.Context, arg WorkspaceNamedEventsParams) ([]*Event, error) {
	rows, err := q.db.QueryContext(ctx, workspaceNamedEvents,
		arg.Cursor,
		arg.Before,
		arg.After,
		arg.Name,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.InternalID,
			&i.AccountID,
			&i.WorkspaceID,
			&i.Source,
			&i.SourceID,
			&i.ReceivedAt,
			&i.EventID,
			&i.EventName,
			&i.EventData,
			&i.EventUser,
			&i.EventV,
			&i.EventTs,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
