package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/util"
)

// @Tags Tools
// @Summary 帧同步 - js
// @Description demo
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param id path int true "项目ID"
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /tools/frame/sync/js/demo [get]
func FrameSyncJsDemo(c *gin.Context) {
	httpresponse.OkWithDetailed("aaaa", "成功", c)
}

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

	list := global.V.MyService.User.GetConstList()

	httpresponse.OkWithDetailed(list, "成功", c)

}

// @Tags Tools
// @Summary 常量列表 - 生成mysql导入脚本
// @Description 给后台使用，生成到MYSQL数据库中，便于后台统一使用
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} httpresponse.Response
// @Router /tools/const/init/db [get]
func ConstInitDb(c *gin.Context) {

	list := global.V.MyService.User.GetConstList()
	sqlTemp := "INSERT INTO `sys_dictionaries` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `type`, `status`, `desc`) VALUES (#id#, NULL,NULL, NULL, '#name#', '#key#', '1', '')"
	subSqlTemp := "INSERT INTO `sys_dictionary_details` (`id`, `created_at`, `updated_at`, `deleted_at`, `label`, `value`, `status`, `sort`, `sys_dictionary_id`) VALUES (NULL,  NULL,NULL , NULL, '#name#', '#value#', '1', '0', '#link_id#')"
	sqlStr := ""
	id := 10
	for _, v := range list {
		sql1 := strings.Replace(sqlTemp, "#id#", strconv.Itoa(id), -1)
		sql1 = strings.Replace(sql1, "#name#", v.Name, -1)
		sql1 = strings.Replace(sql1, "#key#", v.Key, -1)
		util.MyPrint(sql1)
		sqlStr += sql1 + ";    "
		//sqlList = append(sqlList, sql1)

		for k, sub := range v.List {
			sql2_sub := strings.Replace(subSqlTemp, "#link_id#", strconv.Itoa(id), -1)
			sql2_sub = strings.Replace(sql2_sub, "#name#", k, -1)
			sql2_sub = strings.Replace(sql2_sub, "#value#", strconv.Itoa(sub), -1)
			sqlStr += sql2_sub + ";    "
			util.MyPrint("    " + sql2_sub)
		}

		id++
	}
	httpresponse.OkWithDetailed(sqlStr, "成功", c)
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
