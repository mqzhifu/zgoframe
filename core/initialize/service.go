package initialize

import (
	"zgoframe/core/global"
	"zgoframe/service"
)
//内部服务，按说：一个项目里最多也就1-2个服务，其它的服务应该在其它项目，并且访问的时候通过HTTP/TCP，这里方便使用，先统计把其它服务当成一个类使用
func InitMyService( ){
	netWayOption := InitGateway()
	MyServiceOptions := service.MyServiceOptions {
		Gorm :global.V.Gorm,
		Zap :global.V.Zap,
		MyEmail :global.V.Email,
		MyRedis :global.V.Redis,
		NetWayOption :netWayOption,
		GrpcManager :global.V.GrpcManager,
		ProjectManager : global.V.ProjectMng,
		ConfigCenterDataDir :global.C.Http.StaticPath + global.C.ConfigCenter.DataPath,
		ConfigCenterPersistenceType	:global.C.ConfigCenter.PersistenceType,
		OpDirName: global.C.System.OpDirName,
		ServiceList: global.V.ServiceManager.Pool,
		HttpPort:global.C.Http.Port,
		GatewayStatus: global.C.Gateway.Status,
	}

	global.V.MyService = service.NewService(MyServiceOptions)
}
