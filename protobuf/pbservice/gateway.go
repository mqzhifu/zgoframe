package pbservice

import (
	"context"
	"zgoframe/protobuf/pb"
)

type Gateway struct{}

func (gateway *Gateway) CS_Login(ctx context.Context, login *pb.Login) (*pb.LoginRes, error) {
	loginRes := &pb.LoginRes{}
	return loginRes, nil
}
func (gateway *Gateway) CS_Ping(ctx context.Context, ping *pb.Ping) (*pb.Pong, error) {
	pong := &pb.Pong{}
	return pong, nil
}

//func (gateway *Gateway)CS_Pong(ctx context.Context,pong *pb.Pong) (*pb.Pong,error){
//    pong := &pb.Pong{}
//    return pong,nil
//}
func (gateway *Gateway) CS_Heartbeat(ctx context.Context, heartbeat *pb.Heartbeat) (*pb.Empty, error) {
	empty := &pb.Empty{}
	return empty, nil
}
func (gateway *Gateway) SC_Login(ctx context.Context, loginRes *pb.LoginRes) (*pb.Empty, error) {
	empty := &pb.Empty{}
	return empty, nil
}
func (gateway *Gateway) SC_Ping(ctx context.Context, ping *pb.Ping) (*pb.Pong, error) {
	pong := &pb.Pong{}
	return pong, nil
}

//func (gateway *Gateway)SC_Pong(ctx context.Context,pong *pb.Pong) (*pb.Pong,error){
//    pong := &pb.Pong{}
//    return pong,nil
//}
func (gateway *Gateway) SC_Heartbeat(ctx context.Context, heartbeat *pb.Heartbeat) (*pb.Empty, error) {
	empty := &pb.Empty{}
	return empty, nil
}
func (gateway *Gateway) SC_KickOff(ctx context.Context, kickOff *pb.KickOff) (*pb.Empty, error) {
	empty := &pb.Empty{}
	return empty, nil
}

//func (gateway *Gateway)SC_ProjectPush(ctx context.Context,projectPush *pb.ProjectPush) (*pb.ProjectPush,error){
//    projectPush := &pb.ProjectPush{}
//    return projectPush,nil
//}
