package util

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"zgoframe/protobuf/pb"
)

type MyGrpcClient struct {
	ServiceName string
	ProjectId 		int
	ServiceId 	int
	Ip 			string
	Port 		string
	Listen  	net.Listener
	Log 		*zap.Logger
	ClientConn 	*grpc.ClientConn
	//GrpcClient  interface{}
	//一个grpc连接，上面可以挂载若干个服务
	GrpcClientList   map[string]interface{}
}

type MyGrpcService struct {
	ServiceName string
	//AppId 		int
	//ServiceId 	int
	ProjectId 	int
	ListenIp	string
	OutIp		string
	Port 		string
	Log 		*zap.Logger

	Listen  	net.Listener
	GrpcServer *grpc.Server
}

func (myGrpcServer *MyGrpcService)ServerStart(){
	go myGrpcServer.GrpcServer.Serve(myGrpcServer.Listen)
}


//func  (myGrpcClient *MyGrpcClient)MountClientToConnect(serviceName string){
//	var grpcClient interface{}
//	switch serviceName {
//	case "Zgoframe":
//		grpcClient = pb.NewZgoframeClient(myGrpcClient.ClientConn)
//	case "Sync":
//		//grpcClient = pb.NewSyncClient(myGrpcClient.ClientConn)
//	}
//
//	myGrpcClient.GrpcClientList[serviceName] = grpcClient
//}

//将一个服务挂载到一个grpc连接上
func  (myGrpcClient *MyGrpcClient)MountClientToConnect(serviceName string){
	myGrpcClient.GrpcClientList[serviceName] = GetGrpcClientByServiceName(serviceName,myGrpcClient.ClientConn)
}

func GetGrpcClientByServiceName(serviceName string,clientConn *grpc.ClientConn)interface{}{
	var incClient interface{}
	switch serviceName {
	case "FrameSync":
		incClient = pb.NewFrameSyncClient(clientConn)
	case "Gateway":
		incClient = pb.NewGatewayClient(clientConn)
	case "Zgoframe":
		incClient = pb.NewZgoframeClient(clientConn)
	}
	return incClient
}

func(grpcManager *GrpcManager) CallServiceFuncZgoframe(funcName string,balanceFactor string,postData []byte)( data interface{},err error){
	//获取GRPC一个连接
	grpcClient,err := grpcManager.GetZgoframeClient("zgoframe",balanceFactor)
	if err != nil{
		return data,err
	}

	ctx := context.Background()
	switch funcName {
	case "FrameSync":
		requestUser := pb.RequestUser{}
		err := json.Unmarshal(postData,&requestUser)
		if err != nil{
			return data,err
		}
		data ,err = grpcClient.SayHello(ctx,&requestUser)
	}

	return data,err
}

//func  (grpcManager *GrpcManager)CallGrpcResMap(funcName string,data []byte)(responseUser pb.ResponseUser,err error){
//	//switch funcName {
//	//case "FrameSync":
//		//responseUser := pb.ResponseUser{}
//		err = proto.Unmarshal(data,&responseUser)
//		if err != nil{
//			return responseUser,err
//		}
//		//data ,err = grpcClient.SayHello(ctx,&requestUser)
//	//}
//
//	return responseUser,err
//}


//#client#start

//#client#end