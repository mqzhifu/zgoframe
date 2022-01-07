package util

import (
	"zgoframe/protobuf/pb"
)

func (grpcManager *GrpcManager)GetZgoframeClient(name string)(pb.ZgoframeClient,error){
	client, err := grpcManager.GetClientByLoadBalance(name,0)
	if err != nil{
		return nil,err
	}

	return client.(pb.ZgoframeClient),nil
}