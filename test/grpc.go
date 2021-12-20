package test

import (
	"context"
	"errors"
	"fmt"
	"time"
	"zgoframe/core/global"
	"zgoframe/protobuf/pb"
	"zgoframe/protobuf/pbservice"
	"zgoframe/util"
)

func Grpc(){
	//StartClient()
	client2()
	//StartService()
}


func StartService()error{

	//global.C.Grpc.ServicePackagePrefix +" . " +
	serviceName :=  global.V.Service.Name
	//serviceName := "pb.First"

	ip := "127.0.0.1"
	listenIp := "127.0.0.1"
	port := "6666"
	//ip := "8.142.177.235"
	//listenIp := "0.0.0.0"
	//port := "7777"


	node := util.ServiceNode{
		ServiceId: global.C.System.ServiceId,
		ServiceName: serviceName,
		Ip:ip ,
		ListenIp:listenIp,
		Port:port ,
		Protocol: util.SERVICE_PROTOCOL_GRPC,
		IsSelfReg: true,
	}

	global.V.ServiceDiscovery.Register(node)
	//serviceNode ,err := global.V.ServiceManager.GetServiceNodeByServiceName(serviceName)
	//if err != nil{
	//	util.ExitPrint("GetServiceNodeByServiceName err")
	//}

	MyGrpcServer,err := global.V.Grpc.GetServer(serviceName,node.ListenIp,node.Port)
	//grpcOption := util.GrpcOption{
	//	AppId 		: global.V.App.Id,
	//	ListenIp	: global.C.Grpc.Ip,
	//	OutIp		: global.C.Grpc.Ip,
	//	Port 		: global.C.Grpc.Port,
	//	Log			: global.V.Zap,
	//}
	//global.V.Grpc,_ =  util.NewMyGrpc(grpcOption)
	//grpcInc,listen,err := global.V.Grpc.GetServer()
	if err != nil{
		util.MyPrint("GetServer err:",err)
		return errors.New(err.Error())
	}
	//挂载服务的handler
	pb.RegisterZgoframeServer(MyGrpcServer.GrpcServer, &pbservice.Zgoframe{})
	fmt.Println("grpc ServerStart...")
	MyGrpcServer.ServerStart()
	fmt.Println("GrpcServer.Serve:",err)
	return nil
}

func  StartClient()error{
	//util.ExitPrint(global.V.ServiceManager.GetByName("zgoframe"))
	//grpcClientConn,err := global.V.Grpc.GetClient(global.C.Grpc.Ip,global.C.Grpc.Port)
	//dns := global.C.Grpc.Ip+ ":4141"
	//dns := global.C.Grpc.Ip+ ":6666"
	//grpcClientConn, err := grpc.Dial(dns,grpc.WithInsecure())
	serviceName :=  global.V.Service.Name
	serviceNode ,err := global.V.ServiceDiscovery.GetLoadBalanceServiceNodeByServiceName(serviceName,"")
	if err != nil{
		util.ExitPrint("GetServiceNodeByServiceName err:",err)
	}

	fmt.Println("serviceNode:",serviceNode)
	grpcClientConn, err := global.V.Grpc.GetClient(serviceName,global.V.App.Id,serviceNode.Ip,serviceNode.Port)
	//grpcClientConn, err := grpc.Dial(dns,grpc.WithInsecure(),grpc.WithUnaryInterceptor(clientInterceptorBack))
	//util.MyPrint("client grp dns:",dns , " err:",err)
	if err != nil{
		util.MyPrint("grpc GetClient err:",err)
		return  err
	}

	pbServiceFirst := pb.NewZgoframeClient(grpcClientConn)
	RequestRegPlayer := pb.RequestUser{}
	RequestRegPlayer.Id = 123123
	RequestRegPlayer.Nickname = "xiaoz"
	res ,err:= pbServiceFirst.SayHello(context.Background(),&RequestRegPlayer)
	util.MyPrint("grpc return:",res , " err:",err)


	global.V.ServiceDiscovery.ShowJsonByService()
	global.V.ServiceDiscovery.ShowJsonByNodeServer()

	return nil
}

func client2()error{

	go clientSend()

	return nil
}

func clientSend(){
	for{
		serviceName :=  global.V.Service.Name
		grpcClientConn, err := global.V.Grpc.GetClientByLoadBalance(serviceName,0)
		if err != nil{
			util.MyPrint("grpc GetClient err:",err)
			return
		}

		pbServiceFirst := pb.NewZgoframeClient(grpcClientConn)
		RequestRegPlayer := pb.RequestUser{}
		RequestRegPlayer.Id = 123123
		RequestRegPlayer.Nickname = "xiaoz"
		res ,err:= pbServiceFirst.SayHello(context.Background(),&RequestRegPlayer)
		util.MyPrint("grpc return:",res , " err:",err)

		time.Sleep(time.Second * 1)
		util.MyPrint("sleep 1 second...")
	}
}


//===================================================
//var ip = "127.0.0.1"
//var port = "1111"
//func TestGrpcServer(){
//	grpcOption := GrpcOption{
//		ListenIp:ip,
//		OutIp: ip,
//		Port: port,
//		//Role: ROLE_SERVER,
//	}
//	grpcClass := NewMyGrpc(grpcOption)
//	grpcClass.GetServer()
//}

func TestGrpcClient(){
	//grpcOption := GrpcOption{
	//	AppId: 1,
	//	ListenIp:ip,
	//	OutIp: ip,
	//	Port: port,
	//	//Role: ROLE_SERVER,
	//}
	//grpcClass := NewMyGrpc(grpcOption)
	//c ,err := grpcClass.GetClient()
	//if err != nil{
	//	MyPrint("grpcClass.GetClient() err:",err.Error())
	//	return
	//}
	//channel := pb.Channel{Id: 999,Name: "tttt"}
	//
	//channelMap := make(map[uint64]*pb.Channel)
	//channelMap[333] = &channel
	//player := pb.Player{
	//	Id: 1,
	//	RoleName: []byte("aaaaaa"),
	//	Nickname: "xiaoz",
	//	Status: 1,
	//	Score: 11.222,
	//	PhoneType: pb.PhoneType_home,
	//	Sex: true,
	//	Level: 10,
	//}
	//
	//playerList := []*pb.Player{&player}
	//
	//requestRegPlayer := pb.RequestRegPlayer{
	//	AddTime : 11234,
	//	PlayerInfo: &player,
	//	PlayerList :playerList,
	//	ChannelMap :channelMap,
	//}
	//
	////if !ok {
	////	ExitPrint("metadata.FromIncomingContext err:",md,ok)
	////}
	//r, err := c.SayHello(context.Background(), &requestRegPlayer)
	//
	//
	//if err != nil {
	//	MyPrint("could not service: %v", err)
	//}
	//
	//MyPrint("FINAL:",r.Rs)

}

