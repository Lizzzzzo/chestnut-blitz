local stockKey = KEYS[1]
local userKey = KEYS[2]
local userID = ARGV[1]

if redis.call('SISMEMBER', userKey, userID) == 1 then
	return -1
end

local stock = redis.call('GET', stockKey)
if not stock or tonumber(stock) <= 0 then
	return 0
end

redis.call('DECR', stockKey)
redis.call('SADD', userKey, userID)

return 1