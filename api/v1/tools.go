package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/service/cicd"
	"zgoframe/util"
	myservice "zgoframe/service"
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
	httpresponse.OkWithAll("aaaa", "成功", c)
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

	httpresponse.OkWithAll(rs, "成功", c)
}

// @Tags Tools
// @Summary 项目列表
// @Description 每个项目的详细信息
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} []model.Project
// @Router /tools/project/list [post]
func ProjectList(c *gin.Context) {
	var a model.Project
	c.ShouldBind(&a)

	rs := global.V.ProjectMng.Pool

	httpresponse.OkWithAll(rs, "成功", c)
}

// @Tags Tools
// @Summary 所有常量列表
// @Description 所有常量列表，方便调用与调试
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} []httpresponse.ConstInfo "常量列表"
// @Router /tools/const/list [get]
func ConstList(c *gin.Context) {
	var a model.Project
	c.ShouldBind(&a)

	//list := global.V.MyService.User.GetConstList()
	list := GetConstList()

	httpresponse.OkWithAll(list, "成功", c)

}

// @Tags Tools
// @Summary 常量列表 - 生成mysql导入脚本
// @Description 给后台使用，生成到MYSQL数据库中，便于后台统一使用
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {string} string "sql script"
// @Router /tools/const/init/db [get]
func ConstInitDb(c *gin.Context) {
	/*
	delete from sys_dictionaries where id > 6;
	delete from sys_dictionary_details where sys_dictionary_id > 6;
	*/

	//list := global.V.MyService.User.GetConstList()
	list := GetConstList()

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
	httpresponse.OkWithAll(sqlStr, "成功", c)
}

// @Tags Tools
// @Summary 基数据 - 生成mysql脚本，导入到DB中，供后台UI可视化查看
// @Description tables: project instance server
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {string} string "sql script"
// @Router /tools/test/init/db [get]
func ConstInitTestDb(c *gin.Context) {
	envList := util.GetConstListEnv()
	//id := 1
	ipList := make(map[int]string)
	ipList[1] = "127.0.0.1"//本地
	//ipList[2] = "1.1.1.1"//开发
	//ipList[3] = "2.2.2.2"//测试
	ipList[4] = "8.142.177.235"//预发布
	//ipList[5] = "3.3.3.3"//线上

	serverSql := ""
	for k,v:=range envList{
		serverInsertSql := "INSERT INTO `server` (`id`, `name`, `platform`, `out_ip`, `inner_ip`, `env`, `status`, `ext`, `charge_user_name`, `start_time`, `end_time`, `price`, `created_at`, `updated_at`, `deleted_at`,`state`) "
		serverInsertSql += "VALUES  ("+strconv.Itoa(v)+",'"+k+"', 1,   '"+ipList[v]+"', '127.0.0.1', "+strconv.Itoa(v)+", '1', '', '小z', '1650006845', '1650006845', '100', '1650006845', '0', NULL,1);   "
		serverSql += serverInsertSql
	}

	instanceSql := ""
	for _,envId:=range envList{
		if envId == 1 || envId == 4 {
			for _,instance := range cicd.ThirdInstance{
				if !CheckInAllowInstance(instance){
					continue
				}
				instanceInsertSql := "INSERT INTO `instance` (`id`, `platform`, `name`, `host`, `port`, `env`, `user`, `ps`, `ext`, `status`, `charge_user_name`, `start_time`, `end_time`, `price`, `created_at`, `updated_at`, `deleted_at`) "
				instanceInsertSql += "VALUES                  (NULL, '1', '"+instance+"', '"+ipList[envId] +"', '3306', '"+strconv.Itoa(envId)+"', 'aaaa', 'bbbb', '', '1', '小z', '1650006845', '1650006845', '200', '1650006845', '0', NULL);"
				instanceSql += instanceInsertSql
			}
		}
	}

	rs := make(map[string]string)
	rs["serverSql"] = serverSql
	rs["instanceSql"] = instanceSql


	httpresponse.OkWithAll(rs, "成功", c)
}

