package pbservice

import (
	"context"
	"zgoframe/protobuf/pb"
)
type TwinAgora struct{}


func (twinAgora *TwinAgora)CS_CallPeople(ctx context.Context,callPeopleReq *pb.CallPeopleReq) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)SC_CallPeople(ctx context.Context,callPeopleRes *pb.CallPeopleRes) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)CS_CancelCallPeople(ctx context.Context,cancelCallPeopleReq *pb.CancelCallPeopleReq) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)SC_CancelCallPeople(ctx context.Context,cancelCallPeopleReq *pb.CancelCallPeopleReq) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)CS_PeopleEntry(ctx context.Context,peopleEntry *pb.PeopleEntry) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)SC_PeopleEntry(ctx context.Context,peopleEntry *pb.PeopleEntry) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)CS_PeopleLeave(ctx context.Context,peopleLeaveRes *pb.PeopleLeaveRes) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)SC_PeopleLeave(ctx context.Context,peopleLeaveRes *pb.PeopleLeaveRes) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)SC_CallReply(ctx context.Context,callReply *pb.CallReply) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)CS_CallPeopleAccept(ctx context.Context,callVote *pb.CallVote) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)SC_CallPeopleAccept(ctx context.Context,callVote *pb.CallVote) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)CS_CallPeopleDeny(ctx context.Context,callVote *pb.CallVote) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)SC_CallPeopleDeny(ctx context.Context,callVote *pb.CallVote) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)CS_RoomHeartbeat(ctx context.Context,roomHeartbeatReq *pb.RoomHeartbeatReq) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)SC_PushMsg(ctx context.Context,pushMsg *pb.PushMsg) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}

