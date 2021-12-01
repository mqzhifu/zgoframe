package util

import (
	"context"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
	"time"
	"zgoframe/protobuf/pb"
)

type MyGrpc struct {
	Option GrpcOption
	Listen  net.Listener
	//ServerStartUp int	//标识：服务器已启动，避免重复启动
	//ClientStartUp int
}

type GrpcOption struct {
	AppId 		int
	ListenIp	string
	OutIp		string
	Port 		string
	Log 		*zap.Logger
}
//拦截器
//var serverInterceptor grpc.UnaryServerInterceptor
//var clientInterceptor grpc.UnaryClientInterceptor

func NewMyGrpc(grpcOption GrpcOption)(*MyGrpc,error){
	//这里其实除了初始化变量外，只是创建一个TCP SOCKET，给后面的GRPC用

	myGrpc := new(MyGrpc)
	myGrpc.Option = grpcOption

	dns := myGrpc.GetDns()
	myGrpc.Option.Log.Info("grpc GetServer:"+dns)
	listen, err := net.Listen("tcp", dns)
	if err != nil {
		MyPrint("failed to listen:",err.Error())
		return nil,errors.New("failed to listen:"+ err.Error())
	}
	myGrpc.Listen = listen

	return myGrpc,nil
}

func (myGrpc *MyGrpc)GetDns()string{
	dns := myGrpc.Option.ListenIp + ":" + myGrpc.Option.Port
	//MyPrint("dns:"+dns)
	return dns
}

func  (myGrpc *MyGrpc)StartServer(grpcInc *grpc.Server,listen net.Listener){
	reflection.Register(grpcInc)
	err := grpcInc.Serve(listen)
	if  err != nil{
		myGrpc.Option.Log.Info("failed to serve:" + err.Error())
	}
}
func (myGrpc *MyGrpc)Shutdown(){
	myGrpc.Listen.Close()
}
func (myGrpc *MyGrpc)GetServer()(*grpc.Server,net.Listener,error){
	var opts []grpc.ServerOption//grpc为使用的第三方的grpc包
	//opts = append(opts, grpc.UnaryInterceptor(serverInterceptorBack))
	grpcInc := grpc.NewServer(opts...) //创建一个grpc 实例
	return grpcInc,myGrpc.Listen,nil
}

func m1(){

}

func m2(){

}

func serverInterceptorBack(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error){
	fmt.Println("gRPC method: %s, %v", info.FullMethod, req)

	md ,ok := metadata.FromIncomingContext(ctx)
	common := pb.Common{
		ServerReceiveTime: time.Now().Unix(),
	}
	if !ok{
		MyPrint("metadata.FromIncomingContext err:",ok)
	}else{
		fmt.Println(md)
		//clientReqTime, _ := strconv.ParseInt( md["client_req_time"][0], 10, 64)
		//common.RequestId = md["request_id"][0]
		//common.TraceId = md["trace_id"][0]
		//common.ClientReqTime = clientReqTime
	}

	//fmt.Println(md["trace_id"])
	//rid,_ := strconv.Atoi( )
	resp, err = handler(ctx, req)
	common.ServerResponseTime = time.Now().Unix()

	ServerReceiveTimeString := strconv.FormatInt(common.ServerReceiveTime,10)
	ServerResponseTimeString := strconv.FormatInt(common.ServerResponseTime,10)
	md = metadata.Pairs(
		"trace_id", md["trace_id"][0],
		"request_id",md["request_id"][0],
		"client_req_time",md["client_req_time"][0],
		"app_id",md["app_id"][0],
		"server_receive_time",ServerReceiveTimeString,
		"server_response_time",ServerResponseTimeString,
	)
	//ctx = metadata.NewOutgoingContext(ctx, md)
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

func (myGrpc *MyGrpc)clientInterceptorBack(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error{
	//MyPrint("req:",req,"reply:",reply,"opts:",opts)
	var header  metadata.MD

	opts = []grpc.CallOption{grpc.Header(&header)}

	nowString:=strconv.FormatInt(GetNowTimeSecondToInt64(),10)
	md := metadata.Pairs("trace_id", "tt111","request_id",MakeRequestId(),"client_req_time",nowString,"app_id",strconv.Itoa(myGrpc.Option.AppId))
	ctx = metadata.NewOutgoingContext(ctx, md)

	invoker(ctx,method,req,reply,cc,opts...)

	md ,ok := metadata.FromIncomingContext(ctx)
	MyPrint("md:",md,ok)
	MyPrint("method:",method,"req:",req,"reply:",reply,"opts:",header)


	return nil
}

func (myGrpc *MyGrpc) GetClient(ip string,port string)(*grpc.ClientConn,error){
	dns := ip + ":" + port
	conn, err := grpc.Dial(dns,grpc.WithInsecure(),grpc.WithUnaryInterceptor(myGrpc.clientInterceptorBack))
	//conn, err := grpc.Dial(dns,grpc.WithInsecure())
	if err != nil {
		MyPrint("did not connect: %v", err)
		return nil,errors.New("did not connect: %v"+err.Error())
	}

	return conn,nil
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
