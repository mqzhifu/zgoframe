package initialize

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"strings"
	"zgoframe/core/global"
	"zgoframe/model"
	"zgoframe/protobuf/pb"
	"zgoframe/protobuf/pbservice"
	"zgoframe/util"
)

type Initialize struct {
	Option InitOption
}

type InitOption struct {
	Env 				string
	ConfigType 			string
	ConfigFileName 		string
	ConfigSourceType 	string
	EtcdConfigFindUrl	string
	RootDir				string
	RootDirName 		string
	RootCtx 			context.Context
	RootCancelFunc		context.CancelFunc
	RootQuitFunc		func(source int)
}

func NewInitialize(option InitOption)*Initialize{
	initialize := new(Initialize)
	initialize.Option = option
	return initialize
}

//初始化-入口
func (initialize * Initialize)Start()error{
	//createDbTable()
	//初始化 : 配置信息
	viperOption := ViperOption{
		ConfigFileName	: initialize.Option.ConfigFileName,
		ConfigFileType	: initialize.Option.ConfigType,
		SourceType		: initialize.Option.ConfigSourceType,
		EtcdUrl			: initialize.Option.EtcdConfigFindUrl,
		ENV				: initialize.Option.Env,
	}

	util.MyPrint("config option~~ ")
	util.PrintStruct(initialize.Option,":")

	myViper,config,err := GetNewViper(viperOption)
	if err != nil{
		util.MyPrint("GetNewViper err:",err)
		return err
	}
	global.V.Vip = myViper	//全局变量管理者
	global.C = config		//全局变量容器
	//---config end -----

	//mysql
	//这里按说不应该先初始化MYSQL，应该最早初始化LOG类，并且不一定所有项目都用MYSQL，但是项目是基于多APP/PROJECT的模式，强依赖app_id
	if global.C.Mysql.Status != global.CONFIG_STATUS_OPEN{
		util.MyPrint("not open mysql db, need read app_id from db.")
		return err
	}
	global.V.Gorm ,err = GetNewGorm()
	if err != nil{
		util.MyPrint("GetGorm err:",err)
		return err
	}
	model.Db = global.V.Gorm
	//初始化APP信息，所有项目都需要有AppId，因为要做验证，同时目录名也包含在里面
	err = InitApp()
	if err !=nil{
		return err
	}
	//项目目录名，必须跟APP-INFO里的key相同
	initialize.Option.RootDirName,err = InitPath(initialize.Option.RootDir)
	if err !=nil{
		return err
	}
	global.V.RootDir = initialize.Option.RootDir
	//预/报警->推送器，这里是推送到3方，如：prometheus
	//这个是必须优先zap日志类优化处理，因为zap里有用
	if global.C.Alert.Status == global.CONFIG_STATUS_OPEN{
		global.V.AlertPush = util.NewAlertPush(global.C.Alert.Ip,global.C.Alert.Port,global.C.Alert.Uri)
	}
	//日志
	global.V.Zap , err  = GetNewZapLog(global.V.AlertPush,"main","main",1)
	if err != nil{
		util.MyPrint("GetNewZapLog err:",err)
		return err
	}
	//错误 文案 管理
	global.V.Err ,err  = util.NewErrMsg(global.V.Zap,  global.C.Http.StaticPath + global.C.System.ErrorMsgFile )
	if err != nil{
		return err
	}
	//基础类：用于恢复一个挂了的协程
	global.V.RecoverGo = util.NewRecoverGo(global.V.Zap)
	//redis
	if global.C.Redis.Status == global.CONFIG_STATUS_OPEN{
		global.V.Redis ,err = GetNewRedis()
		if err != nil{
			util.MyPrint("GetRedis err:",err)
			return err
		}
		//TestRedis()
	}
	//Http log zap 这里单独再开个zap 实例，用于专门记录http 请求
	HttpZap , err  := GetNewZapLog(global.V.AlertPush,"http","http",0)
	if err != nil{
		util.MyPrint("GetNewZapLog err:",err)
		return err
	}
	//http server
	if global.C.Http.Status == global.CONFIG_STATUS_OPEN{
		global.V.Gin ,err = GetNewHttpGIN(HttpZap)
		if err != nil{
			util.MyPrint("GetNewHttpGIN err:",err)
			return err
		}
	}
	//etcd
	if global.C.Etcd.Status  == global.CONFIG_STATUS_OPEN{
		global.V.Etcd ,err = GetNewEtcd(initialize.Option.Env)
		if err != nil{
			util.MyPrint("GetNewEtcd err:",err)
			return err
		}
	}
	//service 服务发现
	if global.C.Service.Status  == global.CONFIG_STATUS_OPEN{
		if global.C.Etcd.Status != global.CONFIG_STATUS_OPEN{
			return errors.New("Service need Etcd open!")
		}
		global.V.Service = GetNewService()
	}
	//metrics
	if global.C.Metrics.Status == global.CONFIG_STATUS_OPEN{
		global.V.Metric =  util.NewMyMetrics(global.V.Zap)

		if global.C.Http.Status != global.CONFIG_STATUS_OPEN{
			return errors.New("metrics need gin open!")
		}
		global.V.Gin.GET("/metrics", gin.WrapH(promhttp.Handler()))
		//测试
		//global.V.Gin.GET("/metrics/count", func(c *gin.Context) {
		//	global.V.Metric.CounterInc("paySuccess")
		//})
		//
		//global.V.Gin.GET("/metrics/gauge", func(c *gin.Context) {
		//	global.V.Metric.CounterInc("payUser")
		//})
		//global.V.Metric.Test()
	}
	//初始化-protobuf 映射文件
	dir := initialize.Option.RootDir + "/protobuf"
	//将rpc service 中的方法，转化成ID（由PHP生成 的ID map）
	global.V.ProtobufMap = util.NewProtobufMap(global.V.Zap,dir)
	//websocket
	//if global.C.Websocket.Status == global.CONFIG_STATUS_OPEN{
	//	if global.C.Http.Status != global.CONFIG_STATUS_OPEN{
	//		return errors.New("Websocket need gin open!")
	//	}
	//	initSocket()
	//}
	//grpc
	if global.C.Grpc.Status == global.CONFIG_STATUS_OPEN{
		grpcOption := util.GrpcOption{
			AppId 		: global.V.App.Id,
			ListenIp	: global.C.Grpc.Ip,
			OutIp		: global.C.Grpc.Ip,
			Port 		: global.C.Grpc.Port,
			Log			: global.V.Zap,
		}
		global.V.Grpc,_ =  util.NewMyGrpc(grpcOption)
		//

	}

	if global.C.Email.Status == global.CONFIG_STATUS_OPEN {
		emailOption := util.EmailOption{
			Host: global.C.Email.Host,
			Port: global.C.Email.Port,
			FromEmail: global.C.Email.From,
			Password: global.C.Email.Ps,
			Log: global.V.Zap,
		}

		global.V.Email = util.NewMyEmail(emailOption)
	}
	//预/报警,这个是真正的报警，如：邮件 SMS 等
	if global.C.Alert.Status == global.CONFIG_STATUS_OPEN {
		global.V.AlertHook = util.NewAlertHook(-1,"程序出错了：#body#","报错",util.ALERT_METHOD_SYNC,global.V.Zap)
		global.V.AlertHook.Email = global.V.Email
		//global.V.AlertHook.Alert("Aaaa")
		//util.ExitPrint(123123123)
	}
	//TestGorm()
	global.C.System.ENV = initialize.Option.Env
	//启动http
	if global.C.Http.Status == global.CONFIG_STATUS_OPEN{
		RegGinHttpRoute()//这里注册项目自己的http 路由策略
		StartHttpGin()
	}

	//_ ,cancelFunc := context.WithCancel(option.RootCtx)
	//进程通信相关
	ProcessPathFileName := "/tmp/"+global.V.App.Name+".pid"
	global.V.Process = util.NewProcess(ProcessPathFileName,initialize.Option.RootCancelFunc,global.V.Zap,initialize.Option.RootQuitFunc)
	global.V.Process.InitProcess()

	return nil
}



