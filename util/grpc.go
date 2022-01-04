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
	AppId 		int
	ServiceId 	int
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
	var incClient interface{}
	switch serviceName {
	case "Zgoframe":
		incClient = myGrpcClient.ZgoframeMountToConnect()
	case "Sync":
		incClient = myGrpcClient.SyncMountToConnect()
	}

	//myGrpcClient.GrpcClient = incClient
	myGrpcClient.GrpcClientList[serviceName] = incClient
}

func  (myGrpcClient *MyGrpcClient)ZgoframeMountToConnect()pb.ZgoframeClient{
	serviceClient := pb.NewZgoframeClient(myGrpcClient.ClientConn)
	return serviceClient
}

func  (myGrpcClient *MyGrpcClient)SyncMountToConnect()pb.SyncClient{
	serviceClient := pb.NewSyncClient(myGrpcClient.ClientConn)
	return serviceClient
}





//#client#start

//#client#end