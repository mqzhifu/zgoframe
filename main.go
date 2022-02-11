package main

import (
	"context"
	_ "embed"
	"flag"
	"os"
	"os/user"
	"strconv"
	"zgoframe/core"
	"zgoframe/core/global"
	"zgoframe/core/initialize"
	_ "zgoframe/docs"
	"zgoframe/util"
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
	etcdUrl 			:= flag.String("etl", "http://127.0.0.1/getEtcdCluster/Ip/Port", "get etcd config url")
	//当前环境,env:local test pre dev online
	env 				:= flag.String("e", "", "must require , env:local test pre dev online")
	//DEBUG模式
	debug 				:= flag.Int("debug", 0, "startup debug mode level")
	//是否为CICD模式
	//deploy 				:= flag.String("dep", "", "deploy")//部署模式下，启动程序只是为了测试脚本正常，因为之后，要立刻退出
	//开启测试模式
	//testFlag 			:= flag.String("t", "", "testFlag:empty or 1")
	//解析命令行参数
	flag.Parse()
	//检测环境变量值ENV是否正常
	if !util.CheckEnvExist(*env){
		msg := "argv env , is err :"
		util.MyPrint(  msg,envList)
		panic(msg  + *env)
	}

	imUser , _ := user.Current()
	util.MyPrint("exec script user info , name: "+imUser.Name + " uid: " + imUser.Uid  +  " , gid :"+ imUser.Gid + " ,homeDir:" +imUser.HomeDir)

	//u2, _ := user.Lookup(imUser.Name)
	//util.ExitPrint(u2.Name)

	pwd, _ := os.Getwd()//当前路径
	util.MyPrint("exec script pwd:"+pwd)
	//开始初始化模块
	//主协程的 context
	util.MyPrint("create main :cancel context")
	mainCxt,mainCancelFunc := context.WithCancel(context.Background())
	//初始化模块需要的参数
	initOption := initialize.InitOption  {
		Env 				:*env,
		Debug				:*debug,
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
	go core.DoMySelf()
	//监听外部进程信号
	go global.V.Process.DemonSignal()
	util.MyPrint("wait mainCxt.done...")
	select {
	case <-mainCxt.Done():
		QuitAll(1)
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

