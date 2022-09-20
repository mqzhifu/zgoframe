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
func (twinAgora *TwinAgora)CS_CancelCallPeople(ctx context.Context,callPeopleReq *pb.CallPeopleReq) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)SC_CancelCallPeople(ctx context.Context,callPeopleReq *pb.CallPeopleReq) (*pb.Empty,error){
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
func (twinAgora *TwinAgora)CS_CallPeopleAccept(ctx context.Context,callPeopleReq *pb.CallPeopleReq) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)CS_CallPeopleDeny(ctx context.Context,callPeopleReq *pb.CallPeopleReq) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)SC_CallReply(ctx context.Context,callReply *pb.CallReply) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}

