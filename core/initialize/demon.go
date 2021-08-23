package initialize

import (
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"zgoframe/core/global"
	"context"
	"zgoframe/util"
	"errors"
)

func getPidFileName()string{
	return  "/tmp/"+global.V.App.Name+".pid"
}

func InitProcess(){
	//主进程的ID号，存储文件
	pid ,err := initPid(getPidFileName())
	if err != nil{
		global.V.Zap.Error("initPid,err:" + err.Error())
		return
	}

	global.V.Zap.Warn("mainPid:"+strconv.Itoa(pid))

}

func HttpQuit(c *gin.Context){
	Quit()
}

func GetConfig(c *gin.Context){
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

func QuitAll(){
	global.V.Zap.Warn("main quit")
	Quit()

	pid ,err := delPid(getPidFileName())
	util.MyPrint("del pid:",pid,err)
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