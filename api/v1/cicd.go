package v1

import (
	"github.com/gin-gonic/gin"
	httpresponse "zgoframe/http/response"
)

// @Tags Cicd
// @Summary superVisor 列表
// @Description demo
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {boolean} true "数据过长，先用bool替代"
// @Router /cicd/superVisor/list [get]
func CicdSuperVisorList(c *gin.Context) {
	//list, err := global.V.MyService.Cicd.GetSuperVisorList()
	//if err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//} else {
	//	httpresponse.OkWithAll(list, "成功", c)
	//}

}

// @Tags Cicd
// @Summary 本地编译-同步远端
// @Description 本地编译
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} util.Service
// @Router /cicd/local/all/server/service/list [get]
func CicdLocalAllServerServiceList(c *gin.Context) {
	//list, _ := global.V.MyService.Cicd.LocalAllServerServiceList()
	//httpresponse.OkWithAll(list, "成功", c)
}

// @Tags Cicd
// @Summary 服务 列表
// @Description 在当前服务器上，从<部署目录>中检索出每个服务（目录名），分析出：哪些服务~已经部署
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} util.Service
// @Router /cicd/service/list [get]
func CicdServiceList(c *gin.Context) {
	//list := global.V.MyService.Cicd.GetServiceList()
	//httpresponse.OkWithAll(list, "成功", c)
}

// @Tags Cicd
// @Summary 服务器 列表
// @Description 获取所有服务器列表，并做ping，确定状态
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} util.Server
// @Router /cicd/server/list [get]
func CicdServerList(c *gin.Context) {
	//list := global.V.MyService.Cicd.GetServerList()
	//httpresponse.OkWithAll(list, "成功", c)
}

// @Tags Cicd
// @Summary 部署/发布 列表
// @Description 部署/发布 列表
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} model.CicdPublish
// @Router /cicd/publish/list [get]
func CicdPublishList(c *gin.Context) {
	//list := global.V.MyService.Cicd.GetPublishList()
	//httpresponse.OkWithAll(list, "成功", c)
}

// @Tags Cicd
// @Summary 发布项目
// @Description 发布项目
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param id path string true "publish id"
// @Produce  application/json
// @Success 200 {bool} bool "true:成功 false:失败"
// @Router /cicd/service/publish/{id}/{flag} [get]
func CicdServicePublish(c *gin.Context) {
	//idStr := c.Param("id")
	//flagStr := c.Param("flag")
	//if idStr == "" {
	//	httpresponse.FailWithMessage("id empty 1", c)
	//	return
	//}
	//id, err := strconv.Atoi(idStr)
	//if err != nil {
	//	httpresponse.FailWithMessage("id empty 2", c)
	//	return
	//}
	//if id == 0 {
	//	httpresponse.FailWithMessage("id empty 3", c)
	//	return
	//}
	//
	//flag, _ := strconv.Atoi(flagStr)
	//err = global.V.MyService.Cicd.Deploy.Publish(id, flag)
	//if err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//} else {
	//	httpresponse.OkWithAll("bbb", "成功", c)
	//}

}

// @Tags Cicd
// @Summary 部署一个服务
// @Description demo
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.CicdDeploy true "用户信息"
// @Produce  application/json
// @Success 200 {bool} bool "true:成功 false:失败"
// @Router /cicd/service/deploy [post]
func CicdServiceDeploy(c *gin.Context) {
	//var form request.CicdDeploy
	//c.ShouldBind(&form)
	//
	//util.MyPrint("CicdServiceDeploy form:", form)
	////这里因为是HTTP连接，而后端处理一次时间接近1分钟，HTTP可能多次重复请求，开个协程
	//go global.V.MyService.Cicd.Deploy.ApiDeployOneService(form)
	//httpresponse.OkWithAll("aaaa", "成功", c)

}

// @Tags Cicd
// @Summary 操作进程
// @Description 通过superVisor
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.CicdDeploy true "用户信息"
// @Produce  application/json
// @Success 200 {bool} bool "true:成功 false:失败"
// @Router /cicd/superVisor/process [post]
func CicdSuperVisorProcess(c *gin.Context) {
	//var form request.CicdSuperVisor
	//c.ShouldBind(&form)
	//
	//util.MyPrint("CicdSuperVisorProcess form:", form)
	//err := global.V.MyService.Cicd.SuperVisorProcess(form)
	//if err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//} else {
	//	httpresponse.OkWithAll("aaaa", "成功", c)
	//}
}

// @Tags Cicd
// @Summary ping
// @Description 测试对端有没有开启服务
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {bool} bool "true:成功 false:失败"
// @Router /cicd/ping [get]
func CicdPing(c *gin.Context) {
	httpresponse.OkWithAll("aaaa", "成功", c)
}

// @Tags Cicd
// @Summary 本机部署的项目
// @Description 本机部署的项目
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /cicd/local/deploy/dir/list [get]
func CicdLocalDeployDirList(c *gin.Context) {
	//list := global.V.MyService.Cicd.GetHasDeployService()
	//httpresponse.OkWithAll(list, "成功", c)
}

// @Tags Cicd
// @Summary 同步 本机部署的项目 -> 到目标机器
// @Description scp
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.CicdSync true "用户信息"
// @Produce  application/json
// @Success 200 {bool} bool "true:成功 false:失败"
// @Router /cicd/local/sync/target [get]
func CicdLocalSyncTarget(c *gin.Context) {
	//var form request.CicdSync
	//c.ShouldBind(&form)
	//
	//list := global.V.MyService.Cicd.LocalSyncTarget(form)
	//httpresponse.OkWithAll(list, "成功", c)
}
