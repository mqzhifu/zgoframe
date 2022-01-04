package pbservice

import (
	"context"
	"zgoframe/protobuf/pb"
)
type Zgoframe struct{}


func (zgoframe *Zgoframe)SayHello(ctx context.Context,requestUser *pb.RequestUser) (*pb.ResponseUser,error){
    responseUser := &pb.ResponseUser{}
    return responseUser,nil
}
func (zgoframe *Zgoframe)Comm(ctx context.Context,requestUser *pb.RequestUser) (*pb.ResponseUser,error){
    responseUser := &pb.ResponseUser{}
    return responseUser,nil
}

type Sync struct{}


func (sync *Sync)Heartbeat(ctx context.Context,cSHeartbeat *pb.CSHeartbeat) (*pb.SCHeartbeat,error){
    sCHeartbeat := &pb.SCHeartbeat{}
    return sCHeartbeat,nil
}

