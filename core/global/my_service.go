package global

import (
	"zgoframe/service"
	"zgoframe/service/cicd"
	"zgoframe/service/config_center"
	"zgoframe/service/frame_sync"
	gamematch "zgoframe/service/game_match"
	"zgoframe/service/gateway"
	"zgoframe/service/msg_center"
	"zgoframe/service/seed_business"
	"zgoframe/service/user_center"
	"zgoframe/util"
)

type MyService struct {
	User                  *user_center.User           //用户中心
	Sms                   *msg_center.Sms             //短信服务
	Email                 *msg_center.Email           //电子邮件服务
	RoomManage            *frame_sync.RoomManager     //房间服务
	FrameSync             *frame_sync.FrameSync       //帧同步服务
	Match                 *gamematch.GameMatch        //匹配服务
	Gateway               *gateway.Gateway            //网关服务
	TwinAgora             *seed_business.TwinAgora    //广播120远程专家指导
	ConfigCenter          *config_center.ConfigCenter //配置中心
	Cicd                  *cicd.CicdManager           //自动部署
	Mail                  *msg_center.Mail            //站内信
	GameMatch             *gamematch.GameMatch
	RequestServiceAdapter *service.RequestServiceAdapter //请求3方服务 适配器
}

var GateDefaultProtocol = int32(util.PROTOCOL_WEBSOCKET)
var GateDefaultContentType = int32(util.CONTENT_TYPE_JSON)

//内部服务，按说：一个项目里最多也就1-2个服务，其它的服务应该在其它项目，并且访问的时候通过HTTP/TCP，这里方便使用，先把其它服务当成一个类使用
func NewMyService() *MyService {
	var err error
	myService := new(MyService)
	//创建一个 请求3方服务 的适配器，服务之间的请求/调用
	myService.RequestServiceAdapter = service.NewRequestServiceAdapter(V.ServiceDiscovery, V.GrpcManager, service.REQ_SERVICE_METHOD_INNER, C.System.ProjectId, V.Zap)
	//用户服务
	myService.User = user_center.NewUser(V.Gorm, V.Redis)
	//站内信服务
	myService.Mail = msg_center.NewMail(V.Gorm, V.Zap)
	//短信服务
	myService.Sms = msg_center.NewSms(V.Gorm)
	//电子邮件服务
	myService.Email = msg_center.NewEmail(V.Gorm, V.Email)
	//配置中心服务
	configCenterOption := config_center.ConfigCenterOption{
		EnvList:            util.GetConstListEnv(),
		Gorm:               V.Gorm,
		Redis:              V.Redis,
		ProjectManager:     V.ProjectMng,
		PersistenceType:    service.PERSISTENCE_TYPE_FILE,
		PersistenceFileDir: C.Http.StaticPath + "/" + C.ConfigCenter.DataPath,
		Log:                V.Zap,
	}
	myService.ConfigCenter, err = config_center.NewConfigCenter(configCenterOption)
	if err != nil {
		util.ExitPrint("NewConfigCenter err:" + err.Error())
	}
	//远程呼叫专家
	twinAgoraOption := seed_business.TwinAgoraOption{
		Log:                   V.Zap,
		Gorm:                  V.Gorm,
		StaticPath:            C.Http.StaticPath,
		RequestServiceAdapter: myService.RequestServiceAdapter,
	}
	myService.TwinAgora, err = seed_business.NewTwinAgora(twinAgoraOption)
	if err != nil {
		util.ExitPrint(err)
	}
	//长连接通信 - 配置
	netWayOption := util.NetWayOption{
		ListenIp:            C.Gateway.ListenIp,     //程序启动时监听的IP
		OutIp:               C.Gateway.OutIp,        //对外访问的IP
		OutDomain:           C.Gateway.OutDomain,    //对外的域名,WS在线上得用wss
		WsPort:              C.Gateway.WsPort,       //监听端口号
		TcpPort:             C.Gateway.TcpPort,      //监听端口号
		UdpPort:             C.Gateway.UdpPort,      //UDP端口号
		WsUri:               C.Gateway.WsUri,        //接HOST的后面的URL地址
		DefaultProtocolType: GateDefaultProtocol,    //兼容协议：ws tcp udp
		DefaultContentType:  GateDefaultContentType, //默认内容格式 ：json protobuf
		LoginAuthType:       "jwt",                  //登陆验证类型-jwt
		LoginAuthSecretKey:  C.Jwt.Key,              //登陆验证-key
		MaxClientConnNum:    1024,                   //客户端最大连接数
		MsgContentMax:       10240,                  //一条消息内容最大值
		IOTimeout:           3,                      //read write sock fd 超时时间
		ConnTimeout:         60,                     //一个FD超时时间
		ClientHeartbeatTime: 3,                      //客户端心跳时间(秒)
		ServerHeartbeatTime: 5,                      //服务端心跳时间(秒)
		ProtoMap:            V.ProtoMap,             //protobuf 映射表
		GrpcManager:         V.GrpcManager,
		Log:                 V.Zap,
	}
	//网关
	if C.Gateway.Status == "open" {
		gateway := gateway.NewGateway(V.GrpcManager, V.Zap, myService.RequestServiceAdapter)
		gateway.MyServiceList.FrameSync = myService.FrameSync
		gateway.MyServiceList.Match = myService.Match
		gateway.MyServiceList.RoomManage = myService.RoomManage
		gateway.MyServiceList.TwinAgora = myService.TwinAgora

		myService.Gateway = gateway

		_, err := gateway.StartSocket(netWayOption)
		if err != nil {
			util.ExitPrint("InitGateway err:" + err.Error())
		}
	}
	CreateGame(myService)

	myService.Cicd, err = InitCicd()

	//myService.RegisterService()

	return myService
}

