package main

import (
	_ "embed"
	"context"
	"flag"
	"os"
	"strconv"
	"time"
	"zgoframe/core/global"
	"zgoframe/core/initialize"
	_ "zgoframe/docs"
	"zgoframe/util"
	"zgoframe/test"
)

var initializeVar *initialize.Initialize
// @title z golang 框架
// @version 0.1 测试版
// @description 拼装一个GO的基础框架方便日常使用，主要是想把经常用的类统一化，像：log 链路追踪 etcd等，保证项目高可用

func main(){
	//获取<环境变量>枚举值
	envList := util.GetEnvList()

	//配置读取源类型，1 文件  2 etcd
	configSourceType 	:= flag.String("cs", global.DEFAULT_CONFIG_SOURCE_TYPE, "configSource:file or etcd")
	//配置文件的类型
	configFileType 		:= flag.String("ct", global.DEFAULT_CONFIT_TYPE, "configFileType")
	//配置文件的名称
	configFileName 		:= flag.String("cfn", global.DEFAULT_CONFIG_FILE_NAME, "configFileName")
	//获取etcd 配置信息的URL
	etcdUrl 			:= flag.String("etl", "", "get etcd config url")
	//当前环境
	env 				:= flag.String("e", "must require", "env:local test pre dev online")
	//是否为CICD模式
	deploy 				:= flag.String("dep", "", "deploy")
	//开启测试模式
	//testFlag 			:= flag.String("t", "", "testFlag:empty or 1")
	//解析命令行参数
	flag.Parse()

	//test(*testFlag)
	//检测环境变量值ENV是否正常
	if !util.CheckEnvExist(*env){
		util.MyPrint(  "env is err , list:",envList)
		panic("env is err : "+*env)
	}

	pwd, _ := os.Getwd()//当前路径
	util.MyPrint("now pwd:"+pwd)
	//开始初始化模块
	//主协程的 context
	mainCxt,mainCancelFunc := context.WithCancel(context.Background())
	initOption := initialize.InitOption  {
		Env 				:*env,
		ConfigType 			:*configFileType,
		ConfigFileName 		:*configFileName,
		ConfigSourceType 	:*configSourceType,
		EtcdConfigFindUrl	:*etcdUrl,
		RootDir				:pwd,
		RootCtx				:mainCxt,
		RootCancelFunc		:mainCancelFunc,
		RootQuitFunc 		:QuitAll,
	}
	//开始正式全局初始化
	initializeVar = initialize.NewInitialize(initOption)
	err := initializeVar.Start()
	if err != nil{
		util.MyPrint("initialize.Init err:",err)
		panic("initialize.Init err:"+err.Error())
		return
	}
	//测试
	test.Index()
	//deploy 是用来方便布置的，看着挺恶心，我回头再优化
	if *deploy == ""{
		go global.V.Process.DemonSignal()
		util.MyPrint("wait mainCxt.done...")
		select {
		case <-mainCxt.Done():
			QuitAll(1)
		}
	}else{
		util.MyPrint("deploy: sleep 5 ,and auto quit...")
		time.Sleep(5)
	}

	util.MyPrint("main end.")
}

func QuitAll(source int){
	defer func() {
		global.V.Process.DelPid()
	}()

	global.V.Zap.Warn("main quit , source : " + strconv.Itoa(source))
	initializeVar.Quit()


	util.MyPrint("QuitAll finish.")
}

