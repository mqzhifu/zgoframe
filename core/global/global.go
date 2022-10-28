//全局容器，1 配置信息中的变量 2 公共初始好的类包
package global

import "C"
import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"os/user"
	"zgoframe/model"
	"zgoframe/util"
)

type Global struct {
	Vip              *viper.Viper
	Zap              *zap.Logger
	Redis            *util.MyRedis
	RedisGo          *util.MyRedisGo
	Gin              *gin.Engine
	Gorm             *gorm.DB   //多数据库模式下，有一个库肯定会被经常访问，这里加一个快捷链接
	GormList         []*gorm.DB //所有数据库，连接成功后的列表
	Project          model.Project
	ProjectMng       *util.ProjectManager
	Etcd             *util.MyEtcd
	HttpServer       *http.Server
	Metric           *util.MyMetrics
	GrpcManager      *util.GrpcManager
	AlertPush        *util.AlertPush //报警推送： prometheus
	AlertHook        *util.AlertHook //报警：邮件 手机
	Websocket        *util.Websocket
	ConnMng          *util.ConnManager
	RecoverGo        *util.RecoverGo
	ProtoMap         *util.ProtoMap
	Process          *util.Process
	Err              *util.ErrMsg
	Email            *util.MyEmail
	DocsManager      *util.FileManager
	ImgManager       *util.FileManager
	VideoManager     *util.FileManager
	NetWay           *util.NetWay
	ServiceManager   *util.ServiceManager   //管理已注册的服务
	ServiceDiscovery *util.ServiceDiscovery //管理服务发现，会用到上面的ServiceManager
	AliOss           *util.AliOss           //阿里网盘
	MyService        *MyService             //内部快捷服务
}

//main主协程的一些参数
type MainEnvironment struct {
	RootDir         string             `json:"root_dir"`
	RootDirName     string             `json:"root_dir_name"`
	GoVersion       string             `json:"go_version"` //当前go版本
	ExecUser        *user.User         `json:"-"`          //执行该脚本的用户信息
	Cpu             string             `json:"cpu"`        //cpu信息
	RootCtx         context.Context    `json:"-"`          //main的上下文，级别最高
	RootCancelFunc  context.CancelFunc `json:"-"`          //main的取消函数，该管理如果能读出值，main会主动退出
	RootQuitFunc    func(source int)   `json:"-"`          //这是个函数，子级可直接驱动：退出MAIN
	BuildTime       string             //编译时：时间
	BuildGitVersion string             //编译时：git版本号
}

//指令行 收集的参数
type CmdParameter struct {
	Env              int    `json:"env"`                //当前环境
	ConfigSourceType string `json:"config_source_type"` //文件 | etcd
	ConfigFileType   string `json:"config_file_type"`   //项目的配置：文件名
	ConfigFileName   string `json:"config_file_name"`   //项目的配置：文件名
	EtcdUrl          string `json:"etcd_url"`           //etcd get url
	Debug            int    `json:"debug"`              //debug 模式
	TestFlag         string `json:"test_flag"`          //是否为测试状态
}

var V = New() //动态的容器
var C Config  //静态从配置文件中读取的
var MainEnv MainEnvironment
var MainCmdParameter CmdParameter

func New() *Global {
	global := new(Global)
	return global
}

func AutoCreateUpDbTable() map[string]string {
	mydb := util.NewDbTool(V.Gorm)
	sql := mydb.CreateTable(&model.User{}, &model.UserReg{}, &model.UserLogin{},
		&model.AgoraCloudRecord{}, &model.AgoraCallbackRecord{}, &model.TwinAgoraRoom{},
		&model.GameMatchRule{},
		&model.OperationRecord{}, &model.Project{}, &model.StatisticsLog{},
		&model.CicdPublish{}, &model.Server{}, &model.Instance{},
		&model.SmsRule{}, &model.SmsLog{}, &model.EmailRule{}, &model.EmailLog{}, &model.MailRule{}, &model.MailLog{}, &model.MailGroup{})

	return sql
}

func GetUtilUploadConst() map[string]int {
	return V.VideoManager.GetConstListFileUploadType()
}
