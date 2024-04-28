package util

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Process struct {
	PathFileName             string
	CancelFunc               context.CancelFunc
	InitBaseInfoCallbackFunc func() string
	RootQuitFunc             func(source int)
	Log                      *zap.Logger
}

func NewProcess(ProcessPathFileName string, cancelFunc context.CancelFunc, log *zap.Logger, RootQuitFunc func(source int), InitBaseInfoCallbackFunc func() string) *Process {
	process := new(Process)
	process.PathFileName = ProcessPathFileName
	process.CancelFunc = cancelFunc
	process.RootQuitFunc = RootQuitFunc
	process.Log = log
	process.InitBaseInfoCallbackFunc = InitBaseInfoCallbackFunc
	return process
}

func (process *Process) InitProcess() {
	//主进程的ID号，存储文件
	pid, err := initPid(process.PathFileName)
	if err != nil {
		process.Log.Error("initPid,err:" + err.Error())
		return
	}

	process.Log.Warn("mainPid:" + strconv.Itoa(pid))
}

func (process *Process) DelPid() (int, error) {
	_, err := FileExist(process.PathFileName)
	if err != nil {
		return 0, errors.New(process.PathFileName + " not exist~ " + err.Error())
	}

	b, err := ioutil.ReadFile(process.PathFileName) // just pass the file name
	if err != nil {
		return 0, err
	}

	str := string(b)
	pid, _ := strconv.Atoi(str)

	err = os.Remove(process.PathFileName)

	MyPrint("del pid finish ", pid, err)

	return pid, err
}

// 进程PID保存到文件
func initPid(pathFile string) (int, error) {
	pid := os.Getpid()
	checkFd, err := FileExist(pathFile)
	if err != nil && checkFd != nil {
		return pid, errors.New(pathFile + " has exist~ " + err.Error())
	}

	fd, err := os.OpenFile(pathFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	defer fd.Close()
	if err != nil {
		return pid, errors.New(pathFile + " " + err.Error())
	}

	_, err = io.WriteString(fd, strconv.Itoa(pid))
	if err != nil {
		return pid, errors.New(pathFile + " " + err.Error())
	}

	return pid, err
}

// 信号 处理
func (process *Process) DemonSignal() {
	process.Log.Warn("SIGNAL init : ")
	c := make(chan os.Signal)
	//syscall.SIGHUP :ssh 挂断会造成这个信号被捕获，先注释掉吧
	//signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	prefix := "SIGNAL-DEMON :"
	for {
		sign := <-c
		process.Log.Warn(prefix)
		switch sign {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			process.Log.Warn(prefix + "SIGINT | SIGTERM | SIGQUIT  , exit!!!")
			process.CancelFunc()
			goto end
		//case syscall.SIGUSR1:
		//	   process.Log.Warn(prefix + " usr1!!!")
		//case syscall.SIGUSR2:
		//	   process.Log.Warn(prefix + " usr2!!!")

		default:
			process.Log.Warn(prefix + " unknow!!!")
		}
		time.Sleep(time.Second * 1)
	}
end:
	process.Log.Warn("DemonSignal DONE.")
}
