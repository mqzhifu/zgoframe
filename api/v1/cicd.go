package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Tags Cicd
// @Summary superVisor 列表
// @Description demo
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /cicd/superVisor/list [get]
func CicdSuperVisorList(c *gin.Context) {
	list,err := global.V.MyService.Cicd.GetSuperVisorList()
	if err != nil{
		httpresponse.FailWithMessage(err.Error(),c)
	}else{
		httpresponse.OkWithDetailed(list, "成功", c)
	}

}

// @Tags Cicd
// @Summary 服务 列表
// @Description 在当前服务器上，从<部署目录>中检索出每个服务（目录名），分析出：哪些服务~已经部署
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /cicd/service/list [get]
func CicdServiceList(c *gin.Context) {
	list := global.V.MyService.Cicd.GetServiceList()
	httpresponse.OkWithDetailed(list, "成功", c)
}

// @Tags Cicd
// @Summary 服务器 列表
// @Description 获取所有服务器列表，并做ping，确定状态
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /cicd/server/list [get]
func CicdServerList(c *gin.Context) {
	list := global.V.MyService.Cicd.GetServerList()
	httpresponse.OkWithDetailed(list, "成功", c)
}

// @Tags Cicd
// @Summary 部署/发布 列表
// @Description 部署/发布 列表
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /cicd/publish/list [get]
func CicdPublishList(c *gin.Context) {
	list := global.V.MyService.Cicd.GetPublishList()
	httpresponse.OkWithDetailed(list, "成功", c)
}

// @Tags Cicd
// @Summary 发布项目
// @Description 发布项目
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param id path string true "publish id"
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /cicd/service/publish/{id}/{flag} [get]
func CicdServicePublish(c *gin.Context) {
	idStr  := c.Param("id")
	flagStr  := c.Param("flag")
	if idStr == ""{
		httpresponse.FailWithMessage("id empty 1",c)
		return
	}
	id ,err := strconv.Atoi(idStr)
	if err != nil {
		httpresponse.FailWithMessage("id empty 2",c)
		return
	}
	if id == 0{
		httpresponse.FailWithMessage("id empty 3",c)
		return
	}

	flag ,_ := strconv.Atoi(flagStr)

	err = global.V.MyService.Cicd.Publish(id,flag)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(),c)
	}else{
		httpresponse.OkWithDetailed("bbb", "成功", c)
	}


}


// @Tags Cicd
// @Summary 部署一个服务
// @Description demo
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.CicdDeploy true "用户信息"
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /cicd/service/deploy [post]
func CicdServiceDeploy(c *gin.Context) {
	var form request.CicdDeploy
	c.ShouldBind(&form)

	util.MyPrint("CicdServiceDeploy form:",form)
	err := global.V.MyService.Cicd.ApiDeployOneService(form)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(),c)
	}else{
		httpresponse.OkWithDetailed("aaaa", "成功", c)
	}

}

// @Tags Cicd
// @Summary ping
// @Description 测试对端有没有开启服务
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /cicd/ping [get]
func CicdPing(c *gin.Context) {
	httpresponse.OkWithDetailed("aaaa", "成功", c)
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
	list := global.V.MyService.Cicd.GetHasDeployService()
	httpresponse.OkWithDetailed(list, "成功", c)
}

// @Tags Cicd
// @Summary 同步 本机部署的项目 -> 到目标机器
// @Description scp
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /cicd/local/sync/target [get]
func CicdLocalSyncTarget(c *gin.Context) {
	var form request.CicdSync
	c.ShouldBind(&form)

	list := global.V.MyService.Cicd.LocalSyncTarget(form)
	httpresponse.OkWithDetailed(list, "成功", c)
}

