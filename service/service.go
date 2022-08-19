//微服务  - 具体的业务
package service

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"zgoframe/service/cicd"
	"zgoframe/service/gamematch"
	"zgoframe/util"
)

type Service struct {
	User           *User
	Sms            *Sms
	Email          *Email
	RoomManage     *RoomManager
	FrameSync      *FrameSync
	Gateway        *Gateway
	Match          *Match
	ConfigCenter   *ConfigCenter
	Cicd           *cicd.CicdManager
	ProjectManager *util.ProjectManager
	Mail           *Mail
	GameMatch      *gamematch.Gamematch
}

type MyServiceOptions struct {
	Gorm                        *gorm.DB
	Zap                         *zap.Logger
	MyEmail                     *util.MyEmail
	MyRedis                     *util.MyRedis
	MyRedisGo                   *util.MyRedisGo
	NetWayOption                util.NetWayOption
	GrpcManager                 *util.GrpcManager
	ProjectManager              *util.ProjectManager
	ConfigCenterDataDir         string
	ConfigCenterPersistenceType int
	OpDirName                   string
	UploadDiskPath              string
	DownloadDiskPath            string
	ServiceList                 map[int]util.Service
	HttpPort                    string
	GatewayStatus               string
	Etcd                        *util.MyEtcd
	Metrics                     *util.MyMetrics
	ProjectId                   int
	ServiceDiscovery            *util.ServiceDiscovery
	RootDir                     string
}

func NewService(options MyServiceOptions) *Service {
	service := new(Service)
	//用户服务
	service.User = NewUser(options.Gorm, options.MyRedis)
	//站内信服务
	service.Mail = NewMail(options.Gorm, options.Zap)
	//短信服务
	service.Sms = NewSms(options.Gorm)
	//电子邮件服务
	service.Email = NewEmail(options.Gorm, options.MyEmail)
	//配置中心服务
	configCenterOption := ConfigCenterOption{
		envList:            util.GetConstListEnv(),
		Gorm:               options.Gorm,
		Redis:              options.MyRedis,
		ProjectManager:     options.ProjectManager,
		PersistenceType:    PERSISTENCE_TYPE_FILE,
		PersistenceFileDir: options.ConfigCenterDataDir,
		Log:                options.Zap,
	}
	//配置中心 - 服务
	var err error
	service.ConfigCenter, err = NewConfigCenter(configCenterOption)
	if err != nil {
		util.ExitPrint("NewConfigCenter err:" + err.Error())
	}

	//房间服务 - room要先实例化,math frame_sync 都强依赖room
	roomManagerOption := RoomManagerOption{
		Log:          options.Zap,
		ReadyTimeout: 60,
		RoomPeople:   2,
	}
	service.RoomManage = NewRoomManager(roomManagerOption)
	//匹配服务 , 依赖 RoomManage
	matchOption := MatchOption{
		Log:         options.Zap,
		RoomManager: service.RoomManage,
		//MatchSuccessChan chan *Room
	}
	service.Match = NewMatch(matchOption)
	syncOption := FrameSyncOption{
		Log:        options.Zap,
		RoomManage: service.RoomManage,
		FPS:        options.NetWayOption.FPS,
		MapSize:    options.NetWayOption.MapSize,
	}
	go service.Match.Start()
	//user -> sign ->Match -> Room -> Rsync
	//帧同步服务 - 强-依赖room
	service.FrameSync = NewFrameSync(syncOption)

	service.RoomManage.SetFrameSync(service.FrameSync)
	//网关
	if options.GatewayStatus == "open" {
		gateway := NewGateway(options.GrpcManager, options.Zap)
		var netway *util.NetWay
		netway, err = gateway.StartSocket(options.NetWayOption)
		if err != nil {
			util.ExitPrint("InitGateway err:" + err.Error())
		}

		service.FrameSync.SetNetway(netway)
		gateway.MyService = service
		service.Gateway = gateway
	}
	//
	service.Cicd, err = InitCicd(options)
	//
	//
	//gameMatchOption :=  gamematch.GamematchOption{
	//	Log :options.Zap,
	//	Redis :options.MyRedisGo,
	//	ServiceDiscovery : options.ServiceDiscovery,
	//	Etcd : options.Etcd,
	//	Metrics: options.Metrics,
	//	ProjectId :options.ProjectId,
	//	//HttpdOption: myHttpdOption,
	//}
	//
	//myGamematch,errs := gamematch.NewGameMatch(gameMatchOption)
	//if errs != nil{
	//	util.ExitPrint("NewGamematch : ",errs.Error())
	//}
	//service.GameMatch = myGamematch

	return service
}

func InitCicd(option MyServiceOptions) (*cicd.CicdManager, error) {
	//util.MyPrint(ServiceList)
	/*依赖
	host.toml cicd.sh
	table:  project instance server cicd_publish
	*/

	opDirName := option.OpDirName
	opDirFull := option.RootDir + "/" + opDirName
	cicdConfig := cicd.ConfigCicd{}
	cicdConfig.System.RootDir = option.RootDir
	//运维：服务器的配置信息
	configFile := opDirFull + "/host" + "." + "toml"

	//读取配置文件中的内容
	err := util.ReadConfFile(configFile, &cicdConfig)
	if err != nil {
		util.ExitPrint(err.Error())
	}
	cicdConfig.SuperVisor.ConfTemplateFileName = cicdConfig.SuperVisor.ConfTemplateFile
	cicdConfig.SuperVisor.ConfTemplateFile = option.RootDir + "/" + opDirName + "/" + cicdConfig.SuperVisor.ConfTemplateFile

	option.Zap.Debug("InitCicd HostConfigFile:" + configFile + " ConfTemplateFile:" + cicdConfig.SuperVisor.ConfTemplateFile)
	//util.PrintStruct(cicdConfig , " : ")

	//3方实例
	instanceManager, _ := util.NewInstanceManager(option.Gorm)
	//服务器列表
	serverManger, _ := util.NewServerManger(option.Gorm)
	serverList := serverManger.Pool
	//发布管理
	publicManager := cicd.NewCICDPublicManager(option.Gorm)

	//util.ExitPrint(22)
	op := cicd.CicdManagerOption{
		HttpPort:         option.HttpPort,
		ServerList:       serverList,
		Config:           cicdConfig,
		ServiceList:      option.ServiceList,
		ProjectList:      option.ProjectManager.Pool,
		InstanceManager:  instanceManager,
		PublicManager:    publicManager,
		Log:              option.Zap,
		OpDirName:        opDirName,
		UploadDiskPath:   option.UploadDiskPath,
		DownloadDiskPath: option.DownloadDiskPath,
	}

	cicd, err := cicd.NewCicdManager(op)
	if err != nil {
		util.ExitPrint(err)
	}
	return cicd, err
	//生成 filebeat 配置文件
	//cicd.GenerateAllFilebeat()
	//cicd.GetSuperVisorList()
	//部署所有机器上的所有服务项目
	//cicd.DeployAllService()
	//go cicd.StartHttp(global.C.Http.StaticPath)
}

var GateDefaultProtocol = int32(util.PROTOCOL_WEBSOCKET)
var GateDefaultContentType = int32(util.CONTENT_TYPE_JSON)
