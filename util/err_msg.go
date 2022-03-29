package util

import (
	"errors"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

const (
	CODE_NOT_EXIST = 5555
)

type ErrInfo struct {
	Code int
	Msg  string
}

type ErrMsg struct {
	LangPathFile string
	Log          *zap.Logger
	Pool         map[int]ErrInfo
}

func NewErrMsg(log *zap.Logger, langPathFile string) (*ErrMsg, error) {
	//log.Info("NewErrMsg class")
	errMsg := new(ErrMsg)
	errMsg.Log = log
	errMsg.LangPathFile = langPathFile
	errMsg.Pool = make(map[int]ErrInfo)

	err := errMsg.loadFileContent()

	return errMsg, err
}

func (errMsg *ErrMsg) loadFileContent() error {
	fileContentArr, err := ReadLine(errMsg.LangPathFile)
	var errContent string
	if err != nil {
		errContent = "errMsg loadFileContent ReadLine err :" + err.Error()
		errMsg.Log.Error(errContent)
		return errors.New(errContent)
	}

	if len(fileContentArr) <= 0 {
		errContent = "errMsg loadFileContent len <= 0"
		errMsg.Log.Error(errContent)
		return errors.New(errContent)
	}

	for _, v := range fileContentArr {
		row := strings.Split(v, "|")
		code, _ := strconv.Atoi(row[0])
		errInfo := ErrInfo{
			Code: code,
			Msg:  row[1],
		}

		_, ok := errMsg.Pool[code]
		if ok {
			errContent = "code " + row[0] + " has exist"
			errMsg.Log.Error(errContent)
			return errors.New(errContent)
		}

		errMsg.Pool[code] = errInfo
	}

	return nil
}

// 根据一个CODE，创建一个错误
func (errMsg *ErrMsg) New(code int) error {
	errInfo, ok := errMsg.Pool[code]
	if !ok {
		return errors.New(errMsg.Pool[CODE_NOT_EXIST].Msg)
	}
	return errors.New(errInfo.Msg)
}

// 根据一个CODE，创建一个错误，全不使用配置中的话术
func (errMsg *ErrMsg) NewMsg(code int, msg string) error {
	_, ok := errMsg.Pool[code]
	if !ok {
		return errors.New(errMsg.Pool[CODE_NOT_EXIST].Msg)
	}
	return errors.New(msg)
}

// 根据一个CODE，创建一个错误，并替换里面的动态值
func (errMsg *ErrMsg) NewReplace(code int, replace map[int]string) error {
	errInfo, ok := errMsg.Pool[code]
	if !ok {
		return errors.New(errMsg.Pool[CODE_NOT_EXIST].Msg)
	}
	for k, v := range replace {
		errInfo.Msg = strings.Replace(errInfo.Msg, "{"+strconv.Itoa(k)+"}", v, -1)
	}
	return errors.New(errInfo.Msg)
}
