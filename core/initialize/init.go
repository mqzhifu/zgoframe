package initialize

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/util"
)


func Init(ENV string ,configType string , configFileName string )error{
	myViper,config,err := GetNewViper(configType,configFileName)
	if err != nil{
		util.MyPrint("GetNewViper err:",err)
		return err
	}
	global.V.Vip = myViper
	global.C = config

	if global.C.System.AppId <=0 {
		return errors.New("appid is empty")
	}

	global.V.App ,err  = GetNewApp()
	if err != nil{
		util.MyPrint("GetNewApp err:",err)
		return err
	}

	global.V.Zap , err  = GetNewZapLog()
	if err != nil{
		util.MyPrint("GetNewZapLog err:",err)
		return err
	}

	if global.C.Redis.Status == global.CONFIG_STATUS_OPEN{
		global.V.Redis ,err = GetNewRedis()
		if err != nil{
			util.MyPrint("GetRedis err:",err)
			return err
		}
	}

	if global.C.Mysql.Status == global.CONFIG_STATUS_OPEN{
		global.V.Gorm ,err = GetNewGorm()
		if err != nil{
			util.MyPrint("GetGorm err:",err)
			return err
		}
	}

	if global.C.Http.Status == global.CONFIG_STATUS_OPEN{
		global.V.Gin ,err = GetNewHttpGIN()
		if err != nil{
			util.MyPrint("GetNewHttpGIN err:",err)
			return err
		}
		StartHttpGin()
	}

	if global.C.Etcd.Status  == global.CONFIG_STATUS_OPEN{
		global.V.Etcd ,err = GetNewEtcd()
		if err != nil{
			util.MyPrint("GetNewEtcd err:",err)
			return err
		}
	}

	if global.C.Service.Status  == global.CONFIG_STATUS_OPEN{
		//global.V.Service ,err = GetNewService()
		//if err != nil{
		//	util.MyPrint("GetNewEtcd err:",err)
		//	return err
		//}
		global.V.Service  = GetNewService()
	}


	global.V.Metric =  util.NewMyMetrics()

	global.C.System.ENV = ENV

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
		Log			: global.V.Zap,
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

func GetNewViper(ConfigType string,ConfigName string)(myViper *viper.Viper,config global.Config,err error){
	util.MyPrint("ConfigType:",ConfigType ," , ConfigName:",ConfigName)
	myViper = viper.New()
	myViper.SetConfigType(ConfigType)
	//myViper.SetConfigName(ConfigName + "." + ConfigType)
	myViper.SetConfigFile(ConfigName + "." + ConfigType)
	err = myViper.ReadInConfig()
	if err != nil{
		util.MyPrint("myViper.ReadInConfig() err :",err)
		return myViper,config,err
	}

	//config := Config{}
	err = myViper.Unmarshal(&config)
	if err != nil{
		util.MyPrint(" myViper.Unmarshal err:",err)
		return myViper,config,err
	}

	if config.Viper.Watch == global.CONFIG_STATUS_OPEN{
		util.MyPrint("viper watch open")
		myViper.WatchConfig()
		handleFunc := func(in fsnotify.Event) {
			util.MyPrint("myViper.WatchConfig onChange:",in.Name ,in.String())

			//if err := viper.Unmarshal(Conf); err != nil {
			//	panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
			//}
		}
		myViper.OnConfigChange(handleFunc)
		viper.OnConfigChange(handleFunc)
	}

	return myViper,config,nil
}


