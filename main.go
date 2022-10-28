package main

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download
//go:generate go get -u github.com/swaggo/swag/cmd/swag@v1.7.9
//go:generate $HOME/go/bin/swag init --parseDependency --parseInternal --parseDepth 3

import (
	"context"
	_ "embed"
	"flag"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"zgoframe/core"
	"zgoframe/core/global"
	"zgoframe/core/initialize"
	_ "zgoframe/docs"
	"zgoframe/util"
)

//go build -ldflags "-X main.BuildGitVersion='1.0.9' -X 'main.BUILD_TIME=`date`' " -o zgo
var (
	BuildTime       string
	BuildGitVersion string
)

var initializeVar *initialize.Initialize

// @title z golang 框架
// @version 0.5 测试版Alta
// @description restful api 工具，模拟客户端请求，方便调试/测试
// @description 注：这只是一个工具，不是万能的，像：动态枚举类型、公共请求header、动态常量等,详细的请去 <a href="http://godoc.seedreality.com" target="_black">godoc</a> 里去查看
// @description 注：所有 header 遵循HTTP标准，即：自定义的header中每个key 以大写X开头，单词以中划线分隔，每个单词首字母大写
// @description 注：header的格式定义，参考结构：request.HeaderRequest，也可以调用调用接口获得:GET /tools/header/struct
// @description 注：所有的请求都需要包含header+body,header主要用于：基础数据收集+基础数据验证
// @description 注：99%的请求内容格式均是JSON(暂不支持兼容json+html-form)，只有上传文件例外(html+form)
// @description 注：所有接口均支持：跨域请求
// @description 注：所有接口的响应格式均是json格式 ，包含3个值: code data msg ,具体参考 model.httpresponse.Response
// @description 测试/开发人员：用户已上传的图片，查看，<a href="http://static.seedreality.com/upload/" target="_blank">点这里</a>
// @description 测试/开发人员：配置中心的文件，查看，<a href="http://static.seedreality.com/data/config/" target="_blank">点这里</a>
// @description 后台UI：<a href="http://admin.seedreality.com" target="_blank">点这里</a>
// @description <a href="http://static.seedreality.com/html/cicd.html" target="_blank">测试cicd</a> <a href="http://static.seedreality.com/html/frame_sync.html" target="_blank">测试帧同步</a> <a href="http://static.seedreality.com/html/file_upload.html" target="_blank">测试多文件上传</a>
// @license.name 小z
// @contact.name 小z
// @contact.email 78878296@qq.com
// @tag.name Base
// @tag.description 基础操作（不需要登陆，但是会验证头信息 , X-SourceType X-Access X-Project 等）
// @tag.name User
// @tag.description 用户相关操作(需要登陆，头里加X-Token = jwt)
// @tag.name System
// @tag.description 系统管理(需要二次认证)，管理员使用，普通用户不要访问
// @tag.name TwinAgora
// @tag.description 数字孪生 - agora
// @tag.name Cicd
// @tag.description 自动化部署与持续集成
// @tag.name Mail
// @tag.description 站内信/内部邮件通知
// @tag.name Callback
// @tag.description 3方回调/推送通知
// @tag.name GameMatch
// @tag.description 游戏匹配机制
// @tag.name persistence
// @tag.description 持久化(文件/日志收集
// @tag.name file
// @tag.description 文件系统，如：上传/下载文件，文件包括：图片、视频、文件流等. 注：文件上传目前仅支持HTTP协议，也就是form+multipart/form-data模式.不支持：分片传输，断点续传等功能
// @tag.name ConfigCenter
// @tag.description 配置中心，它有几个维度注意下： 环境->项目->文件->模块，项目这个维度http-header头中是公共的且已处理，余下3个请求的时候都要带上。目前仅支持：toml格式，后期可加ymal和ini
// @securityDefinitions.apikey ApiKeyAuth
// @name xa
// @name X-Token
// @in header

