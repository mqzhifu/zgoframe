package util

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"zgoframe/protobuf/pb"
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
	//GrpcClient  interface{}
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


//
func  (myGrpcClient *MyGrpcClient)MountClientToConnect(serviceName string){
	var grpcClient interface{}
	switch serviceName {
	case "Zgoframe":
		grpcClient = pb.NewZgoframeClient(myGrpcClient.ClientConn)
	case "Sync":
		//grpcClient = pb.NewSyncClient(myGrpcClient.ClientConn)
	}

	myGrpcClient.GrpcClientList[serviceName] = grpcClient
}

//#client#start

//#client#end