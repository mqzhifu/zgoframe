package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func GrabOrder(Router *gin.RouterGroup) {
	configCenterRouter := Router.Group("grab/order")
	{
		configCenterRouter.GET("get/pay/category", v1.GetPayCategory)
		configCenterRouter.GET("get/data", v1.GrabOrderGetData)
		configCenterRouter.POST("create", v1.GrabOrderCreate)
	}
}
