package test

import (
	"os"
	"zgoframe/core/global"
	"zgoframe/util"
)

func Cicd(){
	//path1 := "/data/www/golang/src/logslave.go"
	//path2 := "/data/www/golang/src/metaverse-api"
	//file_fd, err1 := os.Stat(path1)
	//dir_fd, err2 := os.Stat(path2)
	//util.MyPrint(err1,err2," ",file_fd.IsDir() , " ",dir_fd.IsDir())
	//util.ExitPrint(33)

	/*依赖
		host.toml cicd.sh
		table:  project instance server cicd_publish
	*/

	opDirName := "operation"
	pwd , _ := os.Getwd()//当前路径]
	opDirFull := pwd + "/" + opDirName
	util.MyPrint(opDirFull,opDirName)

	cicdConfig := util.ConfigCicd{}
	configFile := opDirFull + "/host" + "." +"toml"

	//读取配置文件中的内容
	err := util.ReadConfFile(configFile,&cicdConfig)
	if err != nil{
		util.ExitPrint(err.Error())
	}

	util.PrintStruct(cicdConfig, " : ")
	//3方实例
	instanceManager ,_:= util.NewInstanceManager(global.V.Gorm)
	//服务器列表
	serverManger,_ := util.NewServerManger(global.V.Gorm)
	serverList := serverManger.Pool
	//发布管理
	publicManager := util.NewCICDPublicManager(global.V.Gorm)

	op := util.CicdManagerOption{
		HttpPort		: "1111",
		ServerList 		: serverList,
		Config			: cicdConfig,
		ServiceList		: global.V.ServiceManager.Pool,
		InstanceManager : instanceManager,
		PublicManager 	: publicManager,
		Log				: global.V.Zap,
		OpDirName		: opDirName,

	}

	cicd ,err := util.NewCicdManager(op)
	if err != nil{
		util.ExitPrint(err)
	}
	//cicd.GetSuperVisorList()
	cicd.DeployAllService()
	//go cicd.StartHttp(global.C.Http.StaticPath)
}


