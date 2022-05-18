//http 中间件
package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	httpresponse "zgoframe/http/response"
)

func ErrAbortWithResponse(errCode int,c *gin.Context){
	httpresponse.ErrWithAllByCode(errCode,c)
	c.Abort()
}
