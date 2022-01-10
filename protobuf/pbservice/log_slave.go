package pbservice

import (
	"context"
	"zgoframe/protobuf/pb"
)
type LogSlave struct{}


func (logSlave *LogSlave)Push(ctx context.Context,slavePushMsg *pb.SlavePushMsg) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}

