package test

import (
	"zgoframe/core/global"
	"zgoframe/util"
)

func Email(){
	toEmail := "mqzhifu@sina.com"
	err := global.V.Email.SendOneEmailSync(toEmail,"alert","test_test_test_test_test_test_test_test_test")
	util.ExitPrint(err)
}
