package global

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"time"
	"zgoframe/util"
)

type Global struct {
	Vip        *viper.Viper
	Zap        *zap.Logger
	Redis      *redis.Client
	Gin        *gin.Engine
	Gorm       *gorm.DB
	App        util.App
	Etcd       *util.MyEtcd
	HttpServer *http.Server
	Service    *util.Service
	Metric 		*util.MyMetrics
	Grpc 		*util.MyGrpc
	Alert 		*util.Alert
}

func New()*Global {
	global  := new(Global)
	return global
}

var V = New()
var C Config

type MODEL struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

const (
	DEFAULT_CONFIT_TYPE  = "toml"
	DEFAULT_CONFIG_FILE_NAME = "config"

	CONFIG_STATUS_OPEN = "open"
)
