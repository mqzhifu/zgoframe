package global

import (
	"zgoframe/service"
	"zgoframe/service/cicd"
	"zgoframe/service/config_center"
	"zgoframe/service/frame_sync"
	gamematch "zgoframe/service/game_match"
	"zgoframe/service/gateway"
	"zgoframe/service/grab_order"
	"zgoframe/service/msg_center"
	"zgoframe/service/seed_business"
	"zgoframe/service/user_center"
	"zgoframe/util"
)

type MyService struct {
	User          *user_center.User // 用户中心
	Sms           *msg_center.Sms   // 短信服务
	Email         *msg_center.Email // 电子邮件服务
	AliSms        *util.AliSms
	Gateway       *gateway.Gateway            // 网关服务
	TwinAgora     *seed_business.TwinAgora    // 广州 120远程专家指导
	ConfigCenter  *config_center.ConfigCenter // 配置中心
	Cicd          *cicd.CicdManager           // 自动部署
	Mail          *msg_center.Mail            // 站内信
	GameMatch     *gamematch.GameMatch
	ServiceBridge *service.Bridge
	FrameSync     *frame_sync.FrameSync
	Alert         *msg_center.Alert
	GrabOrder     *grab_order.GrabOrder
	// StaticFileSystem *util.StaticFileSystem
	// Match                 *gamematch.GameMatch        //匹配服务
	// RequestServiceAdapter *service.RequestServiceAdapter //请求3方服务 适配器
	// RoomManage            *frame_sync.RoomManager     //房间服务
}

var GateDefaultProtocol = int32(util.PROTOCOL_WEBSOCKET)
var GateDefaultContentType = int32(util.CONTENT_TYPE_PROTOBUF)

