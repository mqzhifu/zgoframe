package service

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"zgoframe/util"
)

type Service struct {
	User       *User
	SendSms    *SendSms
	SendEmail  *SendEmail
	RoomManage *RoomManager
	FrameSync  *FrameSync
}

func NewService(gorm *gorm.DB, zap *zap.Logger, myEmail *util.MyEmail, myRedis *util.MyRedis) *Service {
	service := new(Service)
	service.User = NewUser(gorm, myRedis)
	service.SendSms = NewSendSms(gorm)
	service.SendEmail = NewSendEmail(gorm, myEmail)

	//room要先实例化,math frame_sync 都强依赖room
	roomManagerOption := RoomManagerOption{
		Log:          zap,
		ReadyTimeout: 60,
		RoomPeople:   4,
	}
	service.RoomManage = NewRoomManager(roomManagerOption)
	matchOption := MatchOption{
		Log:         zap,
		RoomManager: service.RoomManage,
		//MatchSuccessChan chan *Room
	}
	NewMatch(matchOption)
	syncOption := FrameSyncOption{
		Log:        zap,
		RoomManage: service.RoomManage,
	}
	//强-依赖room
	service.FrameSync = NewFrameSync(syncOption)

	service.RoomManage.SetFrameSync(service.FrameSync)

	return service
}
