package pbservice

import (
	"context"
	"zgoframe/protobuf/pb"
)
type Gateway struct{}


func (gateway *Gateway)ClientLogin(ctx context.Context,requestLogin *pb.RequestLogin) (*pb.ResponseLoginRes,error){
    responseLoginRes := &pb.ResponseLoginRes{}
    return responseLoginRes,nil
}
func (gateway *Gateway)ClientPing(ctx context.Context,requestClientPing *pb.RequestClientPing) (*pb.ResponseServerPong,error){
    responseServerPong := &pb.ResponseServerPong{}
    return responseServerPong,nil
}
func (gateway *Gateway)ClientPong(ctx context.Context,requestClientPong *pb.RequestClientPong) (*pb.ResponseServerPong,error){
    responseServerPong := &pb.ResponseServerPong{}
    return responseServerPong,nil
}
func (gateway *Gateway)ClientHeartbeat(ctx context.Context,requestClientHeartbeat *pb.RequestClientHeartbeat) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (gateway *Gateway)ServerPing(ctx context.Context,responseServerPing *pb.ResponseServerPing) (*pb.RequestClientPong,error){
    requestClientPong := &pb.RequestClientPong{}
    return requestClientPong,nil
}
func (gateway *Gateway)ServerPong(ctx context.Context,responseServerPong *pb.ResponseServerPong) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (gateway *Gateway)ServerHeartbeat(ctx context.Context,requestClientHeartbeat *pb.RequestClientHeartbeat) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (gateway *Gateway)ServerLogin(ctx context.Context,responseLoginRes *pb.ResponseLoginRes) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (gateway *Gateway)KickOff(ctx context.Context,responseKickOff *pb.ResponseKickOff) (*pb.Empty,error){
    empty := &pb.Empty{}
    return empty,nil
}
func (gateway *Gateway)ProjectPush(ctx context.Context,requestProjectPush *pb.RequestProjectPush) (*pb.ResponseProjectPush,error){
    responseProjectPush := &pb.ResponseProjectPush{}
    return responseProjectPush,nil
}

