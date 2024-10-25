// 微服务  - 具体的业务
package service

import (
	"fmt"
	"zgoframe/core/global"
	"zgoframe/service/bridge"
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
	User                  *user_center.User             // 用户中心
	Sms                   *msg_center.Sms               // 短信服务
	Email                 *msg_center.Email             // 电子邮件服务
	AliSms                *util.AliSms                  // 阿里-短信业务
	Gateway               *gateway.Gateway              // 网关服务
	TwinAgora             *seed_business.TwinAgora      // 120远程专家指导
	ConfigCenter          *config_center.ConfigCenter   // 配置中心
	Cicd                  *cicd.CicdManager             // 自动部署
	Mail                  *msg_center.Mail              // 站内信
	GameMatch             *gamematch.GameMatch          // 游戏匹配
	FrameSync             *frame_sync.FrameSync         // 游戏帧同步
	Alert                 *msg_center.Alert             // 报警
	GrabOrder             *grab_order.GrabOrder         // 抢单服务
	ServiceBridge         *bridge.Bridge                // 服务之间通信-桥连接
	RequestServiceAdapter *bridge.RequestServiceAdapter // 请求3方服务 适配器
	// StaticFileSystem *util.StaticFileSystem
	// Match                 *gamematch.GameMatch        //匹配服务
	// RoomManage            *frame_sync.RoomManager     //房间服务
}

var GateDefaultProtocol = int32(util.PROTOCOL_WEBSOCKET)
var GateDefaultContentType = int32(util.CONTENT_TYPE_PROTOBUF)

