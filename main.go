package main

import (
	"context"
	"flag"
	"os"
	"strings"
	"time"
	"zgoframe/core/global"
	"zgoframe/core/initialize"
	_ "zgoframe/docs"
	"zgoframe/util"
)

// @title z golang 框架
// @version 0.1
// @description 拼装一个GO的基础框架方便日常使用

var mainDirName string
func main(){
	util.LogLevelFlag = util.LOG_LEVEL_DEBUG

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
	//开启测试模式
	//testFlag 			:= flag.String("t", "", "testFlag:empty or 1")
	//是否为CICD模式
	deploy 				:= flag.String("dep", "", "deploy")
	//解析命令行参数
	flag.Parse()

	//test(*testFlag)
	//return
	//检测环境变量是否正常
	if !util.CheckEnvExist(*env){
		util.ExitPrint(  "env is err , list:",envList)
	}

	pwd, _ := os.Getwd()
	pwdArr:=strings.Split(pwd,"/")
	mainDirName = pwdArr[len(pwdArr)-1]
	//开始初始化模块
	err := initialize.Init(*env,*configFileType,*configFileName,*configSourceType,*etcdUrl,mainDirName)
	if err != nil{
		util.MyPrint("nitialize.Init err:",err)
	}else{
		//主协程的 context
		mainCxt := context.Background()
		cancelCTX ,cancelFunc := context.WithCancel(mainCxt)
		//进程通信相关
		initialize.InitProcess()

		//这里才是 应用层代码需要DIY自己东西的入口
		applicationDoSomething()

		if *deploy == ""{
			initialize.DemonSignal(cancelFunc)
			select {
			case <-cancelCTX.Done():
				initialize.QuitAll()
			}
		}else{
			util.MyPrint("deploy: sleep 5 ,and auto quit...")
			time.Sleep(5)
		}

	}
	util.MyPrint("main end.")
}

func applicationDoSomething(){

}
