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
	Option 		GrpcManagerOption
	ClientList 	map[string]*MyGrpcClient
	ServiceList map[string]*MyGrpcService
}

type GrpcManagerOption struct {
	ProjectId 			int				//在做链路追踪时，需要有项目ID
	Log 				*zap.Logger
	ServiceDiscovery *ServiceDiscovery	//获取一个grpc连接时，需要使用
}

func NewGrpcManager(grpcManagerOption GrpcManagerOption)(*GrpcManager,error){
	//projectId int,Log *zap.Logger,serviceId int

	grpcManager := new(GrpcManager)
	grpcManager.Option = grpcManagerOption
	//存储自己创建的服务，注册注册
	grpcManager.ServiceList = make(map[string]*MyGrpcService)
	//存储它人创建的服务，调用其它人的服务
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
func (grpcManager *GrpcManager)GetClientByLoadBalance(serviceName string, balanceFactor string)(inter interface{},err error){
	serviceNode ,err :=  grpcManager.Option.ServiceDiscovery.GetLoadBalanceServiceNodeByServiceName(serviceName,balanceFactor)
	if err != nil{
		return inter,err
	}

	grpcManager.Option.Log.Info("serviceNode LoadBalance , ip:"+serviceNode.Ip + " , port:"+serviceNode.Port)

	grpcClientConn, err := grpcManager.GetClient(serviceName,serviceNode.Ip,serviceNode.Port)
	return grpcClientConn, err
}

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
		//clientReqTime, _ := strconv.ParseInt( md["client_req_time"][0], 10, 64)
		//clientprojectId ,_:= strconv.Atoi(md["app_id"][0])
		//clientHeader :=ClientHeader{
		//	RequestTime: clientReqTime,
		//	TraceId: md["trace_id"][0],
		//	RequestId: md["request_id"][0],
		//	TargetServiceName: md["service_name"][0],
		//	projectId:clientprojectId ,
		//	Protocol: md["protocol"][0],
		//}
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
		//common.RequestId = md["request_id"][0]
		//common.TraceId = md["trace_id"][0]
		//common.ClientReqTime = clientReqTime
	}
	//fmt.Println(md["trace_id"])
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
	//headerInfo["app_id"] = strconv.Itoa(serverHeader.projectId)
	//headerInfo["protocol"] = "grpc"
	//headerInfo["service_name"] = serverHeader.ServiceName
	//headerInfo["target_app_id"] = strconv.Itoa(serverHeader.TargetprojectId)
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
	//	"trace_id", traceId,"request_id",MakeRequestId(),"client_req_time",nowString,"app_id",strconv.Itoa(myGrpcClient.projectId),"protocol","grpc")

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
	//headerInfo["trace_id"] = MakeTraceId()
	//headerInfo["client_req_time"] = strconv.FormatInt(GetNowTimeSecondToInt64(),10)
	//headerInfo["request_id"] = MakeRequestId()
	//headerInfo["app_id"] = strconv.Itoa(myGrpcClient.projectId)
	//headerInfo["protocol"] = "grpc"
	//headerInfo["service_name"] = myGrpcClient.ServiceName

	MyPrint("grpc client ready header:",headerInfo)

	md := metadata.New(headerInfo)
	ctx = metadata.NewOutgoingContext(ctx, md)

	err = invoker(ctx,method,req,reply,cc,opts...)

	MyPrint("clientInterceptorBack receive , method:",method,"req:",req,"reply:",reply,"opts:",opts,  " cc:",cc.GetState() , " err:",err)
	//md ,ok := metadata.FromIncomingContext(ctx)

	//MyPrint("md:",md,ok)
	//MyPrint("trailer:",trailer)

	//projectId ,_ := strconv.Atoi( header.Get("app_id")[0])
	//serverHeader := ServerHeader{
	//	projectId : projectId,
	//	ServiceName : header.Get("service_name")[0],
	//	ReceiveTime: time.Now().Unix(),
	//	RequestId: MakeRequestId(),
	//	Protocol: "grpc",
	//}
	serviceResponseHeaderDiyString  :=  header.Get(SERVICE_HEADER_KEY)[0]
	serverHeader := ServiceServerHeader{}
	json.Unmarshal([]byte(serviceResponseHeaderDiyString),&serverHeader)
	MyPrint(serverHeader)


	return nil
}

