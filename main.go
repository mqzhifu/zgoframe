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
// @version 0.5 测试版Alta
// @description restful api 工具，模拟客户端请求，方便调试/测试
// @description 注：这只是一个工具，不是万能的，像：动态枚举类型、公共请求header、动态常量等,详细的请去 <a href="http://127.0.0.1:6060" target="_black">godoc</a> 里去查看
// @description 注：所有 header 遵循HTTP标准，即：自定义的header中每个key 以大写X开头，单词以中划线分隔，每个单词首字母大写
// @description 注：header的格式定义，参考结构：request.HeaderRequest，也可以调用调用接口获得:GET /tools/header/struct
// @description 注：所有的请求都需要包含header+body,header主要用于：基础数据收集+基础数据验证
// @description 注：99%的请求内容格式均是JSON(暂不支持兼容json+html)，只有上传文件例外(html+form)
// @description 注：所有接口的响应格式均是json格式 ，包含3个值: code data msg ,具体参考 model.httpresponse.Response
// @description 测试/开发人员：用户已上传的图片，查看，<a href="/static/upload/" target="_blank">点这里</a>
// @description 测试/开发人员：配置中心的文件，查看，<a href="/static/data/config/" target="_blank">点这里</a>
// @description <a href="/static/html/cicd.html" target="_blank">测试cicd</a> <a href="/static/html/frame_sync.html" target="_blank">测试帧同步</a> <a href="/static/html/file_upload.html" target="_blank">测试多文件上传</a>
// @license.name 小z
// @contact.name 小z
// @contact.email 78878296@qq.com
// @tag.name Base
// @tag.description 基础操作（不需要登陆，但是会验证头信息 , X-SourceType X-Access X-Project 等）
// @tag.name User
// @tag.description 用户相关操作(需要登陆，头里加X-Token = jwt)
// @tag.name System
// @tag.description 系统管理(需要二次认证)，管理员使用，普通用户不要访问
// @tag.name Cicd
// @tag.description 自动化部署与持续集成
// @tag.name Mail
// @tag.description 站内信/内部邮件通知
// @tag.name GameMatch
// @tag.description 游戏匹配机制
// @tag.name persistence
// @tag.description 持久化(文件/日志收集)，注：文件上传目前仅支持HTTP协议，也就是form+multipart/form-data模式.不支持：分片传输，断点续传等功能
// @tag.name ConfigCenter
// @tag.description 配置中心，它有几个维度注意下： 环境->项目->文件->模块，项目这个维度http-header头中是公共的且已处理，余下3个请求的时候都要带上。目前仅支持：toml格式，后期可加ymal和ini
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
