package httpmiddleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
)

//限制访问
func RateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果ip请求连接数在两秒内超过5次，返回429并抛出error
		ltime := global.C.Http.ReqLimitTimes
		if ltime <= 0 {
			ltime = 100
		}
		if !LimiterAllow(c.ClientIP(), int64(ltime), 2*time.Second) {
			err := errors.New("too many requests")
			global.V.Zap.Error("RateMiddleware", zap.Any("err", err))
			httpresponse.FailWithMessage(err.Error(), c)
			c.Abort()
			return
		}
		c.Next()
	}
}


// 通过redis的value判断第几次访问并返回是否允许访问
func LimiterAllow(key string, events int64, per time.Duration) bool {
	//curr := global.GVA_REDIS.LLen(key).Val()
	//if curr >= events {
	//	return false
	//}
	//
	//if v := global.GVA_REDIS.Exists(key).Val(); v == 0 {
	//	pipe := global.GVA_REDIS.TxPipeline()
	//	pipe.RPush(key, key)
	//	//设置过期时间
	//	pipe.Expire(key, per)
	//	_, _ = pipe.Exec()
	//} else {
	//	global.GVA_REDIS.RPushX(key, key)
	//}

	return true
}