func createDbTable(){
	mydb := util.NewDbTool(global.V.Gorm)
	mydb.CreateTable(&model.User{},&model.SmsLog{},&model.SmsRule{},&model.App{},&model.UserReg{} , &model.OperationRecord{})
	util.ExitPrint("init done.")
}

func (initialize *Initialize)StartService()error{
	grpcInc,listen,err := global.V.Grpc.GetServer()
	if err != nil{
		return errors.New(err.Error())
	}

	//挂载服务的handler
	pb.RegisterFirstServer(grpcInc, &pbservice.First{})
	pb.RegisterSecondServer(grpcInc, &pbservice.Second{})
	// 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系
	go global.V.Grpc.StartServer(grpcInc,listen)

	return nil
}

func (initialize *Initialize)StartClient()error{
	grpcClientConn,err := global.V.Grpc.GetClient(global.C.Grpc.Ip,global.C.Grpc.Port)
	if err != nil{
		util.MyPrint(err)
		return errors.New(err.Error())
	}
	pbServiceFirst := pb.NewFirstClient(grpcClientConn)
	RequestRegPlayer := pb.RequestRegPlayer{}
	RequestRegPlayer.AddTime = 123123
	res ,_:= pbServiceFirst.SayHello(context.Background(),&RequestRegPlayer)
	util.MyPrint("grpc return:",res)
	return nil
}