func CheckInAllowInstance(name string)bool{
	allowInstance := []string{"mysql","redis","prometheus","kibana","grafana","etcd","es","http"}
	for _,v:= range allowInstance {
		if v == name{
			return true
		}

	}
	return false

}



// @Tags Tools
// @Summary header头-结构体
// @Description 日常header里放一诸如验证类的东西，统一公示出来，仅是说明，方便测试/前端查看，方便使用
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

	httpresponse.OkWithAll(myheader, "成功lalalalala", c)
}



var ConstDataList = []httpresponse.ConstInfo{}

func AddConstLis(row httpresponse.ConstInfo){
	ConstDataList = append(ConstDataList,row)
}

func  GetConstList() []httpresponse.ConstInfo {
	AddConstLis(httpresponse.ConstInfo{
			List: util.GetConstListEnv(),
			Name: "env-环境",
			Key:  "ENV",
		})

	AddConstLis( httpresponse.ConstInfo{
		List: util.GetConstListEnv(),
		Name: "env-环境",
		Key:  "ENV",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListProjectType(),
		Name: "项目类型",
		Key:  "PROJECT_TYPE",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListProjectStatus(),
		Name: "项目状态",
		Key:  "PROJECT_STATUS",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListPlatform(),
		Name: "平台类型",
		Key:  "PLATFORM",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListUserTypeThird(),
		Name: "用户类型3方",
		Key:  "USER_TYPE_THIRD",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListUserRegType(),
		Name: "用户注册类型",
		Key:  "USER_REG_TYPE",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListUserSex(),
		Name: "用户性别",
		Key:  "USER_SEX",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListUserStatus(),
		Name: "用户状态",
		Key:  "USER_STATUS",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListUserTypeThirdCN(),
		Name: "用户类型-中国",
		Key:  "USER_TYPE_THIRD_CN",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListUserTypeThirdNotCN(),
		Name: "用户类型-外国",
		Key:  "USER_TYPE_THIRD_NOT_CN",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListUserGuest(),
		Name: "游客分类",
		Key:  "USER_GUEST",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListUserRobot(),
		Name: "机器人",
		Key:  "USER_ROBOT",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListUserTest(),
		Name: "测试账号",
		Key:  "USER_TEST",
	})

	AddConstLis(httpresponse.ConstInfo{
		List: model.GetConstListPurpose(),
		Name: "目的",
		Key:  "PURPOSE",
	})

	AddConstLis(   httpresponse.ConstInfo{
		List: model.GetConstListAuthCodeStatus(),
		Name: "验证码状态",
		Key:  "AUTH_CODE_STATUS",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListRuleType(),
		Name: "配置规则类型",
		Key:  "RULE_TYPE",
	})

	//
	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListSmsChannel(),
		Name: "短信渠道",
		Key:  "SMS_CHANNEL",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListServerPlatform(),
		Name: "服务器平台",
		Key:  "SERVER_PLATFORM",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: model.GetConstListCicdPublishStatus(),
		Name: "CICD发布状态",
		Key:  "CICD_PUBLISH_STATUS",
	})
	AddConstLis( httpresponse.ConstInfo{
		List: myservice.GetConstListMailBoxType(),
		Name: "站内信,信件箱类型",
		Key:  "MAIL_BOX",
	})

	AddConstLis( httpresponse.ConstInfo{
		List: myservice.GetConstListMailPeople(),
		Name: "站内信,接收人群类型",
		Key:  "MAIL_PEOPLE",
	})


	AddConstLis(httpresponse.ConstInfo{
		List: myservice.GetConstListConfigPersistenceType(),
		Name: "配置中心持久化类型",
		Key:  "CONFIG_PERSISTENCE_TYPE",
	})

	return ConstDataList
}
