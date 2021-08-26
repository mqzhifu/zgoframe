package util

import (
	"go.uber.org/zap"
	"runtime"
	"strconv"
	"sync"
	"context"
)


type RecoverGo struct{
	RetryTimes  map[string]int
	RetryTimesRWLock *sync.RWMutex
	Log	*zap.Logger
}

func NewRecoverGo(log	*zap.Logger)*RecoverGo{
	recoverGo := new(RecoverGo)
	recoverGo.RetryTimes  = make(map[string]int)
	recoverGo.RetryTimesRWLock = &sync.RWMutex{}
	recoverGo.Log = log
	return recoverGo
}

func (recoverGo *RecoverGo) RecoverGoRoutine(back func(ctx context.Context),ctx context.Context,err interface{}){
	pc, file, lineNo, ok := runtime.Caller(3)
	if !ok {
		recoverGo.Log.Error("runtime.Caller ok is false :")
	}
	funcName := runtime.FuncForPC(pc).Name()
	recoverGo.Log.Info(" RecoverGoRoutine  panic in defer  :"+ funcName + " "+file + " "+ strconv.Itoa(lineNo))
	recoverGo.RetryTimesRWLock.RLock()
	retryTimes , ok := recoverGo.RetryTimes[funcName]
	recoverGo.RetryTimesRWLock.RUnlock()
	if ok{
		if retryTimes > 3{
			recoverGo.Log.Error("retry than max times")
			panic(err)
			return
		}else{
			recoverGo.RetryTimesRWLock.Lock()
			recoverGo.RetryTimes[funcName]++
			recoverGo.RetryTimesRWLock.Unlock()
			recoverGo.Log.Info("RecoverGoRoutineRetryTimes = " + strconv.Itoa(recoverGo.RetryTimes[funcName]))
		}
	}else{
		recoverGo.Log.Info("RecoverGoRoutineRetryTimes = 1")
		recoverGo.RetryTimesRWLock.Lock()
		recoverGo.RetryTimes[funcName] = 1
		recoverGo.RetryTimesRWLock.Unlock()
	}
	go back(ctx)
}
