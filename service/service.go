//微服务  - 具体的业务
package service

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"zgoframe/util"
)

type Service struct {
	User           *User
	SendSms        *SendSms
	SendEmail      *SendEmail
	RoomManage     *RoomManager
	FrameSync      *FrameSync
	Gateway        *Gateway
	GameMatch      *Match
	ConfigCenter   *ConfigCenter
	ProjectManager *util.ProjectManager
	Mail 		*Mail
}

type MyServiceOptions struct {
	Gorm *gorm.DB
	Zap *zap.Logger
	MyEmail *util.MyEmail
	MyRedis *util.MyRedis
	NetWayOption util.NetWayOption
	GrpcManager *util.GrpcManager
	ProjectManager *util.ProjectManager
	ConfigCenterDataDir string
	ConfigCenterPersistenceType	int
}

func NewService(options MyServiceOptions) *Service {
	service := new(Service)
	service.User = NewUser(options.Gorm, options.MyRedis)
	service.Mail = NewMail(options.Gorm,options.Zap)
	service.SendSms = NewSendSms(options.Gorm)
	service.SendEmail = NewSendEmail(options.Gorm, options.MyEmail)

	configCenterOption := ConfigCenterOption{
		envList:            util.GetConstListEnv(),
		Gorm:               options.Gorm,
		Redis:              options.MyRedis,
		ProjectManager:     options.ProjectManager,
		PersistenceType:    PERSISTENCE_TYPE_FILE,
		PersistenceFileDir: options.ConfigCenterDataDir,
	}
	var err error
	service.ConfigCenter, err = NewConfigCenter(configCenterOption)
	if err != nil {
		util.ExitPrint("NewConfigCenter err:" + err.Error())
	}

	//room要先实例化,math frame_sync 都强依赖room
	roomManagerOption := RoomManagerOption{
		Log:          options.Zap,
		ReadyTimeout: 60,
		RoomPeople:   4,
	}
	service.RoomManage = NewRoomManager(roomManagerOption)
	matchOption := MatchOption{
		Log:         options.Zap,
		RoomManager: service.RoomManage,
		//MatchSuccessChan chan *Room
	}
	service.GameMatch = NewMatch(matchOption)
	syncOption := FrameSyncOption{
		Log:        options.Zap,
		RoomManage: service.RoomManage,
	}
	//强-依赖room
	service.FrameSync = NewFrameSync(syncOption)

	//var err error
	var netway *util.NetWay
	service.RoomManage.SetFrameSync(service.FrameSync)
	//service.Gateway, netway, err = CreateGateway(netWayOption, grpcManager, zap)

	//util.ExitPrint(netWayOption)
	gateway := NewGateway(options.GrpcManager, options.Zap)
	netway, err = gateway.StartSocket(options.NetWayOption)
	if err != nil {
		util.ExitPrint("InitGateway err:" + err.Error())
	}

	service.FrameSync.SetNetway(netway)
	gateway.MyService = service
	return service
}

var GateDefaultProtocol = int32(util.PROTOCOL_WEBSOCKET)
var GateDefaultContentType = int32(util.CONTENT_TYPE_JSON)
