package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func InitTwinAgoraRouter(Router *gin.RouterGroup) {
	TwinAgora := Router.Group("twin/agora")
	{
		TwinAgora.POST("rtc/get/token", v1.TwinAgoraRTCGetToken) // 设置/修改密码
		TwinAgora.POST("rtm/get/token", v1.TwinAgoraRTMGetToken) // 设置/修改密码
	}
}

