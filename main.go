package main

import (
	"context"
	_ "embed"
	"flag"
	"os"
	"os/user"
	"strconv"
	"zgoframe/core"
	"zgoframe/core/global"
	"zgoframe/core/initialize"
	_ "zgoframe/docs"
	"zgoframe/util"
)

var initializeVar *initialize.Initialize

// @title z golang 框架
// @version 0.1 测试版
// @description restful api 工具，模拟客户端请求，方便调试/测试
// @description 注：这只是一个工具，不是万能的，像：动态枚举类型、公共请求header、动态常量等
// @description 详细的请去 <a href="http://127.0.0.1:6060" target="_black">godoc</a> 里去查看
// @license.name header.BaseInfo-demo:{   "app_version": "v1.1.1",   "device": "iphone",   "device_id": "aaaaaaaa",   "device_version": "12",   "dpi": "390x844",   "ip": "127.0.0.1",   "lat": "21.1111",   "lon": "32.4444",   "os": 1,   "os_version": "11",   "referer": "" }
// @tag.name Base
// @tag.description 不需要登陆，但是会验证头信息 , X-SourceType X-Access X-Project 等，(注：header 中的每个key X开头)
// @tag.name User
// @tag.description 用户相关操作(需要登陆，头里加X-Token = jwt)
// @tag.name System
// @tag.description 系统管理(需要二次认证)
// @tag.name Cicd
// @tag.description 自动化部署与持续集成
// @tag.name Mail
// @tag.description 站内信/内部邮件通知
// @tag.name ConfigCenter
// @tag.description 配置中心
// @securityDefinitions.apikey ApiKeyAuth
// @name xa
// @name X-Token
// @in header

func main() {
	prefix := "main "
	//获取<环境变量>枚举值
	envList := util.GetConstListEnv()
	envListStr := util.ConstListEnvToStr()
	//配置读取源类型，1 文件  2 etcd
	configSourceType := flag.String("cs", global.DEFAULT_CONFIG_SOURCE_TYPE, "configSource:file or etcd")
	//配置文件的类型
	configFileType := flag.String("ct", global.DEFAULT_CONFIT_TYPE, "configFileType")
	//配置文件的名称
	configFileName := flag.String("cfn", global.DEFAULT_CONFIG_FILE_NAME, "configFileName")
	//获取etcd 配置信息的URL
	etcdUrl := flag.String("etl", "http://127.0.0.1/getEtcdCluster/Ip/Port", "get etcd config url")
	//当前环境,env:local test pre dev online
	env := flag.Int("e", 0, "must require , "+envListStr)
	//DEBUG模式
	debug := flag.Int("debug", 0, "startup debug mode level")
	//是否为CICD模式
	//deploy 				:= flag.String("dep", "", "deploy")//部署模式下，启动程序只是为了测试脚本正常，因为之后，要立刻退出
	//开启自动测试模式
	testFlag 			:= flag.String("t", "", "testFlag:empty or 1")
	//解析命令行参数
	flag.Parse()
	//检测环境变量值ENV是否正常
	if !util.CheckEnvExist(*env) {
		msg := prefix + " argv env , is err :"
		util.MyPrint(msg, envList)
		panic(msg + strconv.Itoa(*env))
	}

	imUser, _ := user.Current()
	util.MyPrint(prefix + "exec script user info , name: " + imUser.Name + " uid: " + imUser.Uid + " , gid :" + imUser.Gid + " ,homeDir:" + imUser.HomeDir)

	pwd, _ := os.Getwd() //当前路径
	util.MyPrint(prefix + "exec script pwd:" + pwd)
	//开始初始化模块
	//主协程的 context
	util.MyPrint(prefix + "create cancel context")
	mainCxt, mainCancelFunc := context.WithCancel(context.Background())
	//初始化模块需要的参数
	initOption := initialize.InitOption{
		Env:               *env,
		Debug:             *debug,
		ConfigType:        *configFileType,
		ConfigFileName:    *configFileName,
		ConfigSourceType:  *configSourceType,
		EtcdConfigFindUrl: *etcdUrl,
		RootDir:           pwd,
		RootCtx:           mainCxt,
		RootCancelFunc:    mainCancelFunc,
		RootQuitFunc:      QuitAll,
		TestFlag :		   *testFlag,
	}
	//开始正式全局初始化
	initializeVar = initialize.NewInitialize(initOption)
	err := initializeVar.Start()
	if err != nil {
		util.MyPrint(prefix+"initialize.Init err:", err)
		panic(prefix + "initialize.Init err:" + err.Error())
		return
	}
	//执行用户自己的一些功能
	go core.DoMySelf(*testFlag)
	//监听外部进程信号
	go global.V.Process.DemonSignal()
	util.MyPrint(prefix + "wait mainCxt.done...")
	select {
	case <-mainCxt.Done():
		QuitAll(1)
	}

	util.MyPrint(prefix + "end.")
}

func QuitAll(source int) {
	defer func() {
		global.V.Process.DelPid()
	}()

	global.V.Zap.Warn("main quit , source : " + strconv.Itoa(source))
	initializeVar.Quit()

	util.MyPrint("main QuitAll finish.")
}
