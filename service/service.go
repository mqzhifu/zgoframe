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
	Gateway    *Gateway
}

func NewService(gorm *gorm.DB, zap *zap.Logger, myEmail *util.MyEmail, myRedis *util.MyRedis, netWayOption util.NetWayOption, grpcManager *util.GrpcManager) *Service {
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

	var err error
	var netway *util.NetWay
	service.RoomManage.SetFrameSync(service.FrameSync)
	service.Gateway, netway, err = CreateGateway(netWayOption, grpcManager, zap)
	if err != nil {
		util.ExitPrint("InitGateway err:" + err.Error())
	}

	service.FrameSync.SetNetway(netway)

	return service
}

var GateDefaultProtocol = int32(util.PROTOCOL_WEBSOCKET)
var GateDefaultContentType = int32(util.CONTENT_TYPE_JSON)

func CreateGateway(netWayOption util.NetWayOption, grpcManager *util.GrpcManager, zap *zap.Logger) (*Gateway, *util.NetWay, error) {
	gateway := NewGateway(grpcManager, zap)
	netway, err := gateway.StartSocket(netWayOption)
	return gateway, netway, err
}
