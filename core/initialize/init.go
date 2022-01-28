package initialize

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strings"
	"zgoframe/core/global"
	"zgoframe/model"
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
	global.C = config		//全局变量
	//---config end -----

	//初始化：mysql
	//这里按说不应该先初始化MYSQL，应该最早初始化LOG类，并且不一定所有项目都用MYSQL，但是项目是基于多APP/SERVICE，强依赖 project_id
	if global.C.Mysql.Status != global.CONFIG_STATUS_OPEN{
		errMsg := "not open mysql db, need read app_id from db."
		util.MyPrint(errMsg)
		return errors.New(errMsg)
	}
	global.V.Gorm ,err = GetNewGorm()
	if err != nil{
		util.MyPrint("GetGorm err:",err)
		return err
	}
	model.Db = global.V.Gorm
	//初始化APP信息，所有项目都需要有AppId或serviceId，因为要做验证，同时目录名也包含在里面
	err = InitProject()
	if err !=nil {
		return err
	}
	//项目目录名，必须跟APP-INFO里的key相同
	initialize.Option.RootDirName,err = InitPath(initialize.Option.RootDir)
	if err !=nil{
		return err
	}
	global.V.RootDir = initialize.Option.RootDir
	util.MyPrint("global.V.RootDir:",global.V.RootDir)
	//预/报警->推送器，这里是推送到3方，如：prometheus,ps:这个是必须优先zap日志类优化处理，因为zap里的<钩子>有用到
	if global.C.Alert.Status == global.CONFIG_STATUS_OPEN{
		global.V.AlertPush = util.NewAlertPush(global.C.Alert.Ip,global.C.Alert.Port,global.C.Alert.Uri)
	}
	//日志
	configZap := global.C.Zap
	configZap.FileName = "main"
	configZap.ModuleName = "main"
	global.V.Zap , err  = GetNewZapLog(global.V.AlertPush,configZap,global.V.Project.Id)
	if err != nil{
		util.MyPrint("GetNewZapLog err:",err)
		return err
	}
	//错误码 文案 管理
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
	}
	configZap = global.C.Zap
	configZap.FileName = "http"
	configZap.ModuleName = "http"
	//Http log zap 这里单独再开个zap 实例，用于专门记录http 请求
	HttpZap , err  := GetNewZapLog(global.V.AlertPush,configZap,global.V.Project.Id)
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
	global.V.ServiceManager,_ = util.NewServiceManager(global.V.Gorm)
	//service 服务发现，这里有个顺序，必须先实现化完成:serviceManager
	if global.C.ServiceDiscovery.Status  == global.CONFIG_STATUS_OPEN{
		if global.C.Etcd.Status != global.CONFIG_STATUS_OPEN{
			return errors.New("ServiceDiscovery need Etcd open!")
		}
		global.V.ServiceDiscovery ,err = GetNewServiceDiscovery()
		if err != nil{
			return err
		}
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
	dir := initialize.Option.RootDir + "/" + global.C.Protobuf.BasePath + "/" + global.C.Protobuf.PbServicePath
	//将rpc service 中的方法，转化成ID（由PHP生成 的ID map）
	global.V.ProtobufMap ,err = util.NewProtobufMap(global.V.Zap,dir,global.C.Protobuf.IdMapFileName,global.V.ProjectMng)
	if err != nil{
		util.MyPrint("GetNewViper err:",err)
		return err
	}
	//websocket
	//if global.C.Websocket.Status == global.CONFIG_STATUS_OPEN{
	//	if global.C.Http.Status != global.CONFIG_STATUS_OPEN{
	//		return errors.New("Websocket need gin open!")
	//	}
	//	initSocket()
	//}
	//grpc
	if global.C.Grpc.Status == global.CONFIG_STATUS_OPEN{
		grpcManagerOption := util.GrpcManagerOption{
			//AppId: global.V.App.Id,
			//ServiceId: global.V.Service.Id,
			ProjectId: global.V.Project.Id,
			Log: global.V.Zap,

		}
		if global.C.ServiceDiscovery.Status == global.CONFIG_STATUS_OPEN{
			grpcManagerOption.ServiceDiscovery = global.V.ServiceDiscovery
		}
		global.V.GrpcManager,_ =  util.NewGrpcManager(grpcManagerOption)
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
	global.C.System.ENV = initialize.Option.Env
	//启动http
	if global.C.Http.Status == global.CONFIG_STATUS_OPEN{
		RegGinHttpRoute()//这里注册项目自己的http 路由策略
		StartHttpGin()
	}

	autoCreateUpDbTable()
	//_ ,cancelFunc := context.WithCancel(option.RootCtx)
	//进程通信相关
	ProcessPathFileName := "/tmp/"+global.V.Project.Name+".pid"
	global.V.Process = util.NewProcess(ProcessPathFileName,initialize.Option.RootCancelFunc,global.V.Zap,initialize.Option.RootQuitFunc)
	global.V.Process.InitProcess()

	return nil
}

func autoCreateUpDbTable(){
	//mydb := util.NewDbTool(global.V.Gorm)
	//mydb.CreateTable(&model.User{},&model.SmsLog{},&model.SmsRule{},&model.App{},&model.UserReg{} , &model.OperationRecord{})
	//util.ExitPrint("init done.")
}



func (initialize * Initialize)Quit(){
	global.V.Zap.Warn("init quit start:")
	if global.C.Http.Status == global.CONFIG_STATUS_OPEN{
		HttpServerShutdown()
	}

	if global.C.Redis.Status == global.CONFIG_STATUS_OPEN{
		RedisShutdown()
	}
	//这个得优于etcd先关
	if global.C.Grpc.Status == global.CONFIG_STATUS_OPEN{
		global.V.GrpcManager.Shutdown()
	}
	//这个得优于etcd先关
	if global.C.ServiceDiscovery.Status == global.CONFIG_STATUS_OPEN{
		global.V.ServiceDiscovery.Shutdown()
	}

	if global.C.Etcd.Status == global.CONFIG_STATUS_OPEN{
		global.V.Etcd.Shutdown()
	}
	//global.V.Websocket.Shutdown()

	GormShutdown()
	ViperShutdown()

	global.V.Zap.Warn("init quit finish.")
}



//=======================================================================================
func InitPath(rootDir string)(rootDirName string,err error){
	pwdArr:=strings.Split(rootDir,"/")//切割路径字符串
	rootDirName = pwdArr[len(pwdArr)-1]//获取路径数组最后一个元素：当前路径的文件夹名
	//option.RootDirName = rootDirName
	//global.V.RootDir = option.RootDir
	//这里要求，DB中项目记录里：name 与项目目录名必须一致，防止有人错用/盗用projectId
	projectNameByte := util.CamelToSnake2([]byte(global.V.Project.Name ))
	projectName := util.StrFirstToLower(string(projectNameByte))
	if rootDirName != projectName {
		//方便测试，先注释掉
		return rootDirName,errors.New("mainDirName != app name , rootDirName : "+rootDirName + " , ProjectName:"+ projectName )
	}

	//if global.C.System.ProjectId > 0 {
	//	if rootDirName != global.V.App.Key {
	//		return rootDirName,errors.New("mainDirName != app name , "+rootDirName + " "+  global.V.App.Key)
	//	}
	//}else{
	//	if rootDirName != global.V.Service.Key {
	//		return rootDirName,errors.New("mainDirName != serviceName name , "+rootDirName + " "+  global.V.Service.Key)
	//	}
	//}

	return rootDirName,nil
}

func GetNewEtcd(env string)(myEtcd *util.MyEtcd,err error){
	option := util.EtcdOption{
		ProjectName		: global.V.Project.Name,
		ProjectENV		: env,
		//ProjectKey		: global.V.Project.Key,
		FindEtcdUrl : global.C.Etcd.Url,
		Username	: global.C.Etcd.Username,
		Password	: global.C.Etcd.Password,
		Ip			: global.C.Etcd.Ip,
		Port		: global.C.Etcd.Port,
		Log			: global.V.Zap,
	}
	myEtcd,err  = util.NewMyEtcdSdk(option)
	//util.ExitPrint(err)
	return myEtcd,err
}

func GetNewServiceDiscovery()(serviceDiscovery *util.ServiceDiscovery,err error) {
	serviceOption := util.ServiceDiscoveryOption{
		Log		: global.V.Zap,
		Etcd	: global.V.Etcd,
		//Prefix	: "/service",
		Prefix: global.C.ServiceDiscovery.Prefix,
		DiscoveryType: util.SERVICE_DISCOVERY_ETCD,
		ServiceManager: global.V.ServiceManager,
	}
	serviceDiscovery ,err = util.NewServiceDiscovery(serviceOption)
	return serviceDiscovery,err
}