// 内部服务，按说：一个项目里最多也就1-2个服务，其它的服务应该在其它项目，并且访问的时候通过HTTP/TCP，这里方便使用，先把其它服务当成一个类使用
func NewMyService() *MyService {
	var err error
	myService := new(MyService)
	// 创建一个 请求3方服务 的适配器，服务之间的请求/调用
	// myService.RequestServiceAdapter = service.NewRequestServiceAdapter(V.ServiceDiscovery, V.GrpcManager, service.REQ_SERVICE_METHOD_INNER, C.System.ProjectId, V.Zap)
	ServiceBridgeOp := service.BridgeOption{
		ProtoMap:         V.ProtoMap,
		ProjectId:        V.Project.Id,
		ServiceDiscovery: V.ServiceDiscovery,
		GrpcManager:      V.GrpcManager,
		Flag:             service.REQ_SERVICE_METHOD_NATIVE,
		Log:              V.Zap,
	}
	myService.AliSms = V.AliSms
	// 服务之间互相调用
	myService.ServiceBridge, _ = service.NewBridge(ServiceBridgeOp)
	// 预警推送
	alertOption := msg_center.AlertOption{
		SendMsgChannel:    C.Alert.SendMsgChannel,
		MsgTemplateRuleId: C.Alert.MsgTemplateRuleId,
		SendSync:          C.Alert.SendSync,
		Log:               V.Zap,
		Sms:               myService.Sms,
		SmsReceiver:       C.Alert.SmsReceiver,
		Email:             myService.Email,
		EmailReceiver:     C.Alert.EmailReceiver,
		SendUid:           C.Alert.SendUid,
	}

	myService.Alert, _ = msg_center.NewAlert(alertOption)

	if C.Service.User == "open" {
		// 用户服务
		myService.User = user_center.NewUser(V.Gorm, V.Redis, V.ProjectMng)
	}
	if C.Service.Email == "open" {
		// 站内信服务
		myService.Mail = msg_center.NewMail(V.Gorm, V.Zap)
	}
	if C.Service.Sms == "open" {
		// 短信服务
		myService.Sms = msg_center.NewSms(V.Gorm, V.AliSms, V.Zap)
	}
	if C.Service.Email == "open" {
		// 电子邮件服务
		myService.Email = msg_center.NewEmail(V.Gorm, V.Email)
	}

	if C.Service.Email == "open" {
		// 配置中心服务
		configCenterOption := config_center.ConfigCenterOption{
			EnvList:            util.GetConstListEnv(),
			Gorm:               V.Gorm,
			Redis:              V.Redis,
			ProjectManager:     V.ProjectMng,
			PersistenceType:    service.PERSISTENCE_TYPE_FILE,
			PersistenceFileDir: C.Http.StaticPath + "/" + C.ConfigCenter.DataPath,
			StaticFileSystem:   V.StaticFileSystem,
			Log:                V.Zap,
		}
		myService.ConfigCenter, err = config_center.NewConfigCenter(configCenterOption)
		if err != nil {
			util.ExitPrint("NewConfigCenter err:" + err.Error())
		}
	}
	if C.Service.TwinAgora == "open" {
		// 远程呼叫专家
		twinAgoraOption := seed_business.TwinAgoraOption{
			Log:        V.Zap,
			Gorm:       V.Gorm,
			StaticPath: C.Http.StaticPath,
			ProtoMap:   V.ProtoMap,
			// RequestServiceAdapter: myService.RequestServiceAdapter,
			ServiceBridge:    myService.ServiceBridge,
			StaticFileSystem: V.StaticFileSystem,
		}
		myService.TwinAgora, err = seed_business.NewTwinAgora(twinAgoraOption)
		if err != nil {
			util.ExitPrint(err)
		}
	}
	if C.Service.GrabOrder == "open" {
		myService.GrabOrder = grab_order.NewGrabOrder(V.Gorm)
	}
	// 网关
	if C.Gateway.Status == "open" {
		// 长连接通信 - 配置
		netWayOption := util.NetWayOption{
			ListenIp:            C.Gateway.ListenIp,     // 程序启动时监听的IP
			OutIp:               C.Gateway.OutIp,        // 对外访问的IP
			OutDomain:           C.Gateway.OutDomain,    // 对外的域名,WS在线上得用wss
			WsPort:              C.Gateway.WsPort,       // 监听端口号
			TcpPort:             C.Gateway.TcpPort,      // 监听端口号
			UdpPort:             C.Gateway.UdpPort,      // UDP端口号
			WsUri:               C.Gateway.WsUri,        // 接HOST的后面的URL地址
			DefaultProtocolType: GateDefaultProtocol,    // 兼容协议：ws tcp udp
			DefaultContentType:  GateDefaultContentType, // 默认内容格式 ：json protobuf
			LoginAuthType:       "jwt",                  // 登陆验证类型-jwt
			LoginAuthSecretKey:  C.Jwt.Key,              // 登陆验证-key
			MaxClientConnNum:    1024,                   // 客户端最大连接数
			MsgContentMax:       10240,                  // 一条消息内容最大值
			IOTimeout:           3,                      // read write sock fd 超时时间
			ConnTimeout:         60,                     // 一个FD超时时间
			ClientHeartbeatTime: 3,                      // 客户端心跳时间(秒)
			ServerHeartbeatTime: 5,                      // 服务端心跳时间(秒)
			ProtoMap:            V.ProtoMap,             // protobuf 映射表
			GrpcManager:         V.GrpcManager,
			Log:                 V.Zap,
			Gorm:                V.Gorm,
		}

		// gateway := gateway.NewGateway(V.GrpcManager, V.Zap, myService.RequestServiceAdapter)
		// gateway.MyServiceList.GameMatch = myService.GameMatch
		// gateway.MyServiceList.FrameSync = myService.FrameSync
		// gateway.MyServiceList.TwinAgora = myService.TwinAgora
		gateway := gateway.NewGateway(V.GrpcManager, V.Zap, myService.ServiceBridge)
		myService.Gateway = gateway

		_, err := gateway.StartSocket(netWayOption)
		if err != nil {
			util.ExitPrint("InitGateway err:" + err.Error())
		}
	}
	if C.Service.FrameSync == "open" {
		// 帧同步 - 房间服务 - room要先实例化,math frame_sync 都强依赖room
		frameSyncOption := frame_sync.FrameSyncOption{
			LockMode:      service.LOCK_MODE_PESSIMISTIC,
			Store:         1,
			Log:           V.Zap,
			ServiceBridge: myService.ServiceBridge,
			// RequestServiceAdapter: myService.RequestServiceAdapter,
			OffLineWaitTime: 10,
			Gorm:            V.Gorm,
			ProtoMap:        V.ProtoMap,
		}
		myService.FrameSync = frame_sync.NewFrameSync(frameSyncOption)
	}

	if C.Service.GameMatch == "open" {
		// 匹配服务 , 依赖 RoomManage
		// matchOption := service.MatchOption{
		//	RequestServiceAdapter: myService.RequestServiceAdapter,
		//	Log:                   V.Zap,
		//	RoomManager:           myService.RoomManage,
		//	//MatchSuccessChan chan *Room
		// }
		// 匹配服务，这个是假的，或者说简易版本，用于快速测试
		// myService.Match = service.NewMatch(matchOption)

		// 这个是真的匹配服务
		gmOp := gamematch.GameMatchOption{
			// RequestServiceAdapter:  myService.RequestServiceAdapter,
			ServiceBridge:          myService.ServiceBridge,
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
			FrameSync:              myService.FrameSync,
			RedisIdSeparator:       ",",
			RedisPayloadSeparation: "%",
			ProtoMap:               V.ProtoMap,
			StaticFileSystem:       V.StaticFileSystem,
		}
		myService.GameMatch, err = gamematch.NewGameMatch(gmOp)
		if err != nil {
			util.ExitPrint("NewGameMatch err:", err)
		}
	}

	if C.Cicd.Status == "open" {
		myService.Cicd, err = InitCicd()
	}

	myService.RegisterService()

	return myService
}