// 内部服务，按说：一个项目里最多也就1-2个服务，其它的服务应该在其它项目，并且访问的时候通过HTTP/TCP，这里方便使用，先把其它服务当成一个类使用
func NewMyService() *MyService {
	global.V.Base.Zap.Info("NewMyService start")
	var err error
	myService := new(MyService)
	// 创建一个 请求3方服务 的适配器，服务之间的请求/调用
	myService.RequestServiceAdapter = bridge.NewRequestServiceAdapter(global.V.Util.ServiceDiscovery, global.V.Util.GrpcManager, bridge.REQ_SERVICE_METHOD_NATIVE, global.C.System.ProjectId, global.V.Base.Zap)
	ServiceBridgeOp := bridge.BridgeOption{
		ProtoMap:         global.V.Util.ProtoMap,
		ProjectId:        global.V.Util.Project.Id,
		ServiceDiscovery: global.V.Util.ServiceDiscovery,
		GrpcManager:      global.V.Util.GrpcManager,
		Flag:             bridge.REQ_SERVICE_METHOD_NATIVE,
		Log:              global.V.Base.Zap,
	}
	// 服务之间互相调用
	myService.ServiceBridge, _ = bridge.NewBridge(ServiceBridgeOp)

	myService.AliSms = global.V.Util.AliSms
	// 预警推送
	alertOption := msg_center.AlertOption{
		SendMsgChannel:    global.C.Alert.SendMsgChannel,
		MsgTemplateRuleId: global.C.Alert.MsgTemplateRuleId,
		SendSync:          global.C.Alert.SendSync,
		Log:               global.V.Base.Zap,
		Sms:               myService.Sms,
		SmsReceiver:       global.C.Alert.SmsReceiver,
		Email:             myService.Email,
		EmailReceiver:     global.C.Alert.EmailReceiver,
		SendUid:           global.C.Alert.SendUid,
	}
	global.V.Base.Zap.Info("NewAlert:")
	myService.Alert, _ = msg_center.NewAlert(alertOption)
	// 用户中心服务
	if global.C.Service.User == "open" {
		myService.User = user_center.NewUser(global.V.Base.Gorm, global.V.Base.Redis, global.V.Util.ProjectMng, global.V.Base.Zap)
	}
	// 站内信服务
	if global.C.Service.Mail == "open" {
		myService.Mail = msg_center.NewMail(global.V.Base.Gorm, global.V.Base.Zap)
	}
	// 短信服务
	if global.C.Service.Sms == "open" {
		myService.Sms = msg_center.NewSms(global.V.Base.Gorm, global.V.Util.AliSms, global.V.Base.Zap)
	}
	// 电子邮件服务
	if global.C.Service.Email == "open" {
		fmt.Println("Service.Email:")
		myService.Email = msg_center.NewEmail(global.V.Base.Gorm, global.V.Util.Email)
	}
	//抢单-服务
	if global.C.Service.GrabOrder == "open" {
		fmt.Println("Service.GrabOrder:")
		myService.GrabOrder = grab_order.NewGrabOrder(global.V.Base.Gorm, global.V.Base.Redis)
	}
	// 配置中心服务
	if global.C.Service.ConfigCenter == "open" {
		configCenterOption := config_center.ConfigCenterOption{
			EnvList:            util.GetConstListEnv(),
			Gorm:               global.V.Base.Gorm,
			Redis:              global.V.Base.Redis,
			ProjectManager:     global.V.Util.ProjectMng,
			PersistenceType:    config_center.PERSISTENCE_TYPE_FILE,
			PersistenceFileDir: global.C.Http.StaticPath + "/" + global.C.ConfigCenter.DataPath,
			StaticFileSystem:   global.V.Util.StaticFileSystem,
			Log:                global.V.Base.Zap,
		}
		myService.ConfigCenter, err = config_center.NewConfigCenter(configCenterOption)
		if err != nil {
			fmt.Println("Service.ConfigCenter:")
			util.ExitPrint("NewConfigCenter err:" + err.Error())
		}
	}
	if global.C.Service.TwinAgora == "open" {
		// 远程呼叫专家
		twinAgoraOption := seed_business.TwinAgoraOption{
			Log:                   global.V.Base.Zap,
			Gorm:                  global.V.Base.Gorm,
			StaticPath:            global.C.Http.StaticPath,
			ProtoMap:              global.V.Util.ProtoMap,
			RequestServiceAdapter: myService.RequestServiceAdapter,
			ServiceBridge:         myService.ServiceBridge,
			StaticFileSystem:      global.V.Util.StaticFileSystem,
		}
		myService.TwinAgora, err = seed_business.NewTwinAgora(twinAgoraOption)
		if err != nil {
			util.ExitPrint(err)
		}
	}
	// 网关
	if global.C.Gateway.Status == "open" {
		// 长连接通信 - 配置
		netWayOption := util.NetWayOption{
			ListenIp:            global.C.Gateway.ListenIp,  // 程序启动时监听的IP
			OutIp:               global.C.Gateway.OutIp,     // 对外访问的IP
			OutDomain:           global.C.Gateway.OutDomain, // 对外的域名,WS在线上得用wss
			WsPort:              global.C.Gateway.WsPort,    // 监听端口号
			TcpPort:             global.C.Gateway.TcpPort,   // 监听端口号
			UdpPort:             global.C.Gateway.UdpPort,   // UDP端口号
			WsUri:               global.C.Gateway.WsUri,     // 接HOST的后面的URL地址
			DefaultProtocolType: GateDefaultProtocol,        // 兼容协议：ws tcp udp
			DefaultContentType:  GateDefaultContentType,     // 默认内容格式 ：json protobuf
			LoginAuthType:       "jwt",                      // 登陆验证类型-jwt
			LoginAuthSecretKey:  global.C.Jwt.Key,           // 登陆验证-key
			MaxClientConnNum:    1024,                       // 客户端最大连接数
			MsgContentMax:       10240,                      // 一条消息内容最大值
			IOTimeout:           3,                          // read write sock fd 超时时间
			ConnTimeout:         60,                         // 一个FD超时时间
			ClientHeartbeatTime: 3,                          // 客户端心跳时间(秒)
			ServerHeartbeatTime: 5,                          // 服务端心跳时间(秒)
			ProtoMap:            global.V.Util.ProtoMap,     // protobuf 映射表
			Log:                 global.V.Base.Zap,          // 日志
			Gorm:                global.V.Base.Gorm,         // 数据库，用于持久化
			GrpcManager:         global.V.Util.GrpcManager,
		}

		gatewayInc := gateway.NewGateway(netWayOption, myService.ServiceBridge, myService.RequestServiceAdapter)
		myService.Gateway = gatewayInc

		_, err := gatewayInc.StartSocket()
		if err != nil {
			util.ExitPrint("InitGateway err:" + err.Error())
		}
	}
	//帧同步
	if global.C.Service.FrameSync == "open" {
		fmt.Println("Service.FrameSync:")
		frameSyncOption := frame_sync.FrameSyncOption{
			LockMode:              frame_sync.LOCK_MODE_PESSIMISTIC,
			Store:                 1,
			Log:                   global.V.Base.Zap,
			RequestServiceAdapter: myService.RequestServiceAdapter,
			ServiceBridge:         myService.ServiceBridge,
			OffLineWaitTime:       10,
			Gorm:                  global.V.Base.Gorm,
			ProtoMap:              global.V.Util.ProtoMap,
		}
		myService.FrameSync = frame_sync.NewFrameSync(frameSyncOption)
	}

	if global.C.Service.GameMatch == "open" {
		fmt.Println("Service.GameMatch:")
		// 匹配服务 , 依赖 RoomManage
		// matchOption := service.MatchOption{
		//	RequestServiceAdapter: myService.RequestServiceAdapter,
		//	Log:                   global.V.Base.Zap,
		//	RoomManager:           myService.RoomManage,
		//	//MatchSuccessChan chan *Room
		// }
		// 匹配服务，这个是假的，或者说简易版本，用于快速测试
		// myService.Match = service.NewMatch(matchOption)

		// 这个是真的匹配服务
		gmOp := gamematch.GameMatchOption{
			// RequestServiceAdapter:  myService.RequestServiceAdapter,
			ServiceBridge:          myService.ServiceBridge,
			Log:                    global.V.Base.Zap,
			Redis:                  global.V.Base.RedisGo,
			Gorm:                   global.V.Base.Gorm,
			Metrics:                global.V.Util.Metric,
			ServiceDiscovery:       global.V.Util.ServiceDiscovery,
			RuleDataSourceType:     gamematch.GAME_MATCH_DATA_SOURCE_TYPE_DB,
			StaticPath:             global.C.Http.StaticPath,
			RedisPrefix:            "gm",
			RedisKeySeparator:      "_",
			RedisTextSeparator:     "#",
			FrameSync:              myService.FrameSync,
			RedisIdSeparator:       ",",
			RedisPayloadSeparation: "%",
			ProtoMap:               global.V.Util.ProtoMap,
			StaticFileSystem:       global.V.Util.StaticFileSystem,
		}
		myService.GameMatch, err = gamematch.NewGameMatch(gmOp)
		if err != nil {
			util.ExitPrint("NewGameMatch err:", err)
		}
	}

	if global.C.Cicd.Status == "open" {
		fmt.Println("Service.Cicd:")
		myService.Cicd, err = InitCicd()
	}

	myService.RegisterService()

	return myService
}

