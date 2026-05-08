-- Extend: renew a distributed lock's TTL
-- Only the lock owner (matched by token) can extend it.
--
-- KEYS[1]: lock key
-- ARGV[1]: unique token (owner identifier)
-- ARGV[2]: new TTL in milliseconds
-- Returns: 1 if extended, 0 if not owner

local value = redis.call("get", KEYS[1])
if value == false then
    return 0
end

if value == ARGV[1] then
    redis.call("pexpire", KEYS[1], ARGV[2])
    return 1
end

return 0
