package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func Gateway(Router *gin.RouterGroup) {

	GatewayRouter := Router.Group("gateway")
	{
		GatewayRouter.POST("service/:service_name/:func_name", v1.GatewayService)
		GatewayRouter.GET("proto", v1.GatewayProto)
		GatewayRouter.GET("action/map", v1.ActionMap)
		GatewayRouter.GET("config", v1.GatewayConfig)
		GatewayRouter.GET("fd/list", v1.GatewayFDList)
		GatewayRouter.POST("send/msg", v1.GatewaySendMsg)
		GatewayRouter.GET("total", v1.GatewayTotal)

	}
}
