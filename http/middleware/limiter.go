package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"zgoframe/core/global"
	"errors"
	httpresponse "zgoframe/http/response"
)

//限制访问
func RateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果ip请求连接数在10秒内超过N次，返回429并抛出error
		ltime := global.C.Http.ReqLimitTimes
		if ltime > 0 {
			if !LimiterAllow(c.ClientIP(), int64(ltime), 10*time.Second) {
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
func LimiterAllow(key string, events int64, per time.Duration) bool {
	curr := global.V.Redis.LLen(key).Val()
	if curr >= events {
		return false
	}

	if v := global.V.Redis.Exists(key).Val(); v == 0 {
		pipe := global.V.Redis.TxPipeline()
		pipe.RPush(key, key)
		//设置过期时间
		pipe.Expire(key, per)
		_, _ = pipe.Exec()
	} else {
		global.V.Redis.RPushX(key, key)
	}

	return true
}