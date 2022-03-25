package service

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	User      *User
	SendSms   *SendSms
	SendEmail *SendEmail
}

func NewService(gorm *gorm.DB, zap *zap.Logger) *Service {
	service := new(Service)
	service.User = NewUser(gorm)
	service.SendSms = NewSendSms(gorm)
	return service
}
