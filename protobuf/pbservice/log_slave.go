package pbservice

import (
	"context"
	"fmt"
	"zgoframe/protobuf/pb"
)
type LogSlave struct{}


func (logSlave *LogSlave)Push(ctx context.Context,slavePushMsg *pb.SlavePushMsg) (*pb.Empty,error){
	fmt.Println("grpc service LogSlave received:",slavePushMsg)
    empty := &pb.Empty{}
    return empty,nil
}

