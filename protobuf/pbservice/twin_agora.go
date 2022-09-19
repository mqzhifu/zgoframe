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
func (twinAgora *TwinAgora)SC_PeopleEntry(ctx context.Context,peopleEntryRes *pb.PeopleEntryRes) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (twinAgora *TwinAgora)SC_PushMsg(ctx context.Context,pushMsgRes *pb.PushMsgRes) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}

