//全局初始化
package initialize

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os/user"
	"strings"
	"zgoframe/core"
	"zgoframe/core/global"
	"zgoframe/model"
	"zgoframe/util"
)

type Initialize struct {
	Option InitOption
}

type MainEnvironment struct {
	RootDir     string `json:"root_dir"`
	RootDirName string `json:"root_dir_name"`

	GoVersion       string             `json:"go_version"`
	ExecUser        *user.User         `json:"-"`
	Cpu             string             `json:"cpu"`
	RootCtx         context.Context    `json:"-"`
	RootCancelFunc  context.CancelFunc `json:"-"`
	RootQuitFunc    func(source int)   `json:"-"`
	BuildTime       string
	BuildGitVersion string
}

type CmdParameter struct {
	Env              int    `json:"env"`
	ConfigSourceType string `json:"config_source_type"`
	ConfigFileType   string `json:"config_file_type"`
	ConfigFileName   string `json:"config_file_name"`
	EtcdUrl          string `json:"etcd_url"`
	Debug            int    `json:"debug"`
	TestFlag         string `json:"test_flag"`
}

type InitOption struct {
	CmdParameter
	MainEnvironment
	//Env               int    `json:"env"`
	//Debug             int    `json:"debug"`
	//ConfigType        string `json:"config_type"`
	//ConfigFileName    string `json:"config_file_name"`
	//ConfigSourceType  string `json:"config_source_type"`
	//EtcdConfigFindUrl string `json:"etcd_config_find_url"`
	//TestFlag          string `json:"-"`

	//RootDir     string `json:"root_dir"`
	//RootDirName string `json:"root_dir_name"`
	//
	//GoVersion      string             `json:"go_version"`
	//Cpu            string             `json:"cpu"`
	//RootCtx        context.Context    `json:"-"`
	//RootCancelFunc context.CancelFunc `json:"-"`
	//RootQuitFunc   func(source int)   `json:"-"`
}

func NewInitialize(option InitOption) *Initialize {
	initialize := new(Initialize)
	initialize.Option = option
	return initialize
}

