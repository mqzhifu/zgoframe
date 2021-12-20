package test

import (
	"zgoframe/core/global"
	"zgoframe/util"
)

func test_cicd(){

	//config := Config{}

	cicdConfig := util.CicdConfig{}
	configFile := "host" + "." +"toml"
	util.ReadConfFile(configFile,&cicdConfig)
	util.PrintStruct(cicdConfig, " : ")


	instanceManager ,_:= util.NewInstanceManager(global.V.Gorm)

	hostServer,_ := util.NewHostServer(global.V.Gorm)
	hostServerList := hostServer.Pool

	publicManager := util.NewCICDPublicManager(global.V.Gorm)

	op := util.CicdManagerOption{
		ServerList :hostServerList,
		Config: cicdConfig,
		ServiceList: global.V.ServiceManager.Pool,
		InstanceManager :instanceManager,
		PublicManager :publicManager,
	}

	cicd := util.NewCicdManager(op)
	cicd.Init()
}


