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



type MyGrpcClient struct {
	ServiceName string
	AppId 		int
	ServiceId 	int
	Ip 			string
	Port 		string
	Listen  	net.Listener
	Log 		*zap.Logger
	ClientConn 	*grpc.ClientConn
}

type MyGrpcServer struct {
	ServiceName string
	AppId 		int
	ServiceId 	int
	ListenIp	string
	OutIp		string
	Port 		string
	Log 		*zap.Logger

	Listen  	net.Listener
	GrpcServer *grpc.Server
}

type GrpcManager struct {
	Option 		GrpcManagerOption
	ClientList 	map[string]*MyGrpcClient
	ServiceList map[string]*MyGrpcServer
	//ServerStartUp int	//标识：服务器已启动，避免重复启动
	//ClientStartUp int
}

type GrpcManagerOption struct {
	AppId 				int
	ServiceId 			int
	Log 				*zap.Logger
	ServiceDiscovery *ServiceDiscovery
}

func NewGrpcManager(grpcManagerOption GrpcManagerOption)(*GrpcManager,error){
	//AppId int,Log *zap.Logger,serviceId int

	grpcManager := new(GrpcManager)
	grpcManager.Option = grpcManagerOption
	//grpcManager.Log = Log
	//grpcManager.AppId = AppId
	//grpcManager.ServiceId = serviceId

	//存储自己创建的服务
	grpcManager.ServiceList = make(map[string]*MyGrpcServer)
	//存储它人创建的服务
	grpcManager.ClientList = make(map[string]*MyGrpcClient)

	return grpcManager,nil
}

func  (grpcManager *GrpcManager) Shutdown(){
	for _,client:=range grpcManager.ClientList{
		client.ClientConn.Close()
	}

	for _,server :=range grpcManager.ServiceList{
		server.GrpcServer.Stop()
		server.Listen.Close()
	}
}

func (myGrpcServer *MyGrpcServer)ServerStart(){
	go myGrpcServer.GrpcServer.Serve(myGrpcServer.Listen)
}

func (grpcManager *GrpcManager) GetServer(serviceName string ,ip string ,port string)(*MyGrpcServer,error){
	dns := ip +":"+port
	grpcManager.Option.Log.Debug("GetServer serviceName:"+serviceName + " dns:" + dns)

	listen, err := net.Listen("tcp", dns)
	if err != nil {
		MyPrint("failed to listen:",err.Error())
		return nil,errors.New("failed to listen:"+ err.Error())
	}

	myServer := MyGrpcServer{
		ServiceName	: serviceName,
		AppId		: grpcManager.Option.AppId,
		ServiceId 	: grpcManager.Option.ServiceId,
		ListenIp	: ip,
		OutIp		: ip,
		Port		: port,
		Listen		: listen,
		Log			: grpcManager.Option.Log,
	}

	var opts []grpc.ServerOption//grpc为使用的第三方的grpc包
	opts = append(opts, grpc.UnaryInterceptor(myServer.serverInterceptorBack))
	grpcInc := grpc.NewServer(opts...) //创建一个grpc 实例

	myServer.GrpcServer = grpcInc
	grpcManager.ServiceList[serviceName] = &myServer

	return &myServer,nil
}
func (grpcManager *GrpcManager)GetClientByLoadBalance(serviceName string, appId int)(clientConn *grpc.ClientConn,err error){
	serviceNode ,err :=  grpcManager.Option.ServiceDiscovery.GetLoadBalanceServiceNodeByServiceName(serviceName,"")
	if err != nil{
		return clientConn,err
	}

	grpcManager.Option.Log.Info("serviceNode LoadBalance , ip:"+serviceNode.Ip + " , port:"+serviceNode.Port)

	grpcClientConn, err := grpcManager.GetClient(serviceName,appId,serviceNode.Ip,serviceNode.Port)
	return grpcClientConn, err
}

func (grpcManager *GrpcManager) GetClient(serviceName string, appId int,ip string,port string)(*grpc.ClientConn,error){
	dns := ip + ":" + port

	client ,ok := grpcManager.ClientList[dns]
	if ok {
		return client.ClientConn,nil
	}

	myClient := MyGrpcClient{
		ServiceName: serviceName,
		AppId: appId,
		Port: port,
		Ip: ip,
		Log: grpcManager.Option.Log,
	}

	conn, err := grpc.Dial(dns,grpc.WithInsecure(),grpc.WithUnaryInterceptor(myClient.clientInterceptorBack))
	//conn, err := grpc.Dial(dns,grpc.WithInsecure())
	if err != nil {
		MyPrint("did not connect: %v", err)
		return nil,errors.New("did not connect: %v"+err.Error())
	}

	myClient.ClientConn = conn

	//grpcManager.ClientList[serviceName] = &myClient

	return conn,nil
}

