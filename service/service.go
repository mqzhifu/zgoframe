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
}

func NewService(gorm *gorm.DB, zap *zap.Logger, myEmail *util.MyEmail, myRedis *util.MyRedis, netWayOption util.NetWayOption, grpcManager *util.GrpcManager, projectManager *util.ProjectManager, dataDir string) *Service {
	service := new(Service)
	service.User = NewUser(gorm, myRedis)
	service.SendSms = NewSendSms(gorm)
	service.SendEmail = NewSendEmail(gorm, myEmail)

	configCenterOption := ConfigCenterOption{
		envList:            util.GetEnvList(),
		Gorm:               gorm,
		Redis:              myRedis,
		ProjectManager:     projectManager,
		PersistenceType:    PERSITENCE_TYPE_FILE,
		PersistenceFileDir: dataDir + "/data/config/",
	}
	var err error
	service.ConfigCenter, err = NewConfigCenter(configCenterOption)
	if err != nil {
		util.ExitPrint("NewConfigCenter err:" + err.Error())
	}

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
	service.GameMatch = NewMatch(matchOption)
	syncOption := FrameSyncOption{
		Log:        zap,
		RoomManage: service.RoomManage,
	}
	//强-依赖room
	service.FrameSync = NewFrameSync(syncOption)

	//var err error
	var netway *util.NetWay
	service.RoomManage.SetFrameSync(service.FrameSync)
	//service.Gateway, netway, err = CreateGateway(netWayOption, grpcManager, zap)

	gateway := NewGateway(grpcManager, zap)
	netway, err = gateway.StartSocket(netWayOption)
	if err != nil {
		util.ExitPrint("InitGateway err:" + err.Error())
	}

	service.FrameSync.SetNetway(netway)
	gateway.MyService = service
	return service
}

var GateDefaultProtocol = int32(util.PROTOCOL_WEBSOCKET)
var GateDefaultContentType = int32(util.CONTENT_TYPE_JSON)
