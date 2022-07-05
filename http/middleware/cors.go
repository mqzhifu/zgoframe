package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zgoframe/util"
)

//跨域请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {

		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		//c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id,X-Source-Type,X-Project-Id,X-Access")
		c.Header("Access-Control-Allow-Headers","*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		method := c.Request.Method
		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			util.MyPrint("http-middleware cross-domain : OPTIONS hit.")
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
