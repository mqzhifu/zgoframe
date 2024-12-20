// 全局初始化
package initialize

import (
	"errors"
	"fmt"
	"strings"
	v1 "zgoframe/api/v1"
	"zgoframe/core"
	"zgoframe/core/global"
	"zgoframe/http/router"
	"zgoframe/util"
)

type Initialize struct{}

func NewInitialize() *Initialize {
	initialize := new(Initialize)
	return initialize
}

// 框架-初始化-入口
func (initialize *Initialize) Start() error {
	//---read config file start -----
	prefix := "initialize ," //输出日志的前缀
	//开始：读取配置文件
	myViper, config, err := GetNewViper(prefix)
	if err != nil {
		fmt.Println(prefix+"GetNewViper err:", err)
		return err
	}
	global.V.Base.Vip = myViper //全局变量管理者
	global.C = config           //全局变量
	//--- read config file end -----

	global.V.Util.StaticFileSystem = util.NewStaticFileSystem(global.V.Base.StaticFileSys, global.MainCmdParameter.BuildStatic)
	//创建全局 日志
	err = createLogByCategory(prefix) //创建 main 和 http 日志-类
	if err != nil {
		return err
	}

	//邮件与短信优先初始化，是一但有报警，就可以直接发邮件/短信

	//邮件模块
	if global.C.Email.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		emailOption := util.EmailOption{
			Host:      global.C.Email.Host,
			Port:      global.C.Email.Port,
			FromEmail: global.C.Email.From,
			Password:  global.C.Email.Ps,
			AuthCode:  global.C.Email.AuthCode,
			Log:       global.V.Base.Zap,
		}

		global.V.Util.Email, err = util.NewMyEmail(emailOption)
		if err != nil {
			return err
		}
	}

	//短信模块
	if global.C.AliSms.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		op := util.AliSmsOp{
			AccessKeyId:     global.C.AliSms.AccessKeyId,
			AccessKeySecret: global.C.AliSms.AccessKeySecret,
			Endpoint:        global.C.AliSms.Endpoint,
		}
		global.V.Util.AliSms, err = util.NewAliSms(op)
		if err != nil {
			util.MyPrint(prefix+"util.NewAliSms err:", err)
			return err
		}
	}

	//预/报警,这个是真正的直接报警，如：邮件 SMS 等，不是推送3方
	//ps:不推荐这么用，最好都统一推送3方报警机制
	//if global.C.Alert.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
	//	global.V.Base.AlertHook = util.NewAlertHook(-1, "程序出错了：#body#", "报错", util.ALERT_METHOD_SYNC, global.V.Base.Zap)
	//	global.V.Base.AlertHook.Email = global.V.Base.Email
	//	global.V.Base.AlertHook.Alert("Aaaa")
	//}

	//预警/报警->推送器，这里是推送到3方服务，如：prometheus，而不是直接发邮件/短信
	//ps:这个要优先zap日志类优化处理，因为zap里的<钩子>有用到,主要是日志里自动触发报警，略方便
	if global.C.AlertPush.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		global.V.Util.AlertPush, err = util.NewAlertPush(global.C.AlertPush.Host, global.C.AlertPush.Port, global.C.AlertPush.Uri, prefix)
		if err != nil {
			return err
		}
	}
	//实例化gorm db，因为要使用 project/app 记录
	global.V.Base.GormList, err = GetNewGorm(prefix)
	if err != nil {
		return err
	}
	if len(global.V.Base.GormList) <= 0 {
		return errors.New("至少有一个数据库需要被连接")
	}
	//默认取第一个DB配置
	global.V.Base.Gorm = global.V.Base.GormList[0]
	//初始化APP信息，所有项目都需要有AppId或serviceId.因为:
	//1. header 要做验证
	//2. CICD 目录名也包含在里面
	//3. 日志里要输出
	err = InitProject(prefix)
	if err != nil {
		global.V.Base.Zap.Error(prefix + err.Error())
		return err
	}
	//gorm 和 project 初始化(成功)完成后，给日志:增加公共输出项：projectId
	global.V.Base.Zap = LoggerWithProject(global.V.Base.Zap, global.V.Util.Project.Id)
	global.V.Base.HttpZap = LoggerWithProject(global.V.Base.HttpZap, global.V.Util.Project.Id)
	//项目目录名，必须跟PROJECT里的key相同(key由驼峰转为下划线模式)
	global.MainEnv.RootDirName, err = InitPath(global.MainEnv.RootDir)
	if err != nil {
		global.V.Base.Zap.Error(prefix + err.Error())
		return err
	}
	//项目的根目录
	global.V.Base.Zap.Info(prefix + "global.V.Base.RootDir: " + global.MainEnv.RootDir)
	//错误码 文案 管理（还未用起来，后期优化）
	errorMsgFileContentDir := global.C.Http.StaticPath + "/" + global.C.System.ErrorMsgFile
	errorMsgFileContent, err := global.V.Util.StaticFileSystem.GetStaticFileContentLine(errorMsgFileContentDir)
	if err != nil {
		global.V.Base.Zap.Error(prefix + err.Error())
		return err
	}
	global.V.Util.Err, err = util.NewErrMsg(global.V.Base.Zap, errorMsgFileContentDir, errorMsgFileContent)
	if err != nil {
		global.V.Base.Zap.Error(prefix + err.Error())
		return err
	}
	//基础类：用于恢复一个挂了的协程,避免主进程被panic fatal 带挂了，同时有重试次数控制
	global.V.Util.RecoverGo = util.NewRecoverGo(global.V.Base.Zap, 3)
	//redis
	if global.C.Redis.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		global.V.Base.Redis, err = GetNewRedis(prefix)
		if err != nil {
			global.V.Base.Zap.Error(prefix + " GetRedis " + err.Error())
			return err
		}
		//这个是另外一个redis sdk库，算是备用吧
		redisGoOption := util.RedisGoOption{
			Host: global.C.Redis.Ip,
			Port: global.C.Redis.Port,
			Ps:   global.C.Redis.Password,
			Log:  global.V.Base.Zap,
		}
		global.V.Base.RedisGo, _ = util.NewRedisConnPool(redisGoOption)
	}
	//http server
	if global.C.Http.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		global.V.Base.Gin, err = GetNewHttpGIN(global.V.Base.HttpZap, prefix)
		if err != nil {
			global.V.Base.Zap.Error(prefix + "GetNewHttpGIN err:" + err.Error())
			return err
		}
		//给 HTTP的 日志文件 增加 projectId 项
		global.V.Base.HttpZap = LoggerWithProject(global.V.Base.HttpZap, global.V.Util.Project.Id)
	}
	//grpc
	if global.C.Grpc.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		grpcManagerOption := util.GrpcManagerOption{
			//AppId: global.V.Base.App.Id,
			//ServiceId: global.V.Base.Service.Id,
			ProjectId: global.V.Util.Project.Id,
			Log:       global.V.Base.Zap,
		}
		if global.C.ServiceDiscovery.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
			grpcManagerOption.ServiceDiscovery = global.V.Util.ServiceDiscovery
		}
		global.V.Util.GrpcManager, _ = util.NewGrpcManager(grpcManagerOption)
	}

	initialize.InitMetric(prefix)
	initialize.InitEtcd(prefix)
	initialize.InitEls(prefix)
	initialize.InitProtobuf(prefix)
	initialize.InitAliOss(prefix)

	InitFileManager()

	//启动http
	if global.C.Http.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		router.RegGinHttpRoute() //这里注册项目自己的http 路由策略
		StartHttpGin()
	}
	//实例化所有微服务
	v1.ApiServices()
	//service.NewMyService()//v1.api 里面进行初始化，也就是谁用：谁就初始化，避免循环引包
	initialize.InitServiceDiscovery(prefix)
	//_, cancelFunc := context.WithCancel(option.RootCtx)
	//进程通信相关
	ProcessPathFileName := "/tmp/" + global.V.Util.Project.Name + ".pid"
	global.V.Util.Process = util.NewProcess(ProcessPathFileName, global.MainEnv.RootCancelFunc, global.V.Base.Zap, global.MainEnv.RootQuitFunc, initialize.OutHttpGetBaseInfo)
	global.V.Util.Process.InitProcess()

	if global.MainCmdParameter.TestFlag != "" {
		core.DoTestAction(global.MainCmdParameter.TestFlag)
		return nil
	}

	return nil
}

