package service

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/protobuf/pb"
	"zgoframe/protobuf/pbservice"
	"zgoframe/util"
	"fmt"
)

func Receive(c *gin.Context){

}

func ProcessOne(data []byte){

}


func LogSlaveInit(){
	ip := "localhost"
	listenIp := "localhost"
	port := "7777"


	serviceName := "LogSlave"

	project ,empty := global.V.ServiceManager.GetByName(serviceName)
	if empty{
		util.ExitPrint("project err:empty")
	}

	node := util.ServiceNode{
		ProjectId	: project.Id,
		Ip			: ip ,
		ListenIp	: listenIp,
		Port		: port ,
		Protocol	: util.SERVICE_PROTOCOL_GRPC,
		IsSelfReg	: true,
	}
	//注册一个服务(不牵扯GRPC)
	err := global.V.ServiceDiscovery.Register(node)
	if err != nil{
		util.ExitPrint("erviceDiscovery.Registe failed:"+err.Error())
	}

	MyGrpcService,err := global.V.GrpcManager.CreateService(serviceName,node.Ip,node.Port)
	if err != nil{
		util.MyPrint("GetServer err:",err)
		//return errors.New(err.Error())
	}
	//挂载服务的handler
	pb.RegisterLogSlaveServer(MyGrpcService.GrpcServer,&pbservice.LogSlave{})
	//pb.RegisterZgoframeServer(MyGrpcService.GrpcServer, &pbservice.Zgoframe{})
	//pb.RegisterSyncServer(MyGrpcService.GrpcServer, &pbservice.Sync{})
	fmt.Println("grpc ServerStart...")
	MyGrpcService.ServerStart()
	fmt.Println("GrpcServer.Serve:",err)
}