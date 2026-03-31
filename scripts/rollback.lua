local stockKey = KEYS[1]
local userKey = KEYS[2]
local userID = ARGV[1]

local stock = redis.call('GET', stockKey)
if not stock then 
	return -1
end

if redis.call('SISMEMBER', userKey, userID) == 1 then
	redis.call('INCR', stockKey)
	redis.call('SREM', userKey, userID)
	return 1
else
	return 0
end