//这里测试一下，服务注册到ETCD
func (myService *MyService) RegisterService() {
	if C.ServiceDiscovery.Status == "open" && C.Etcd.Status == "open" {
		var node util.ServiceNode
		ip := "127.0.0.1"
		listenIp := "127.0.0.1"
		port := "1111"

		node = util.ServiceNode{
			//ServiceId: global.C.System.ProjectId,
			ServiceId:   1, //游戏匹配服务
			ServiceName: "GameMatch",
			Ip:          ip,
			ListenIp:    listenIp,
			Port:        port,
			Protocol:    util.SERVICE_PROTOCOL_HTTP,
			IsSelfReg:   true,
		}
		V.ServiceDiscovery.Register(node)

		node = util.ServiceNode{
			//ServiceId: global.C.System.ProjectId,
			ServiceId:   6, //游戏匹配服务
			ServiceName: "Zgoframe",
			Ip:          ip,
			ListenIp:    listenIp,
			Port:        port,
			Protocol:    util.SERVICE_PROTOCOL_HTTP,
			IsSelfReg:   true,
		}

		V.ServiceDiscovery.Register(node)
		node = util.ServiceNode{
			//ServiceId: global.C.System.ProjectId,
			ServiceId:   2, //游戏匹配服务
			ServiceName: "FrameSync",
			Ip:          ip,
			ListenIp:    listenIp,
			Port:        port,
			Protocol:    util.SERVICE_PROTOCOL_HTTP,
			IsSelfReg:   true,
		}
		V.ServiceDiscovery.Register(node)

	}
}

