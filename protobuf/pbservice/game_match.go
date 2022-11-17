package pbservice

import (
	"context"
	"zgoframe/protobuf/pb"
)
type GameMatch struct{}


func (gameMatch *GameMatch)CS_PlayerMatchSign(ctx context.Context,gameMatchSign *pb.GameMatchSign) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (gameMatch *GameMatch)CS_PlayerMatchSignCancel(ctx context.Context,gameMatchPlayerCancel *pb.GameMatchPlayerCancel) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (gameMatch *GameMatch)SC_GameMatchOptResult(ctx context.Context,gameMatchOptResult *pb.GameMatchOptResult) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}

