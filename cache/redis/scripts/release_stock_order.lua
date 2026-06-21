-- ReleaseStockOrder: atomically release stock for a specific order.
-- Idempotent: if order was not reserved (or already released), returns 0 without side effects.
--
-- KEYS[1]: stock key (e.g. seckill:stock:{activity}:{sku})
-- KEYS[2]: reserved orders set key (e.g. seckill:reserved:{activity}:{sku})
-- ARGV[1]: quantity to release
-- ARGV[2]: orderNo
-- Returns: remaining stock after release, 0 if order was not reserved (idempotent no-op)

-- Only release if this order was previously reserved.
if redis.call('SISMEMBER', KEYS[2], ARGV[2]) == 0 then
    return 0
end

local qty = tonumber(ARGV[1])
redis.call('INCRBY', KEYS[1], qty)
redis.call('SREM', KEYS[2], ARGV[2])
return tonumber(redis.call('GET', KEYS[1]))
