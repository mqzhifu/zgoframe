//全局容器，1 配置信息中的变量 2 公共初始好的类包
package global

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"zgoframe/model"
	"zgoframe/util"
	"zgoframe/service"
)

type Global struct {
	//App        util.App
	//AppMng		*util.AppManager
	RootDir          string
	Vip              *viper.Viper
	Zap              *zap.Logger
	Redis            *util.MyRedis
	Gin              *gin.Engine
	Gorm             *gorm.DB		//多数据库模式下，有一个库肯定会被经常访问，这里加一个快捷链接
	GormList 		 []*gorm.DB		//所有数据库，连接成功后的列表
	Project          util.Project
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

	MyService        *service.Service//内部快捷服务

	//Service          util.Service
	ServiceManager   *util.ServiceManager   //管理已注册的服务
	ServiceDiscovery *util.ServiceDiscovery //管理服务发现，会用到上面的ServiceManager

	//ConnProtocol *util.ConnProtocol
}

func New() *Global {
	global := new(Global)
	return global
}

var V = New()
var C Config

const (
	DEFAULT_CONFIT_TYPE        = "toml"
	DEFAULT_CONFIG_FILE_NAME   = "config"
	DEFAULT_CONFIG_SOURCE_TYPE = "file"

	CONFIG_STATUS_OPEN = "open"
)

func AutoCreateUpDbTable()map[string]string {
	mydb := util.NewDbTool(V.Gorm)
	sql := mydb.CreateTable(&model.User{}, &model.UserReg{}, &model.UserLogin{},
		&model.OperationRecord{}, &model.Project{},&model.StatisticsLog{},
		&model.CicdPublish{}, &model.Server{}, &model.Instance{},
		&model.SmsRule{}, &model.SmsLog{}, &model.EmailRule{}, &model.EmailLog{}, &model.MailRule{}, &model.MailLog{}, &model.MailGroup{})

	return sql
	//util.ExitPrint("init done.")
}