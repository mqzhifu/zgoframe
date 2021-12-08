package util

import (
	"context"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"strconv"
	"time"
)

type ClientHeader struct {
	TraceId string
	RequestId string
	Protocol string
	AppId int
	RequestTime int64
	TargetServiceName string
	ServerReceiveTime int64
	ServerResponseTime int64
}

type ServerHeader struct {
	TraceId string
	RequestId string
	Protocol string
	AppId int
	TargetAppId int
	ServiceName string
	ReceiveTime int64
	ResponseTime int64
}

type MyGrpcClient struct {
	ServiceName string
	AppId int
	Ip 		string
	Port 		string
	Listen  net.Listener
	ClientConn *grpc.ClientConn
}

type MyGrpcServer struct {
	ServiceName string
	AppId 		int
	ListenIp	string
	OutIp		string
	Port 		string
	Listen  net.Listener
	GrpcServer *grpc.Server
}

type GrpcManager struct {
	//Option GrpcOption
	AppId 	int
	Log 		*zap.Logger
	ClientList 	map[string]*MyGrpcClient
	ServiceList map[string]*MyGrpcServer
	//ServerStartUp int	//标识：服务器已启动，避免重复启动
	//ClientStartUp int
}

//type GrpcOption struct {
//	Log 		*zap.Logger
//}

