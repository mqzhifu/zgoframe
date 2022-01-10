package global

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"zgoframe/util"
)

type Global struct {
	//App        util.App
	//AppMng		*util.AppManager
	RootDir 		string
	Vip        		*viper.Viper
	Zap        		*zap.Logger
	Redis      		*util.MyRedis
	Gin        		*gin.Engine
	Gorm       		*gorm.DB
	Service			util.Service

	Project 		util.Project
	ProjectMng  	*util.ProjectManager
	Etcd       		*util.MyEtcd
	HttpServer 		*http.Server
	Metric 			*util.MyMetrics
	GrpcManager		*util.GrpcManager
	AlertPush		*util.AlertPush	//报警推送： prometheus
	AlertHook 		*util.AlertHook	//报警：邮件 手机
	Websocket  		*util.Websocket
	ConnMng 		*util.ConnManager
	RecoverGo		*util.RecoverGo
	ProtobufMap 	*util.ProtobufMap
	Process 		*util.Process
	Err 			*util.ErrMsg
	Email 			*util.MyEmail
	ServiceManager  *util.ServiceManager		//管理已注册的服务
	ServiceDiscovery	*util.ServiceDiscovery	//管理服务发现，会用到上面的ServiceManager

	Gateway			*util.Gateway
	//ConnProtocol *util.ConnProtocol
}

func New()*Global {
	global  := new(Global)
	return global
}

var V = New()
var C Config

const (
	DEFAULT_CONFIT_TYPE  = "toml"
	DEFAULT_CONFIG_FILE_NAME = "config"
	DEFAULT_CONFIG_SOURCE_TYPE = "file"

	CONFIG_STATUS_OPEN = "open"
)
