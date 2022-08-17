package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func CallbackRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("callback")
	{
		//图形 - 验证码
		BaseRouter.POST("agora/rtc", v1.AgoraCallbackRTC)
	}
}
