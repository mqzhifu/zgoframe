package initialize

import (
	"zgoframe/core/global"
	"zgoframe/service"
	"zgoframe/util"
)
//内部服务，按说：一个项目里最多也就1-2个服务，其它的服务应该在其它项目，并且访问的时候通过HTTP/TCP，这里方便使用，先统计把其它服务当成一个类使用
func InitMyService( redisGo *util.MyRedisGo){
	netWayOption := InitGateway()
	MyServiceOptions := service.MyServiceOptions {
		Gorm 			: global.V.Gorm,
		Zap 			: global.V.Zap,
		MyEmail 		: global.V.Email,
		MyRedis 		: global.V.Redis,
		NetWayOption 	: netWayOption,
		GrpcManager 	: global.V.GrpcManager,
		ProjectManager 	: global.V.ProjectMng,
		//ConfigCenterDataDir :global.C.Http.StaticPath + global.C.ConfigCenter.DataPath,
		ConfigCenterDataDir : global.C.ConfigCenter.DataPath,
		ConfigCenterPersistenceType	:global.C.ConfigCenter.PersistenceType,
		OpDirName		: global.C.System.OpDirName,//用于CICD
		ServiceList		: global.V.ServiceManager.Pool,
		HttpPort		: global.C.Http.Port,
		GatewayStatus	: global.C.Gateway.Status,
		MyRedisGo		: redisGo,
		Etcd			: global.V.Etcd,
		Metrics			: global.V.Metric,
		ServiceDiscovery: global.V.ServiceDiscovery,
		ProjectId 		: global.C.System.ProjectId,
		UploadDiskPath  : global.V.RootDir + "/" + global.C.Upload.Path,
	}

	global.V.MyService = service.NewService(MyServiceOptions)
}