func main() {
	//编译打进去的两个参数：BuildTime 编译时间，编译的git版本号
	util.MyPrint("code , BuildTime:", BuildTime, " BuildGitVersion:", BuildGitVersion)
	//日志前缀
	prefix := "main "
	//处理指令行的参数
	cmdParameter := processCmdParameter(prefix)
	util.MyPrint("cmdParameter")
	util.PrintStruct(cmdParameter, ":")
	//util.MyPrint(prefix+" cmd parameter:", cmdParameter)
	//获取当前脚本执行用户信息
	imUser, _ := user.Current()
	util.MyPrint(prefix + "exec script user info , name: " + imUser.Name + " uid: " + imUser.Uid + " , gid :" + imUser.Gid + " ,homeDir:" + imUser.HomeDir)
	//当前脚本执行的路径
	pwd, _ := os.Getwd()
	util.MyPrint(prefix + "exec script pwd:" + pwd)
	//开始初始化模块
	//main主协程的 context
	util.MyPrint(prefix + "create cancel context")
	ttt()
	mainCxt, mainCancelFunc := context.WithCancel(context.Background())
	mainEnvironment := global.MainEnvironment{
		RootDir:         pwd,
		GoVersion:       runtime.Version(),
		Cpu:             runtime.GOARCH,
		RootCtx:         mainCxt,
		RootCancelFunc:  mainCancelFunc,
		RootQuitFunc:    QuitAll,
		ExecUser:        imUser,
		BuildTime:       BuildTime,
		BuildGitVersion: BuildGitVersion,
	}
	global.MainEnv = mainEnvironment
	global.MainCmdParameter = cmdParameter
	//开始正式全局初始化
	initializeVar = initialize.NewInitialize()
	err := initializeVar.Start()
	if err != nil {
		util.MyPrint(prefix+"initialize.Init err:", err)
		panic(prefix + "initialize.Init err:" + err.Error())
	}

	//执行用户自己的一些功能
	go core.DoMySelf()
	//监听外部进程信号
	go global.V.Process.DemonSignal()
	util.MyPrint(prefix + "wait mainCxt.done...")
	select {
	case <-mainCxt.Done(): //阻塞
		QuitAll(1)
	}

	util.MyPrint(prefix + "end.")
}

//处理指令行参数
func processCmdParameter(prefix string) global.CmdParameter {
	//获取<环境变量>枚举值
	envList := util.GetConstListEnv()
	envListStr := util.ConstListEnvToStr()
	//当前环境,env:local test pre dev online
	env := flag.Int("e", 0, "must require , "+envListStr)
	//配置读取源类型，1 文件  2 etcd
	configSourceType := flag.String("cs", core.DEFAULT_GLOBAL_CONFIG_TYPE_FILE, "configSource:file or etcd")
	//配置文件的类型:toml yaml
	configFileType := flag.String("ct", core.DEFAULT_GLOBAL_CONFIG_FILE_TYPE, "configFileType")
	//配置文件的名称
	configFileName := flag.String("cfn", core.DEFAULT_GLOBAL_CONFIG_FILE_NAME, "configFileName")
	//获取etcd 配置信息的URL,也可以把配置文件中的信息存于ETCD中，通过URL请求ETCD获取
	etcdUrl := flag.String("etl", "http://127.0.0.1/getEtcdCluster/Ip/Port", "get etcd config url")
	//DEBUG模式
	debug := flag.Int("debug", 0, "startup debug mode level")
	//开启自动测试模式
	testFlag := flag.String("t", "", "testFlag:empty or 1")
	//解析命令行参数
	flag.Parse()
	//检测环境变量值ENV是否正常
	if !util.CheckEnvExist(*env) {
		msg := prefix + " argv env , is err :"
		util.MyPrint(msg, envList)
		panic(msg + strconv.Itoa(*env))
	}

	cmdParameter := global.CmdParameter{
		Env:              *env,
		ConfigSourceType: *configSourceType,
		ConfigFileType:   *configFileType,
		ConfigFileName:   *configFileName,
		EtcdUrl:          *etcdUrl,
		Debug:            *debug,
		TestFlag:         *testFlag,
	}

	return cmdParameter
}
func ttt() {
	//ss := util.NewConstHandle()
	//util.MyPrint(ss.EnumConstPool)
	//util.ExitPrint(33)
}
func QuitAll(source int) {
	defer func() {
		global.V.Process.DelPid()
	}()

	global.V.Zap.Warn("main quit , source : " + strconv.Itoa(source))
	initializeVar.Quit()

	util.MyPrint("main QuitAll finish.")
}
