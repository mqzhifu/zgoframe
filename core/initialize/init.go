package initialize

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/model"
	"zgoframe/util"
)

type InitOption struct {
	Env 				string
	ConfigType 			string
	ConfigFileName 		string
	ConfigSourceType 	string
	EtcdConfigFindUrl	string
	RootDirName 		string
	RootCtx 			context.Context
	RootCancelFunc		context.CancelFunc
	RootQuitFunc		func(source int)
}


func Init(option InitOption)error{
	//createDbTable()

	//初始化配置信息
	viperOption := ViperOption{
		ConfigFileName: option.ConfigFileName,
		ConfigFileType:  option.ConfigType,
		SourceType: option.ConfigSourceType,
		EtcdUrl: option.EtcdConfigFindUrl,
		ENV: option.Env,
	}
	util.MyPrint("config flow : ")
	util.PrintStruct(option,":")

	myViper,config,err := GetNewViper(viperOption)
	if err != nil{
		util.MyPrint("GetNewViper err:",err)
		return err
	}
	global.V.Vip = myViper
	global.C = config
	//---config end -----

	//mysql
	//这里按说不应该先初始化MYSQL，而且不一定所有项目都用MYSQL，但是项目是基于多APP/PROJECT的模式，强依赖app_id
	//if global.C.Mysql.Status == global.CONFIG_STATUS_OPEN{
		global.V.Gorm ,err = GetNewGorm()
		if err != nil{
			util.MyPrint("GetGorm err:",err)
			return err
		}
	//}

	//初始化APP信息，所有项目都需要有AppId
	if global.C.System.AppId <=0 {
		return errors.New("appid is empty")
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

	//这里要求，项目表里配置的key与项目目录名必须一致.
	if option.RootDirName != global.V.App.Key{
		return errors.New("mainDirName != app name , "+option.RootDirName + " "+  global.V.App.Name)
	}
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
	//基础类：用于恢复一个挂了的协程
	global.V.RecoverGo = util.NewRecoverGo(global.V.Zap)
	//redis
	if global.C.Redis.Status == global.CONFIG_STATUS_OPEN{
		global.V.Redis ,err = GetNewRedis()
		if err != nil{
			util.MyPrint("GetRedis err:",err)
			return err
		}
	}
	//Http log zap
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
		global.V.Etcd ,err = GetNewEtcd()
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
	pwd,_ := os.Getwd()
	dir := pwd + "/protobuf"
	global.V.ProtobufMap = util.NewProtobufMap(global.V.Zap,dir)
	//websocket
	if global.C.Websocket.Status == global.CONFIG_STATUS_OPEN{
		if global.C.Http.Status != global.CONFIG_STATUS_OPEN{
			return errors.New("Websocket need gin open!")
		}
		initSocket()
	}
	//grpc
	if global.C.Grpc.Status == global.CONFIG_STATUS_OPEN{
		grpcOption := util.GrpcOption{
			AppId 		: global.V.App.Id,
			ListenIp	: global.C.Grpc.Ip,
			OutIp		: global.C.Grpc.Ip,
			Port 		: global.C.Grpc.Port,
			Log			: global.V.Zap,
		}
		global.V.Grpc =  util.NewMyGrpc(grpcOption)
		//grpcInc,listen,err := global.V.Grpc.GetServer()
		//if err != nil{
		//	return errors.New(err.Error())
		//}
		////挂载服务的handler
		//pb.RegisterFirstServer(grpcInc, &pbservice.First{})
		//pb.RegisterSecondServer(grpcInc, &pbservice.Second{})
		//// 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系
		//go global.V.Grpc.StartServer(grpcInc,listen)
		//
		//grpcClientConn,err := global.V.Grpc.GetClient("127.0.0.1","6666")
		//if err != nil{
		//	return errors.New(err.Error())
		//}
		//pbServiceFirst := pb.NewFirstClient(grpcClientConn)
	}
	//预/报警,这个是真正的报警，如：邮件 SMS 等
	global.V.AlertHook = util.NewAlertHook(global.V.Zap)


	global.C.System.ENV = option.Env
	//启动http
	if global.C.Http.Status == global.CONFIG_STATUS_OPEN{
		StartHttpGin()
	}
	//_ ,cancelFunc := context.WithCancel(option.RootCtx)
	//进程通信相关
	ProcessPathFileName := "/tmp/"+global.V.App.Name+".pid"
	global.V.Process = util.NewProcess(ProcessPathFileName,option.RootCancelFunc,global.V.Zap,option.RootQuitFunc)
	global.V.Process.InitProcess()

	return nil
}

func Quit(){
	global.V.Zap.Warn("init quit start:")
	HttpServerShutdown()
	global.V.Zap.Warn("init quit mmmmm:")
	global.V.Redis.Close()
	db , _ := global.V.Gorm.DB()
	db.Close()
	global.V.Vip.WatchRemoteConfig()
	global.V.Zap.Warn("init quit finish.")
}

func createDbTable(){
	mydb := util.NewDb(global.V.Gorm)
	mydb.CreateTable(&model.User{},&model.SmsLog{},&model.SmsRule{},&model.App{},&model.UserReg{} , &model.OperationRecord{})
	util.ExitPrint("init done.")
}


//初始化app管理容器
func GetNewApp()(m *util.AppManager,e error){
	appM,err := util.NewAppManager(global.V.Gorm)
	if err != nil{
		return m,err
	}

	return appM,nil
}

func GetNewEtcd()(myEtcd *util.MyEtcd,err error){
	option := util.EtcdOption{
		AppName		: global.V.App.Name,
		AppENV		: global.C.System.ENV,
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
		Prefix	: global.V.App.Name,
	}
	myService := util.NewService(serviceOption)

	return myService
}




