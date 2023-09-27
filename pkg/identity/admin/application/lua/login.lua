local login = function(merchantId, email, token, value, expire)
    local ns = "token:admin:" .. merchantId .. ":"
    local result = redis.call("GET", ns .. email)
    if result ~= false then
        redis.call("DEL", ns .. result)
        redis.call("DEL", ns .. email)
    end
    redis.call("SET", ns .. email, token, "EX", expire)
    redis.call("SET", ns .. token, value, "EX", expire)
    return true
end

return login(KEYS[1], ARGV[1], ARGV[2], ARGV[3], ARGV[4])