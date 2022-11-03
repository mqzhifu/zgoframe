package global

import (
	"zgoframe/service"
	"zgoframe/service/cicd"
	gamematch "zgoframe/service/game_match"
	"zgoframe/util"
)

type MyService struct {
	User         *service.User         //用户中心
	Sms          *service.Sms          //短信服务
	Email        *service.Email        //电子邮件服务
	RoomManage   *service.RoomManager  //房间服务
	FrameSync    *service.FrameSync    //帧同步服务
	Gateway      *service.Gateway      //网关服务
	Match        *service.Match        //匹配服务
	TwinAgora    *service.TwinAgora    //广播120远程专家指导
	ConfigCenter *service.ConfigCenter //配置中心
	Cicd         *cicd.CicdManager     //自动部署
	Mail         *service.Mail         //站内信
	GameMatch    *gamematch.GameMatch
}

var GateDefaultProtocol = int32(util.PROTOCOL_WEBSOCKET)
var GateDefaultContentType = int32(util.CONTENT_TYPE_JSON)

//内部服务，按说：一个项目里最多也就1-2个服务，其它的服务应该在其它项目，并且访问的时候通过HTTP/TCP，这里方便使用，先把其它服务当成一个类使用
func NewMyService() *MyService {
	var err error
	myService := new(MyService)
	//用户服务
	myService.User = service.NewUser(V.Gorm, V.Redis)
	//站内信服务
	myService.Mail = service.NewMail(V.Gorm, V.Zap)
	//短信服务
	myService.Sms = service.NewSms(V.Gorm)
	//电子邮件服务
	myService.Email = service.NewEmail(V.Gorm, V.Email)
	//配置中心服务
	configCenterOption := service.ConfigCenterOption{
		EnvList:            util.GetConstListEnv(),
		Gorm:               V.Gorm,
		Redis:              V.Redis,
		ProjectManager:     V.ProjectMng,
		PersistenceType:    service.PERSISTENCE_TYPE_FILE,
		PersistenceFileDir: C.Http.StaticPath + "/" + C.ConfigCenter.DataPath,
		Log:                V.Zap,
	}
	myService.ConfigCenter, err = service.NewConfigCenter(configCenterOption)
	if err != nil {
		util.ExitPrint("NewConfigCenter err:" + err.Error())
	}
	//房间服务 - room要先实例化,math frame_sync 都强依赖room
	roomManagerOption := service.RoomManagerOption{
		Log:          V.Zap,
		ReadyTimeout: 60,
		RoomPeople:   2,
	}
	myService.RoomManage = service.NewRoomManager(roomManagerOption)
	//匹配服务 , 依赖 RoomManage
	matchOption := service.MatchOption{
		Log:         V.Zap,
		RoomManager: myService.RoomManage,
		//MatchSuccessChan chan *Room
	}
	//远程呼叫专家
	myService.TwinAgora, _ = service.NewTwinAgora(V.Gorm, V.Zap, C.Http.StaticPath)

	//长连接通信 - 配置
	netWayOption := util.NetWayOption{
		ListenIp:            C.Gateway.ListenIp, //程序启动时监听的IP
		OutIp:               C.Gateway.OutIp,    //对外访问的IP
		OutDomain:           C.Gateway.OutDomain,
		WsPort:              C.Gateway.WsPort,       //监听端口号
		TcpPort:             C.Gateway.TcpPort,      //监听端口号
		UdpPort:             C.Gateway.UdpPort,      //UDP端口号
		WsUri:               C.Gateway.WsUri,        //接HOST的后面的URL地址
		DefaultProtocolType: GateDefaultProtocol,    //兼容协议：ws tcp udp
		DefaultContentType:  GateDefaultContentType, //默认内容格式 ：json protobuf
		LoginAuthType:       "jwt",                  //jwt
		LoginAuthSecretKey:  C.Jwt.Key,
		MaxClientConnNum:    10,    //客户端最大连接数
		MsgContentMax:       10240, //一条消息内容最大值
		IOTimeout:           3,     //read write sock fd 超时时间
		ConnTimeout:         60,    //一个FD超时时间
		ClientHeartbeatTime: 3,
		ServerHeartbeatTime: 5,
		GrpcManager:         V.GrpcManager,
		Log:                 V.Zap,
		ProtoMap:            V.ProtoMap,
		//ProtobufMapPath		string		`json:"portobuf_map_path"`//协议号对应的函数名
		//两种快速关闭方式，也可以直接调用 shutdown 函数
		//OutCxt 				context.Context `json:"-"`			//调用方的CTX，用于所有协程的退出操作
		//CloseChan 			chan int		`json:"-"`
		//FPS:     10,
		//MapSize: 10,
	}
	//匹配服务
	myService.Match = service.NewMatch(matchOption)
	syncOption := service.FrameSyncOption{
		Log:        V.Zap,
		RoomManage: myService.RoomManage,
		FPS:        10,
		MapSize:    5,
	}
	go myService.Match.Start()
	//user -> sign ->Match -> Room -> Rsync
	//帧同步服务 - 强-依赖room
	myService.FrameSync = service.NewFrameSync(syncOption)
	myService.RoomManage.SetFrameSync(myService.FrameSync)
	//网关
	if C.Gateway.Status == "open" {
		gateway := service.NewGateway(V.GrpcManager, V.Zap)
		var netway *util.NetWay
		myService.FrameSync.SetNetway(netway)
		gateway.MyServiceList.FrameSync = myService.FrameSync
		gateway.MyServiceList.Match = myService.Match
		gateway.MyServiceList.RoomManage = myService.RoomManage
		gateway.MyServiceList.TwinAgora = myService.TwinAgora
		myService.Gateway = gateway

		netway, err = gateway.StartSocket(netWayOption)
		if err != nil {
			util.ExitPrint("InitGateway err:" + err.Error())
		}

	}
	myService.Cicd, err = InitCicd()

	//type GameMatchOption struct {
	//	Log                *zap.Logger            //log 实例
	//	Redis              *util.MyRedisGo        //redis 实例
	//	Gorm               *gorm.DB               //mysql 实例
	//	Service            *util.Service          //服务 实例
	//	Metrics            *util.MyMetrics        //统计 实例
	//	ServiceDiscovery   *util.ServiceDiscovery //服务发现 实例
	//	StaticPath         string                 //静态文件公共目录
	//	RuleDataSourceType int                    //rule的数据来源类型
	//	RedisPrefix        string                 //redis公共的前缀，主要是怕key重复
	//	RedisTextSeparator string                 //结构体不能直接存到redis中，得手动分隔存进去。不存JSON是因为浪费空间
	//	RedisKeySeparator  string                 //redis key 的分隔符号
	//	ProjectId          int
	//	//Etcd             *util.MyEtcd
	//}

	gmOp := gamematch.GameMatchOption{
		Log:     V.Zap,
		Redis:   V.RedisGo,
		Gorm:    V.Gorm,
		Metrics: V.Metric,
		//Service:            V.ServiceManager,
		ServiceDiscovery:       V.ServiceDiscovery,
		RuleDataSourceType:     service.GAME_MATCH_DATA_SOURCE_TYPE_DB,
		StaticPath:             C.Http.StaticPath,
		RedisPrefix:            "gm",
		RedisKeySeparator:      "_",
		RedisTextSeparator:     "#",
		RedisIdSeparator:       ",",
		RedisPayloadSeparation: "%",
	}
	myService.GameMatch, err = gamematch.NewGameMatch(gmOp)
	if err != nil {
		util.ExitPrint("NewGameMatch err:", err)
	}

	//这个是真的匹配服务，上面是个假的DEMO类型的匹配服务
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

	myService.RegisterService()

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
