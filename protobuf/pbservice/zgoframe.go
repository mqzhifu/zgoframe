package pb

import (
	"context"
	"zgoframe/protobuf/pb"
)
type Zgoframe struct{}


func (zgoframe *Zgoframe)SayHello(ctx context.Context,requestUser *pbRequestUser.) (*pbResponseUser.,error){
    responseUser := &pbResponseUser.{}
    return responseUser,nil
}
func (zgoframe *Zgoframe)Comm(ctx context.Context,requestUser *pbRequestUser.) (*pbResponseUser.,error){
    responseUser := &pbResponseUser.{}
    return responseUser,nil
}

type Sync struct{}


func (sync *Sync)Heartbeat(ctx context.Context,cSHeartbeat *pbCSHeartbeat.) (*pbSCHeartbeat.,error){
    sCHeartbeat := &pbSCHeartbeat.{}
    return sCHeartbeat,nil
}

