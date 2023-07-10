package v1

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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

type MiguRes struct {
	AppId     string
	Data      string
	DataBytes []byte
	Time      int64
	TimeStr   string
	Sign      string
	FinalData string
	SignLower string
}

// @Tags Tools
// @Summary 测试咪咕
// @Description 120项目API接口
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Produce  application/json
// @Success 200 {object} v1.MiguRes "最终的请求参数信息"
// @Router /tools/test/migu/api/para [GET]
func TestMiguAPIGetPara(c *gin.Context) {
	type DataStruct struct {
		ArSn   string `json:"arSn"`
		CpeMac string `json:"cpeMac"`
	}

	appId := "wechat_625"
	appSecret := "b267a314-2208-4970-a0fc-b9f0e677b437"
	data := DataStruct{ArSn: "T20230101BJ00011", CpeMac: "AC-BD-8I-FE"}
	dataBytes, _ := json.Marshal(&data)
	dataStr := string(dataBytes)

	first16AppSecret := []byte(appSecret)[0:16]
	encrypted := util.AesEncryptCBC(dataBytes, first16AppSecret)
	base64Encrypted := base64.StdEncoding.EncodeToString(encrypted)
	finalData := "{\"data\":" + "\"" + base64Encrypted + "\"" + "}"
	//dataStr = "{"a":1}"
	time := util.GetNowMillisecond()
	timeStr := strconv.FormatInt(time, 10)
	//timeStr := "1676340948931"
	//String plaintext = appId + timestamp + appSecret + jsonString;
	joinStr := appId + timeStr + appSecret + dataStr
	sign := util.SHA1_1(joinStr)
	sigLower := strings.ToLower(sign)
	util.MyPrint("app-id:", appId, "appSecret:", appSecret, "data:", data, "time:", time, "timeStr", timeStr, "sign", sign, " sigLower:", sigLower, "FinalData:", finalData)

	rs := MiguRes{
		AppId:     appId,
		Time:      time,
		TimeStr:   timeStr,
		Sign:      sign,
		Data:      dataStr,
		DataBytes: dataBytes,
		FinalData: finalData,
		SignLower: sigLower,
	}

	httpresponse.OkWithAll(rs, "ok", c)
}

// @Tags Tools
// @Summary 测试咪咕,对方返回的数据信息
// @Description 120项目API接口aaa
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param data body request.SendSMS true "基础信息"
// @Produce  application/json
// @Success 200 {object} v1.MiguRes "最终的请求参数信息"
// @Router /tools/test/migu/api/backdata [POST]
func ReceiveMiguBackData(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		util.MyPrint(" ioutil.ReadAll err:" + err.Error())
	} else {
		util.MyPrint(string(body))
	}
	httpresponse.Ok(c)
}

// @Tags Tools
// @Summary 一个项目的详细信息
// @Description 用于开发工具测试
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
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
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
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
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Produce  application/json
// @Success 200 {object} []httpresponse.ConstInfo "常量列表"
// @Router /tools/const/list [get]
func ConstList(c *gin.Context) {
	//var a model.Project
	//c.ShouldBind(&a)
	//
	//list := GetConstList()
	//
	httpresponse.OkWithAll("aaa", "成功", c)

}

// @Tags Tools
// @Summary 常量列表 - 生成mysql导入GVA中的脚本
// @Description 给后台使用，生成到MYSQL数据库中，便于后台统一使用
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Security ApiKeyAuth
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
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
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
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
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
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
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
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
// @Param X-Client-Req-Time header string true "客户端请求时间 unix" default(1648277052)
// @Param X-Request-Id header string true "追踪ID" default(abcdefg)
// @Param X-Trace-Id header string true "请求ID" default(12345667)
// @Param X-Token header string true "追踪ID" default(12345667)
// @Param X-Sign header string true "签名" default(12345667)
// @Param X-Base-Info header string true "客户端信息" default({'sn':”,'pack_name':”,'app_version':”,'os':”,'os_version':”,'device':”,'device_version':”,'lat':”,'lon':”,'device_id':”,'dpi':”,'ip':”,'referer':”})
// @Success 200 {object} request.TestHeader
// @Router /tools/test/full/header [get]
func TestFullHeaderStruct(c *gin.Context) {
	header, _ := request.GetMyHeader(c)
	httpresponse.OkWithAll(header, "OK~d", c)

}

// @Tags Tools
// @Summary header头-结构体
// @Description 日常header里放一诸如验证类的东西，统一公示出来，仅是说明，方便测试/前端查看，方便使用
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID" Enums(1,2,3,4) default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Second-Auth-Uname header string true "二次验证-用户名" default(test)
// @Param X-Second-Auth-Ps header string true "二次验证-密码" default(qweASD1234560)
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