// 这里测试一下，服务注册到ETCD
func (myService *MyService) RegisterService() {
	if C.ServiceDiscovery.Status == "open" && C.Etcd.Status == "open" {
		var node util.ServiceNode
		ip := "127.0.0.1"
		listenIp := "127.0.0.1"
		port := "1111"

		node = util.ServiceNode{
			// ServiceId: global.C.System.ProjectId,
			ServiceId:   1, // 游戏匹配服务
			ServiceName: "GameMatch",
			Ip:          ip,
			ListenIp:    listenIp,
			Port:        port,
			Protocol:    util.SERVICE_PROTOCOL_HTTP,
			IsSelfReg:   true,
		}
		V.ServiceDiscovery.Register(node)

		node = util.ServiceNode{
			// ServiceId: global.C.System.ProjectId,
			ServiceId:   6, // 游戏匹配服务
			ServiceName: "Zgoframe",
			Ip:          ip,
			ListenIp:    listenIp,
			Port:        port,
			Protocol:    util.SERVICE_PROTOCOL_HTTP,
			IsSelfReg:   true,
		}

		V.ServiceDiscovery.Register(node)
		node = util.ServiceNode{
			// ServiceId: global.C.System.ProjectId,
			ServiceId:   2, // 游戏匹配服务
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

func InitCicd() (*cicd.CicdManager, error) {
	/*依赖
	host.toml cicd.sh
	table:  project instance server cicd_publish
	*/

	// opDirFull := MainEnv.RootDir + "/" + C.System.OpDirName
	cicdConfig := cicd.ConfigCicd{}

	// 运维：服务器的配置信息
	// configFile := opDirFull + "/host" + "." + "toml"
	//
	// // 读取配置文件中的内容
	// err := util.ReadConfFile(configFile, &cicdConfig)
	// if err != nil {
	// 	util.ExitPrint(err.Error())
	// }

	cicdConfig.SuperVisor = cicd.ConfigCicdSuperVisor{
		RpcPort:          C.SuperVisor.RpcPort,
		ConfTemplateFile: C.SuperVisor.ConfTemplateFile,
		// ConfDir:          C.SuperVisor.ConfDir,
	}

	cicdConfig.System = cicd.ConfigCicdSystem{
		Env:                C.Cicd.Env,
		LogDir:             C.Cicd.LogDir,
		WorkBaseDir:        C.Cicd.WorkBaseDir,
		RemoteBaseDir:      C.Cicd.RemoteBaseDir,
		RemoteUploadDir:    C.Cicd.RemoteUploadDir,
		RemoteDownloadDir:  C.Cicd.RemoteDownloadDir,
		MasterDirName:      C.Cicd.MasterDirName,
		GitCloneTmpDirName: C.Cicd.GitCloneTmpDirName,
		// HttpPort:           C.Cicd.HttpPort,
	}
	cicdConfig.System.RootDir = MainEnv.RootDir

	cicdConfig.SuperVisor.ConfTemplateFileName = cicdConfig.SuperVisor.ConfTemplateFile
	cicdConfig.SuperVisor.ConfTemplateFile = MainEnv.RootDir + "/" + C.System.OpDirName + "/" + cicdConfig.SuperVisor.ConfTemplateFile

	V.Zap.Debug(" ConfTemplateFile:" + cicdConfig.SuperVisor.ConfTemplateFile)
	// util.PrintStruct(cicdConfig , " : ")

	// 3方实例
	instanceManager, _ := util.NewInstanceManager(V.Gorm)
	// 服务器列表
	serverManger, _ := util.NewServerManger(V.Gorm)
	serverList := serverManger.Pool
	// 发布管理
	publicManager := cicd.NewCICDPublicManager(V.Gorm)

	// util.ExitPrint(22)
	op := cicd.CicdManagerOption{
		// HttpPort:         C.Http.Port,
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
	// 生成 filebeat 配置文件
	// cicd.GenerateAllFilebeat()
	// cicd.GetSuperVisorList()
	// 部署所有机器上的所有服务项目
	// cicd.DeployAllService()
	// go cicd.StartHttp(global.C.Http.StaticPath)

	return cicd, err

}
