package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func CallbackRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("callback")
	{
		BaseRouter.POST("agora/cloud", v1.AgoraCallbackCloud)
		BaseRouter.POST("agora/cloud/test", v1.AgoraCallbackCloudTest)
		BaseRouter.POST("agora/rtc", v1.AgoraCallbackRTC)
		BaseRouter.POST("agora/rtc/test", v1.AgoraCallbackRTCTest)

	}
}