// 这里测试一下，服务注册到ETCD
func (myService *MyService) RegisterService() {
	if global.C.ServiceDiscovery.Status == "open" && global.C.Etcd.Status == "open" {
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
		global.V.Util.ServiceDiscovery.Register(node)

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

		global.V.Util.ServiceDiscovery.Register(node)
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
		global.V.Util.ServiceDiscovery.Register(node)

	}
}

func InitCicd() (*cicd.CicdManager, error) {
	/*依赖
	host.toml cicd.sh
	table:  project instance server cicd_publish
	*/

	// opDirFull := MainEnglobal.V.Base.RootDir + "/" + C.System.OpDirName
	cicdConfig := cicd.ConfigCicd{}

	// 运维：服务器的配置信息
	// configFile := opDirFull + "/host" + "." + "toml"
	//
	// // 读取配置文件中的内容
	// err := util.ReadConfFile(configFile, &cicdConfig)
	// if err != nil {
	// }

	cicdConfig.SuperVisor = cicd.ConfigCicdSuperVisor{
		RpcPort:          global.C.SuperVisor.RpcPort,
		ConfTemplateFile: global.C.SuperVisor.ConfTemplateFile,
		// ConfDir:          C.SuperVisor.ConfDir,
	}

	cicdConfig.System = cicd.ConfigCicdSystem{
		Env:                global.C.Cicd.Env,
		LogDir:             global.C.Cicd.LogDir,
		WorkBaseDir:        global.C.Cicd.WorkBaseDir,
		RemoteBaseDir:      global.C.Cicd.RemoteBaseDir,
		RemoteUploadDir:    global.C.Cicd.RemoteUploadDir,
		RemoteDownloadDir:  global.C.Cicd.RemoteDownloadDir,
		MasterDirName:      global.C.Cicd.MasterDirName,
		GitCloneTmpDirName: global.C.Cicd.GitCloneTmpDirName,
		// HttpPort:           C.Cicd.HttpPort,
	}

	cicdConfig.System.RootDir = global.MainEnv.RootDir

	cicdConfig.SuperVisor.ConfTemplateFileName = cicdConfig.SuperVisor.ConfTemplateFile
	cicdConfig.SuperVisor.ConfTemplateFile = global.MainEnv.RootDir + "/" + global.C.System.OpDirName + "/" + cicdConfig.SuperVisor.ConfTemplateFile

	global.V.Base.Zap.Debug(" ConfTemplateFile:" + cicdConfig.SuperVisor.ConfTemplateFile)
	// util.PrintStruct(cicdConfig , " : ")

	// 3方实例
	instanceManager, _ := util.NewInstanceManager(global.V.Base.Gorm)
	// 服务器列表
	serverManger, _ := util.NewServerManger(global.V.Base.Gorm)
	serverList := serverManger.Pool
	// 发布管理
	publicManager := cicd.NewCICDPublicManager(global.V.Base.Gorm)

	op := cicd.CicdManagerOption{
		// HttpPort:         C.Http.Port,
		ServerList:       serverList,
		Config:           cicdConfig,
		InstanceManager:  instanceManager,
		PublicManager:    publicManager,
		Log:              global.V.Base.Zap,
		OpDirName:        global.C.System.OpDirName,
		ServiceList:      global.V.Util.ServiceManager.Pool,
		ProjectList:      global.V.Util.ProjectMng.Pool,
		UploadDiskPath:   global.C.Http.StaticPath + "/" + global.C.FileManager.UploadPath,
		DownloadDiskPath: global.C.Http.StaticPath + "/" + global.C.FileManager.DownloadPath,
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
