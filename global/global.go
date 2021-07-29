package global

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Global struct {
	Vip 	*viper.Viper
	Zap  	*zap.Logger
	Redis 	*redis.Client
	Gin		*gin.Engine
	Gorm 	*gorm.DB
}

func New()*Global{
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
