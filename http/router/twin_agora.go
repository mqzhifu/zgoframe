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
		TwinAgora.POST("rtc/get/cloud/record/acquire", v1.TwinAgoraRTCGetCloudRecordAcquire)
		TwinAgora.POST("rtc/cloud/record/start", v1.TwinAgoraRTCCloudRecordStart)
		TwinAgora.POST("rtc/cloud/record/stop", v1.TwinAgoraRTCCloudRecordStop)
		TwinAgora.POST("rtc/cloud/record/query", v1.TwinAgoraRTCCloudRecordQuery)

	}
}
