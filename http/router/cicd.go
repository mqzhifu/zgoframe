package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func InitCicdRouter(Router *gin.RouterGroup) {
	CicdRouter := Router.Group("cicd")
	{
		//
		CicdRouter.GET("superVisor/list", v1.CicdSuperVisorList)

		CicdRouter.GET("service/list", v1.CicdServiceList)

		CicdRouter.GET("server/list", v1.CicdServerList)
		//部署一个服务
		CicdRouter.POST("service/deploy", v1.CicdServiceDeploy)

		CicdRouter.GET("service/deploy/all", v1.CicdServiceDeploy)

		CicdRouter.GET("ping", v1.CicdPing)

	}
}