func (initialize *Initialize) OutHttpGetBaseInfo() string {
	//optionStr, _ := json.Marshal(initialize.Option)
	//return string(optionStr)
	return "img OutHttpGetBaseInfo"
}

func (initialize *Initialize) Quit() {
	global.V.Base.Zap.Warn("init quit start:")
	if global.C.Http.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		HttpServerShutdown()
	}

	if global.C.Redis.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		RedisShutdown()
	}
	//这个得优于etcd先关
	if global.C.Grpc.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		global.V.Util.GrpcManager.Shutdown()
	}
	//这个得优于etcd先关
	if global.C.ServiceDiscovery.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		global.V.Util.ServiceDiscovery.Shutdown()
	}

	if global.C.Etcd.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		global.V.Util.Etcd.Shutdown()
	}

	//global.V.Base.Websocket.Shutdown()

	GormShutdown()
	ViperShutdown()

	global.V.Base.Zap.Warn("init quit finish.")
}

// =======================================================================================
func InitPath(rootDir string) (rootDirName string, err error) {
	pwdArr := strings.Split(rootDir, "/") //切割路径字符串
	rootDirName = pwdArr[len(pwdArr)-1]   //获取路径数组最后一个元素：当前路径的文件夹名
	//这里要求，DB中项目记录里：name 与项目目录名必须一致，防止有人错用/盗用projectId
	projectNameByte := util.CamelToSnake2([]byte(global.V.Util.Project.Name))
	projectName := util.StrFirstToLower(string(projectNameByte))
	if rootDirName != projectName {
		//这里与CICD部署的时候冲突，先注释掉，回头想想怎么解决掉
		return rootDirName, errors.New("mainDirName != app name , rootDirName : " + rootDirName + " , ProjectName:" + projectName)
	}

	return rootDirName, nil
}

func GetNewServiceDiscovery() (serviceDiscovery *util.ServiceDiscovery, err error) {
	serviceOption := util.ServiceDiscoveryOption{
		Log: global.V.Base.Zap,
		//Etcd:           global.V.Base.Etcd,
		Prefix:         global.C.ServiceDiscovery.Prefix,
		DiscoveryType:  util.SERVICE_DISCOVERY_ETCD,
		ServiceManager: global.V.Util.ServiceManager,
		//Prefix	: "/service",
	}
	serviceDiscovery, err = util.NewServiceDiscovery(serviceOption)
	return serviceDiscovery, err
}
