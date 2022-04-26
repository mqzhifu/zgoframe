package gamematch
//
//import (
//	"context"
//	"gamematch/gamematch"
//	_ "net/http/pprof"
//	"os"
//	"os/signal"
//	"strconv"
//	"syscall"
//	"time"
//	"zlib"
//)
//
////var AddRuleFlag = 0
//var SERVICE_MATCH_NAME = ""
//
//var mainOutPrefix = "main "
//var ctx = context.Background()
//var cancelCtx ,cancelFunc = context.WithCancel(ctx)
////===========
//var mylog 		*zlib.Log
//var myMetrics 	*zlib.Metrics
//var myEtcd 		*zlib.MyEtcd
//var myRedis 	*zlib.MyRedis
//var myService 	*zlib.Service
//var myGamematch *gamematch.Gamematch


//func main(){
	//zlib.LogLevelFlag = zlib.LOG_LEVEL_DEBUG
	////处理指令行参数
	//cmdArgsStruct := gamematch.CmdArgs{}
	//cmsArg ,err := zlib.CmsArgs(cmdArgsStruct)
	//if err != nil{
	//	zlib.ExitPrint(mainOutPrefix + " err " +err.Error())
	//}
	//zlib.MyPrint(mainOutPrefix + " argc  ")
	//for k,v := range cmsArg{
	//	msg :=  k + ":"+ v
	//	zlib.MyPrint(msg)
	//}
	//
	//if !zlib.CheckEnvExist(cmsArg["Env"]){
	//	list := zlib.GetEnvList()
	//	zlib.ExitPrint(mainOutPrefix + "env is err , list:",list)
	//}
	////指令行参数处理完成，进入初始化阶段
	//gm := enter(cmsArg)
	//这里有个顺序关系，要后台都正常启动完成后，再开启httpd入口
	//go gm.Startup()
	//myGoroutine.CreateExec(myGamematch,"DemonAll")
	//myGoroutine.CreateExec(myGamematch,"StartHttpd",myHttpdOption)
	//go DemonSignal()



	//go DemonSignal()
	//
	//now := zlib.GetNowTimeSecondToInt()
	//myMetrics.FastLog("InitEndTime",zlib.METRICS_OPT_PLUS,now)
	//mylog.Alert("InitTime:" ,myMetrics.GetInitTime()  )
	//<- cancelCtx.Done()
	//zlib.MyPrint("ExecTime:",myMetrics.GetExecTime())
	//
	//
	//time.Sleep(time.Second * 1)
	//mylog.CloseChan <- 1
