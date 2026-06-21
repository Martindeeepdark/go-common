-- DeductStockOrder: atomically check, deduct stock, and record order reservation.
-- Prevents overselling and duplicate deductions for the same order.
--
-- KEYS[1]: stock key (e.g. seckill:stock:{activity}:{sku})
-- KEYS[2]: reserved orders set key (e.g. seckill:reserved:{activity}:{sku})
-- ARGV[1]: quantity to deduct
-- ARGV[2]: orderNo
-- Returns: remaining stock on success, -1 if key missing, -2 if insufficient, -3 if order already reserved

-- Idempotent: if order already reserved, return current stock (deduct was done before).
if redis.call('SISMEMBER', KEYS[2], ARGV[2]) == 1 then
    local stock = tonumber(redis.call('GET', KEYS[1]))
    if stock == false or stock == nil then
        return -1
    end
    return stock
end

local stock = tonumber(redis.call('GET', KEYS[1]))
if stock == false or stock == nil then
    return -1
end

local qty = tonumber(ARGV[1])
if stock < qty then
    return -2
end

redis.call('DECRBY', KEYS[1], qty)
redis.call('SADD', KEYS[2], ARGV[2])
return stock - qty
