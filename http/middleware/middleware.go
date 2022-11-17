//http 中间件
package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
)

//公共函数 - 结束中间件处理，不执行后面的函数了
func ErrAbortWithResponse(errCode int, c *gin.Context) {
	global.V.Zap.Error("ErrAbortWithResponse:" + strconv.Itoa(errCode))
	httpresponse.ErrWithAllByCode(errCode, c)
	c.Abort()
}
