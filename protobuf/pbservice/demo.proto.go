package pbservice

import (
	"context"
	"zgoframe/protobuf/pb"
)
type First struct{}


func (first *First)SayHello(ctx context.Context,requestRegPlayer *pb.RequestRegPlayer) (*pb.ResponseReg,error){

    responseReg := &pb.ResponseReg{
    	Rs: true,
	}

    return responseReg,nil
}
func (first *First)SayHi(ctx context.Context,requestRegPlayer *pb.RequestRegPlayer) (*pb.ResponseReg,error){
    responseReg := &pb.ResponseReg{}
    return responseReg,nil
}

