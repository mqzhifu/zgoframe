//http 中间件
package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
)

func ErrAbortWithResponse(errCode int, c *gin.Context) {
	global.V.Zap.Error("ErrAbortWithResponse:" + strconv.Itoa(errCode))
	httpresponse.ErrWithAllByCode(errCode, c)
	c.Abort()
}
