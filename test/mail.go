package test

import "zgoframe/core/global"

func Email() {

	to := "mqzhifu@sina.com"
	subject := "报警"
	msg := "程序出错了，请赶快修复"
	global.V.Email.SendOneEmailSync(to, subject, msg)
}
