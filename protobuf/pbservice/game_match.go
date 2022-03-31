package pbservice

import (
	"context"
	"zgoframe/protobuf/pb"
)
type GameMatch struct{}


func (gameMatch *GameMatch)CS_PlayerMatchSign(ctx context.Context,playerMatchSign *pb.PlayerMatchSign) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (gameMatch *GameMatch)CS_PlayerMatchSignCancel(ctx context.Context,playerMatchSignCancel *pb.PlayerMatchSignCancel) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (gameMatch *GameMatch)SC_PlayerMatchSignFailed(ctx context.Context,playerMatchSignFailed *pb.PlayerMatchSignFailed) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (gameMatch *GameMatch)SC_PlayerMatchingFailed(ctx context.Context,playerMatchingFailed *pb.PlayerMatchingFailed) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}

