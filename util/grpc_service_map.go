package util

import (
	"zgoframe/protobuf/pb"
)

func (grpcManager *GrpcManager)GetFrameSyncClient(name string,balanceFactor string)(pb.FrameSyncClient,error){
	client, err := grpcManager.GetClientByLoadBalance(name,balanceFactor)
	if err != nil{
		return nil,err
	}

	return client.(pb.FrameSyncClient),nil
}
func (grpcManager *GrpcManager)GetGatewayClient(name string,balanceFactor string)(pb.GatewayClient,error){
	client, err := grpcManager.GetClientByLoadBalance(name,balanceFactor)
	if err != nil{
		return nil,err
	}

	return client.(pb.GatewayClient),nil
}
func (grpcManager *GrpcManager)GetZgoframeClient(serviceName string,balanceFactor string)(pb.ZgoframeClient,error){
	client, err := grpcManager.GetClientByLoadBalance(serviceName,balanceFactor)
	if err != nil{
		return nil,err
	}

	return client.(pb.ZgoframeClient),nil
}