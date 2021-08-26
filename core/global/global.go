package global

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"zgoframe/util"
)

type Global struct {
	Vip        *viper.Viper
	Zap        *zap.Logger
	Redis      *redis.Client
	Gin        *gin.Engine
	Gorm       *gorm.DB
	App        util.App
	AppMng		*util.AppManager
	Etcd       *util.MyEtcd
	HttpServer *http.Server
	Service    *util.Service
	Metric 		*util.MyMetrics
	Grpc 		*util.MyGrpc
	AlertPush	*util.AlertPush	//报警推送： prometheus
	AlertHook 	*util.AlertHook	//报警：邮件 手机
	Websocket  *util.Websocket
	ConnMng 	*util.ConnManager
	ConnPRotocol *util.ConnProtocol
	RecoverGo	*util.RecoverGo
	ProtobufMap *util.ProtobufMap
	Process 	*util.Process
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
