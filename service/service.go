package service

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"zgoframe/util"
)

type Service struct {
	User      *User
	SendSms   *SendSms
	SendEmail *SendEmail
}

func NewService(gorm *gorm.DB, zap *zap.Logger, myEmail *util.MyEmail, myRedis *util.MyRedis) *Service {
	service := new(Service)
	service.User = NewUser(gorm, myRedis)
	service.SendSms = NewSendSms(gorm)
	service.SendEmail = NewSendEmail(gorm, myEmail)

	return service
}
