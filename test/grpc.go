package test

import (
	"errors"
	"fmt"
	"zgoframe/core/global"
	"zgoframe/protobuf/pb"
	"zgoframe/protobuf/pbservice"
	"zgoframe/util"
	"context"
)

func Grpc(){
	//StartClient()
	client2()
	//StartService()
}


func StartService()error{
	//包前缀 + 服务名
	serviceName :=  global.C.Grpc.ServicePackagePrefix +"." + global.V.Service.Name
	//serviceName := "pb.First"
	ip := "127.0.0.1"
	listenIp := "127.0.0.1"
	port := "6666"
	//ip := "8.142.177.235"
	//listenIp := "0.0.0.0"
	//port := "7777"


	node := util.ServiceNode{
		ProjectId	: global.C.System.ProjectId,
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
	//测试一下刚刚注册的服务，是否成功，从服务管理池中直接寻找
	testServerRegRs := false
	for _,service := range global.V.ServiceManager.Pool{
		if service.Name == serviceName{
			testServerRegRs = true
			break
		}
	}
	if !testServerRegRs{
		util.ExitPrint("reg failed .")
	}
	//服务发现注册成功后，再创建一个grpc server
	MyGrpcService,err := global.V.GrpcManager.CreateService(serviceName,node.Ip,node.Port)
	if err != nil{
		util.MyPrint("GetServer err:",err)
		return errors.New(err.Error())
	}
	//挂载服务的handler
	pb.RegisterZgoframeServer(MyGrpcService.GrpcServer, &pbservice.Zgoframe{})
	//pb.RegisterSyncServer(MyGrpcService.GrpcServer, &pbservice.Sync{})
	fmt.Println("grpc ServerStart...")
	MyGrpcService.ServerStart()
	fmt.Println("GrpcServer.Serve:",err)
	return nil
}

func  StartClient()error{
	serviceName :=  global.V.Service.Key
	zgoframeClient ,err := global.V.GrpcManager.GetZgoframeClient(serviceName,"")
	if err != nil{
		RequestRegPlayer := pb.RequestUser{}
		RequestRegPlayer.Id = 123123
		RequestRegPlayer.Nickname = "xiaoz"
		res ,err:= zgoframeClient.SayHello(context.Background(),&RequestRegPlayer)
		fmt.Println(res,err)
	}


	//serviceNode ,err := global.V.ServiceDiscovery.GetLoadBalanceServiceNodeByServiceName(serviceName,"")
	//if err != nil{
	//	util.ExitPrint("GetServiceNodeByServiceName err:",err)
	//}
	//
	//fmt.Println("serviceNode:",serviceNode)
	//grpcClientConn, err := global.V.Grpc.GetClient(serviceName,global.V.Project.Id,serviceNode.Ip,serviceNode.Port)
	////grpcClientConn, err := grpc.Dial(dns,grpc.WithInsecure(),grpc.WithUnaryInterceptor(clientInterceptorBack))
	////util.MyPrint("client grp dns:",dns , " err:",err)
	//if err != nil{
	//	util.MyPrint("grpc GetClient err:",err)
	//	return  err
	//}
	//
	//pbServiceFirst := pb.NewZgoframeClient(grpcClientConn)
	//RequestRegPlayer := pb.RequestUser{}
	//RequestRegPlayer.Id = 123123
	//RequestRegPlayer.Nickname = "xiaoz"
	//res ,err:= pbServiceFirst.SayHello(context.Background(),&RequestRegPlayer)
	//util.MyPrint("grpc return:",res , " err:",err)
	//
	//
	//global.V.ServiceDiscovery.ShowJsonByService()
	//global.V.ServiceDiscovery.ShowJsonByNodeServer()

	return nil
}

func client2()error{

	go clientSend()

	return nil
}

func clientSend(){
	//for{
	//	serviceName :=  global.V.Service.Name
	//	grpcClientConn, err := global.V.Grpc.GetClientByLoadBalance(serviceName,0)
	//	if err != nil{
	//		util.MyPrint("grpc GetClient err:",err)
	//		return
	//	}
	//
	//	pbServiceFirst := pb.NewZgoframeClient(grpcClientConn)
	//	RequestRegPlayer := pb.RequestUser{}
	//	RequestRegPlayer.Id = 123123
	//	RequestRegPlayer.Nickname = "xiaoz"
	//	res ,err:= pbServiceFirst.SayHello(context.Background(),&RequestRegPlayer)
	//	util.MyPrint("grpc return:",res , " err:",err)
	//
	//	time.Sleep(time.Second * 1)
	//	util.MyPrint("sleep 1 second...")
	//}
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

