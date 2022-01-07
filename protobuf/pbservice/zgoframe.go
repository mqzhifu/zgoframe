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