//游戏类的服务,一个游戏至少得有：房间、匹配、帧同步
func CreateGame(myService *MyService) (err error) {
	//帧同步 - 房间服务 - room要先实例化,math frame_sync 都强依赖room
	roomManagerOption := frame_sync.RoomManagerOption{
		Log:                   V.Zap,
		ReadyTimeout:          60,
		RoomPeople:            2,
		RequestServiceAdapter: myService.RequestServiceAdapter,
	}
	myService.RoomManage = frame_sync.NewRoomManager(roomManagerOption)
	//匹配服务 , 依赖 RoomManage
	//matchOption := service.MatchOption{
	//	RequestServiceAdapter: myService.RequestServiceAdapter,
	//	Log:                   V.Zap,
	//	RoomManager:           myService.RoomManage,
	//	//MatchSuccessChan chan *Room
	//}
	//匹配服务，这个是假的，或者说简易版本，用于快速测试
	//myService.Match = service.NewMatch(matchOption)

	//这个是真的匹配服务
	gmOp := gamematch.GameMatchOption{
		RequestServiceAdapter:  myService.RequestServiceAdapter,
		Log:                    V.Zap,
		Redis:                  V.RedisGo,
		Gorm:                   V.Gorm,
		Metrics:                V.Metric,
		ServiceDiscovery:       V.ServiceDiscovery,
		RuleDataSourceType:     service.GAME_MATCH_DATA_SOURCE_TYPE_DB,
		StaticPath:             C.Http.StaticPath,
		RedisPrefix:            "gm",
		RedisKeySeparator:      "_",
		RedisTextSeparator:     "#",
		FrameSyncRoom:          myService.RoomManage,
		RedisIdSeparator:       ",",
		RedisPayloadSeparation: "%",
	}
	myService.GameMatch, err = gamematch.NewGameMatch(gmOp)
	if err != nil {
		util.ExitPrint("NewGameMatch err:", err)
	}

	//user -> sign ->Match -> Room -> Rsync
	//帧同步服务 - 强-依赖room
	syncOption := frame_sync.FrameSyncOption{
		RequestServiceAdapter: myService.RequestServiceAdapter,
		ProjectId:             C.System.ProjectId,
		Log:                   V.Zap,
		RoomManage:            myService.RoomManage,
		FPS:                   10,
		MapSize:               5,
	}
	myService.FrameSync = frame_sync.NewFrameSync(syncOption)
	myService.RoomManage.SetFrameSync(myService.FrameSync)
	//go myService.Match.Start()
	return nil
}

func InitCicd() (*cicd.CicdManager, error) {
	//util.MyPrint(ServiceList)
	/*依赖
	host.toml cicd.sh
	table:  project instance server cicd_publish
	*/

	opDirFull := MainEnv.RootDir + "/" + C.System.OpDirName
	cicdConfig := cicd.ConfigCicd{}
	cicdConfig.System.RootDir = MainEnv.RootDir
	//运维：服务器的配置信息
	configFile := opDirFull + "/host" + "." + "toml"

	//读取配置文件中的内容
	err := util.ReadConfFile(configFile, &cicdConfig)
	if err != nil {
		util.ExitPrint(err.Error())
	}
	cicdConfig.SuperVisor.ConfTemplateFileName = cicdConfig.SuperVisor.ConfTemplateFile
	cicdConfig.SuperVisor.ConfTemplateFile = MainEnv.RootDir + "/" + C.System.OpDirName + "/" + cicdConfig.SuperVisor.ConfTemplateFile

	V.Zap.Debug("InitCicd HostConfigFile:" + configFile + " ConfTemplateFile:" + cicdConfig.SuperVisor.ConfTemplateFile)
	//util.PrintStruct(cicdConfig , " : ")

	//3方实例
	instanceManager, _ := util.NewInstanceManager(V.Gorm)
	//服务器列表
	serverManger, _ := util.NewServerManger(V.Gorm)
	serverList := serverManger.Pool
	//发布管理
	publicManager := cicd.NewCICDPublicManager(V.Gorm)

	//util.ExitPrint(22)
	op := cicd.CicdManagerOption{
		HttpPort:         C.Http.Port,
		ServerList:       serverList,
		Config:           cicdConfig,
		InstanceManager:  instanceManager,
		PublicManager:    publicManager,
		Log:              V.Zap,
		OpDirName:        C.System.OpDirName,
		ServiceList:      V.ServiceManager.Pool,
		ProjectList:      V.ProjectMng.Pool,
		UploadDiskPath:   C.Http.StaticPath + "/" + C.FileManager.UploadPath,
		DownloadDiskPath: C.Http.StaticPath + "/" + C.FileManager.DownloadPath,
	}

	cicd, err := cicd.NewCicdManager(op)
	if err != nil {
		util.ExitPrint(err)
	}
	//生成 filebeat 配置文件
	//cicd.GenerateAllFilebeat()
	//cicd.GetSuperVisorList()
	//部署所有机器上的所有服务项目
	//cicd.DeployAllService()
	//go cicd.StartHttp(global.C.Http.StaticPath)

	return cicd, err

}
