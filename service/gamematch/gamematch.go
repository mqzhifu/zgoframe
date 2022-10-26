package gamematch

import (
	"errors"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
	"zgoframe/service"
	"zgoframe/util"
)

type Gamematch struct {
	Option GamematchOption //初始化 配置参数

	RuleConfig *RuleConfig //配置信息

	containerSign    map[int]*QueueSign    //容器 - 报名类
	containerSuccess map[int]*QueueSuccess //容器 -
	containerMatch   map[int]*Match        //容器 -
	containerPush    map[int]*Push         //容器 -

	signalChan       map[int]map[string]chan int //一个rule会开N个协作守护协程，之间的通信信号保存在这里
	signalChanRWLock *sync.RWMutex
	HttpdRuleState   map[int]int
	PlayerStatus     *PlayerStatus //控制玩家的状态 容器

}

//玩家结构体，目前暂定就一个ID，其它的字段值，后期得加，主要是用于权重计算
type Player struct {
	Id        int
	MatchAttr map[string]int
	Weight    float32
}

//快捷 全局(包) 变量
var mylog *zap.Logger
var myerr *util.ErrMsg
var myredis *util.MyRedisGo
var myServiceDiscovery *util.ServiceDiscovery
var myetcd *util.MyEtcd
var mymetrics *util.MyMetrics
var playerStatus *PlayerStatus //玩家状态类
var projectId int

type GamematchOption struct {
	Log              *zap.Logger
	Redis            *util.MyRedisGo
	Service          *util.Service
	Etcd             *util.MyEtcd
	Metrics          *util.MyMetrics
	Err              *util.ErrMsg
	ServiceDiscovery *util.ServiceDiscovery
	ProjectId        int
	MonitorRuleIds   []int
}

func NewGameMatch(gamematchOption GamematchOption) (gamematch *Gamematch, errs error) {
	gamematchOption.Log.Info("NewGameMatch : ")
	//初始化全局变量，用于便捷访问
	mylog = gamematchOption.Log
	myredis = gamematchOption.Redis
	myServiceDiscovery = gamematchOption.ServiceDiscovery
	myetcd = gamematchOption.Etcd
	mymetrics = gamematchOption.Metrics
	projectId = gamematchOption.ProjectId
	//初始化-错误码
	container := getErrorCode()
	mylog.Info(" init ErrorCodeList , len : " + strconv.Itoa(len(container)))
	if len(container) == 0 {
		return gamematch, errors.New("getErrorCode len  = 0")
	}
	//初始化-错误/异常 类
	myerr = gamematchOption.Err
	gamematch = new(Gamematch)
	gamematch.Option = gamematchOption
	//初始化 - 信号和管道
	gamematch.signalChan = make(map[int]map[string]chan int)
	gamematch.signalChanRWLock = &sync.RWMutex{}
	//初始化所有 匹配规则
	gamematch.RuleConfig, errs = NewRuleConfig(gamematch, gamematchOption.MonitorRuleIds)
	if errs != nil {
		return gamematch, errs
	}
	//初始化- 容器 : 报名、匹配、推送、报名成功
	gamematch.containerPush = make(map[int]*Push)
	gamematch.containerSign = make(map[int]*QueueSign)
	gamematch.containerSuccess = make(map[int]*QueueSuccess)
	gamematch.containerMatch = make(map[int]*Match)
	//实例化容器
	//每一个RULE都有对应的上面的：4个节点  push sign success match
	for _, rule := range gamematch.RuleConfig.getAll() {
		//fmt.Printf("%+v",rule)
		gamematch.containerPush[rule.Id] = NewPush(rule, gamematch)
		gamematch.containerSign[rule.Id] = NewQueueSign(rule, gamematch)
		gamematch.containerSuccess[rule.Id] = NewQueueSuccess(rule, gamematch)
		gamematch.containerMatch[rule.Id] = NewMatch(rule, gamematch)
	}
	//共计开了多少个容器，间接约等于多少个协程
	containerTotal := len(gamematch.containerPush) + len(gamematch.containerSign) + len(gamematch.containerSuccess) + len(gamematch.containerMatch)
	mylog.Info("rule container total :" + strconv.Itoa(containerTotal))

	playerStatus = NewPlayerStatus()
	gamematch.PlayerStatus = playerStatus
	//初始化 - 每个rule - httpd 状态
	gamematch.HttpdRuleState = make(map[int]int)
	mylog.Info("HTTPD_RULE_STATE_INIT")
	for ruleId, _ := range gamematch.RuleConfig.getAll() {
		//设置状态
		gamematch.HttpdRuleState[ruleId] = service.HTTPD_RULE_STATE_INIT
	}

	return gamematch, nil
}
func (gamematch *Gamematch) GetContainerSignByRuleId(ruleId int) *QueueSign {
	content, ok := gamematch.containerSign[ruleId]
	if !ok {
		mylog.Error("getContainerSignByRuleId is null")
	}
	return content
}
func (gamematch *Gamematch) getContainerSuccessByRuleId(ruleId int) *QueueSuccess {
	content, ok := gamematch.containerSuccess[ruleId]
	if !ok {
		mylog.Error("getContainerSuccessByRuleId is null")
	}
	return content
}
func (gamematch *Gamematch) getContainerPushByRuleId(ruleId int) *Push {
	content, ok := gamematch.containerPush[ruleId]
	if !ok {
		mylog.Error("getContainerPushByRuleId is null")
	}
	return content
}
func (gamematch *Gamematch) getContainerMatchByRuleId(ruleId int) *Match {
	content, ok := gamematch.containerMatch[ruleId]
	if !ok {
		mylog.Error("getContainerMatchByRuleId is null")
	}
	return content
}

