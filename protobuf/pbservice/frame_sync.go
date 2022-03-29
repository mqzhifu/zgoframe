package pbservice

import (
	"context"
	"zgoframe/protobuf/pb"
)
type FrameSync struct{}


func (frameSync *FrameSync)CS_PlayerOperations(ctx context.Context,logicFrame *pb.LogicFrame) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)CS_PlayerResumeGame(ctx context.Context,playerResumeGame *pb.PlayerResumeGame) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)CS_PlayerReady(ctx context.Context,playerReady *pb.PlayerReady) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)CS_PlayerOver(ctx context.Context,playerOver *pb.PlayerOver) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)CS_RoomHistory(ctx context.Context,roomHistory *pb.RoomHistory) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)CS_RoomBaseInfo(ctx context.Context,roomBaseInfo *pb.RoomBaseInfo) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)CS_PlayerMatchSign(ctx context.Context,playerMatchSign *pb.PlayerMatchSign) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)CS_PlayerMatchSignCancel(ctx context.Context,playerMatchSignCancel *pb.PlayerMatchSignCancel) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)SC_PlayerMatchSignFailed(ctx context.Context,playerMatchSignFailed *pb.PlayerMatchSignFailed) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)SC_PlayerMatchingFailed(ctx context.Context,playerMatchingFailed *pb.PlayerMatchingFailed) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)SC_EnterBattle(ctx context.Context,enterBattle *pb.EnterBattle) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)SC_LogicFrame(ctx context.Context,logicFrame *pb.LogicFrame) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)SC_RoomHistory(ctx context.Context,roomHistory *pb.RoomHistory) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)SC_RoomBaseInfo(ctx context.Context,roomBaseInfo *pb.RoomBaseInfo) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)SC_OtherPlayerOffline(ctx context.Context,otherPlayerOffline *pb.OtherPlayerOffline) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)SC_OtherPlayerOver(ctx context.Context,playerOver *pb.PlayerOver) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)SC_OtherPlayerResumeGame(ctx context.Context,playerResumeGame *pb.PlayerResumeGame) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)SC_ReadyTimeout(ctx context.Context,readyTimeout *pb.ReadyTimeout) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)SC_StartBattle(ctx context.Context,startBattle *pb.StartBattle) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)SC_RestartGame(ctx context.Context,restartGame *pb.RestartGame) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)SC_GameOver(ctx context.Context,gameOver *pb.GameOver) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}

