local removeToken = function(merchantId, email)
    local ns = "token:admin:" .. merchantId .. ":"
    local token = redis.call("GET", ns .. email)
    if token ~= false then
        redis.call("DEL", ns .. token)
        redis.call("DEL", ns .. email)
    end
    return true
end

return removeToken(KEYS[1], ARGV[1])