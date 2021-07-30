package util

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strconv"
	"time"
	"errors"
	"fmt"
)

const (
	LEVEL_INFO 		=	 1 << iota
	LEVEL_DEBUG		//2
	LEVEL_ERROR		//4
	LEVEL_PANIC		//8

	LEVEL_EMERGENCY	//16
	LEVEL_ALERT		//32
	LEVEL_CRITICAL	//64
	LEVEL_WARNING	//128
	LEVEL_NOTICE	//256
	LEVEL_TRACE		//512

	LEVEL_ALL 		= LEVEL_INFO | LEVEL_DEBUG | LEVEL_ERROR | LEVEL_PANIC | LEVEL_EMERGENCY |LEVEL_ALERT| LEVEL_CRITICAL |LEVEL_WARNING |LEVEL_NOTICE|LEVEL_TRACE
	LEVEL_DEV 		= LEVEL_INFO | LEVEL_DEBUG | LEVEL_ERROR | LEVEL_PANIC | LEVEL_TRACE
	LEVEL_ONLINE 	= LEVEL_INFO | LEVEL_ERROR | LEVEL_PANIC

	FILE_HASH_NONE = 0
	FILE_HASH_MONTH = 1
	FILE_HASH_DAY = 2
	FILE_HASH_HOUR = 3

	CONTENT_TYPE_STRING = 0
	CONTENT_TYPE_JSON = 1

)

const (
	OUT_TARGET_SC = 	 1 << iota
	OUT_TARGET_FILE
	OUT_TARGET_NET

	OUT_TARGET_ALL = OUT_TARGET_SC|OUT_TARGET_FILE|OUT_TARGET_NET

	OUT_TARGET_NET_TCP  = 1
	OUT_TARGET_NET_UDP = 2
)

var levelContentPrefixes = map[int]string{
	LEVEL_INFO		: "INFO",
	LEVEL_DEBUG		: "DEBUG",
	LEVEL_ERROR		: "ERROR",
	LEVEL_PANIC		: "PANIC",
	LEVEL_EMERGENCY	: "EMERG",
	LEVEL_ALERT		: "ALERT",
	LEVEL_CRITICAL	: "CRITI",
	LEVEL_WARNING	: "WARNI",
	LEVEL_NOTICE	: "NOTIC",
	LEVEL_TRACE		: "TRACE",
}

type Msg struct {
	AppId 		int		`json:"appId"`
	ModuleId	int		`json:"moduleId"`
	LevelPrefix string	`json:"levelPrefix"`
	Content 	string	`json:"content"`
	Header 		string	`json:"header"`
}


type Log struct {
	Option LogOption
	InChan 	chan Msg	`json:"-"`
	CloseChan chan int	`json:"-"`
}

