package pbservice

import (
	"context"
	"fmt"
	"zgoframe/protobuf/pb"
)
type First struct{}


func (first *First)SayHello(ctx context.Context,requestRegPlayer *pb.RequestRegPlayer) (*pb.ResponseReg,error){
	fmt.Println("SayHello receive :",requestRegPlayer.AddTime)
    responseReg := &pb.ResponseReg{}
	responseReg.Rs = true
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

