package pbservice

import (
	"context"
	"zgoframe/protobuf/pb"
)
type First struct{}


func (first *First)SayHello(ctx context.Context,requestRegPlayer *pb.RequestRegPlayer) (*pb.ResponseReg,error){
    responseReg := &pb.ResponseReg{}
    return responseReg,nil
}
func (first *First)SayHi(ctx context.Context,requestRegPlayer *pb.RequestRegPlayer) (*pb.ResponseReg,error){
    responseReg := &pb.ResponseReg{}
    return responseReg,nil
}

type Second struct{}


func (second *Second)SayHello2(ctx context.Context,requestRegPlayer *pb.RequestRegPlayer) (*pb.ResponseReg,error){
    responseReg := &pb.ResponseReg{}
    return responseReg,nil
}
func (second *Second)SayHi2(ctx context.Context,requestRegPlayer *pb.RequestRegPlayer) (*pb.ResponseReg,error){
    responseReg := &pb.ResponseReg{}
    return responseReg,nil
}

