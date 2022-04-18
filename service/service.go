//微服务  - 具体的业务
package service

import (
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
	GameMatch      *Match
	ConfigCenter   *ConfigCenter
	Cicd 		*cicd.CicdManager
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
	OpDirName 	string
	ServiceList map[int]util.Service
}

func NewService(options MyServiceOptions) *Service {
	service := new(Service)
	service.User = NewUser(options.Gorm, options.MyRedis)
	service.Mail = NewMail(options.Gorm,options.Zap)
	service.Sms = NewSms(options.Gorm)
	service.Email = NewEmail(options.Gorm, options.MyEmail)

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

	service.Cicd ,err = InitCicd(options.Gorm,options.Zap,options.OpDirName,options.ServiceList)

	service.FrameSync.SetNetway(netway)
	gateway.MyService = service
	return service
}

func InitCicd(gorm *gorm.DB,zap *zap.Logger,opDir string,ServiceList map[int]util.Service)(*cicd.CicdManager,error){

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
		HttpPort		: cicdConfig.System.HttpPort,
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
