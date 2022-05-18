package util

import (
	"errors"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

const (
	CODE_NOT_EXIST = 5555
	ERR_separate = "-_-"
)

type ErrInfo struct {
	Code int
	Msg  string
}

func (errInfo *ErrInfo) Error()string{
	return errInfo.Msg
}

func (errInfo *ErrInfo) GetCode()int{
	return errInfo.Code
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

func (errMsg ErrMsg) loadFileContent() error {
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
		e := errMsg.Pool[CODE_NOT_EXIST]
		return & e
	//	//return errors.New(errMsg.Pool[CODE_NOT_EXIST].Msg)
	//	return & errMsg.Pool[CODE_NOT_EXIST]
	}
	//return errors.New(errInfo.Msg)
	//e := errors.New(errInfo.Msg)
	//e.Error()
	errInfo.Msg = strconv.Itoa(code) + ERR_separate + errInfo.Msg
	return & errInfo
}

//// 根据一个CODE，创建一个错误
//func (errMsg *ErrMsg) New(code int) error {
//	errInfo, ok := errMsg.Pool[code]
//	if !ok {
//		//return errors.New(errMsg.Pool[CODE_NOT_EXIST].Msg)
//		return & errMsg.Pool[CODE_NOT_EXIST]
//	}
//	//return errors.New(errInfo.Msg)
//	//e := errors.New(errInfo.Msg)
//	//e.Error()
//	return errInfo
//}


// 根据一个CODE，创建一个错误，全不使用配置中的话术
func (errMsg *ErrMsg) NewMsg(code int, msg string) error {
	errInfo, ok := errMsg.Pool[code]
	if !ok {
		//return errors.New(errMsg.Pool[CODE_NOT_EXIST].Msg)
		e := errMsg.Pool[CODE_NOT_EXIST]
		return &e
	}
	errInfo.Msg = strconv.Itoa(code) + ERR_separate + msg
	//return errors.New(msg)
	return &errInfo
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
	errInfo.Msg = strconv.Itoa(code) + ERR_separate + errInfo.Msg
	return & errInfo
}

func (errMsg *ErrMsg)MakeOneStringReplace(str string)map[int]string{
	msg := make(map[int]string)
	msg[0] = str
	return msg
}

func (errMsg *ErrMsg)SplitMsg(msg string)(code int ,eMsg string ,err error){
	list := strings.Split(msg,ERR_separate)
	MyPrint("errMsg SplitMsg:"+msg , " list:",list)
	if len(list )== 2 {
		code ,_ = strconv.Atoi(list[0])
		eMsg = list[1]
		return code,msg,nil
	}
	return code,eMsg,errors.New("len != 2"+strconv.Itoa(len(list )))
}

