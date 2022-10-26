package v1

import (
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/service/cicd"
	"zgoframe/util"
)

// @Tags Tools
// @Summary 一个项目的详细信息
// @Description 用于开发工具测试
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(aaaa)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(bbbb)
// @Param id path int true "项目ID"
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /tools/project/info/{id} [get]
func ProjectOneInfo(c *gin.Context) {
	var a model.Project
	c.ShouldBind(&a)

	id := util.Atoi(c.Param("id"))
	if id <= 0 {
		httpresponse.FailWithMessage("id <= 0", c)
		return
	}

	info, empty := global.V.ProjectMng.GetById(id)
	if empty {
		httpresponse.FailWithMessage("id not found in db.", c)
		return
	} else {
		httpresponse.OkWithAll(info, "成功", c)
	}

}

// @Tags Tools
// @Summary 项目列表
// @Description 每个项目的详细信息
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(aaaa)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(bbbb)
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
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(aaaa)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(bbbb)
// @Produce  application/json
// @Success 200 {object} []httpresponse.ConstInfo "常量列表"
// @Router /tools/const/list [get]
func ConstList(c *gin.Context) {
	//var a model.Project
	//c.ShouldBind(&a)
	//
	//list := GetConstList()
	//
	//httpresponse.OkWithAll(list, "成功", c)

}

// @Tags Tools
// @Summary 常量列表 - 生成mysql导入GVA中的脚本
// @Description 给后台使用，生成到MYSQL数据库中，便于后台统一使用
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {string} string "sql script"
// @Router /tools/const/init/db [get]
func ConstInitDb(c *gin.Context) {
	constHandle := util.NewConstHandle()
	enumConstPool := constHandle.EnumConstPool

	sqlTemp := "INSERT INTO `sys_dictionaries` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `type`, `status`, `desc`) VALUES (#id#, '2022-04-04 00:01:01',NULL, NULL, '#name#', '#key#', '1', '')"
	subSqlTemp := "INSERT INTO `sys_dictionary_details` (`id`, `created_at`, `updated_at`, `deleted_at`, `label`, `value`, `status`, `sort`, `sys_dictionary_id`) VALUES (NULL,  '2022-04-04 00:00:01',NULL , NULL, '#name#', '#value#', '1', '0', '#link_id#')"
	sqlStr := ""
	id := 10
	for _, EnumConst := range enumConstPool {
		sql1 := strings.Replace(sqlTemp, "#id#", strconv.Itoa(id), -1)
		sql1 = strings.Replace(sql1, "#name#", EnumConst.Desc, -1)
		sql1 = strings.Replace(sql1, "#key#", EnumConst.CommonPrefix, -1)
		//util.MyPrint(sql1)
		sqlStr += sql1 + ";  \n"
		//sqlList = append(sqlList, sql1)

		for _, constItem := range EnumConst.ConstList {
			value := ""
			if EnumConst.Type == "int" {
				aa := constItem.Value.(int)
				value = strconv.Itoa(aa)
			} else {
				value = constItem.Value.(string)
			}
			sql2_sub := strings.Replace(subSqlTemp, "#link_id#", strconv.Itoa(id), -1)
			sql2_sub = strings.Replace(sql2_sub, "#name#", constItem.Desc, -1)
			sql2_sub = strings.Replace(sql2_sub, "#value#", value, -1)
			sqlStr += sql2_sub + ";    \n"
			//util.MyPrint("    " + sql2_sub)
		}

		id++
	}
	util.MyPrint(sqlStr)
	httpresponse.OkWithAll(sqlStr, "成功", c)
}

// @Tags Tools
// @Summary 生成mysql数据-脚本，可导入到DB中
// @Description tables: project instance server
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {string} string "sql script"
// @Router /tools/init/db/data [get]
func InitDbData(c *gin.Context) {
	tables := []string{"sms_rule", "mail_rule", "project", "server", "instance"}
	tablesStr := ""
	for _, v := range tables {
		tablesStr += v + " "
	}
	outFile := global.C.Http.DiskStaticPath + "/data/db_data.sql"
	config := global.C.Mysql[0]
	shell := "/soft/mysql/bin/mysqldump -h " + config.Ip + " -u" + config.Username + " -p" + config.Password + " " + config.DbName + " " + tablesStr + " > " + outFile

	httpresponse.OkWithMessage(shell, c)

}