func  (myGrpcServer *MyGrpcServer) serverInterceptorBack(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error){
	fmt.Println("gRPC serverInterceptorBack ", ctx, req , info)

	md ,ok := metadata.FromIncomingContext(ctx)
	//common := pb.Common{
	//	ServerReceiveTime: time.Now().Unix(),
	//}
	//fmt.Println("server md:",md)
	serverHeader := ServiceServerHeader{
		AppId : strconv.Itoa( myGrpcServer.AppId),
		ServiceName : myGrpcServer.ServiceName,
		ReceiveTime: strconv.FormatInt( time.Now().Unix(),10),
		RequestId: MakeRequestId(),
		Protocol: "grpc",
	}
	if !ok{
		myGrpcServer.Log.Info("grpc server receive  metadata empty")
	}else{
		//clientReqTime, _ := strconv.ParseInt( md["client_req_time"][0], 10, 64)
		//clientAppId ,_:= strconv.Atoi(md["app_id"][0])
		//clientHeader :=ClientHeader{
		//	RequestTime: clientReqTime,
		//	TraceId: md["trace_id"][0],
		//	RequestId: md["request_id"][0],
		//	TargetServiceName: md["service_name"][0],
		//	AppId:clientAppId ,
		//	Protocol: md["protocol"][0],
		//}
		clientHeader := ServiceClientHeader{}
		err = json.Unmarshal([]byte(md[SERVICE_HEADER_KEY][0]),&clientHeader)
		MyPrint("json.Unmarshal err:",err)

		MyPrint("grpc server receive clientHeader:",clientHeader)

		serverHeader.TraceId = clientHeader.TraceId
		serverHeader.TargetAppId = clientHeader.AppId

		//common.RequestId = md["request_id"][0]
		//common.TraceId = md["trace_id"][0]
		//common.ClientReqTime = clientReqTime
	}

	//fmt.Println(md["trace_id"])
	//rid,_ := strconv.Atoi( )
	resp, err = handler(ctx, req)
	//common.ServerResponseTime = time.Now().Unix()

	serverHeader.ResponseTime = strconv.FormatInt( time.Now().Unix(),10)

	//serverHeader.ResponseTime = time.Now().Unix()
	////ServerReceiveTimeString := strconv.FormatInt(common.ServerReceiveTime,10)
	////ServerResponseTimeString := strconv.FormatInt(common.ServerResponseTime,10)
	headerInfo := make(map[string]string)
	////headerInfo["aaaaa"] = "bbbb"
	//headerInfo["trace_id"] = serverHeader.TraceId
	//headerInfo["receive_time"] = strconv.FormatInt(serverHeader.ReceiveTime,10)
	//headerInfo["response_time"] = strconv.FormatInt(serverHeader.ResponseTime,10)
	//headerInfo["request_id"] = serverHeader.RequestId
	//headerInfo["app_id"] = strconv.Itoa(serverHeader.AppId)
	//headerInfo["protocol"] = "grpc"
	//headerInfo["service_name"] = serverHeader.ServiceName
	//headerInfo["target_app_id"] = strconv.Itoa(serverHeader.TargetAppId)
	jsonServerHeader,_ := json.Marshal(serverHeader)
	headerInfo[SERVICE_HEADER_KEY]= string(jsonServerHeader)
	mdData := metadata.New(headerInfo)
	//md = metadata.Pairs(
	//	//"trace_id", md["trace_id"][0],
	//	//"request_id",md["request_id"][0],
	//	//"client_req_time",md["client_req_time"][0],
	//	//"app_id",md["app_id"][0],
	//	//"server_receive_time",ServerReceiveTimeString,
	//	//"server_response_time",ServerResponseTimeString,
	//)

	MyPrint("serverHeader:",headerInfo)

	ctx = metadata.NewOutgoingContext(ctx, mdData)
	//grpc.SetHeader(ctx,md)
	//grpc.SetTrailer(ctx,md)
	grpc.SendHeader(ctx,mdData)

	return resp, err
}





func (myGrpcClient *MyGrpcClient) clientInterceptorBack(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error{
	MyPrint("grpc clientInterceptorBack ,  method",method, " req:" , req)
	//MyPrint("req:",req,"reply:",reply,"opts:",opts)
	var header  ,trailer metadata.MD

	opts = []grpc.CallOption{grpc.Header(&header),grpc.Trailer(&trailer)}

	//md := metadata.Pairs(
	//	"trace_id", traceId,"request_id",MakeRequestId(),"client_req_time",nowString,"app_id",strconv.Itoa(myGrpcClient.AppId),"protocol","grpc")

	clientHeader := ServiceClientHeader{
		TraceId: MakeTraceId(),
		RequestTime: strconv.FormatInt(GetNowTimeSecondToInt64(),10),
		AppId: strconv.Itoa(myGrpcClient.AppId),
		Protocol: "grpc",
		RequestId: MakeRequestId(),
		TargetServiceName:  myGrpcClient.ServiceName,
	}
	clientHeaderJson,err := json.Marshal(clientHeader)
	MyPrint("json.Marshal:",err)
	headerInfo := make(map[string]string)
	headerInfo[SERVICE_HEADER_KEY] = string( clientHeaderJson)
	//headerInfo["trace_id"] = MakeTraceId()
	//headerInfo["client_req_time"] = strconv.FormatInt(GetNowTimeSecondToInt64(),10)
	//headerInfo["request_id"] = MakeRequestId()
	//headerInfo["app_id"] = strconv.Itoa(myGrpcClient.AppId)
	//headerInfo["protocol"] = "grpc"
	//headerInfo["service_name"] = myGrpcClient.ServiceName

	MyPrint("grpc client ready header:",headerInfo)

	md := metadata.New(headerInfo)
	ctx = metadata.NewOutgoingContext(ctx, md)

	invoker(ctx,method,req,reply,cc,opts...)

	//md ,ok := metadata.FromIncomingContext(ctx)

	//MyPrint("md:",md,ok)
	//MyPrint("trailer:",trailer)

	//appId ,_ := strconv.Atoi( header.Get("app_id")[0])
	//serverHeader := ServerHeader{
	//	AppId : appId,
	//	ServiceName : header.Get("service_name")[0],
	//	ReceiveTime: time.Now().Unix(),
	//	RequestId: MakeRequestId(),
	//	Protocol: "grpc",
	//}
	serviceResponseHeaderDiyString  :=  header.Get(SERVICE_HEADER_KEY)[0]
	serverHeader := ServiceServerHeader{}
	json.Unmarshal([]byte(serviceResponseHeaderDiyString),&serverHeader)
	MyPrint("clientInterceptorBack receive , method:",method,"req:",req,"reply:",reply,"opts:",serverHeader)


	return nil
}

