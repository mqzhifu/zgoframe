package main

import (
	"context"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
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


	deploy 	:= flag.String("dep", "", "deploy")

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

		pathFilePid := "/tmp/zgoframe.pid"
		pid ,err := initPid(pathFilePid)
		if err != nil{
			global.V.Zap.Error("initPid,err:" + err.Error())
			return
		}

		global.V.Zap.Warn("mainPid:"+strconv.Itoa(pid))


		global.V.Gin.GET("/sys/quit",httpQuit)
		global.V.Gin.GET("/sys/config",getConfig)

		if *deploy == ""{
			DemonSignal(cancelFunc)
			select {
			case <-cancelCTX.Done():
				Quit()
			}
		}else{
			util.MyPrint("deploy: sleep 5 ,and auto quit...")
			time.Sleep(5)
		}


		pid ,err = delPid(pathFilePid)
		util.MyPrint("del pid:",pid,err)

	}
	util.MyPrint("main end.")
}
func httpQuit(c *gin.Context){
	Quit()
}

func getConfig(c *gin.Context){
	Quit()
}

func delPid(pathFile string )(int,error){
	if !util.CheckFileIsExist(pathFile){
		return 0,errors.New(pathFile + " not exist~ ")
	}

	b, err := ioutil.ReadFile(pathFile) // just pass the file name
	if err != nil {
		return 0,err
	}

	str := string(b)
	pid ,_ := strconv.Atoi(str)

	err = os.Remove(pathFile)

	return pid,err
}
//进程PID保存到文件
func initPid(pathFile string )(int,error){
	pid := os.Getpid()
	if util.CheckFileIsExist(pathFile){
		return pid,errors.New(pathFile + " has exist~ ")
	}

	fd, err  := os.OpenFile(pathFile, os.O_WRONLY | os.O_CREATE | os.O_TRUNC , 0777)
	defer fd.Close()
	if err != nil{
		return pid,errors.New(pathFile + " " + err.Error())
	}

	_, err = io.WriteString(fd, strconv.Itoa(pid))
	if err != nil{
		return pid,errors.New(pathFile + " " + err.Error())
	}

	return pid,err
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