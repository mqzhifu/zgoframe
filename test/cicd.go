package test

import (
	"zgoframe/core/global"
	"zgoframe/util"
)

func test_cicd(){

	//config := Config{}

	cicdConfig := util.ConfigCicd{}
	configFile := "host" + "." +"toml"
	util.ReadConfFile(configFile,&cicdConfig)
	util.PrintStruct(cicdConfig, " : ")


	instanceManager ,_:= util.NewInstanceManager(global.V.Gorm)

	serverManger,_ := util.NewServerManger(global.V.Gorm)
	serverList := serverManger.Pool

	publicManager := util.NewCICDPublicManager(global.V.Gorm)

	op := util.CicdManagerOption{
		ServerList :serverList,
		Config: cicdConfig,
		ServiceList: global.V.ServiceManager.Pool,
		InstanceManager :instanceManager,
		PublicManager :publicManager,
	}

	cicd := util.NewCicdManager(op)

	//cicd.GetSuperVisorList()
	//cicd.Init()
	go cicd.StartHttp(global.C.Http.StaticPath)
}


