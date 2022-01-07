package pbservice

import (
	"context"
	"zgoframe/protobuf/pb"
)
type FrameSync struct{}


func (frameSync *FrameSync)PlayerOperations(ctx context.Context,requestPlayerOperations *pb.RequestPlayerOperations) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)PlayerResumeGame(ctx context.Context,requestPlayerResumeGame *pb.RequestPlayerResumeGame) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)PlayerReady(ctx context.Context,requestPlayerReady *pb.RequestPlayerReady) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)PlayerOver(ctx context.Context,requestPlayerOver *pb.RequestPlayerOver) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)RoomHistory(ctx context.Context,requestRoomHistory *pb.RequestRoomHistory) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)GetRoom(ctx context.Context,requestGetRoom *pb.RequestGetRoom) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)PlayerMatchSign(ctx context.Context,requestPlayerMatchSign *pb.RequestPlayerMatchSign) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)PlayerMatchSignCancel(ctx context.Context,requestPlayerMatchSignCancel *pb.RequestPlayerMatchSignCancel) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)EnterBattle(ctx context.Context,responseEnterBattle *pb.ResponseEnterBattle) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)PushLogicFrame(ctx context.Context,responsePushLogicFrame *pb.ResponsePushLogicFrame) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)OtherPlayerOffline(ctx context.Context,responseOtherPlayerOffline *pb.ResponseOtherPlayerOffline) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)ReadyTimeout(ctx context.Context,responseReadyTimeout *pb.ResponseReadyTimeout) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)PushRoomHistory(ctx context.Context,responsePushRoomHistory *pb.ResponsePushRoomHistory) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)GameOver(ctx context.Context,responseGameOver *pb.ResponseGameOver) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)PushRoomInfo(ctx context.Context,responsePushRoomInfo *pb.ResponsePushRoomInfo) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)StartBattle(ctx context.Context,responseStartBattle *pb.ResponseStartBattle) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)OtherPlayerOver(ctx context.Context,responseOtherPlayerOver *pb.ResponseOtherPlayerOver) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)RestartGame(ctx context.Context,responseRestartGame *pb.ResponseRestartGame) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)PlayerMatchSignFailed(ctx context.Context,responsePlayerMatchSignFailed *pb.ResponsePlayerMatchSignFailed) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)PlayerMatchingFailed(ctx context.Context,responsePlayerMatchingFailed *pb.ResponsePlayerMatchingFailed) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (frameSync *FrameSync)OtherPlayerResumeGame(ctx context.Context,responseOtherPlayerResumeGame *pb.ResponseOtherPlayerResumeGame) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}

