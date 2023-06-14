package httpmiddleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
	"zgoframe/core/global"
)

// 对API的访问次数、频繁，做限制,防止恶意DDos
func Limiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqIp := c.ClientIP()
		secondPeriod := 10
		// 如果ip请求连接数在10秒内超过N次，返回429并抛出error
		maxTimes := global.C.Http.ReqLimitTimes
		showMsg := "http middleware <Limiter> , ReqLimitTimes:" + strconv.Itoa(maxTimes) + " , reqIp:" + reqIp + " , secondPeriod:" + strconv.Itoa(secondPeriod)
		if maxTimes > 0 {
			//N秒允许访问多少次
			rs, nowTimes := LimiterAllow(reqIp, maxTimes, secondPeriod)
			showMsg += " , nowTimes:" + strconv.Itoa(nowTimes)
			if !rs {
				err := errors.New("too many requests")
				global.V.Zap.Error("RateMiddleware", zap.Any("err", err))
				ErrAbortWithResponse(5208, c)
				//httpresponse.FailWithMessage(err.Error(), c)
				//c.Abort()
				return
			}
		} else {
			showMsg += " , no need process."
		}

		global.V.Zap.Debug(showMsg)
		//global.V.Zap.Debug("middle Limiter finish.")
		c.Next()
		//fmt.Println("RateMiddleware after")
	}
}

// 通过redis的value判断第几次访问并返回是否允许访问
func LimiterAllow(ip string, maxTimes int, second int) (bool, int) {
	element, _ := global.V.Redis.GetElementByIndex("limiter", ip)
	nowTimesStr, err := global.V.Redis.Get(element)
	nowTimes, _ := strconv.Atoi(nowTimesStr)

	if err == redis.Nil {
		global.V.Redis.SetEX(element, strconv.Itoa(1), second)
		return true, nowTimes
	} else if err != nil {
		return false, nowTimes
	} else {
		if nowTimes >= maxTimes {
			return false, nowTimes
		}
		global.V.Redis.Incr(element)
		return true, nowTimes
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
