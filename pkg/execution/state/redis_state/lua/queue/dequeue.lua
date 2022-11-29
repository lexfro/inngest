--[[

Output:
  0: Successfully leased item
  1: Queue item not found

]]

local queueKey      = KEYS[1]
local queueIndexKey = KEYS[2]
local partitionKey  = KEYS[3]

local queueID = ARGV[1]

-- Fetch this item to see if it was in progress prior to deleting.
local item = cjson.decode(redis.call("HGET", queueKey, queueID))
if item == nil then
	return 1
end

redis.call("HDEL", queueKey, queueID)
redis.call("ZREM", queueIndexKey, queueID)
redis.call("HINCRBY", partitionKey, "len", -1) -- len of enqueued items decreases

if item.leaseID ~= nil then
	-- Remove total number in progress, if there's a lease.
	redis.call("HINCRBY", partitionKey, "n", -1)
end

return 0