//给一个 协程 管道  发信号
func (gamematch *Gamematch) notifyRoutine(sign chan int, signType int) {
	mylog.Warn("send routine : " + strconv.Itoa(signType))
	sign <- signType
}
func (gamematch *Gamematch) Quit(source int) {
	gamematch.closeDemonRoutine()
	//os.Exit(state)
}

//启动后台守护-协程
func (gamematch *Gamematch) Startup() {
	queueList := gamematch.RuleConfig.getAll()
	queueLen := len(queueList)
	mylog.Info("start Startup ,  rule total : " + strconv.Itoa(queueLen))
	if queueLen <= 0 {
		mylog.Error(" RuleConfig list is empty!!!")
		util.ExitPrint(" RuleConfig list is empty!!!")
		return
	}
	//开始每个rule
	for _, rule := range queueList {
		if rule.Id == 0 { //0是特殊管理，仅给HTTPD使用
			continue
		}
		gamematch.startOneRuleDomon(rule)
	}
	//后台守护协程均已开启完毕，可以开启前端HTTPD入口了
	//gamematch.StartHttpd(gamematch.Option.HttpdOption)
}

//睡眠 - 协程
func mySleepSecond(second time.Duration, msg string) {
	mylog.Info(msg + " sleep second " + strconv.Itoa(int(second)))
	time.Sleep(second * time.Second)
}

func (gamematch *Gamematch) DelOneRuleById(ruleId int) {
	gamematch.RuleConfig.delOne(ruleId)
}

//删除全部数据
func (gamematch *Gamematch) DelAll() {
	mylog.Warn(" action :  DelAll")
	keys := service.RedisPrefix + "*"
	myredis.RedisDelAllByPrefix(keys)
}

//开启一条rule的所有守护协程，
//虽然有4个，但是只有match是最核心、最复杂的，另外3个算是辅助
func (gamematch *Gamematch) startOneRuleDomon(rule Rule) {
	//zlib.AddRoutineList("startOneRuleDomon push")
	//zlib.AddRoutineList("startOneRuleDomon signTimeout")
	//zlib.AddRoutineList("startOneRuleDomon successTimeout")
	//zlib.AddRoutineList("startOneRuleDomon matching")

	//报名超时,这里注释掉了，因为：是并发协程开启，如果其它协程先执行检测到有要处理的数据，这个时候此协程检查出有超时的数据开始删除操作，结果其它协程执行到一半发现数据丢了
	//queueSign := gamematch.GetContainerSignByRuleId(rule.Id)
	//go gamematch.StartOneGoroutineDemon(rule.Id,"signTimeout",queueSign.Log,queueSign.CheckTimeout)
	//gamematch.Option.Goroutine.CreateExec(gamematch,"MyDemon",rule.Id,"signTimeout",queueSign.Log,queueSign.CheckTimeout)
	//报名成功，但3方迟迟推送失败，无人接收，超时
	queueSuccess := gamematch.getContainerSuccessByRuleId(rule.Id)
	go gamematch.StartOneGoroutineDemon(rule.Id, "successTimeout", queueSuccess.Log, queueSuccess.CheckTimeout)
	//gamematch.Option.Goroutine.CreateExec(gamematch,"MyDemon",rule.Id,"successTimeout",queueSuccess.Log,queueSuccess.CheckTimeout)
	//推送
	push := gamematch.getContainerPushByRuleId(rule.Id)
	go gamematch.StartOneGoroutineDemon(rule.Id, "push", push.Log, push.checkStatus)
	//gamematch.Option.Goroutine.CreateExec(gamematch,"MyDemon",rule.Id,"push",push.Log,push.checkStatus)
	//匹配
	match := gamematch.containerMatch[rule.Id]
	go gamematch.StartOneGoroutineDemon(rule.Id, "matching", match.Log, match.matching)
	//gamematch.Option.Goroutine.CreateExec(gamematch,"MyDemon",rule.Id,"matching",match.Log,match.matching)

	mylog.Info("start httpd ,up state...")
	gamematch.HttpdRuleState[rule.Id] = service.HTTPD_RULE_STATE_OK
}

