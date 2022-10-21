package httpmiddleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
)

//对API的访问次数、频繁，做限制,防止恶意DDos
func Limiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		global.V.Zap.Debug("middle Limiter start:")
		//fmt.Println("RateMiddleware pref")
		// 如果ip请求连接数在10秒内超过N次，返回429并抛出error
		maxTimes := global.C.Http.ReqLimitTimes
		if maxTimes > 0 {
			//N秒允许访问多少次
			if !LimiterAllow(c.ClientIP(), maxTimes, 10) {
				err := errors.New("too many requests")
				global.V.Zap.Error("RateMiddleware", zap.Any("err", err))
				httpresponse.FailWithMessage(err.Error(), c)
				c.Abort()
				return
			}
		}
		//global.V.Zap.Debug("middle Limiter finish.")
		c.Next()
		//fmt.Println("RateMiddleware after")
	}
}

// 通过redis的value判断第几次访问并返回是否允许访问
func LimiterAllow(ip string, maxTimes int, second int) bool {
	element, _ := global.V.Redis.GetElementByIndex("limiter", ip)
	nowTimes, err := global.V.Redis.Get(element)

	if err == redis.Nil {
		global.V.Redis.SetEX(element, strconv.Itoa(1), second)
		return true
	} else if err != nil {
		return false
	} else {
		t, _ := strconv.Atoi(nowTimes)
		if t >= maxTimes {
			return false
		}
		global.V.Redis.Incr(element)
		return true
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

}
