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
		TwinAgora.POST("cloud/record/create/acquire", v1.TwinAgoraCloudRecordCreateAcquire)
		TwinAgora.POST("cloud/record/start", v1.TwinAgoraCloudRecordStart)
		TwinAgora.GET("cloud/record/stop/:rid", v1.TwinAgoraCloudRecordStop)
		TwinAgora.GET("cloud/record/query/:rid", v1.TwinAgoraCloudRecordQuery)
		TwinAgora.GET("cloud/record/oss/files/:rid", v1.TwinAgoraCloudRecordOssFiles)

	}
}
