package frame_sync

import (
	"errors"
	"time"
	"zgoframe/protobuf/pb"
	"zgoframe/service"
	"zgoframe/util"
)

func (frameSync *FrameSync) ListeningBridgeMsg() {
	for {
		select {
		case msg := <-frameSync.Option.ServiceBridge.NativeServiceList.FrameSync:
			frameSync.NativeServiceFuncRouter(msg)
		default:
			time.Sleep(time.Millisecond * service.BRIDGE_SLEEP_TIME)
		}
	}
}

//帧同步的路由
func (frameSync *FrameSync) NativeServiceFuncRouter(msg pb.Msg) (data []byte, err error) {
	prefix := "RouterServiceSync "
	requestLogicFrame := pb.LogicFrame{}
	requestPlayerResumeGame := pb.PlayerResumeGame{}
	requestPlayerReady := pb.PlayerReady{}
	requestPlayerOver := pb.PlayerOver{}
	requestRoomHistory := pb.ReqRoomHistory{}
	requestRoomBaseInfo := pb.RoomBaseInfo{}
	//requestPlayerMatchSign := pb.GameMatchSign{}
	//requestPlayerMatchSignCancel := pb.GameMatchPlayerCancel{}

	reqPlayerBase := pb.PlayerBase{}

	reqFDCreateEvent := pb.FDCreateEvent{}
	reqHeartbeat := pb.Heartbeat{}
	requestFDCloseEvent := pb.FDCloseEvent{}

	protoServiceFunc, _ := frameSync.Option.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	switch protoServiceFunc.FuncName {
	//case "CS_PlayerMatchSign":
	//	err = util.ConnParserContentMsg(msg, &requestPlayerMatchSign, conn.UserId)
	//case "CS_PlayerMatchSignCancel":
	//	err = util.ConnParserContentMsg(msg, &requestPlayerMatchSignCancel, conn.UserId)
	case "CS_PlayerOperations":
		err = util.ConnParserContentMsg(msg, &requestLogicFrame)
	case "CS_PlayerResumeGame":
		err = util.ConnParserContentMsg(msg, &requestPlayerResumeGame)
	case "CS_PlayerReady":
		err = util.ConnParserContentMsg(msg, &requestPlayerReady)
	case "CS_PlayerOver":
		err = util.ConnParserContentMsg(msg, &requestPlayerOver)
	case "CS_RoomHistory":
		err = util.ConnParserContentMsg(msg, &requestRoomHistory)
	case "CS_RoomBaseInfo":
		err = util.ConnParserContentMsg(msg, &requestRoomBaseInfo)
	case "FdClose":
		err = util.ConnParserContentMsg(msg, &requestFDCloseEvent)
	case "CS_Heartbeat":
		err = util.ConnParserContentMsg(msg, &reqHeartbeat)
	case "CS_PlayerState":
		err = util.ConnParserContentMsg(msg, &reqPlayerBase)
	case "FdCreate":
		err = util.ConnParserContentMsg(msg, &reqFDCreateEvent)
	default:
		return data, errors.New(frameSync.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "1"))
	}
	if err != nil {
		frameSync.Option.Log.Error(prefix + " , ParserContentMsg err:" + err.Error())
	}

	switch protoServiceFunc.FuncName {
	case "CS_PlayerOperations":
		//requestLogicFrame.SourceUid = conn.UserId
		err = frameSync.ReceivePlayerOperation(requestLogicFrame)
	case "CS_PlayerResumeGame":
		//requestPlayerResumeGame.SourceUid = conn.UserId
		err = frameSync.PlayerResumeGame(requestPlayerResumeGame)
	case "CS_PlayerReady":
		//requestPlayerReady.SourceUid = conn.UserId
		err = frameSync.PlayerReady(requestPlayerReady)
	case "CS_PlayerOver":
		//requestPlayerOver.SourceUid = conn.UserId
		err = frameSync.PlayerOver(requestPlayerOver)
	case "CS_RoomHistory":
		//requestRoomHistory.SourceUid = conn.UserId
		err = frameSync.RoomHistory(requestRoomHistory)
	case "CS_RoomBaseInfo":
		//requestRoomBaseInfo.SourceUid = conn.UserId
		err = frameSync.RoomManage.GetRoom(requestRoomBaseInfo)
	case "CS_PlayerState":
		frameSync.GetPlayerBase(reqPlayerBase)
	case "FdClose":
		err = frameSync.CloseFD(requestFDCloseEvent)
	case "CS_Heartbeat":
		err = frameSync.Heartbeat(reqHeartbeat)
	case "FdCreate":
		err = frameSync.CreateFD(reqFDCreateEvent)
	default:
		return data, errors.New(frameSync.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "2"))
	}

	return data, err
}

func (frameSync *FrameSync) MakeRouterErrNotFound(prefix string, funcName string, index string) string {
	errMsg := prefix + " , FuncName not found-" + index + " :" + funcName
	frameSync.Option.Log.Error(prefix + " , FuncName not found-" + index + " :" + funcName)
	return errMsg
}
