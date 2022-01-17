package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"strconv"
	"time"
)
//帮助管理GRPC连接、请求、响应
type GrpcManager struct {
	Option 			GrpcManagerOption
	ClientList 		map[string]*MyGrpcClient	//map[dns]*MyGrpcClient
	ServiceList 	map[string]*MyGrpcService	//map[dns]*MyGrpcService
}
//实例化参数
type GrpcManagerOption struct {
	ProjectId 			int				//在做链路追踪时，需要有项目ID
	Log 				*zap.Logger
	ServiceDiscovery *ServiceDiscovery	//获取一个grpc连接时，需要使用
}
//创建一个实例
func NewGrpcManager(grpcManagerOption GrpcManagerOption)(*GrpcManager,error){
	grpcManager 		   := new(GrpcManager)
	grpcManager.Option 		= grpcManagerOption
	//存储自己创建的服务
	grpcManager.ServiceList =	 make(map[string]*MyGrpcService)
	//存储它人创建的服务
	grpcManager.ClientList 		= make(map[string]*MyGrpcClient)
	//监听服务发现发生变更
	go grpcManager.WatchServiceChange()

	return grpcManager,nil
}
//服务发现里的数据发生变更，要连带着更新GRPC里的连接信息
func  (grpcManager *GrpcManager)WatchServiceChange(){
	grpcManager.Option.Log.Info("start WatchServiceChange.")
	for changeService := range grpcManager.Option.ServiceDiscovery.WatchMsg{
		grpcManager.Option.Log.Info(" grpc receive serviceChange :"+changeService.Name + " " +changeService.Ip + " " + changeService.Port + " " +changeService.Action)
		dns := changeService.Name + ":" +changeService.Ip
		myGrpcClient , myGrpcClientOk := grpcManager.ClientList[dns]
		if changeService.Action == "PUT"{
			if myGrpcClientOk {
				for k ,_ := range myGrpcClient.GrpcClientList{
					if k == changeService.Name{
						grpcManager.Option.Log.Error("grpc watch don't repeat")
						return
					}
				}
			}
			grpcManager.GetClient(changeService.Name,changeService.NewIp,changeService.NewPort)
		}else if  changeService.Action == "DELETE" {
			if !myGrpcClientOk{
				grpcManager.Option.Log.Error("grpc watch DELETE failed, not in map")
				return
			}
			_ ,GrpcClientListOk := myGrpcClient.GrpcClientList[changeService.Name]
			if !GrpcClientListOk {
				err := errors.New("myGrpcClient CloseOneService no search ")
				grpcManager.Option.Log.Error(err.Error())
				return
			}
			delete(myGrpcClient.GrpcClientList,changeService.Name)

			if len(grpcManager.ClientList[dns].GrpcClientList) == 0 {
				grpcManager.ClientList[dns].ClientConn.Close()
				delete(grpcManager.ClientList , dns)
			}
			return
		}else{
			grpcManager.Option.Log.Error("grpc watch action error.")
			continue
		}
	}
}
//关闭
func  (grpcManager *GrpcManager) Shutdown(){
	//这里主要就是把tcp 连接 给关闭了

	for _,client:=range grpcManager.ClientList{
		client.ClientConn.Close()
	}

	for _,server :=range grpcManager.ServiceList{
		server.GrpcServer.Stop()
		server.Listen.Close()
	}
}
//创建/注册一个服务
func (grpcManager *GrpcManager) CreateService(serviceName string ,ip string ,port string)(*MyGrpcService,error){
	_,empty := grpcManager.Option.ServiceDiscovery.option.ServiceManager.GetByName(serviceName)
	if empty{
		errrorMsg := "serviceName not in serviceManager(db)."+ serviceName
		return nil,errors.New(errrorMsg)
	}
	dns := ip +":"+port
	grpcManager.Option.Log.Debug("GetServer serviceName:"+serviceName + " dns:" + dns)
	//创建一个TCP 连接
	listen, err := net.Listen("tcp", dns)
	if err != nil {
		MyPrint("failed to listen:",err.Error())
		return nil,errors.New("failed to listen:"+ err.Error())
	}
	//实例化一个自己的grpc service
	myGrpcService := MyGrpcService{
		ServiceName	: serviceName,
		ProjectId 	: grpcManager.Option.ProjectId,
		ListenIp	: ip,
		OutIp		: ip,
		Port		: port,
		Listen		: listen,
		Log			: grpcManager.Option.Log,
	}

	var opts []grpc.ServerOption//grpc为使用的第三方的grpc包
	//设定GRPC 公共拦截器
	opts = append(opts, grpc.UnaryInterceptor(myGrpcService.serverInterceptorBack))
	grpcInc := grpc.NewServer(opts...) //创建一个grpc 实例

	myGrpcService.GrpcServer = grpcInc
	grpcManager.ServiceList[serviceName] = &myGrpcService

	return &myGrpcService,nil
}
//获取一个grpc client 连接，同时支持自动负载
func (grpcManager *GrpcManager)GetClientByLoadBalance(serviceName string, balanceFactor string)(inter interface{},err error){
	serviceNode ,err :=  grpcManager.Option.ServiceDiscovery.GetLoadBalanceServiceNodeByServiceName(serviceName,balanceFactor)
	if err != nil{
		return inter,err
	}

	grpcManager.Option.Log.Info("serviceNode LoadBalance , ip:"+serviceNode.Ip + " , port:"+serviceNode.Port)

	grpcClientConn, err := grpcManager.GetClient(serviceName,serviceNode.Ip,serviceNode.Port)
	return grpcClientConn, err
}
//获取一个grpc client 连接
func (grpcManager *GrpcManager) GetClient(serviceName string,ip string,port string)(interface{},error){
	serviceInfo,empty := grpcManager.Option.ServiceDiscovery.option.ServiceManager.GetByName(serviceName)
	if empty{
		errrorMsg := "serviceName not in serviceManager(db)."+ serviceName
		return nil,errors.New(errrorMsg)
	}

	dns := ip + ":" + port
	myGrpcClient ,ok := grpcManager.ClientList[dns]
	if ok {
		client ,ok := myGrpcClient.GrpcClientList[serviceName]
		if ok {
			return client,nil
		}
		err := myGrpcClient.MountClientToConnect(serviceName)
		//grpcManager.Option.Log.Info(" use has exist inc :"+dns)
		return client,err
	}

	myClient := MyGrpcClient{
		ServiceName: serviceName,
		ProjectId: serviceInfo.Id,
		Port: port,
		Ip: ip,
		Log: grpcManager.Option.Log,
		GrpcClientList: make(map[string]interface{}),
	}

	conn, err := grpc.Dial(dns,grpc.WithInsecure(),grpc.WithUnaryInterceptor(myClient.clientInterceptorBack))
	//conn, err := grpc.Dial(dns,grpc.WithInsecure())
	if err != nil {
		MyPrint("did not connect: %v", err)
		return nil,errors.New("did not connect: %v"+err.Error())
	}

	myClient.ClientConn = conn
	grpcManager.ClientList[dns] = &myClient

	err = myClient.MountClientToConnect(serviceName)
    //grpcManager.GrpcClientList[serviceName] = myClient.MountClientToConnect(serviceName)
	return conn,err
}
//server端接收拦截器
func  (myGrpcService *MyGrpcService) serverInterceptorBack(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error){
	fmt.Println("gRPC serverInterceptorBack ", ctx, req , info)

	md ,ok := metadata.FromIncomingContext(ctx)

	traceId := MakeTraceId()

	serverHeader := ServiceServerHeader{
		ProjectId : strconv.Itoa( myGrpcService.ProjectId),
		ServiceName : myGrpcService.ServiceName,
		ReceiveTime: strconv.FormatInt( time.Now().Unix(),10),
		RequestId: MakeRequestId(),
		Protocol: "grpc",
	}
	if !ok{
		myGrpcService.Log.Info("grpc server receive  metadata empty")
	}else{
		clientHeader := ServiceClientHeader{}
		err = json.Unmarshal([]byte(md[SERVICE_HEADER_KEY][0]),&clientHeader)
		if err != nil{
			myGrpcService.Log.Error("json.Unmarshal err:" + err.Error())
		}
		MyPrint("grpc server receive clientHeader:",clientHeader)

		if clientHeader.TraceId != ""{
			traceId = clientHeader.TraceId
		}

		serverHeader.TraceId = traceId
		serverHeader.TargetProjectId = clientHeader.ProjectId
	}
	resp, err = handler(ctx, req)
	serverHeader.ResponseTime = strconv.FormatInt( time.Now().Unix(),10)
	headerInfo := make(map[string]string)
	jsonServerHeader,_ := json.Marshal(serverHeader)
	headerInfo[SERVICE_HEADER_KEY]= string(jsonServerHeader)
	mdData := metadata.New(headerInfo)

	MyPrint("serverHeader:",headerInfo)

	ctx = metadata.NewOutgoingContext(ctx, mdData)
	grpc.SendHeader(ctx,mdData)

	return resp, err
}
//客户端请求拦截器
func (myGrpcClient *MyGrpcClient) clientInterceptorBack(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error{
	MyPrint("grpc clientInterceptorBack ,  method",method, " req:" , req)
	var header  ,trailer metadata.MD

	opts = []grpc.CallOption{grpc.Header(&header),grpc.Trailer(&trailer)}

	clientHeader := ServiceClientHeader{
		TraceId: MakeTraceId(),
		RequestTime: strconv.FormatInt(GetNowTimeSecondToInt64(),10),
		ProjectId: strconv.Itoa(myGrpcClient.ProjectId),
		Protocol: "grpc",
		RequestId: MakeRequestId(),
		TargetServiceName:  myGrpcClient.ServiceName,
	}
	clientHeaderJson,err := json.Marshal(clientHeader)
	MyPrint("json.Marshal:",err)
	headerInfo := make(map[string]string)
	headerInfo[SERVICE_HEADER_KEY] = string( clientHeaderJson)

	MyPrint("grpc client ready header:",headerInfo)

	md := metadata.New(headerInfo)
	ctx = metadata.NewOutgoingContext(ctx, md)

	err = invoker(ctx,method,req,reply,cc,opts...)

	MyPrint("clientInterceptorBack receive , method:",method,"req:",req,"reply:",reply,"opts:",opts,  " cc:",cc.GetState() , " err:",err)
	serviceResponseHeaderDiyString  :=  header.Get(SERVICE_HEADER_KEY)[0]
	serverHeader := ServiceServerHeader{}
	json.Unmarshal([]byte(serviceResponseHeaderDiyString),&serverHeader)
	MyPrint(serverHeader)


	return nil
}