type LogOption struct {
	AppId 			int
	ModuleId		int
	OutTarget 		int		//输出目标介质
	OutContentType  int		//输出的内容格式类型
	Level 			int		//当前类可以输出的等级类型，如DEBUG模式，所有日志类型均输出

	OutFileBasePath 	string	//输出到文件的：基础目录
	OutFilePath 		string	//输出到文件的：最终目录 = 基础目录 + 当前项目类型
	OutFilePathFile 	string 	//最终真实输出的： 路径+文件名+扩展名
	OutFileFileName 	string	//输出到文件的：文件名
	OutFileFileExtName 	string	//输出到文件的：文件扩展名
	OutFileHashType		int		//文件的保存形式，使用HASH模式

	OutFileFileFd		*os.File	`json:"-"`
	OutFileFileFdOpenTime string


}
func NewLog( logOption LogOption)(log *Log ,errs error){
	//MyPrint("New log class ,OutFilePath : ",logOption.OutFilePath , " level : ",logOption.Level ," target : ",logOption.Target)

	if logOption.Level == 0{
		return nil,errors.New(" level is empty ")
	}

	if logOption.OutTarget == 0 {
		return nil,errors.New(" target is empty ")
	}

	log = new(Log)
	log.InChan = make(chan Msg,10000)
	log.CloseChan = make(chan int)
	log.Option = logOption

	//如果要输出到文件中，要做判断，并提前找到FD
	if log.checkTargetIncludeByBit(OUT_TARGET_FILE){
		err := log.OpenNewFd()
		if err != nil{
			return nil,err
		}
	}

	log.Debug(logOption)

	go log.loopRealWriteMsg()
	return log,nil
}
func  (log *Log)OpenNewFd( )(err error ){
	if log.Option.OutFilePath == ""{
		return errors.New(" OutFilePath is empty ")
	}

	if log.Option.OutFileFileName == ""{
		return errors.New(" OutFilePath is empty ")
	}

	errs := log.checkOutFilePathPower(log.Option.OutFilePath)
	if errs != nil{
		return errs
	}

	//log.GetPathFile()
	loutFilePathFileTmp := log.Option.OutFilePath + "/" + log.Option.OutFileFileName
	switch log.Option.OutFileHashType {
	case FILE_HASH_NONE:
		break
	case FILE_HASH_HOUR:
		dateHour := GetNowDateHour()
		log.Option.OutFileFileFdOpenTime = dateHour
		loutFilePathFileTmp += dateHour
	case FILE_HASH_DAY:
		dateDay := GetNowDate()
		log.Option.OutFileFileFdOpenTime = dateDay
		loutFilePathFileTmp += dateDay
	case FILE_HASH_MONTH:
		date := GetNowDateMonth()
		log.Option.OutFileFileFdOpenTime =date
		loutFilePathFileTmp += date
	}
	pathFile := loutFilePathFileTmp + "." + log.Option.OutFileFileExtName
	fd, err  := os.OpenFile(pathFile, os.O_WRONLY | os.O_CREATE | os.O_APPEND , 0777)
	if err != nil{
		return errors.New(" log out file , OpenFile :  " + err.Error())
	}

	log.Option.OutFileFileFd = fd
	log.Option.OutFilePathFile = pathFile
	if log.Option.OutFileFileExtName == ""{
		log.Option.OutFileFileExtName = "log"
	}

	return nil
}

func  (log *Log)checkFileFdTimeout()bool{
	timeStr := ""
	switch log.Option.OutFileHashType {
	case FILE_HASH_NONE:
		return false
	case FILE_HASH_HOUR:
		timeStr = GetNowDateHour()
	case FILE_HASH_DAY:
		timeStr = GetNowDate()
	case FILE_HASH_MONTH:
		timeStr = GetNowDateMonth()
	}
	if timeStr == log.Option.OutFileFileFdOpenTime{
		return false
	}else{
		return true
	}
}

//permission
func  (log *Log)  checkOutFilePathPower(path string)error{
	if path == ""{
		return errors.New(" checkOutFilePathPower ("+path+") : path is empty")
	}

	fd,e :=  os.Stat(path)
	if e != nil{
		return errors.New(" checkOutFilePathPower ("+path+"): os.Stat : "+ e.Error())
	}

	if !fd.IsDir(){
		return errors.New(" checkOutFilePathPower ("+path+"): path is not a dir ")
	}
	perm := fd.Mode().Perm().String()
	//MyPrint(perm,os.FileMode(0755).String())
	//log.Debug(fd.Mode(),fd.Mode().Perm())
	if perm < os.FileMode(0755).String(){
		return errors.New(" checkOutFilePathPower ("+path+"):path permission 0777 ")
	}
	return nil
}
func  linkOutFilePath(basePath string,category string)string{
	return basePath + "/" + category + "/"
}
func (log *Log) Info(content ...interface{}){
	log.Out(LEVEL_INFO,content)
}

func (log *Log) Debug(content ...interface{}){
	log.Out(LEVEL_DEBUG,content...)
}

func (log *Log) Error(content ...interface{}){
	log.Out(LEVEL_ERROR,content...)
}

func (log *Log) Notice(content ...interface{}){
	log.Out(LEVEL_NOTICE,content...)
}

func (log *Log) Warning(content ...interface{}){
	log.Out(LEVEL_WARNING,content...)
}

func (log *Log) Alert(content ...interface{}){
	log.Out(LEVEL_ALERT,content...)
}

func (log *Log) Panic(content ...interface{}){
	log.Out(LEVEL_PANIC,content...)
}

