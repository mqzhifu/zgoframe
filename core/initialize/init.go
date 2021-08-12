package initialize

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/util"
)

func Init(ENV string ,configType string , configFileName string,configSourceType string ,etcdUrl string)error{
	viperOption := ViperOption{
		ConfigFileName: configFileName,
		ConfigFileType:  configType,
		SourceType: configSourceType,
		EtcdUrl: etcdUrl,
		ENV: ENV,
	}
	//初始化配置信息
	myViper,config,err := GetNewViper(viperOption)
	if err != nil{
		util.MyPrint("GetNewViper err:",err)
		return err
	}
	global.V.Vip = myViper
	global.C = config
	//初始化APP信息，所有项目都需要有APPID
	if global.C.System.AppId <=0 {
		return errors.New("appid is empty")
	}

	global.V.App ,err  = GetNewApp()
	if err != nil{
		util.MyPrint("GetNewApp err:",err)
		return err
	}
	//预警器
	if global.C.Alert.Status == global.CONFIG_STATUS_OPEN{
		global.V.Alert = util.NewAlert(global.C.Alert.Ip,global.C.Alert.Port,global.C.Alert.Uri)
	}
	//日志
	global.V.Zap , err  = GetNewZapLog(global.V.Alert)
	if err != nil{
		util.MyPrint("GetNewZapLog err:",err)
		return err
	}
	//redis
	if global.C.Redis.Status == global.CONFIG_STATUS_OPEN{
		global.V.Redis ,err = GetNewRedis()
		if err != nil{
			util.MyPrint("GetRedis err:",err)
			return err
		}
	}
	//mysql
	if global.C.Mysql.Status == global.CONFIG_STATUS_OPEN{
		global.V.Gorm ,err = GetNewGorm()
		if err != nil{
			util.MyPrint("GetGorm err:",err)
			return err
		}
	}
	//http server
	if global.C.Http.Status == global.CONFIG_STATUS_OPEN{
		global.V.Gin ,err = GetNewHttpGIN()
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

	if global.C.Service.Status  == global.CONFIG_STATUS_OPEN{
		global.V.Service = GetNewService()
		global.V.Service  = GetNewService()
	}
	//metrics
	if global.C.Metrics.Status == global.CONFIG_STATUS_OPEN{
		global.V.Metric =  util.NewMyMetrics()

		if global.C.Http.Status != global.CONFIG_STATUS_OPEN{
			return errors.New("metrics nee gin open!")
		}
		global.V.Gin.GET("/metrics", gin.WrapH(promhttp.Handler()))
		//global.V.Gin.GET("/metrics/count", func(c *gin.Context) {
		//	global.V.Metric.CounterInc("paySuccess")
		//})
		//
		//global.V.Gin.GET("/metrics/gauge", func(c *gin.Context) {
		//	global.V.Metric.CounterInc("payUser")
		//})
		//global.V.Metric.Test()
	}

	//hook := util.NewAlertHook()

	global.C.System.ENV = ENV

	if global.C.Http.Status == global.CONFIG_STATUS_OPEN{
		StartHttpGin()
	}

	return nil
}

func Quit(){
	HttpServerShutdown()
	global.V.Redis.Close()
	db , _ := global.V.Gorm.DB()
	db.Close()
	global.V.Vip.WatchRemoteConfig()
}

func GetNewApp()(util.App,error){
	appM := util.NewAppManager()
	app ,err := appM.GetById(global.C.System.AppId)
	if err {
		return app,errors.New("AppId not match : " + strconv.Itoa(global.C.System.AppId) )
	}
	return app,nil
}

func GetNewEtcd()(myEtcd *util.MyEtcd,err error){
	option := util.EtcdOption{
		AppName		: global.V.App.Name,
		AppENV		: global.C.System.ENV,
		FindEtcdUrl :	global.C.Etcd.Url,
		Username	: global.C.Etcd.Username,
		Password	: global.C.Etcd.Port,
		Ip			: global.C.Etcd.Ip,
		Port		: global.C.Etcd.Port,
		Log: global.V.Zap,
	}
	myEtcd,err  = util.NewMyEtcdSdk(option)
	return myEtcd,err
}

func GetNewService()*util.Service {
	serviceOption := util.ServiceOption{
		Log: global.V.Zap,
		Etcd: global.V.Etcd,
		Prefix: global.V.App.Name,
	}
	myService := util.NewService(serviceOption)
	return myService
}




