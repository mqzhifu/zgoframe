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
	/*
	delete from sys_dictionaries where id > 6;
	delete from sys_dictionary_details where sys_dictionary_id > 6;
	*/

	list := global.V.MyService.User.GetConstList()
	sqlTemp := "INSERT INTO `sys_dictionaries` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `type`, `status`, `desc`) VALUES (#id#, '2022-04-04 00:01:01',NULL, NULL, '#name#', '#key#', '1', '')"
	subSqlTemp := "INSERT INTO `sys_dictionary_details` (`id`, `created_at`, `updated_at`, `deleted_at`, `label`, `value`, `status`, `sort`, `sys_dictionary_id`) VALUES (NULL,  '2022-04-04 00:00:01',NULL , NULL, '#name#', '#value#', '1', '0', '#link_id#')"
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
// @Summary 基数据 - 生成mysql导入脚本
// @Description tables: project instance server
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} httpresponse.Response
// @Router /tools/test/init/db [get]
func ConstInitTestDb(c *gin.Context) {
	envList := util.GetConstListEnv()
	//id := 1
	ipList := make(map[int]string)
	ipList[1] = "127.0.0.1"
	ipList[2] = "1.1.1.1"
	ipList[3] = "2.2.2.2"
	ipList[4] = "8.142.177.235"
	ipList[5] = "3.3.3.3"

	serverSql := ""
	for k,v:=range envList{
		serverInsertSql := "INSERT INTO `server` (`id`, `name`, `platform`, `out_ip`, `inner_ip`, `env`, `status`, `ext`, `charge_user_name`, `start_time`, `end_time`, `price`, `created_at`, `updated_at`, `deleted_at`,`state`) "
		serverInsertSql += "VALUES  ("+strconv.Itoa(v)+",'"+k+"', 1,   '"+ipList[v]+"', '127.0.0.1', "+strconv.Itoa(v)+", '1', '', '小z', '1650006845', '1650006845', '100', '1650006845', '0', NULL,1);   "
		serverSql += serverInsertSql
	}

	instanceSql := ""
	for _,envId:=range envList{
		for _,instance := range util.ThirdInstance{
			if !CheckInAllowInstance(instance){
				continue
			}
			instanceInsertSql := "INSERT INTO `instance` (`id`, `platform`, `name`, `host`, `port`, `env`, `user`, `ps`, `ext`, `status`, `charge_user_name`, `start_time`, `end_time`, `price`, `created_at`, `updated_at`, `deleted_at`) "
			instanceInsertSql += "VALUES                  (NULL, '1', '"+instance+"', '"+ipList[4] +"', '3306', '"+strconv.Itoa(envId)+"', 'aaaa', 'bbbb', '', '1', '小z', '1650006845', '1650006845', '200', '1650006845', '0', NULL);"
			instanceSql += instanceInsertSql
		}
	}

	rs := make(map[string]string)
	rs["serverSql"] = serverSql
	rs["instanceSql"] = instanceSql


	httpresponse.OkWithDetailed(rs, "成功", c)
}

func CheckInAllowInstance(name string)bool{
	allowInstance := []string{"mysql","redis","prometheus","kibana","grafana","etcd","es"}
	for _,v:= range allowInstance {
		if v == name{
			return true
		}

	}
	return false

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
