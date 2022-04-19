package v1

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
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
// @Summary 部署一个服务
// @Description demo
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /cicd/service/deploy/{name} [get]
func CicdServiceDeploy(c *gin.Context) {
	//global.V.MyService.Cicd.DeployAllService()
	httpresponse.OkWithDetailed("aaaa", "成功", c)
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

