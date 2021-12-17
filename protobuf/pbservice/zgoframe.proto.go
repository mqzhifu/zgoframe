package pbservice

import (
	"context"
	"fmt"
	"zgoframe/protobuf/pb"
)
type Zgoframe struct{}


func (zgoframe *Zgoframe)SayHello(ctx context.Context,requestUser *pb.RequestUser) (*pb.ResponseUser,error){
	fmt.Println("service Zgoframe SayHello:",requestUser)
    responseUser := &pb.ResponseUser{}
    return responseUser,nil
}