//}
//每个rule的协程包括：HTTPD、报名超时、匹配成功超时、PUSH推送
//func enter(cmsArg map[string]string)*gamematch.Gamematch{
	//zlib.MyPrint("enter initialize")
	////获取app project 项目信息
	//appM  := zlib.NewAppManager()
	//app,empty := appM.GetById(gamematch.APP_ID)
	//if empty{
	//	zlib.PanicPrint(mainOutPrefix + " err " + ": appId is empty "+ strconv.Itoa(app.Id))
	//}
	////该项目的英文名称
	//SERVICE_MATCH_NAME = app.Key
	////创建全局日志类
	//logLevel,_ := strconv.Atoi(cmsArg["LogLevel"])
	//logOutFilePath :=  zlib.BasePathPlusTypeStr(cmsArg["LogBasePath"],appM.GetTypeName(app.Type))
	//logOption := zlib.LogOption{
	//	AppId			: app.Id,
	//	ModuleId		: 1,
	//	OutFilePath 	: logOutFilePath,
	//	OutFileFileName	: SERVICE_MATCH_NAME,
	//	Level 			: logLevel,
	//	OutTarget 		: zlib.OUT_TARGET_ALL,
	//	//OutContentType	: zlib.CONTENT_TYPE_JSON,
	//	OutContentType	: zlib.CONTENT_TYPE_STRING,
	//	OutFileHashType	: zlib.FILE_HASH_DAY,
	//	OutFileFileExtName : "log",
	//}
	//
	//newlog,errs  := zlib.NewLog(logOption)
	//if errs != nil{
	//	zlib.ExitPrint("new log err",errs.Error())
	//}
	//mylog = newlog
	////全局标量统计类
	//MetricsOption :=zlib.MetricsOption{
	//	Log: mylog,
	//}
	//myMetrics  = zlib.NewMetrics(MetricsOption)
	////启动 全局标量统计类 接收统计信息
	//go myMetrics.Start(ctx)
	////记录一下程序的启动时间
	//myMetrics.FastLog("StartUpTime",zlib.METRICS_OPT_PLUS,zlib.GetNowTimeSecondToInt())
	////短连接地址，主要获取配置信息，这里主要是获取ETCD的配置信息
	//url := cmsArg["BaseUrl"]
	//etcdOption := zlib.EtcdOption{
	//	FindEtcdUrl: url,
	//	Log : mylog,
	//	AppName: SERVICE_MATCH_NAME,
	//	AppENV : cmsArg["Env"],
	//}
	//myEtcd,errs = zlib.NewMyEtcdSdk(etcdOption)
	//if errs != nil{
	//	zlib.ExitPrint("NewMyEtcdSdk err",errs.Error())
	//}
	////实例化-<redis>-组件
	//redisOption := zlib.RedisOption{
	//	Host: myEtcd.GetAppConfByKey("redis_host"),
	//	Port: myEtcd.GetAppConfByKey("redis_port"),
	//	Ps: myEtcd.GetAppConfByKey("redis_ps"),
	//	Log: mylog,
	//}
	//myRedis , errs = zlib.NewRedisConnPool(redisOption)
	//if errs != nil{
	//	zlib.ExitPrint("new redis err",errs.Error())
	//}
	////实例化-<服务发现>-组件
	//serviceOption := zlib.ServiceOption{
	//	Etcd: myEtcd,
	//	Log: mylog,
	//	Prefix: gamematch.SERVICE_PREFIX,
	//	TestHttpGamematchPushReceiveHsot: gamematch.TEST_HTTP_PUSH_RECEIVE_HOST,
	//	//Goroutine: myGoroutine,
	//}
	////获取当前机器IP，这里获取的好像是个肉网IP
	//localIp,err := zlib.GetLocalIp()
	//if err !=nil{
	//	zlib.ExitPrint("GetLocalIp err : ",err)
	//}
	////myHost := "192.168.31.148"
	////myHost := "192.168.192.170"
	//myHost := localIp
	//myPort := myEtcd.GetAppConfByKey("reg_self_port")
	////创建 服务发现 类
	//myService = zlib.NewService(serviceOption)
	////将自己注册成一个服务
	//err = myService.RegOneDynamic(SERVICE_MATCH_NAME,myHost+":"+myPort)
	//if err !=nil{
	//	mylog.Error("myService.RegOneDynamic err :",err.Error())
	//	zlib.ExitPrint("err")
	//}
	//myHttpdOption :=  gamematch.HttpdOption{
	//	Host: myHost,
	//	Port: myPort,
	//	Log : mylog,
	//}
	//最后，终于，实例化：匹配机制
	//gamematchOption := gamematch.GamematchOption  {
	//	Log :mylog,
	//	Redis :myRedis,
	//	Service :myService,
	//	Etcd : myEtcd,
	//	Metrics: myMetrics,
	//	PidFilePath :myEtcd.GetAppConfByKey("pid_file_path"),
	//	CmsArg: cmsArg,
	//	HttpdOption: myHttpdOption,
	//	//Goroutine: myGoroutine,
	//}
	//
	//myGamematch,errs = gamematch.NewGamematch(gamematchOption)
	//if errs != nil{
	//	zlib.ExitPrint("NewGamematch : ",errs.Error())
	//}

	//return myGamematch
//}

//func Quit(){
//	myMetrics.FastLog("ShutdownStartTime",zlib.METRICS_OPT_PLUS,zlib.GetNowTimeSecondToInt())
//
//	myGamematch.Quit(1)
//	myService.Shutdown()
//	myEtcd.Shutdown()
//	myRedis.Shutdown()
//
//	cancelFunc()
//}
////信号 处理
//func DemonSignal(){
//	mylog.Warning("SIGNAL init : ")
//	c := make(chan os.Signal)
//	//syscall.SIGHUP :ssh 挂断会造成这个信号被捕获，先注释掉吧
//	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
//	prefix := "SIGNAL-DEMON :"
//	for{
//		sign := <- c
//		mylog.Warning(prefix,sign)
//		switch sign {
//			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
//				mylog.Warning(prefix+" exit!!!")
//				Quit()
//				goto end
//			case syscall.SIGUSR1:
//				mylog.Warning(prefix+" usr1!!!")
//			case syscall.SIGUSR2:
//				mylog.Warning(prefix+" usr2!!!")
//			default:
//				mylog.Warning(prefix+" unknow!!!")
//		}
//		time.Sleep(time.Second * 1)
//	}
//	end:
//		mylog.Alert("DemonSignal DONE.")
//}







