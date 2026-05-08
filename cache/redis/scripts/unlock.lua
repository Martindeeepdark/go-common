-- Unlock: safely release a distributed lock
-- Only the lock owner (matched by token) can release it.
--
-- KEYS[1]: lock key
-- ARGV[1]: unique token (owner identifier)
-- Returns: 1 if unlocked, 0 if not owner

local value = redis.call("get", KEYS[1])
if value == false then
    return 0
end

if value == ARGV[1] then
    redis.call("del", KEYS[1])
    return 1
end

return 0