func NewGrpcManager(AppId 	int,Log 		*zap.Logger)(*GrpcManager,error){

	grpcManager := new(GrpcManager)
	grpcManager.Log = Log
	grpcManager.AppId = AppId
	grpcManager.ServiceList = make(map[string]*MyGrpcServer)
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

func (grpcManager *GrpcManager) GetServer(serviceName string ,appId int,ip string ,port string)(*MyGrpcServer,error){
	dns := ip +":"+port
	grpcManager.Log.Debug("GetServer serviceName:"+serviceName + " dns:" + dns)

	listen, err := net.Listen("tcp", dns)
	if err != nil {
		MyPrint("failed to listen:",err.Error())
		return nil,errors.New("failed to listen:"+ err.Error())
	}

	myServer := MyGrpcServer{
		ServiceName: serviceName,
		AppId: grpcManager.AppId,
		ListenIp: ip,
		OutIp: ip,
		Port: port,
		Listen: listen,
	}

	var opts []grpc.ServerOption//grpc为使用的第三方的grpc包
	opts = append(opts, grpc.UnaryInterceptor(myServer.serverInterceptorBack))
	grpcInc := grpc.NewServer(opts...) //创建一个grpc 实例

	myServer.GrpcServer = grpcInc
	grpcManager.ServiceList[serviceName] = &myServer

	return &myServer,nil
}

func  (myGrpcServer *MyGrpcServer) serverInterceptorBack(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error){
	fmt.Println("gRPC method: %s, %v", info.FullMethod, req)

	md ,ok := metadata.FromIncomingContext(ctx)
	//common := pb.Common{
	//	ServerReceiveTime: time.Now().Unix(),
	//}
	//fmt.Println("server md:",md)
	serverHeader := ServerHeader{
		AppId : myGrpcServer.AppId,
		ServiceName : myGrpcServer.ServiceName,
		ReceiveTime: time.Now().Unix(),
		RequestId: MakeRequestId(),
		Protocol: "grpc",
	}
	if !ok{
		MyPrint("grpc server receive  metadata.FromIncomingContext err:",ok)
	}else{
		//fmt.Println(md)
		clientReqTime, _ := strconv.ParseInt( md["client_req_time"][0], 10, 64)
		clientAppId ,_:= strconv.Atoi(md["app_id"][0])
		clientHeader :=ClientHeader{
			RequestTime: clientReqTime,
			TraceId: md["trace_id"][0],
			RequestId: md["request_id"][0],
			TargetServiceName: md["service_name"][0],
			AppId:clientAppId ,
			Protocol: md["protocol"][0],
		}

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

	serverHeader.ResponseTime = time.Now().Unix()
	//ServerReceiveTimeString := strconv.FormatInt(common.ServerReceiveTime,10)
	//ServerResponseTimeString := strconv.FormatInt(common.ServerResponseTime,10)
	headerInfo := make(map[string]string)
	headerInfo["trace_id"] = serverHeader.TraceId
	headerInfo["receive_time"] = strconv.FormatInt(serverHeader.ReceiveTime,10)
	headerInfo["response_time"] = strconv.FormatInt(serverHeader.ResponseTime,10)
	headerInfo["request_id"] = serverHeader.RequestId
	headerInfo["app_id"] = strconv.Itoa(serverHeader.AppId)
	headerInfo["protocol"] = "grpc"
	headerInfo["service_name"] = serverHeader.ServiceName
	headerInfo["target_app_id"] = strconv.Itoa(serverHeader.TargetAppId)

	mdData := metadata.New(headerInfo)
	//md = metadata.Pairs(
	//	//"trace_id", md["trace_id"][0],
	//	//"request_id",md["request_id"][0],
	//	//"client_req_time",md["client_req_time"][0],
	//	//"app_id",md["app_id"][0],
	//	//"server_receive_time",ServerReceiveTimeString,
	//	//"server_response_time",ServerResponseTimeString,
	//)

	MyPrint("serverHeader:",serverHeader)

	ctx = metadata.NewOutgoingContext(ctx, mdData)
	grpc.SendHeader(ctx,md)

	return resp, err
}

func MakeTraceId()string{
	id :=  uuid.NewV4()
	return id.String()
}

func MakeRequestId()string{
	id :=  uuid.NewV4()
	return id.String()
}



func (myGrpcClient *MyGrpcClient) clientInterceptorBack(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error{
	//MyPrint("req:",req,"reply:",reply,"opts:",opts)
	var header  metadata.MD

	opts = []grpc.CallOption{grpc.Header(&header)}

	//md := metadata.Pairs(
	//	"trace_id", traceId,"request_id",MakeRequestId(),"client_req_time",nowString,"app_id",strconv.Itoa(myGrpcClient.AppId),"protocol","grpc")


	headerInfo := make(map[string]string)
	headerInfo["trace_id"] = MakeTraceId()
	headerInfo["client_req_time"] = strconv.FormatInt(GetNowTimeSecondToInt64(),10)
	headerInfo["request_id"] = MakeRequestId()
	headerInfo["app_id"] = strconv.Itoa(myGrpcClient.AppId)
	headerInfo["protocol"] = "grpc"
	headerInfo["service_name"] = myGrpcClient.ServiceName

	MyPrint("grpc client header:",headerInfo)

	md := metadata.New(headerInfo)
	ctx = metadata.NewOutgoingContext(ctx, md)



	invoker(ctx,method,req,reply,cc,opts...)

	md ,ok := metadata.FromIncomingContext(ctx)
	MyPrint("md:",md,ok)
	MyPrint("grpc client receive , method:",method,"req:",req,"reply:",reply,"opts:",header)


	return nil
}

func (grpcManager *GrpcManager) GetClient(serviceName string, appId int,ip string,port string)(*grpc.ClientConn,error){
	dns := ip + ":" + port

	mygrpc := MyGrpcClient{
		ServiceName: serviceName,
		AppId: appId,
		Port: port,
		Ip: ip,
	}

	conn, err := grpc.Dial(dns,grpc.WithInsecure(),grpc.WithUnaryInterceptor(mygrpc.clientInterceptorBack))
	//conn, err := grpc.Dial(dns,grpc.WithInsecure())
	if err != nil {
		MyPrint("did not connect: %v", err)
		return nil,errors.New("did not connect: %v"+err.Error())
	}

	mygrpc.ClientConn = conn

	grpcManager.ClientList[serviceName] = &mygrpc

	return conn,nil
}
