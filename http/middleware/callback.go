package httpmiddleware

import (
	"github.com/gin-gonic/gin"
)

//跨域请求
func Callback() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
