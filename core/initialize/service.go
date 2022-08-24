package initialize

////内部服务，按说：一个项目里最多也就1-2个服务，其它的服务应该在其它项目，并且访问的时候通过HTTP/TCP，这里方便使用，先把其它服务当成一个类使用
//func InitMyService() {
//	netWayOption := InitGateway()
//	//fm := global.GetUploadObj(1, "")
//	MyServiceOptions := service.MyServiceOptions{
//		Gorm:             global.V.Gorm,
//		Zap:              global.V.Zap,
//		MyEmail:          global.V.Email,
//		MyRedis:          global.V.Redis,
//		MyRedisGo:        global.V.RedisGo,
//		ServiceDiscovery: global.V.ServiceDiscovery,
//		Etcd:             global.V.Etcd,
//		GrpcManager:      global.V.GrpcManager,
//		ProjectManager:   global.V.ProjectMng,
//		ServiceList:      global.V.ServiceManager.Pool,
//		Metrics:          global.V.Metric,
//
//		NetWayOption:                netWayOption,
//		ConfigCenterPersistenceType: global.C.ConfigCenter.PersistenceType,
//		HttpPort:                    global.C.Http.Port,
//		GatewayStatus:               global.C.Gateway.Status,
//
//		ProjectId:           global.C.System.ProjectId,
//		OpDirName:           global.C.System.OpDirName, //用于CICD
//		ConfigCenterDataDir: global.C.Http.StaticPath + "/" + global.C.ConfigCenter.DataPath,
//		UploadDiskPath:      global.C.Http.StaticPath + "/" + global.C.FileManager.UploadPath,
//		DownloadDiskPath:    global.C.Http.StaticPath + "/" + global.C.FileManager.DownloadPath,
//		RootDir:             global.MainEnv.RootDir,
//	}
//	myGlobal := service.Global{}
//	global.V.MyService = service.NewService(MyServiceOptions)
//	global.V.MyService = global.NewMyService()
//}