func (log *Log) OutScreen(a ...interface{}){
	if a[0] == "[INFO]"{
		fmt.Printf("%c[1;40;33m%s%c[0m", 0x1B, a[0], 0x1B)
	}else if a[0] == "[ERROR]" {
		fmt.Printf("%c[1;40;31m%s%c[0m", 0x1B, a[0], 0x1B)
	}else if a[0] == "[NOTIC]" {
		fmt.Printf("%c[1;40;34m%s%c[0m", 0x1B, a[0], 0x1B)
	}else if a[0] == "[ALERT]" {
		fmt.Printf("%c[1;40;36m%s%c[0m", 0x1B, a[0], 0x1B)
	}else if a[0] == "[WARNI]" {
		fmt.Printf("%c[1;40;35m%s%c[0m", 0x1B, a[0], 0x1B)
	}else if a[0] == "[PANIC]" {
		fmt.Printf("%c[1;40;41m%s%c[0m", 0x1B, a[0], 0x1B)
	}else{
		fmt.Printf("%c[1;40;32m%s%c[0m", 0x1B, a[0], 0x1B)
	}

	newlist := append(a[:0], a[(0+1):]...)
	fmt.Println(newlist...)
}
//func (log *Log) GetPathFile()string{
//	return log.option.OutFilePath + "/" + log.option.OutFileFileName
//}

func (log *Log) OutFile(content string){
	_, err := io.WriteString(log.Option.OutFileFileFd, content)
	if err != nil{
		log.Error("OutFile io.WriteString : ",err.Error())
	}
}

func (log *Log) CloseFileFd()error{
	if !log.checkTargetIncludeByBit(OUT_TARGET_FILE){
		return errors.New("checkTargetIncludeByBit OUT_TARGET_FILE :false")
	}
	err := log.Option.OutFileFileFd.Close()
	return err
}

func (log *Log)getHeaderContentStr()string{
	timeStr:=time.Now().Format("2006-01-02 15:04:05")
	//unixstamp := GetNowTimeSecondToInt()
	//msTimeStr :=  GetNowMillisecond()
	//uuid4 := getUuid4()
	pid  := os.Getpid()
	str := timeStr + "[" + strconv.Itoa(pid) + "]"
	//str :=   "[" + strconv.Itoa(pid) + "]"
	//str := strconv.Itoa(int(msTimeStr)) + "[" + strconv.Itoa(pid) + "]"
	return str
}

func (log *Log) checkTargetIncludeByBit(flag int)bool{
	if log.Option.OutTarget & flag == flag {
		return true
	}
	return false
}

func (log *Log) checkLevelIncludeByBit(level int)bool{
	//MyPrint(log.option.Level,level)
	if log.Option.Level & level == level {
		return true
	}
	return false
}

func  (log *Log)Out(level int ,argcs ...interface{}){
	if !log.checkLevelIncludeByBit(level){
		return
	}
	contentLevelPrefix := "[" + levelContentPrefixes[level] +"]"
	content := ""
	for _,argc := range argcs{
		content += " " + log.String(argc)
	}

	msg := Msg{
		LevelPrefix: contentLevelPrefix,
		Content: content,
		Header: log.getHeaderContentStr(),
	}

	log.InChan <-msg
}
func  (log *Log)SlaveOut(level int ,msg Msg){
	contentLevelPrefix := "[" + levelContentPrefixes[level] +"]"
	msg.LevelPrefix = contentLevelPrefix

	log.InChan <-msg
}


func  (log *Log)loopRealWriteMsg(){
	MyPrint("zlog : loopRealWriteMsg start")
	isBreak := 0
	for{
		select {
		case msg := <- log.InChan:
			log.Write(msg)
		case <- log.CloseChan:
			isBreak = 1
		}
		if isBreak == 1{
			goto end
		}
	}
end:
	MyPrint("ctx.done() zlib.log - loopRealWriteMsg")
	log.CloseFileFd()
}

func  (log *Log)Write(msg Msg){
	if log.checkTargetIncludeByBit(OUT_TARGET_FILE){
		if msg.AppId == 0{
			msg.AppId = log.Option.AppId
		}

		if msg.ModuleId == 0{
			msg.ModuleId = log.Option.ModuleId
		}

		str := ""
		switch log.Option.OutContentType {
		case CONTENT_TYPE_JSON:
			strByte,_ := json.Marshal(&msg)
			str = string(strByte)
		case CONTENT_TYPE_STRING:
			str = msg.LevelPrefix + msg.Header + msg.Content
		}
		if log.checkFileFdTimeout(){
			log.OpenNewFd()
		}
		log.OutFile(str  + "\n")
	}

	if log.checkTargetIncludeByBit(OUT_TARGET_SC){
		log.OutScreen(msg.LevelPrefix,msg.Header,msg.Content)
	}

	if log.checkTargetIncludeByBit(OUT_TARGET_NET){

	}
}
//https://github.com/gogf/gf/tree/master/os/glog
type apiString interface {
	String() string
}

