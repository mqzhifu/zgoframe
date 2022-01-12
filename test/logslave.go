package test

import (
	"zgoframe/core/global"
	"zgoframe/protobuf/pb"
	"zgoframe/protobuf/pbservice"
	"zgoframe/util"
	"fmt"
)

func LogSlave()  {
	ip := "127.0.0.1"
	listenIp := "127.0.0.1"
	port := "7777"


	serviceName := global.V.Project.Key
	project ,empty := global.V.ServiceManager.GetByName(serviceName)
	if empty{
		util.ExitPrint("project err:empty")
		return
	}

	node := util.ServiceNode{
		ProjectId	: project.Id,
		Ip			: ip ,
		ListenIp	: listenIp,
		Port		: port ,
		Protocol	: util.SERVICE_PROTOCOL_GRPC,
		IsSelfReg	: true ,
	}
	//注册一个服务发现(不牵扯GRPC)
	err := global.V.ServiceDiscovery.Register(node)
	if err != nil{
		util.ExitPrint("serviceDiscovery.Register failed:"+err.Error())
		return
	}
	//创建一个grpc 服务
	MyGrpcService,err := global.V.GrpcManager.CreateService(serviceName,node.Ip,node.Port)
	if err != nil{
		util.MyPrint("GetServer err:",err)
		return
	}
	//挂载服务的handler
	pb.RegisterLogSlaveServer(MyGrpcService.GrpcServer,&pbservice.LogSlave{})
	fmt.Println("grpc ServerStart...")
	err = MyGrpcService.ServerStart()
	fmt.Println("GrpcServer.Serve:",err)
}