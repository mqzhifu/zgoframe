package test

import (
	"zgoframe/core/global"
	"zgoframe/service/cicd"
)

func Cicd() {
	global.V.MyService.Cicd.Deploy.AllService(cicd.DEPLOY_TARGET_TYPE_REMOTE)
}
