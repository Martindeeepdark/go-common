-- DeductStock: atomically check and deduct stock
-- Prevents overselling by combining GET + check + DECRBY in one atomic operation.
--
-- KEYS[1]: stock key
-- ARGV[1]: quantity to deduct
-- Returns: remaining stock on success, -1 if key missing, -2 if insufficient

local stock = tonumber(redis.call('GET', KEYS[1]))
if stock == false or stock == nil then
    return -1
end

local qty = tonumber(ARGV[1])
if stock < qty then
    return -2
end

redis.call('DECRBY', KEYS[1], qty)
return stock - qty
