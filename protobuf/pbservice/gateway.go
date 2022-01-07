package pb

import (
	"context"
	"zgoframe/protobuf/pb"
)
type FrameSync struct{}


func (frameSync *FrameSync)ClientLogin(ctx context.Context,requestLogin *pbRequestLogin.) (*pbResponseLoginRes.,error){
    responseLoginRes := &pbResponseLoginRes.{}
    return responseLoginRes,nil
}
func (frameSync *FrameSync)ClientPing(ctx context.Context,requestClientPing *pbRequestClientPing.) (*pbResponseServerPong.,error){
    responseServerPong := &pbResponseServerPong.{}
    return responseServerPong,nil
}
func (frameSync *FrameSync)ClientPong(ctx context.Context,requestClientPong *pbRequestClientPong.) (*pbResponseServerPong.,error){
    responseServerPong := &pbResponseServerPong.{}
    return responseServerPong,nil
}
func (frameSync *FrameSync)ClientHeartbeat(ctx context.Context,requestClientHeartbeat *pbRequestClientHeartbeat.) (*pbEmpty.,error){
    empty := &pbEmpty.{}
    return empty,nil
}
func (frameSync *FrameSync)ServerPing(ctx context.Context,responseServerPing *pbResponseServerPing.) (*pbRequestClientPong.,error){
    requestClientPong := &pbRequestClientPong.{}
    return requestClientPong,nil
}
func (frameSync *FrameSync)ServerPong(ctx context.Context,responseServerPong *pbResponseServerPong.) (*pbEmpty.,error){
    empty := &pbEmpty.{}
    return empty,nil
}
func (frameSync *FrameSync)ServerHeartbeat(ctx context.Context,requestClientHeartbeat *pbRequestClientHeartbeat.) (*pbEmpty.,error){
    empty := &pbEmpty.{}
    return empty,nil
}
func (frameSync *FrameSync)ServerLogin(ctx context.Context,responseLoginRes *pbResponseLoginRes.) (*pbEmpty.,error){
    empty := &pbEmpty.{}
    return empty,nil
}
func (frameSync *FrameSync)KickOff(ctx context.Context,responseKickOff *pbResponseKickOff.) (*pbEmpty.,error){
    empty := &pbEmpty.{}
    return empty,nil
}

