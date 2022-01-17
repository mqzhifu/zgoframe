package util

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
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
//启动一个grpc server ，阻塞模式(s端使用)
func (myGrpcServer *MyGrpcService)ServerStart()error{
	return myGrpcServer.GrpcServer.Serve(myGrpcServer.Listen)
}
//将一个服务挂载到一个grpc连接上(c端使用)
func  (myGrpcClient *MyGrpcClient)MountClientToConnect(serviceName string)error{
	client,err:= myGrpcClient.GetGrpcClientByServiceName(serviceName,myGrpcClient.ClientConn)
	if err != nil{
		return err
	}
	myGrpcClient.GrpcClientList[serviceName] = client
	return nil
}

//func  (myGrpcClient *MyGrpcClient)CloseOneService(serviceName string)error{
//	_ ,ok := myGrpcClient.GrpcClientList[serviceName]
//	if ok {
//		delete(myGrpcClient.GrpcClientList,serviceName)
//		return nil
//	}
//	return errors.New("myGrpcClient CloseOneService no search ")
//}

//以下都是动态脚本生成的了=====================================================================


// //根据服务名获取一个GRPC-CLIENT 连接(c端使用)
// func  (myGrpcClient *MyGrpcClient) GetGrpcClientByServiceName(serviceName string,clientConn *grpc.ClientConn)interface{}{
// 	var incClient interface{}
// 	switch serviceName {
// 	case "FrameSync":
// 		incClient = pb.NewFrameSyncClient(clientConn)
// 	case "Gateway":
// 		incClient = pb.NewGatewayClient(clientConn)
// 	case "Zgoframe":
// 		incClient = pb.NewZgoframeClient(clientConn)
// 	}
// 	return incClient
// }
// //动态调用一个GRPC-SERVER 的一个方法(c端使用)
// func (grpcManager *GrpcManager) CallGrpc(serviceName string,funcName string,balanceFactor string,requestData []byte)( resData interface{},err error){
// 	switch serviceName {
// 	case "Zgoframe":
// 		resData , err = grpcManager.CallServiceFuncZgoframe(funcName,balanceFactor,requestData)
// 	}
//
// 	return resData,err
// }
//
//
// func (grpcManager *GrpcManager) CallServiceFuncZgoframe(funcName string,balanceFactor string,postData []byte)( data interface{},err error){
// 	//获取GRPC一个连接
// 	grpcClient,err := grpcManager.GetZgoframeClient("zgoframe",balanceFactor)
// 	if err != nil{
// 		return data,err
// 	}
//
// 	ctx := context.Background()
// 	switch funcName {
// 	case "FrameSync":
// 		requestUser := pb.RequestUser{}
// 		err := json.Unmarshal(postData,&requestUser)
// 		if err != nil{
// 			return data,err
// 		}
// 		data ,err = grpcClient.SayHello(ctx,&requestUser)
// 	}
//
// 	return data,err
// }
//
// //以下均是快捷方法，快速获取一个grpc连接的client.(如果return interface 就不需要下面这些方法了，只是方便调用)
// func (grpcManager *GrpcManager)GetFrameSyncClient(name string,balanceFactor string)(pb.FrameSyncClient,error){
// 	client, err := grpcManager.GetClientByLoadBalance(name,balanceFactor)
// 	if err != nil{
// 		return nil,err
// 	}
//
// 	return client.(pb.FrameSyncClient),nil
// }
// func (grpcManager *GrpcManager)GetGatewayClient(name string,balanceFactor string)(pb.GatewayClient,error){
// 	client, err := grpcManager.GetClientByLoadBalance(name,balanceFactor)
// 	if err != nil{
// 		return nil,err
// 	}
//
// 	return client.(pb.GatewayClient),nil
// }
// func (grpcManager *GrpcManager)GetZgoframeClient(serviceName string,balanceFactor string)(pb.ZgoframeClient,error){
// 	client, err := grpcManager.GetClientByLoadBalance(serviceName,balanceFactor)
// 	if err != nil{
// 		return nil,err
// 	}
//
// 	return client.(pb.ZgoframeClient),nil
// }

