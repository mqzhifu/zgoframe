// 全局容器，1 配置信息中的变量 2 公共初始好的类包
package global

import (
	"context"
	"embed"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"os/user"
	"zgoframe/model"
	"zgoframe/util"
)

type Base struct {
	Vip            *viper.Viper
	Zap            *zap.Logger
	HttpZap        *zap.Logger
	Redis          *util.MyRedis
	RedisGo        *util.MyRedisGo
	Gin            *gin.Engine
	Gorm           *gorm.DB   // 多数据库模式下，有一个库肯定会被经常访问，这里加一个快捷链接
	GormList       []*gorm.DB // 所有数据库，连接成功后的列表
	HttpServer     *http.Server
	StaticFileSys  embed.FS // embed.FS 静态文件
	ES8TypedClient *elasticsearch.TypedClient
	ES8Client      *elasticsearch.Client
}

//type Service struct {
//	MyService        *MyService             // 内部快捷服务
//}

type Util struct {
	Project          model.Project
	ProjectMng       *util.ProjectManager
	Etcd             *util.MyEtcd
	Metric           *util.MyMetrics
	GrpcManager      *util.GrpcManager
	AlertPush        *util.AlertPush // 报警推送： prometheus
	Websocket        *util.Websocket
	ConnMng          *util.ConnManager
	RecoverGo        *util.RecoverGo
	ProtoMap         *util.ProtoMap
	Process          *util.Process
	Err              *util.ErrMsg
	Email            *util.MyEmail
	ImageSlice       *util.ImageSlice
	DocsManager      *util.FileManager
	PackagesManager  *util.FileManager
	ImgManager       *util.FileManager
	VideoManager     *util.FileManager
	NetWay           *util.NetWay
	AliSms           *util.AliSms
	ServiceManager   *util.ServiceManager   // 管理已注册的服务
	ServiceDiscovery *util.ServiceDiscovery // 管理服务发现，会用到上面的ServiceManager
	AliOss           *util.AliOss           // 阿里网盘
	StaticFileSystem *util.StaticFileSystem // 兼容，管理静态文件读取
	// AlertHook        *util.AlertHook //报警：邮件 手机
}

// 所有容器挂 在这个上面
type Container struct {
	Base *Base
	Util *Util
	//Service *MyService
}

var V = NewContainer() // 动态的容器
var B Base
var C Config // 静态从配置文件中读取的
var MainEnv MainEnvironment
var MainCmdParameter CmdParameter

// main主协程的一些参数-环境参数
type MainEnvironment struct {
	RootDir         string             `json:"root_dir"` //main.go 文件路径
	RootDirName     string             `json:"root_dir_name"`
	GoVersion       string             `json:"go_version"` // 当前go版本
	ExecUser        *user.User         `json:"-"`          // 执行该脚本的用户信息
	Cpu             string             `json:"cpu"`        // cpu信息
	RootCtx         context.Context    `json:"-"`          // main的上下文，级别最高
	RootCancelFunc  context.CancelFunc `json:"-"`          // main的取消函数，该管理如果能读出值，main会主动退出
	RootQuitFunc    func(source int)   `json:"-"`          // 这是个函数，子级可直接驱动：退出MAIN
	BuildTime       string             // 编译时：时间
	BuildGitVersion string             // 编译时：git版本号
}

// 指令行 收集的参数
type CmdParameter struct {
	Env              int    `json:"env"`                // 当前环境
	ConfigSourceType string `json:"config_source_type"` // 文件 | etcd
	ConfigFileType   string `json:"config_file_type"`   // 项目的配置：文件名
	ConfigFileName   string `json:"config_file_name"`   // 项目的配置：文件名
	EtcdUrl          string `json:"etcd_url"`           // etcd get url
	Debug            int    `json:"debug"`              // debug 模式
	TestFlag         string `json:"test_flag"`          // 是否为测试状态
	BuildStatic      string `json:"build_static"`       // 编译时：把静态文件一并打包进二进制文件中，牵扯到：读文件时，是从编译包里读还是从硬盘文件中读。默认为关闭，启动进程时 指令行操作
}

func NewContainer() *Container {
	container := new(Container)
	container.Base = new(Base)
	container.Util = new(Util)
	return container
}

func AutoCreateUpDbTable() map[string]string {
	dbTool := util.NewDbTool(V.Base.Gorm)
	sql := dbTool.CreateTable(&model.User{}, &model.UserReg{}, &model.UserLogin{},
		&model.GameMatchRule{}, &model.GameMatchSuccess{}, &model.GameMatchGroup{}, &model.GameMatchPush{}, &model.GameSyncRoom{},
		&model.OperationRecord{}, &model.Project{}, &model.ProjectPushMsg{}, &model.StatisticsLog{},
		&model.CicdPublish{}, &model.Server{}, &model.Instance{}, &model.ConnRecord{},
		&model.SmsRule{}, &model.SmsLog{}, &model.EmailRule{}, &model.EmailLog{}, &model.MailRule{}, &model.MailLog{}, &model.MailGroup{},
		&model.Goods{}, &model.Orders{}, &model.PayOrder{}, &model.PayOrderMatch{}, &model.UserTotal{})
	//&model.AgoraCloudRecord{}, &model.AgoraCallbackRecord{}, &model.TwinAgoraRoom{},

	return sql
}

func GetUtilUploadConst() map[string]int {
	return V.Util.VideoManager.GetConstListFileUploadType()
}