// @Tags Tools
// @Summary 生成mysql表结构-脚本，可导入到DB中
// @Description tables: project instance server
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {string} string "sql script"
// @Router /tools/init/db/structure [get]
func InitDbStructure(c *gin.Context) {
	sql := global.AutoCreateUpDbTable()
	path := global.C.Http.DiskStaticPath + "/data/db_structure.sql"
	newConfigFile, _ := os.Create(path)
	for _, v := range sql {
		v += " ;\n"
		newConfigFile.Write([]byte(v))
	}
	//util.(path)

	httpresponse.OkWithAll(path, "成功", c)
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
	//envList := util.GetConstListEnv()
	////id := 1
	//ipList := make(map[int]string)
	//ipList[1] = "127.0.0.1"//本地
	////ipList[2] = "1.1.1.1"//开发
	////ipList[3] = "2.2.2.2"//测试
	//ipList[4] = "8.142.177.235"//预发布
	////ipList[5] = "3.3.3.3"//线上
	//
	//serverSql := ""
	//for k,v:=range envList{
	//	serverInsertSql := "INSERT INTO `server` (`id`, `name`, `platform`, `out_ip`, `inner_ip`, `env`, `status`, `ext`, `charge_user_name`, `start_time`, `end_time`, `price`, `created_at`, `updated_at`, `deleted_at`,`state`) "
	//	serverInsertSql += "VALUES  ("+strconv.Itoa(v)+",'"+k+"', 1,   '"+ipList[v]+"', '127.0.0.1', "+strconv.Itoa(v)+", '1', '', '小z', '1650006845', '1650006845', '100', '1650006845', '0', NULL,1);   "
	//	serverSql += serverInsertSql
	//}

	serverMng, _ := util.NewServerManger(global.V.Gorm)
	serverList := serverMng.Pool

	instanceSql := ""
	for _, server := range serverList {
		if server.Env != 5 { //正式环境，都已经配置好了，就不要随便再重新插入了
			for _, instance := range cicd.ThirdInstance {
				if !CheckInAllowInstance(instance) {
					continue
				}
				instanceInsertSql := "INSERT INTO `instance` (`id`, `platform`, `name`, `host`, `port`, `env`, `user`, `ps`, `ext`, `status`, `charge_user_name`, `start_time`, `end_time`, `price`, `created_at`, `updated_at`, `deleted_at`) "
				instanceInsertSql += "VALUES                  (NULL, '1', '" + instance + "', '" + server.OutIp + "', '', '" + strconv.Itoa(server.Env) + "', 'aaaa', 'bbbb', '', '1', '小z', '1650006845', '1650006845', '200', '1650006845', '0', NULL);"
				instanceSql += instanceInsertSql
			}
		}
	}

	rs := make(map[string]string)
	//rs["serverSql"] = serverSql
	rs["instanceSql"] = instanceSql

	httpresponse.OkWithAll(rs, "成功", c)
}

func CheckInAllowInstance(name string) bool {
	allowInstance := []string{"mysql", "redis", "etcd", "oss", "http", "grpc", "ali_email", "email", "gateway", "agora", "domain", "cdn", "alert", "sms", "super_visor"}
	for _, v := range allowInstance {
		if v == name {
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

func AddConstLis(row httpresponse.ConstInfo) {
	ConstDataList = append(ConstDataList, row)
}

//func GetConstList() []httpresponse.ConstInfo {
//	AddConstLis(httpresponse.ConstInfo{
//		List: util.GetConstListEnv(),
//		Name: "env-环境",
//		Key:  "ENV",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListProjectType(),
//		Name: "项目类型",
//		Key:  "PROJECT_TYPE",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListProjectStatus(),
//		Name: "项目状态",
//		Key:  "PROJECT_STATUS",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListPlatform(),
//		Name: "平台类型",
//		Key:  "PLATFORM",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListUserTypeThird(),
//		Name: "用户类型3方",
//		Key:  "USER_TYPE_THIRD",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListUserRegType(),
//		Name: "用户注册类型",
//		Key:  "USER_REG_TYPE",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListUserSex(),
//		Name: "用户性别",
//		Key:  "USER_SEX",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListUserStatus(),
//		Name: "用户状态",
//		Key:  "USER_STATUS",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListUserTypeThirdCN(),
//		Name: "用户类型-中国",
//		Key:  "USER_TYPE_THIRD_CN",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListUserTypeThirdNotCN(),
//		Name: "用户类型-外国",
//		Key:  "USER_TYPE_THIRD_NOT_CN",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListUserGuest(),
//		Name: "游客分类",
//		Key:  "USER_GUEST",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListUserRobot(),
//		Name: "机器人",
//		Key:  "USER_ROBOT",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListUserTest(),
//		Name: "测试账号",
//		Key:  "USER_TEST",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListPurpose(),
//		Name: "本次操作目的",
//		Key:  "PURPOSE",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListAuthCodeStatus(),
//		Name: "验证码状态",
//		Key:  "AUTH_CODE_STATUS",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListRuleType(),
//		Name: "配置规则类型",
//		Key:  "RULE_TYPE",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListSendMsgStatus(),
//		Name: "消息发送状态",
//		Key:  "SEND_MSG_STATUS",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListProjectLanguage(),
//		Name: "项目开发语言",
//		Key:  "PROJECT_LANG",
//	})
//
//	//AddConstLis(httpresponse.ConstInfo{
//	//	List: model.GetConstListFileHashType(),
//	//	Name: "文件hash类型",
//	//	Key:  "FILE_HASH_TYPE",
//	//})
//
//	//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListSmsChannel(),
//		Name: "短信渠道",
//		Key:  "SMS_CHANNEL",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListServerPlatform(),
//		Name: "服务器平台",
//		Key:  "SERVER_PLATFORM",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListCicdPublishStatus(),
//		Name: "CICD发布状态",
//		Key:  "CICD_PUBLISH_STATUS",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: model.GetConstListCicdPublishDeployStatus(),
//		Name: "CICD发布部署状态",
//		Key:  "CICD_PUBLISH_STATUS",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: myservice.GetConstListMailBoxType(),
//		Name: "站内信,信件箱类型",
//		Key:  "MAIL_BOX",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: myservice.GetConstListMailPeople(),
//		Name: "站内信,接收人群类型",
//		Key:  "MAIL_PEOPLE",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: myservice.GetConstListConfigPersistenceType(),
//		Name: "配置中心持久化类型",
//		Key:  "CONFIG_PERSISTENCE_TYPE",
//	})
//
//	AddConstLis(httpresponse.ConstInfo{
//		List: global.GetUtilUploadConst(),
//		Name: "上传，文件类型",
//		Key:  "UPLOAD_FILE_TYPE",
//	})
//
//	return ConstDataList
//}
