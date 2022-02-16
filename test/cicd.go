package test

import (
	"os"
	"zgoframe/core/global"
	"zgoframe/util"
)

func Cicd(){

	testPushCounterName := "testPushCounter"

	global.V.Metric.CreateCounter(testPushCounterName,"im_test_counter")
	global.V.Metric.CounterInc(testPushCounterName)

	testPushGaugeName := "testPushGauge"

	global.V.Metric.CreateGauge(testPushGaugeName,"im_test_gauge")
	global.V.Metric.GaugeSet(testPushGaugeName,0.001)


	push_err := global.V.Metric.PushMetrics()
	util.MyPrint("test pusher err:",push_err)
	return
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

	util.PrintStruct(cicdConfig , " : ")
	//3方实例
	instanceManager ,_:= util.NewInstanceManager(global.V.Gorm)
	//服务器列表
	serverManger,_ := util.NewServerManger(global.V.Gorm)
	serverList := serverManger.Pool
	//发布管理
	publicManager := util.NewCICDPublicManager(global.V.Gorm)

	//util.ExitPrint(22)
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
	//生成 filebeat 配置文件
	cicd.GenerateAllFilebeat()
	//cicd.GetSuperVisorList()
	//部署所有机器上的所有服务项目
	//cicd.DeployAllService()
	//go cicd.StartHttp(global.C.Http.StaticPath)
}


