package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zgoframe/core/global"
	"zgoframe/core/initialize"
	"zgoframe/util"
	_ "zgoframe/docs"

)
// @title z golang 框架
// @version 0.1
// @description 拼装一个GO的基础框架方便日常使用
func main(){
	util.LogLevelFlag = util.LOG_LEVEL_DEBUG

	envList := util.GetEnvList()


	configSourceType 		:= flag.String("cs", "file", "configSource:file or etcd")
	configFileType 		:= flag.String("ct", global.DEFAULT_CONFIT_TYPE, "configFileType")
	configFileName 	:= flag.String("cfn", global.DEFAULT_CONFIG_FILE_NAME, "configFileName")
	etcdUrl 	:= flag.String("etl", "", "get etcd config url")
	env 			:= flag.String("e", "must require", "env:loca test pre dev online")
	testFlag 		:= flag.String("t", "", "testFlag:empty or 1")

	flag.Parse()

	test(*testFlag)
	//return

	if !util.CheckEnvExist(*env){
		util.ExitPrint(  "env is err , list:",envList)
	}

	err := initialize.Init(*env,*configFileType,*configFileName,*configSourceType,*etcdUrl)
	if err != nil{
		util.MyPrint("nitialize.Init err:",err)
	}else{
		mainCxt := context.Background()
		cancelCTX ,cancelFunc := context.WithCancel(mainCxt)

		DemonSignal(cancelFunc)
		select {
		case <-cancelCTX.Done():
			Quit()
		}
	}
	util.MyPrint("main end.")
}

func test(testFlag string ){
	//if testFlag == ""{
	//	util.TestGrpcServer()
	//}else{
	//	util.TestGrpcClient()
	//}
}

func Quit(){
	global.V.Zap.Warn("main quit")
	initialize.Quit()
}
//信号 处理
func DemonSignal(cancelFunc context.CancelFunc){
	global.V.Zap.Warn("SIGNAL init : ")
	c := make(chan os.Signal)
	//syscall.SIGHUP :ssh 挂断会造成这个信号被捕获，先注释掉吧
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	prefix := "SIGNAL-DEMON :"
	for{
		sign := <- c
		global.V.Zap.Warn(prefix)
		switch sign {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			global.V.Zap.Warn(prefix+"SIGINT | SIGTERM | SIGQUIT  , exit!!!")
			cancelFunc()
			goto end
		case syscall.SIGUSR1:
			global.V.Zap.Warn(prefix+" usr1!!!")
		case syscall.SIGUSR2:
			global.V.Zap.Warn(prefix+" usr2!!!")
		default:
			global.V.Zap.Warn(prefix+" unknow!!!")
		}
		time.Sleep(time.Second * 1)
	}
end:
	global.V.Zap.Warn("DemonSignal DONE.")
}