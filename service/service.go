//微服务  - 具体的业务
package service

import (
	"zgoframe/service/gamematch"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"zgoframe/service/cicd"
	"zgoframe/util"
)

type Service struct {
	User           *User
	Sms        *Sms
	Email      *Email
	RoomManage     *RoomManager
	FrameSync      *FrameSync
	Gateway        *Gateway
	Match      *Match
	ConfigCenter   *ConfigCenter
	Cicd 		*cicd.CicdManager
	ProjectManager *util.ProjectManager
	Mail 		*Mail
	GameMatch 	*gamematch.Gamematch
}

type MyServiceOptions struct {
	Gorm *gorm.DB
	Zap *zap.Logger
	MyEmail *util.MyEmail
	MyRedis *util.MyRedis
	MyRedisGo *util.MyRedisGo
	NetWayOption util.NetWayOption
	GrpcManager *util.GrpcManager
	ProjectManager *util.ProjectManager
	ConfigCenterDataDir string
	ConfigCenterPersistenceType	int
	OpDirName 	string
	ServiceList map[int]util.Service
	HttpPort 	string
	GatewayStatus string
	Etcd *util.MyEtcd
	Metrics *util.MyMetrics
	ProjectId int
	ServiceDiscovery *util.ServiceDiscovery
}

func NewService(options MyServiceOptions) *Service {
	service := new(Service)
	//用户服务
	service.User = NewUser(options.Gorm, options.MyRedis)
	//站内信服务
	service.Mail = NewMail(options.Gorm,options.Zap)
	//短信服务
	service.Sms = NewSms(options.Gorm)
	//电子邮件服务
	service.Email = NewEmail(options.Gorm, options.MyEmail)
	////配置中心服务
	//configCenterOption := ConfigCenterOption{
	//	envList:            util.GetConstListEnv(),
	//	Gorm:               options.Gorm,
	//	Redis:              options.MyRedis,
	//	ProjectManager:     options.ProjectManager,
	//	PersistenceType:    PERSISTENCE_TYPE_FILE,
	//	PersistenceFileDir: options.ConfigCenterDataDir,
	//}
	//var err error
	//service.ConfigCenter, err = NewConfigCenter(configCenterOption)
	//if err != nil {
	//	util.ExitPrint("NewConfigCenter err:" + err.Error())
	//}
	//
	////房间服务 - room要先实例化,math frame_sync 都强依赖room
	//roomManagerOption := RoomManagerOption{
	//	Log:          options.Zap,
	//	ReadyTimeout: 60,
	//	RoomPeople:   4,
	//}
	//service.RoomManage = NewRoomManager(roomManagerOption)
	////匹配服务
	//matchOption := MatchOption{
	//	Log:         options.Zap,
	//	RoomManager: service.RoomManage,
	//	//MatchSuccessChan chan *Room
	//}
	//service.Match = NewMatch(matchOption)
	//syncOption := FrameSyncOption{
	//	Log:        options.Zap,
	//	RoomManage: service.RoomManage,
	//}
	////帧同步服务 - 强-依赖room
	//service.FrameSync = NewFrameSync(syncOption)
	//
	//service.RoomManage.SetFrameSync(service.FrameSync)
	////他们3个的关系：
	////room -> match , match -> room ,room ->frame sync ， frame sync -> room
	//
	//
	//if options.GatewayStatus == "open"{
	//	gateway := NewGateway(options.GrpcManager, options.Zap)
	//	var netway *util.NetWay
	//	netway, err = gateway.StartSocket(options.NetWayOption)
	//	if err != nil {
	//		util.ExitPrint("InitGateway err:" + err.Error())
	//	}
	//
	//	service.FrameSync.SetNetway(netway)
	//	gateway.MyService = service
	//}
	//
	//service.Cicd ,err = InitCicd(options.Gorm,options.Zap,options.OpDirName,options.ServiceList,options.HttpPort)
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

func InitCicd(gorm *gorm.DB,zap *zap.Logger,opDir string,ServiceList map[int]util.Service,httpPort string)(*cicd.CicdManager,error){

	/*依赖
	host.toml cicd.sh
	table:  project instance server cicd_publish
	*/

	opDirName := opDir
	pwd , _ := os.Getwd()//当前路径]
	opDirFull := pwd + "/" + opDirName
	util.MyPrint(opDirFull,opDirName)

	cicdConfig := cicd.ConfigCicd{}
	//运维：服务器的配置信息
	configFile := opDirFull + "/host" + "." +"toml"

	//读取配置文件中的内容
	err := util.ReadConfFile(configFile,&cicdConfig)
	if err != nil{
		util.ExitPrint(err.Error())
	}

	cicdConfig.SuperVisor.ConfTemplateFile = pwd + "/" + opDir  +  "/" + cicdConfig.SuperVisor.ConfTemplateFile

	util.PrintStruct(cicdConfig , " : ")
	//3方实例
	instanceManager ,_:= util.NewInstanceManager(gorm )
	//服务器列表
	serverManger,_ := util.NewServerManger( gorm )
	serverList := serverManger.Pool
	//发布管理
	publicManager := cicd.NewCICDPublicManager(gorm )

	//util.ExitPrint(22)
	op := cicd.CicdManagerOption{
		HttpPort		: httpPort,
		ServerList 		: serverList,
		Config			: cicdConfig,
		ServiceList		: ServiceList,
		InstanceManager : instanceManager,
		PublicManager 	: publicManager,
		Log				: zap,
		OpDirName		: opDirName,
	}

	cicd ,err := cicd.NewCicdManager(op)
	if err != nil{
		util.ExitPrint(err)
	}
	return cicd,err
	//生成 filebeat 配置文件
	//cicd.GenerateAllFilebeat()
	//cicd.GetSuperVisorList()
	//部署所有机器上的所有服务项目
	//cicd.DeployAllService()
	//go cicd.StartHttp(global.C.Http.StaticPath)
}

var GateDefaultProtocol = int32(util.PROTOCOL_WEBSOCKET)
var GateDefaultContentType = int32(util.CONTENT_TYPE_JSON)
