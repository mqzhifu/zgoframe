package core

import (
	"zgoframe/core/global"
	"zgoframe/http/request"
	"zgoframe/test"
	"zgoframe/util"
)

func DoMySelf() {

}

//test_command()
//Gateway()
func DoTestAction(flag string) {
	//util.GetHTTPBaseAuth()
	util.MyPrint("=======================", flag)
	switch flag {
	case "alert":
		global.V.MyService.Alert.Send(6, "商品库存不足，请及时补充货源", "warning")
	case "db_table":
		sqlList := global.AutoCreateUpDbTable()
		sqlStrings := ""
		for _, v := range sqlList {
			sqlStrings += v
		}
		util.MyPrint(sqlStrings)
		util.ExitPrint("i want exit 2.")
	case "ali_sms":
		//AliSms := util.NewAliSms()
		//AliSms.Send()
	case "email":
		test.Email()
	case "sms":
		test.Sms()
	case "service_sms":
		SendSMS := request.SendSMS{}
		recordNewId, err := global.V.MyService.Sms.Send(6, SendSMS)
		util.MyPrint(recordNewId, err)
	case "cicd":
		test.Cicd()
	case "ProjectAutoCreateUserDbRecord":
		test.ProjectAutoCreateUserDbRecord()
	case "alert_push":
		//global.V.AlertPush.Push(1, "error", "test push alert info.")
	case "grpc":
		test.Grpc()
	case "super_visor":
		test.Test_command()
		test.Test_supervisor()
	case "gateway":
		test.Gateway()
	default:
		util.ExitPrint("DoTestAction flag no hit , flag:" + flag)
	}

}
