package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func Cicd(Router *gin.RouterGroup) {
	CicdRouter := Router.Group("cicd")
	{
		CicdRouter.GET("superVisor/list", v1.CicdSuperVisorList)                          //每台机器上的 superVisor 状态
		CicdRouter.GET("service/list", v1.CicdServiceList)                                //属于基础数据，查看下 当前所有项目的状态
		CicdRouter.GET("server/list", v1.CicdServerList)                                  // 属于基础数据，获取所有服务器列表，并做ping，确定状态
		CicdRouter.POST("service/deploy", v1.CicdServiceDeploy)                           //部署一个服务
		CicdRouter.GET("publish/list", v1.CicdPublishList)                                // 查看 项目的已部署 日志 列表
		CicdRouter.GET("service/publish/:id/:flag", v1.CicdServicePublish)                //已部署好的项目，正式 发布
		CicdRouter.POST("superVisor/process", v1.CicdSuperVisorProcess)                   //本机操作远端的 superVisor ，管理进程
		CicdRouter.GET("local/all/server/service/list", v1.CicdLocalAllServerServiceList) //查看 所有项目的 可部署 列表
		//CicdRouter.POST("local/sync/target", v1.CicdLocalSyncTarget)

	}
}