type apiError interface {
	Error() string
}

func  (log *Log)String(i interface{}) string {
	if i == nil {
		return ""
	}
	switch value := i.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	//case time.Time:
	//	if value.IsZero() {
	//		return ""
	//	}
	//	return value.String()
	//case *time.Time:
	//	if value == nil {
	//		return ""
	//	}
	//	return value.String()
	//case gtime.Time:
	//	if value.IsZero() {
	//		return ""
	//	}
	//	return value.String()
	//case *gtime.Time:
	//	if value == nil {
	//		return ""
	//	}
	//	return value.String()
	default:
		// Empty checks.
		if value == nil {
			return ""
		}
		if f, ok := value.(apiString); ok {
			// If the variable implements the String() interface,
			// then use that interface to perform the conversion
			return f.String()
		}
		if f, ok := value.(apiError); ok {
			// If the variable implements the Error() interface,
			// then use that interface to perform the conversion
			return f.Error()
		}
		// Reflect checks.
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
		case reflect.String:
			return rv.String()
		}
		if kind == reflect.Ptr {
			return log.String(rv.Elem().Interface())
		}
		// Finally we use json.Marshal to convert.
		if jsonContent, err := json.Marshal(value); err != nil {
			return fmt.Sprint(value)
		} else {
			return string(jsonContent)
		}
	}
}

//func  (log *Log)parseInterfaceValueCovertStr(interValue interface{})string{
//
//	switch f := interValue.(type) {
//		case bool:
//			if f {
//				return "true"
//			}else{
//				return "false"
//			}
//		case float32:
//			return FloatToString(interValue.(float32),3)
//		case float64:
//			return Float64ToString(interValue.(float64),3)
//		//case complex64:
//		//	p.fmtComplex(complex128(f), 64, verb)
//		//case complex128:
//		//	p.fmtComplex(f, 128, verb)
//		case int:
//			return strconv.Itoa(interValue.(int))
//		case int8:
//			return strconv.Itoa(int(interValue.(int8)))
//		case int16:
//			strconv.Itoa(int (interValue.(int16)))
//		case int32:
//			strconv.FormatInt(int64 (interValue.(int32)),10)
//		case int64:
//			strconv.FormatInt(interValue.(int64),10)
//		case uint:
//			strconv.FormatUint(uint64(interValue.(uint)),10)
//		case uint8:
//			strconv.FormatUint(uint64(interValue.(uint8)),10)
//		case uint16:
//			strconv.FormatUint(uint64(interValue.(uint16)),10)
//		case uint32:
//			strconv.FormatUint(uint64(interValue.(uint32)),10)
//		case uint64:
//			strconv.FormatUint(interValue.(uint64),10)
//		case uintptr:
//			p.fmtInteger(uint64(f), unsigned, verb)
//		case string:
//			return interValue.(string)
//		case []byte:
//			return interValue.(string)
//		case reflect.Value:
//			if f.IsValid() && f.CanInterface() {
//				p.arg = f.Interface()
//				if p.handleMethods(verb) {
//					return
//				}
//			}
//			p.printValue(f, verb, 0)
//		default:
//			if !p.handleMethods(verb) {
//				p.printValue(reflect.ValueOf(f), verb, 0)
//			}
//
//	}
//}

//var levelStringMap = map[string]int{
//	"ALL":      LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"DEV":      LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"DEVELOP":  LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"PROD":     LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"PRODUCT":  LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"DEBU":     LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"DEBUG":    LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"INFO":     LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"NOTI":     LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"NOTICE":   LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"WARN":     LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"WARNING":  LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"ERRO":     LEVEL_ERRO | LEVEL_CRIT,
//	"ERROR":    LEVEL_ERRO | LEVEL_CRIT,
//	"CRIT":     LEVEL_CRIT,
//	"CRITICAL": LEVEL_CRIT,
//}