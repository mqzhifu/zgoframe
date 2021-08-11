package util

import (
	"context"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
	"time"
	"zgoframe/protobuf/pb"
	"zgoframe/protobuf/pbservice"
)


const(
	ROLE_CLIENT = 1
	ROLE_SERVER = 2
)
type MyGrpc struct {
	Option GrpcOption
	ServerStartUp int
	ClientStartUp int
}

type GrpcOption struct {
	AppId 		int
	ListenIp	string
	OutIp		string
	Port 		string
	Role 		int
}

var serverInterceptor grpc.UnaryServerInterceptor
var clientInterceptor grpc.UnaryClientInterceptor
func NewMyGrpc(grpcOption GrpcOption)*MyGrpc{
	myGrpc := new(MyGrpc)
	myGrpc.Option = grpcOption
	return myGrpc
}

func (myGrpc *MyGrpc)GetDns()string{
	dns := myGrpc.Option.ListenIp + ":" + myGrpc.Option.Port
	MyPrint("dns:"+dns)
	return dns
}

func (myGrpc *MyGrpc)GetServer()error{
	if myGrpc.ServerStartUp == 1{
		return errors.New("server has start up...")
	}
	dns := myGrpc.GetDns()
	serverInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error){
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

	lis, err := net.Listen("tcp", dns)
	if err != nil {
		MyPrint("failed to listen:",err.Error())
		return errors.New("failed to listen:"+ err.Error())
	}
	defer lis.Close()

	var opts []grpc.ServerOption//grpc为使用的第三方的grpc包
	opts = append(opts, grpc.UnaryInterceptor(serverInterceptor))
	s := grpc.NewServer(opts...) //创建一个grpc 实例
	//挂载服务的handler
	pb.RegisterFirstServer(s, &pbservice.First{})
	// 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		MyPrint("failed to serve:",err.Error())
		return errors.New("failed to serve:"+err.Error())
	}
	return nil
}
func MakeTraceId()string{
	id :=  uuid.NewV4()
	return id.String()
}

func MakeRequestId()string{
	id :=  uuid.NewV4()
	return id.String()
}

func (myGrpc *MyGrpc) GetClient()(pb.FirstClient,error){
	if myGrpc.ClientStartUp == 1{
		return nil,errors.New("client has start up...")
	}


	clientInterceptor = func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error{
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

	dns := myGrpc.GetDns()
	conn, err := grpc.Dial(dns,grpc.WithInsecure(),grpc.WithUnaryInterceptor(clientInterceptor))
	if err != nil {
		MyPrint("did not connect: %v", err)
		return nil,errors.New("did not connect: %v"+err.Error())
	}else{
		MyPrint("tcp conn ok!")
	}

	c := pb.NewFirstClient(conn)



	return c,nil
}

var ip = "127.0.0.1"
var port = "1111"
func TestGrpcServer(){
	grpcOption := GrpcOption{
		ListenIp:ip,
		OutIp: ip,
		Port: port,
		Role: ROLE_SERVER,
	}
	grpcClass := NewMyGrpc(grpcOption)
	grpcClass.GetServer()
}

func TestGrpcClient(){
	grpcOption := GrpcOption{
		AppId: 1,
		ListenIp:ip,
		OutIp: ip,
		Port: port,
		Role: ROLE_SERVER,
	}
	grpcClass := NewMyGrpc(grpcOption)
	c ,err := grpcClass.GetClient()
	if err != nil{
		MyPrint("grpcClass.GetClient() err:",err.Error())
		return
	}
	channel := pb.Channel{Id: 999,Name: "tttt"}

	channelMap := make(map[uint64]*pb.Channel)
	channelMap[333] = &channel
	player := pb.Player{
		Id: 1,
		RoleName: []byte("aaaaaa"),
		Nickname: "xiaoz",
		Status: 1,
		Score: 11.222,
		PhoneType: pb.PhoneType_home,
		Sex: true,
		Level: 10,
	}

	playerList := []*pb.Player{&player}

	requestRegPlayer := pb.RequestRegPlayer{
		AddTime : 11234,
		PlayerInfo: &player,
		PlayerList :playerList,
		ChannelMap :channelMap,
	}

	//if !ok {
	//	ExitPrint("metadata.FromIncomingContext err:",md,ok)
	//}
	r, err := c.SayHello(context.Background(), &requestRegPlayer)


	if err != nil {
		MyPrint("could not service: %v", err)
	}

	MyPrint("FINAL:",r.Rs)

}
