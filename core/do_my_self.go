package core

import (
	"time"
	"zgoframe/core/global"
	"zgoframe/test"
	"zgoframe/util"
)

func DoMySelf(){

	global.V.AlertPush.Push(1,"error","test push alert info.")
	time.Sleep(time.Second * 1)
	util.ExitPrint(22)
	//测试
	//test.Index()
	if global.C.System.ProjectId == 3{
		test.LogSlave()
	}else{
		//test.Grpc()
		//test.Gateway()
		//test.Cicd()
		test.Email()
	}
}