//初始化-入口
func (initialize *Initialize) Start() error {
	prefix := "initialize ,"
	//初始化 : 配置信息
	viperOption := ViperOption{
		ConfigFileName: initialize.Option.ConfigFileName,
		ConfigFileType: initialize.Option.ConfigFileType,
		SourceType:     initialize.Option.ConfigSourceType,
		EtcdUrl:        initialize.Option.EtcdUrl,
		ENV:            initialize.Option.Env,
		PrintPrefix:    prefix,
	}

	util.MyPrint(prefix + "start CoreInitialize : config option~~ ")
	util.PrintStruct(initialize.Option, ":")
	util.MyPrint("-------")

	myViper, config, err := GetNewViper(viperOption)
	if err != nil {
		util.MyPrint(prefix+"GetNewViper err:", err)
		return err
	}
	//util.MyPrint(prefix + "read config info to assignment GlobalVariable , finish. ")
	global.V.Vip = myViper //全局变量管理者
	global.C = config      //全局变量
	//---config end -----

	//预/报警->推送器，这里是推送到3方，如：prometheus
	//ps:这个要优先zap日志类优化处理，因为zap里的<钩子>有用到,主要是日志里自动触发报警，略方便
	if global.C.Alert.Status == global.CONFIG_STATUS_OPEN {
		global.V.AlertPush = util.NewAlertPush(global.C.Alert.Host, global.C.Alert.Port, global.C.Alert.Uri, prefix)
	}
	//创建main日志类
	configZap := global.C.Zap
	configZap.FileName = "main"
	configZap.ModuleName = "main"
	mailZap, configZapReturn, err := GetNewZapLog(global.V.AlertPush, configZap)
	if err != nil {
		util.MyPrint("GetNewZapLog err:", err)
		return err
	}
	global.V.Zap = mailZap
	//初始化：mysql
	//原单库连接可以判断，现改成框架可以多库连接，判断 就意义了，数据结构还麻烦，只是多了一个status字段
	//PS:并不一定所有项目都用MYSQL，但基于<多APP/SERVICE>，强依赖 project_id，另外，日志也需要
	//if global.C.Mysql.Status != global.CONFIG_STATUS_OPEN {
	//	errMsg := "please open mysql db Module, because need project_id from read db."
	//	return errors.New(errMsg)
	//}
	//这个变量，主要是给gorm做日志使用，也就是DB的日志，最终也交由zap来接管
	util.LoggerZap = global.V.Zap
	//实例化gorm db
	global.V.GormList, err = GetNewGorm(prefix)
	if err != nil {
		return err
	}
	if len(global.V.GormList) <= 0 {
		return errors.New("至少有一个数据库需要被连接")
	}
	global.V.Gorm = global.V.GormList[0]
	//DB 快捷变量
	model.Db = global.V.Gorm
	//初始化APP信息，所有项目都需要有AppId或serviceId，因为要做验证，同时目录名也包含在里面
	err = InitProject(prefix)
	if err != nil {
		global.V.Zap.Error(prefix + err.Error())
		return err
	}
	//gorm 和 project 初始化(成功)完成后，给main日志增加公共输出项：projectId
	global.V.Zap = LoggerWithProject(global.V.Zap, global.V.Project.Id)
	//项目目录名，必须跟PROJECT里的key相同(key由驼峰转为下划线模式)
	initialize.Option.RootDirName, err = InitPath(initialize.Option.RootDir)
	if err != nil {
		global.V.Zap.Error(prefix + err.Error())
		return err
	}
	//项目的根目录
	global.V.RootDir = initialize.Option.RootDir
	global.V.Zap.Info(prefix + "global.V.RootDir: " + global.V.RootDir)
	//错误码 文案 管理（还未用起来，后期优化）
	global.V.Err, err = util.NewErrMsg(global.V.Zap, global.C.Http.StaticPath+global.C.System.ErrorMsgFile)
	if err != nil {
		global.V.Zap.Error(prefix + err.Error())
		return err
	}
	//基础类：用于恢复一个挂了的协程,避免主进程被panic fatal 带挂了，同时有重度次数控制
	global.V.RecoverGo = util.NewRecoverGo(global.V.Zap, 3)
	//redis
	var redisGo *util.MyRedisGo
	if global.C.Redis.Status == global.CONFIG_STATUS_OPEN {
		global.V.Redis, err = GetNewRedis(prefix)
		if err != nil {
			global.V.Zap.Error(prefix + " GetRedis " + err.Error())
			return err
		}

		redisGoOption := util.RedisGoOption{
			Host: global.C.Redis.Ip,
			Port: global.C.Redis.Port,
			Ps:   global.C.Redis.Password,
			Log:  global.V.Zap,
		}
		redisGo, _ = util.NewRedisConnPool(redisGoOption)

	}
	//http server
	if global.C.Http.Status == global.CONFIG_STATUS_OPEN {
		configZap = global.C.Zap
		configZap.FileName = "http"
		configZap.ModuleName = "http"
		//Http log zap 这里单独再开个zap 实例，用于专门记录http 请求
		HttpZap, _, err := GetNewZapLog(global.V.AlertPush, configZap)
		if err != nil {
			global.V.Zap.Error(prefix + "GetNewZapLog err:" + err.Error())
			return err
		}

		global.V.Gin, err = GetNewHttpGIN(HttpZap, prefix)
		if err != nil {
			global.V.Zap.Error(prefix + "GetNewHttpGIN err:" + err.Error())
			return err
		}
		HttpZap = LoggerWithProject(HttpZap, global.V.Project.Id)
	}
	//etcd
	if global.C.Etcd.Status == global.CONFIG_STATUS_OPEN {
		global.V.Etcd, err = GetNewEtcd(initialize.Option.Env, configZapReturn, prefix)
		if err != nil {
			global.V.Zap.Error(prefix + "GetNewEtcd err:" + err.Error())
			return err
		}
	}
	//服务管理器，这里跟project manager 有点差不多，不同的只是：project是DB中所有记录,service是type=N的情况
	//ps:之所以单独加一个模块，也是因为service有些特殊的结构变量，与project的结构变量不太一样
	global.V.ServiceManager, _ = util.NewServiceManager(global.V.Gorm)
	//service 服务发现，这里有个顺序，必须先实现化完成:serviceManager
	if global.C.ServiceDiscovery.Status == global.CONFIG_STATUS_OPEN {
		if global.C.Etcd.Status != global.CONFIG_STATUS_OPEN {
			return errors.New("ServiceDiscovery need Etcd open!")
		}
		global.V.ServiceDiscovery, err = GetNewServiceDiscovery()
		if err != nil {
			return err
		}
	}
	//metrics
	if global.C.Metrics.Status == global.CONFIG_STATUS_OPEN {
		myPushGateway := util.PushGateway{
			Status:  global.C.PushGateway.Status,
			Ip:      global.C.PushGateway.Ip,
			Port:    global.C.PushGateway.Port,
			JobName: global.V.Project.Name,
		}
		myMetricsOption := util.MyMetricsOption{
			Log:         global.V.Zap,
			NameSpace:   global.V.Project.Name,
			PushGateway: myPushGateway,
			Env:         global.C.System.ENV,
		}
		global.V.Metric = util.NewMyMetrics(myMetricsOption)

		if global.C.Http.Status != global.CONFIG_STATUS_OPEN {
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
	if global.C.Protobuf.Status == global.CONFIG_STATUS_OPEN {
		global.V.ProtoMap, err = util.NewProtoMap(global.V.Zap, dir, global.C.Protobuf.IdMapFileName, global.V.ProjectMng)
		if err != nil {
			util.MyPrint("GetNewViper err:", err)
			return err
		}
	}

	//websocket
	//if global.C.Websocket.Status == global.CONFIG_STATUS_OPEN{
	//	if global.C.Http.Status != global.CONFIG_STATUS_OPEN{
	//		return errors.New("Websocket need gin open!")
	//	}
	//	initSocket()
	//}

	//grpc
	if global.C.Grpc.Status == global.CONFIG_STATUS_OPEN {
		grpcManagerOption := util.GrpcManagerOption{
			//AppId: global.V.App.Id,
			//ServiceId: global.V.Service.Id,
			ProjectId: global.V.Project.Id,
			Log:       global.V.Zap,
		}
		if global.C.ServiceDiscovery.Status == global.CONFIG_STATUS_OPEN {
			grpcManagerOption.ServiceDiscovery = global.V.ServiceDiscovery
		}
		global.V.GrpcManager, _ = util.NewGrpcManager(grpcManagerOption)
	}
	//邮件模块
	if global.C.Email.Status == global.CONFIG_STATUS_OPEN {
		emailOption := util.EmailOption{
			Host:      global.C.Email.Host,
			Port:      global.C.Email.Port,
			FromEmail: global.C.Email.From,
			Password:  global.C.Email.Ps,
			AuthCode:  global.C.Email.AuthCode,
			Log:       global.V.Zap,
		}

		global.V.Email = util.NewMyEmail(emailOption)
	}
	//预/报警,这个是真正的直接报警，如：邮件 SMS 等，不是推送3方
	//ps:不推荐这么用，最好都统一推送3方报警机制
	if global.C.Alert.Status == global.CONFIG_STATUS_OPEN {
		global.V.AlertHook = util.NewAlertHook(-1, "程序出错了：#body#", "报错", util.ALERT_METHOD_SYNC, global.V.Zap)
		global.V.AlertHook.Email = global.V.Email
		//global.V.AlertHook.Alert("Aaaa")
		//util.ExitPrint(123123123)
	}

	//var netWayOption util.NetWayOption
	//if global.C.Gateway.Status == global.CONFIG_STATUS_OPEN {
	//	netWayOption = InitGateway()
	//	//	global.V.Gateway, err = InitGateway()
	//	//	if err != nil {
	//	//		global.V.Zap.Error(prefix + "InitGateway err:" + err.Error())
	//	//		return err
	//	//	}
	//}
	InitMyService(redisGo)

	global.C.System.ENV = initialize.Option.Env
	//启动http
	if global.C.Http.Status == global.CONFIG_STATUS_OPEN {
		RegGinHttpRoute() //这里注册项目自己的http 路由策略
		StartHttpGin(initialize.Option)
	}

	//_ ,cancelFunc := context.WithCancel(option.RootCtx)
	//进程通信相关
	ProcessPathFileName := "/tmp/" + global.V.Project.Name + ".pid"
	global.V.Process = util.NewProcess(ProcessPathFileName, initialize.Option.RootCancelFunc, global.V.Zap, initialize.Option.RootQuitFunc, initialize.OutHttpGetBaseInfo)
	global.V.Process.InitProcess()

	if initialize.Option.TestFlag != "" {
		core.DoTestAction(initialize.Option.TestFlag)
		return nil
	}

	return nil
}

func (initialize *Initialize) OutHttpGetBaseInfo() string {
	optionStr, _ := json.Marshal(initialize.Option)
	return string(optionStr)
}

func (initialize *Initialize) Quit() {
	global.V.Zap.Warn("init quit start:")
	if global.C.Http.Status == global.CONFIG_STATUS_OPEN {
		HttpServerShutdown()
	}

	if global.C.Redis.Status == global.CONFIG_STATUS_OPEN {
		RedisShutdown()
	}
	//这个得优于etcd先关
	if global.C.Grpc.Status == global.CONFIG_STATUS_OPEN {
		global.V.GrpcManager.Shutdown()
	}
	//这个得优于etcd先关
	if global.C.ServiceDiscovery.Status == global.CONFIG_STATUS_OPEN {
		global.V.ServiceDiscovery.Shutdown()
	}

	if global.C.Etcd.Status == global.CONFIG_STATUS_OPEN {
		global.V.Etcd.Shutdown()
	}
	//global.V.Websocket.Shutdown()

	GormShutdown()
	ViperShutdown()

	global.V.Zap.Warn("init quit finish.")
}

//=======================================================================================
func InitPath(rootDir string) (rootDirName string, err error) {
	pwdArr := strings.Split(rootDir, "/") //切割路径字符串
	rootDirName = pwdArr[len(pwdArr)-1]   //获取路径数组最后一个元素：当前路径的文件夹名
	//option.RootDirName = rootDirName
	//global.V.RootDir = option.RootDir
	//这里要求，DB中项目记录里：name 与项目目录名必须一致，防止有人错用/盗用projectId
	projectNameByte := util.CamelToSnake2([]byte(global.V.Project.Name))
	projectName := util.StrFirstToLower(string(projectNameByte))
	if rootDirName != projectName {
		//这里与CICD部署的时候冲突，先注释掉，回头想想怎么解决掉
		//return rootDirName, errors.New("mainDirName != app name , rootDirName : " + rootDirName + " , ProjectName:" + projectName)
	}

	return rootDirName, nil
}

func GetNewEtcd(env int, configZapReturn global.Zap, prefix string) (myEtcd *util.MyEtcd, err error) {
	//这个是给3方库：clientv3使用的
	//有点操蛋，我回头想想如何优化掉
	zl := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(configZapReturn.LevelInt8)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:      "json",
		EncoderConfig: zap.NewProductionEncoderConfig(),
		//OutputPaths:      []string{"stderr"},
		OutputPaths:      []string{"stdout", configZapReturn.FileName},
		ErrorOutputPaths: []string{"stderr"},
	}

	option := util.EtcdOption{
		ProjectName: global.V.Project.Name,
		ProjectENV:  env,
		//ProjectKey		: global.V.Project.Key,
		FindEtcdUrl: global.C.Etcd.Url,
		Username:    global.C.Etcd.Username,
		Password:    global.C.Etcd.Password,
		Ip:          global.C.Etcd.Ip,
		Port:        global.C.Etcd.Port,
		Log:         global.V.Zap,
		ZapConfig:   zl,
		PrintPrefix: prefix,
	}
	myEtcd, err = util.NewMyEtcdSdk(option)
	//util.ExitPrint(err)
	return myEtcd, err
}

func GetNewServiceDiscovery() (serviceDiscovery *util.ServiceDiscovery, err error) {
	serviceOption := util.ServiceDiscoveryOption{
		Log:  global.V.Zap,
		Etcd: global.V.Etcd,
		//Prefix	: "/service",
		Prefix:         global.C.ServiceDiscovery.Prefix,
		DiscoveryType:  util.SERVICE_DISCOVERY_ETCD,
		ServiceManager: global.V.ServiceManager,
	}
	serviceDiscovery, err = util.NewServiceDiscovery(serviceOption)
	return serviceDiscovery, err
}
