package v1

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Tags ConfigCenter
// @Summary 以模块(文件)为单位，获取该模块(文件)下的所有配置信息
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param name path string true "module name"
// @Success 200 {object} httpresponse.LoginResponse
// @Router /config/center/module/{name} [get]
func ConfigCenterGetByModule(c *gin.Context) {
	projectId := request.GetProjectIdByHeader(c)
	moduleName  := c.Param("name")
	configInfo , err := global.V.MyService.ConfigCenter.GetByModule(global.C.System.ENV,projectId,moduleName)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(),c)
	}else{
		httpresponse.OkWithAll(configInfo,"成功",c)
	}

}

// @Tags ConfigCenter
// @Summary 以以模块(文件)+里面具体的key 为单位，获取配置信息
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.ConfigCenterGetByKeyReq true " "
// @Success 200 {object} httpresponse.LoginResponse
// @Router /config/center/get/key [post]
func ConfigCenterGetByModuleByKey(c *gin.Context) {
	projectId := request.GetProjectIdByHeader(c)

	var form request.ConfigCenterGetByKeyReq
	c.ShouldBind(&form)

	configInfo , err := global.V.MyService.ConfigCenter.GetByKey(global.C.System.ENV,projectId,form.Module,form.Key)
	util.MyPrint(configInfo)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(),c)
	}else{
		httpresponse.OkWithAll(configInfo,"成功",c)
	}
}

// @Tags ConfigCenter
// @Summary 以模块(文件)+里面具体的key 为单位，设置置信息(如果存在，覆盖)
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.ConfigCenterSetByKeyReq true " "
// @Success 200 {object} httpresponse.LoginResponse
// @Router /config/center/set/key [post]
func ConfigCenterSetByModuleByKey(c *gin.Context) {
	projectId := request.GetProjectIdByHeader(c)

	var form request.ConfigCenterSetByKeyReq
	c.ShouldBind(&form)

	err := global.V.MyService.ConfigCenter.SetByKey(global.C.System.ENV,projectId,form.Module,form.Key,form.Value)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(),c)
	}else{
		httpresponse.OkWithAll("ok","成功",c)
	}
}



