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
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param X-Client-Req-Time header string true "客户端请求时间unix" default(1648277052)
// @Param data body request.ConfigCenterOpt true "请求参数"
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /config/center/get/module [POST]
func ConfigCenterGetByModule(c *gin.Context) {
	projectId := request.GetProjectIdByHeader(c)

	var form request.ConfigCenterOpt
	c.ShouldBind(&form)

	// moduleName  := c.Param("name")
	configInfo, err := ApiServices().ConfigCenter.GetByModule(form.Env, projectId, form.Module)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		httpresponse.OkWithAll(configInfo, "成功", c)
	}

}

// @Tags ConfigCenter
// @Summary 以以模块(文件)+里面具体的key 为单位，获取配置信息
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param X-Client-Req-Time header string true "客户端请求时间unix" default(1648277052)
// @Param data body request.ConfigCenterOpt true "请求参数"
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /config/center/get/key [post]
func ConfigCenterGetByModuleByKey(c *gin.Context) {
	projectId := request.GetProjectIdByHeader(c)

	var form request.ConfigCenterOpt
	c.ShouldBind(&form)

	configInfo, err := ApiServices().ConfigCenter.GetByKey(global.MainCmdParameter.Env, projectId, form.Module, form.Key)
	util.MyPrint(configInfo)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		httpresponse.OkWithAll(configInfo, "成功", c)
	}
}

// @Tags ConfigCenter
// @Summary 以模块(文件)+里面具体的key 为单位，设置置信息(如果存在，覆盖)
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param X-Client-Req-Time header string true "客户端请求时间unix" default(1648277052)
// @Param data body request.ConfigCenterOpt true "请求参数"
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /config/center/module/set/key [post]
func ConfigCenterSetByModuleByKey(c *gin.Context) {
	projectId := request.GetProjectIdByHeader(c)

	var form request.ConfigCenterOpt
	c.ShouldBind(&form)

	// util.MyPrint("set key form:",form.Env,form.Module,form.Value)

	err := ApiServices().ConfigCenter.SetByKey(form.Env, projectId, form.Module, form.Key, form.Value)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		httpresponse.OkWithAll("ok", "成功", c)
	}
}

// @Tags ConfigCenter
// @Summary 创建模块(文件)
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param X-Client-Req-Time header string true "客户端请求时间unix" default(1648277052)
// @Param data body request.ConfigCenterOpt true "请求参数"
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /config/center/create/module [post]
func ConfigCenterCreateModule(c *gin.Context) {
	projectId := request.GetProjectIdByHeader(c)

	var form request.ConfigCenterOpt
	c.ShouldBind(&form)

	err := ApiServices().ConfigCenter.CreateModule(form.Env, projectId, form.Module)
	if err == nil {
		httpresponse.OkWithMessage("已创建", c)
	} else {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
}