func (initialize * Initialize)Quit(){
	global.V.Zap.Warn("init quit start:")
	HttpServerShutdown()
	RedisShutdown()
	GormShutdown()
	global.V.Websocket.Shutdown()
	ViperShutdown()
	global.V.Grpc.Shutdown()
	global.V.Etcd.Shutdown()
	global.V.Service.Shutdown()

	global.V.Zap.Warn("init quit finish.")
}


func InitApp()(err error){
	if global.C.System.AppId <=0 {
		return errors.New("appId is empty")
	}

	global.V.AppMng ,err  = GetNewApp()
	if err != nil{
		util.MyPrint("GetNewApp err:",err)
		return err
	}
	//根据APPId去DB中查找详细信息
	app,empty := global.V.AppMng.GetById(global.C.System.AppId)
	if empty {
		return errors.New("AppId not match : " + strconv.Itoa(global.C.System.AppId) )
	}
	global.V.App = app
	util.MyPrint("project app info flow:")
	util.PrintStruct(app,":")

	return nil
}

func InitPath(rootDir string)(rootDirName string,err error){
	pwdArr:=strings.Split(rootDir,"/")//切割路径字符串
	rootDirName = pwdArr[len(pwdArr)-1]//获取路径数组最后一个元素：当前路径的文件夹名
	//option.RootDirName = rootDirName
	//global.V.RootDir = option.RootDir
	//这里要求，项目表里配置的key与项目目录名必须一致.
	if rootDirName != global.V.App.Key{
		return rootDirName,errors.New("mainDirName != app name , "+rootDirName + " "+  global.V.App.Name)
	}
	return rootDirName,nil
}

//初始化app管理容器
func GetNewApp()(m *util.AppManager,e error){
	appM,err := util.NewAppManager(global.V.Gorm)
	if err != nil{
		return m,err
	}

	return appM,nil
}

func GetNewEtcd(env string)(myEtcd *util.MyEtcd,err error){
	option := util.EtcdOption{
		AppName		: global.V.App.Name,
		AppENV		: env,
		AppKey: global.V.App.Key,
		FindEtcdUrl : global.C.Etcd.Url,
		Username	: global.C.Etcd.Username,
		Password	: global.C.Etcd.Password,
		Ip			: global.C.Etcd.Ip,
		Port		: global.C.Etcd.Port,
		Log			: global.V.Zap,
	}
	myEtcd,err  = util.NewMyEtcdSdk(option)
	return myEtcd,err
}

func GetNewService()*util.Service {
	serviceOption := util.ServiceOption{
		Log		: global.V.Zap,
		Etcd	: global.V.Etcd,
		Prefix	: "/service",
	}
	myService := util.NewService(serviceOption)

	return myService
}




