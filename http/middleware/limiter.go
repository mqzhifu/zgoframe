package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
	"errors"
)

//限制访问
func RateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果ip请求连接数在10秒内超过N次，返回429并抛出error
		maxTimes := global.C.Http.ReqLimitTimes
		if maxTimes > 0 {
			//N秒允许访问多少次
			if !LimiterAllow(c.ClientIP(), maxTimes, 10  ) {
				err := errors.New("too many requests")
				global.V.Zap.Error("RateMiddleware", zap.Any("err", err))
				httpresponse.FailWithMessage(err.Error(), c)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
// 通过redis的value判断第几次访问并返回是否允许访问
func LimiterAllow(ip string, maxTimes int, second int) bool {
	element , _ := global.V.Redis.GetElementByIndex("limiter",ip)
	nowTimes, err   := global.V.Redis.Get(element)
	if err != nil && err == redis.Nil{
		global.V.Redis.SetEX(element,strconv.Itoa(1),second)
	}else{
		t, _ := strconv.Atoi(nowTimes)
		if t >=  maxTimes  {
			return false
		}
		global.V.Redis.Incr(element)
	}


	//if v := global.V.Redis.(key).Val(); v == 0 {
	//	pipe := global.V.Redis.TxPipeline()
	//	pipe.RPush(key, key)
	//	//设置过期时间
	//	pipe.Expire(key, per)
	//	_, _ = pipe.Exec()
	//} else {
	//	global.V.Redis.RPushX(key, key)
	//}

	return true
}