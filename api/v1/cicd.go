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
// @Summary 部署/发布 列表
// @Description 查看 项目的已部署 日志 列表
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param X-Client-Req-Time header string true "客户端请求时间unix" default(1648277052)
// @Produce  application/json
// @Success 200 {object} model.CicdPublish
// @Router /cicd/publish/list [get]
func CicdPublishList(c *gin.Context) {
	list := global.V.MyService.Cicd.GetPublishList(20) // 取出最新的20条即可，不然数据太大
	httpresponse.OkWithAll(list, "成功", c)
}

// @Tags Cicd
// @Summary 查看 所有项目的 可部署 列表
// @Description 每台机器上，有多少个服务，具体可以操作部署哪个项目到远端服务器
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param X-Client-Req-Time header string true "客户端请求时间unix" default(1648277052)
// @Produce  application/json
// @Success 200 {object} util.Service
// @Router /cicd/local/all/server/service/list [get]
func CicdLocalAllServerServiceList(c *gin.Context) {
	list, _ := global.V.MyService.Cicd.LocalAllServerServiceList()
	httpresponse.OkWithAll(list, "成功", c)
}

// @Tags Cicd
// @Summary 所有服务(项目) 列表
// @Description 属于基础数据，查看下 当前所有项目的状态
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param X-Client-Req-Time header string true "客户端请求时间unix" default(1648277052)
// @Produce  application/json
// @Success 200 {object} util.Service
// @Router /cicd/service/list [get]
func CicdServiceList(c *gin.Context) {
	list := global.V.MyService.Cicd.GetServiceList()
	httpresponse.OkWithAll(list, "成功", c)
}

// @Tags Cicd
// @Summary 服务器 列表
// @Description 属于基础数据，获取所有服务器列表，并做ping，确定状态
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param X-Client-Req-Time header string true "客户端请求时间unix" default(1648277052)
// @Produce  application/json
// @Success 200 {object} util.Server
// @Router /cicd/server/list [get]
func CicdServerList(c *gin.Context) {
	list := global.V.MyService.Cicd.GetServerList()
	httpresponse.OkWithAll(list, "成功", c)
}

// @Tags Cicd
// @Summary 每台机器上的 superVisor
// @Description 每台机上都会有一个 superVisor 进程，管理着所有 服务/项目，从 superVisor 视角看一下所有 项目的状态
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param X-Client-Req-Time header string true "客户端请求时间unix" default(1648277052)
// @Produce  application/json
// @Success 200 {boolean} true "数据过长，先用bool替代"
// @Router /cicd/superVisor/list [get]
func CicdSuperVisorList(c *gin.Context) {
	list, err := global.V.MyService.Cicd.GetSuperVisorList()
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		httpresponse.OkWithAll(list, "成功", c)
	}

}

// @Tags Cicd
// @Summary 本机操作远端的 superVisor ，管理进程
// @Description 对远端服务器上的 superVisor，管理一个服务，如：停止进程 重启进程 启动进程
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param X-Client-Req-Time header string true "客户端请求时间unix" default(1648277052)
// @Param data body request.CicdDeploy true "用户信息"
// @Produce  application/json
// @Success 200 {bool} bool "true:成功 false:失败"
// @Router /cicd/superVisor/process [post]
func CicdSuperVisorProcess(c *gin.Context) {
	var form request.CicdSuperVisor
	c.ShouldBind(&form)

	util.MyPrint("CicdSuperVisorProcess form:", form)
	err := global.V.MyService.Cicd.SuperVisorProcess(form)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		httpresponse.OkWithAll("aaaa", "成功", c)
	}
}

// @Tags Cicd
// @Summary 发布项目
// @Description 已部署好的项目，正式 发布
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param X-Client-Req-Time header string true "客户端请求时间unix" default(1648277052)
// @Param id path string true "publish id"
// @Param flag path string flase "deployTargetType"
// @Produce  application/json
// @Success 200 {bool} bool "true:成功 false:失败"
// @Router /cicd/service/publish/{id}/{flag} [get]
func CicdServicePublish(c *gin.Context) {
	idStr := c.Param("id")
	flagStr := c.Param("flag")
	if idStr == "" {
		httpresponse.FailWithMessage("id empty 1", c)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpresponse.FailWithMessage("id empty 2", c)
		return
	}
	if id == 0 {
		httpresponse.FailWithMessage("id empty 3", c)
		return
	}

	flag, _ := strconv.Atoi(flagStr)
	err = global.V.MyService.Cicd.Deploy.Publish(id, flag)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		httpresponse.OkWithAll("bbb", "成功", c)
	}

}

// @Tags Cicd
// @Summary 部署一个服务
// @Description 开始把项目部署到指定的服务器上，本机编译，最后再把代码同步到远端服务器上
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param X-Client-Req-Time header string true "客户端请求时间unix" default(1648277052)
// @Param data body request.CicdDeploy true "用户信息"
// @Produce  application/json
// @Success 200 {bool} bool "true:成功 false:失败"
// @Router /cicd/service/deploy [post]
func CicdServiceDeploy(c *gin.Context) {
	var form request.CicdDeploy
	c.ShouldBind(&form)

	util.MyPrint("CicdServiceDeploy form:", form)
	// 这里因为是HTTP连接，而后端处理一次时间接近1分钟，HTTP可能多次重复请求，开个协程
	go global.V.MyService.Cicd.Deploy.ApiDeployOneService(form)
	httpresponse.OkWithAll("aaaa", "成功", c)
	// if err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	// } else {
	//	httpresponse.OkWithAll("aaaa", "成功", c)
	// }

}

// func CicdLocalSyncTarget(c *gin.Context) {
//	var form request.CicdSync
//	c.ShouldBind(&form)
//
//	list := global.V.MyService.Cicd.LocalSyncTarget(form)
//	httpresponse.OkWithAll(list, "成功", c)
// }
