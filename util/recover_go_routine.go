package util

import (
	"go.uber.org/zap"
	"runtime"
	"strconv"
	"sync"
	"context"
)
/*
	实现功能
	1. 一些重点的协程，守护模式，一但挂了，需要重新拉起
	2. 拉起次数可自定义
	ps:实际上看，这个类只是帮忙统一处理了重新次数功能，没其它实际意义

	使用：
	1. 调用者，在内部使用defer recover 捕获到 panic 之后调用此方法
*/
type RecoverGo struct{
	RetryTimes  map[string]int//map[函数名]已重试次数
	RetryMaxTimes int
	RetryTimesRWLock *sync.RWMutex
	Log	*zap.Logger
}

func NewRecoverGo(log	*zap.Logger , retryMaxTimes int)*RecoverGo{
	recoverGo := new(RecoverGo)
	recoverGo.RetryTimes  = make(map[string]int)
	recoverGo.RetryTimesRWLock = &sync.RWMutex{}
	recoverGo.Log = log

	if retryMaxTimes <= 0{
		//黑鹰最大3
		retryMaxTimes = 3
	}

	if retryMaxTimes > 10{
		retryMaxTimes = 10
	}
	recoverGo.RetryMaxTimes = retryMaxTimes

	log.Info("NewRecoverGo , retryMaxTimes:"+strconv.Itoa(retryMaxTimes))

	return recoverGo
}
/*
	callback:回调|启动函数，该函数接收2个参数 context 和 任意个参数
	ctx:启动新函数时，要把这个值传递回去
	err:错误信息
	v2:启动新函数时，要把这个值传递回去
*/
func (recoverGo *RecoverGo) RecoverGoRoutine(callback func(ctx context.Context,v ...interface{}) , ctx context.Context , err interface{},v2 ...interface{}){
	//0 当前函数  1及之后，就是逐级调用的函数信息,注：这里的skip(数字)不能大于调用最大层级
	pc, file, lineNo, ok := runtime.Caller(3)
	if !ok {
		recoverGo.Log.Error("runtime.Caller <ok> is false.")
	}
	funcName := runtime.FuncForPC(pc).Name()//出现panic|fatal 的函数名，也就是调用此函数的函数
	recoverGo.Log.Info(" RecoverGoRoutine  panic in defer  :"+ funcName + " "+file + " "+ strconv.Itoa(lineNo))
	//加锁，避免错误使用，导致多次调用此函数，多次重新启动新的函数协程
	recoverGo.RetryTimesRWLock.RLock()
	retryTimes , ok := recoverGo.RetryTimes[funcName]
	recoverGo.RetryTimesRWLock.RUnlock()
	if ok{//map-key已存在，证明之前已经挂过并创建了map key
		if retryTimes > recoverGo.RetryMaxTimes{
			recoverGo.Log.Error("retry than max times")
			panic(err)
			return
		}else{
			recoverGo.RetryTimesRWLock.Lock()
			recoverGo.RetryTimes[funcName]++
			recoverGo.RetryTimesRWLock.Unlock()
			recoverGo.Log.Info("RecoverGoRoutineRetryTimes = " + strconv.Itoa(recoverGo.RetryTimes[funcName]))
		}
	}else{//map-key不存在，证明之前没有挂过
		recoverGo.Log.Info("RecoverGoRoutineRetryTimes = 1")
		recoverGo.RetryTimesRWLock.Lock()
		recoverGo.RetryTimes[funcName] = 1
		recoverGo.RetryTimesRWLock.Unlock()
	}
	//重新启动该函数
	go callback(ctx,v2)
}
