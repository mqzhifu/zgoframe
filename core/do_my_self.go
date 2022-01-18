package core

import (
	"zgoframe/core/global"
	"zgoframe/test"
)

func DoMySelf(){
	//测试
	//test.Index()
	if global.C.System.ProjectId == 3{
		test.LogSlave()
	}else{
		//test.Grpc()
		test.Gateway()
	}
}