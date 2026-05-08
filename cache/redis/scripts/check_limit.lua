-- CheckLimit: atomically check and increment purchase limit counter
-- Prevents over-purchasing by combining check + INCRBY in one atomic operation.
--
-- KEYS[1]: user purchase limit key
-- ARGV[1]: quantity to purchase this time
-- ARGV[2]: purchase limit (max allowed)
-- Returns: updated count on success, -1 if would exceed limit

local qty     = tonumber(ARGV[1])
local limit   = tonumber(ARGV[2])
local current = tonumber(redis.call('GET', KEYS[1]))

if current == false or current == nil then
    current = 0
end

if current + qty > limit then
    return -1
end

if current == 0 then
    redis.call('SET', KEYS[1], qty)
else
    redis.call('INCRBY', KEYS[1], qty)
end

return current + qty
