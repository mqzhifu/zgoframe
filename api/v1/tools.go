package v1

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/util"
)

// @Tags Tools
// @Summary 一个项目的详细信息
// @Description 用于开发工具测试
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param id path int true "项目ID"
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /tools/project/info/{id} [get]
func ProjectOneInfo(c *gin.Context) {
	var a model.Project
	c.ShouldBind(&a)

	rs := global.V.ProjectMng.Pool

	httpresponse.OkWithDetailed(rs, "成功", c)
}

// @Tags Tools
// @Summary 项目列表
// @Description 每个项目的详细信息
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /tools/project/list [post]
func ProjectList(c *gin.Context) {
	var a model.Project
	c.ShouldBind(&a)

	rs := global.V.ProjectMng.Pool

	httpresponse.OkWithDetailed(rs, "成功", c)
}

// @Tags Tools
// @Summary 所有常量列表
// @Description 所有常量列表，方便调用与调试
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} httpresponse.Response
// @Router /tools/const/list [get]
func ConstList(c *gin.Context) {
	var a model.Project
	c.ShouldBind(&a)

	list := make(map[string]interface{})

	list["PROJECT_TYPE_MAP"] = util.PROJECT_TYPE_MAP
	list["PlatformList"] = request.GetPlatformList()
	list["ThirdTypeList"] = model.GetUserThirdTypeList()
	list["UserRegTypeList"] = model.GetUserRegTypeList()
	list["UserRegTypeList"] = model.GetUserSexList()
	list["UserStatusList"] = model.GetUserStatusList()

	httpresponse.OkWithDetailed(list, "成功", c)

}

// @Tags Tools
// @Summary header头结构体
// @Description 日常header里放一诸如验证类的东西，统一公示出来，仅是说明，方便测试，不是真实API，方便使用
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID" Enums(1,2,3,4) default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Client-Req-Time header string true "客户端请求时间unix" default(1648277052)
// @Success 200 {object} request.TestHeader
// @Router /tools/header/struct [get]
func HeaderStruct(c *gin.Context) {
	myheader := request.TestHeader{
		HeaderRequest:  request.HeaderRequest{},
		HeaderResponse: request.HeaderResponse{},
	}

	httpresponse.OkWithDetailed(myheader, "成功lalalalala", c)
}
