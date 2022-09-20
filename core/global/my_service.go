package global

import (
	"zgoframe/service"
	"zgoframe/service/cicd"
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

	myService.TwinAgora = service.NewTwinAgora(V.Gorm)

	//长连接通信 - 配置
	netWayOption := util.NetWayOption{
		ListenIp: C.Gateway.ListenIp, //程序启动时监听的IP
		OutIp:    C.Gateway.OutIp,    //对外访问的IP

		WsPort:  C.Gateway.WsPort,  //监听端口号
		TcpPort: C.Gateway.TcpPort, //监听端口号
		UdpPort: C.Gateway.UdpPort, //UDP端口号

		WsUri:               C.Gateway.WsUri,        //接HOST的后面的URL地址
		DefaultProtocolType: GateDefaultProtocol,    //兼容协议：ws tcp udp
		DefaultContentType:  GateDefaultContentType, //默认内容格式 ：json protobuf

		LoginAuthType:      "jwt",     //jwt
		LoginAuthSecretKey: C.Jwt.Key, //密钥

		MaxClientConnNum: 10,    //客户端最大连接数
		MsgContentMax:    10240, //一条消息内容最大值
		IOTimeout:        1,     //read write sock fd 超时时间
		ConnTimeout:      60,    //一个FD超时时间
		GrpcManager:      V.GrpcManager,
		Log:              V.Zap,
		ProtoMap:         V.ProtoMap,
		//ProtobufMapPath		string		`json:"portobuf_map_path"`//协议号对应的函数名
		//两种快速关闭方式，也可以直接调用shutdown函数
		//OutCxt 				context.Context `json:"-"`			//调用方的CTX，用于所有协程的退出操作
		//CloseChan 			chan int		`json:"-"`
		FPS:     10,
		MapSize: 10,
	}
	//匹配
	myService.Match = service.NewMatch(matchOption)
	syncOption := service.FrameSyncOption{
		Log:        V.Zap,
		RoomManage: myService.RoomManage,
		FPS:        netWayOption.FPS,
		MapSize:    netWayOption.MapSize,
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
		netway, err = gateway.StartSocket(netWayOption)
		if err != nil {
			util.ExitPrint("InitGateway err:" + err.Error())
		}

		myService.FrameSync.SetNetway(netway)
		gateway.MyServiceList.FrameSync = myService.FrameSync
		gateway.MyServiceList.Match = myService.Match
		gateway.MyServiceList.RoomManage = myService.RoomManage
		gateway.MyServiceList.TwinAgora = myService.TwinAgora
		myService.Gateway = gateway
	}
	myService.Cicd, err = InitCicd()

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

	return myService
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