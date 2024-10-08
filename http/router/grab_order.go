package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func GrabOrder(Router *gin.RouterGroup) {
	configCenterRouter := Router.Group("grab/order")
	{
		configCenterRouter.GET("get/pay/category", v1.GetPayCategory)
		configCenterRouter.GET("get/base/data", v1.GrabOrderGetBaseData)
		configCenterRouter.GET("get/order/bucket/list", v1.GrabOrderBucketList)
		configCenterRouter.GET("get/user/total", v1.GrabOrderGetUserTotal)
		configCenterRouter.GET("get/user/bucket/list", v1.GrabOrderGetUserBucketAmountList)

		configCenterRouter.POST("create", v1.GrabOrderCreate)
		configCenterRouter.POST("user/open", v1.GrabOrderUserOpen)

	}
}
