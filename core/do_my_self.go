package core

import (
	"zgoframe/core/global"
	"zgoframe/test"
	"zgoframe/util"
)

func DoMySelf() {

}

//test_command()
//Gateway()
func DoTestAction(flag string) {
	util.GetHTTPBaseAuth()

	switch flag {
	case "db_table":
		sqlList := global.AutoCreateUpDbTable()
		sqlStrings := ""
		for _, v := range sqlList {
			sqlStrings += v
		}
		util.MyPrint(sqlStrings)
		util.ExitPrint("i want exit 2.")
	case "email":
		test.Email()
	case "cicd":
		test.Cicd()
	case "alert_push":
		global.V.AlertPush.Push(1, "error", "test push alert info.")
	case "grpc":
		test.Grpc()
	case "super_visor":
		test.Test_command()
		test.Test_supervisor()
	case "gateway":
		//test.gateway()
	default:
		util.ExitPrint("DoTestAction flag no hit , flag:" + flag)
	}

}
