package test

import (
	"zgoframe/core/global"
	"zgoframe/util"
)

func cicd(){

	/*依赖
		host.toml cicd.sh
		table:  project instance server cicd_publish
	*/
	cicdConfig := util.ConfigCicd{}
	configFile := "host" + "." +"toml"
	//读取配置文件中的内容
	util.ReadConfFile(configFile,&cicdConfig)
	util.PrintStruct(cicdConfig, " : ")
	//3方实例
	instanceManager ,_:= util.NewInstanceManager(global.V.Gorm)
	//服务器列表
	serverManger,_ := util.NewServerManger(global.V.Gorm)
	serverList := serverManger.Pool
	//发布管理
	publicManager := util.NewCICDPublicManager(global.V.Gorm)

	op := util.CicdManagerOption{
		ServerList 		:serverList,
		Config			: cicdConfig,
		ServiceList		: global.V.ServiceManager.Pool,
		InstanceManager :instanceManager,
		PublicManager 	:publicManager,
	}

	cicd := util.NewCicdManager(op)

	//cicd.GetSuperVisorList()
	cicd.DeployAllService()
	//go cicd.StartHttp(global.C.Http.StaticPath)
}