////请便是方便记日志，每次要写两个FD的日志，太麻烦
//func rootAndSingToLogInfoMsg(sign *QueueSign,a ...interface{}){
//	mylog.Info(a)
//	sign.Log.Info(a)
//}
//func (gamematch *Gamematch) StartHttpd(httpdOption HttpdOption)error{
//	httpd,err  := NewHttpd(httpdOption,gamematch)
//	if err != nil{
//		return err
//	}
//	httpd.Start()
//	return nil
//}
////让出当前协程执行时间
//func myGosched(msg string){
//	mylog.Info(msg + " Gosched ..." )
//	runtime.Gosched()
//}
////死循环
//func deadLoopBlock(sleepSecond time.Duration,msg string){
//	for {
//		mySleepSecond(sleepSecond,  " deadLoopBlock: " +msg)
//	}
//}
////实例化，一个LOG类，基于模块
//func getModuleLogInc( moduleName string)(newLog *zap.Logger ,err error){
//	logOption := myGamematchOption.Option.Log.Option
//	logOption.OutFileFileName = moduleName
//	//logOption := zlib.LogOption{
//	//	OutFilePath : mylog.Op.OutFilePath ,
//	//	OutFileName: moduleName + ".log",
//	//	//Level : zlib.Atoi(myetcd.GetAppConfByKey("log_level")),
//	//	Level:  mylog.Op.Level,
//	//	Target : 6,
//	//}
//	newLog,err = zlib.NewLog(logOption)
//	if err != nil{
//		return newLog,err
//	}
//	return newLog,nil
//}
////实例化，一个LOG类，基于RULE+模块
//func getRuleModuleLogInc(ruleCategory string,moduleName string)*zlib.Log{
//	//dir := myetcd.GetAppConfByKey("log_base_path") + "/" + ruleCategory
//
//	logOption := myGamematchOption.Option.Log.Option
//	logOption.OutFileFileName = moduleName
//
//	dir := logOption.OutFilePath + "/" + ruleCategory + "/"
//	logOption.OutFilePath = dir
//
//	_ ,err := util.PathExists(dir)
//	if err != nil {//证明目录存在
//		//mylog.Debug("dir has exist",dir)
//	}else{
//		err := os.Mkdir(dir, 0777)
//		if err != nil{
//			util.ExitPrint("create dir failed ",err.Error())
//		}else{
//			mylog.Debug("create dir success : ",dir)
//		}
//	}
//
//	//logOption := zlib.LogOption{
//	//	OutFilePath : dir ,
//	//	OutFileName: moduleName + ".log",
//	//	Level : mylog.Op.Level,
//	//	Target : 6,
//	//}
//	newLog,err := zlib.NewLog(logOption)
//	if err != nil{
//		util.ExitPrint(err.Error())
//	}
//	return newLog
//}

//通用 业务型  请求 数据  检查
func (gamematch *Gamematch) BusinessCheckData(form HttpReqBusiness) (errCode int, httpReqBusiness HttpReqBusiness) {
	//mylog.Info(" businessCheckData : ")
	//if postJsonStr == ""{
	//	return 802,httpReqBusiness
	//}
	//var jsonUnmarshalErr error
	//jsonUnmarshalErr = json.Unmarshal([]byte(postJsonStr),&httpReqBusiness)
	//if jsonUnmarshalErr != nil{
	//	mylog.Error(jsonUnmarshalErr.Error())
	//	return 459,httpReqBusiness
	//}
	if form.MatchCode == "" {
		return 450, httpReqBusiness
	}
	rule, err := gamematch.RuleConfig.getByCategory(httpReqBusiness.MatchCode)
	if err != nil {
		return 806, httpReqBusiness
	}
	httpReqBusiness.RuleId = rule.Id
	_, err = gamematch.checkHttpdState(httpReqBusiness.RuleId)
	if err != nil {
		return 804, httpReqBusiness
	}

	return 0, httpReqBusiness
}

func (gamematch *Gamematch) checkHttpdState(ruleId int) (bool, error) {
	//state ,ok := gamematch.HttpdRuleState[ruleId]
	//if !ok {
	//	return false,myerr.NewErrorCode(803)
	//}
	//if state == HTTPD_RULE_STATE_OK{
	//	return true,nil
	//}
	//return false,myerr.NewErrorCode(804)
	return true, nil
}
