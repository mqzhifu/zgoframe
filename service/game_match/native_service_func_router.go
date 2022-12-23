package gamematch

import (
	"errors"
	"time"
	"zgoframe/protobuf/pb"
	"zgoframe/service"
	"zgoframe/util"
)

func (gameMatch *GameMatch) ListeningBridgeMsg() {
	for {
		select {
		case msg := <-gameMatch.Option.ServiceBridge.NativeServiceList.GameMatch:
			gameMatch.NativeServiceFuncRouter(msg)
		default:
			time.Sleep(time.Millisecond * service.BRIDGE_SLEEP_TIME)
		}
	}
}

func (gameMatch *GameMatch) NativeServiceFuncRouter(msg pb.Msg) (data []byte, err error) {
	prefix := "RouterServiceGameMatch"

	requestPlayerMatchSign := pb.GameMatchSign{}
	requestPlayerMatchSignCancel := pb.GameMatchPlayerCancel{}
	protoServiceFunc, _ := gameMatch.Option.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	reqFDCreateEvent := pb.FDCreateEvent{}
	reqHeartbeat := pb.Heartbeat{}
	requestFDCloseEvent := pb.FDCloseEvent{}

	switch protoServiceFunc.FuncName {
	case "CS_PlayerMatchSign":
		err = util.ConnParserContentMsg(msg, &requestPlayerMatchSign)
	case "CS_PlayerMatchSignCancel":
		err = util.ConnParserContentMsg(msg, &requestPlayerMatchSignCancel)
	case "FdClose":
		err = util.ConnParserContentMsg(msg, &requestFDCloseEvent)
	case "CS_Heartbeat":
		err = util.ConnParserContentMsg(msg, &reqHeartbeat)
	case "FdCreate":
		err = util.ConnParserContentMsg(msg, &reqFDCreateEvent)
	default:
		return data, errors.New(gameMatch.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "1"))
	}

	if err != nil {
		gameMatch.Option.Log.Error(prefix + " , ParserContentMsg err:" + err.Error())
		return data, err
	}

	requestPlayerMatchSign.SourceUid = msg.SourceUid
	requestPlayerMatchSignCancel.SourceUid = msg.SourceUid
	reqFDCreateEvent.SourceUid = msg.SourceUid
	reqHeartbeat.SourceUid = msg.SourceUid
	requestFDCloseEvent.SourceUid = msg.SourceUid

	switch protoServiceFunc.FuncName {
	case "CS_PlayerMatchSign":
		_, e := gameMatch.PlayerJoin(requestPlayerMatchSign)
		if e != nil {
			gameMatch.Option.Log.Debug("PlayerJoin return e: " + e.Error())
		}
	case "CS_PlayerMatchSignCancel":
		//requestPlayerMatchSignCancel.SourceUid = conn.UserId
		gameMatch.Cancel(requestPlayerMatchSignCancel)
	case "FdClose":
		util.MyPrint("FdClose")
		//err = frameSync.CloseFD(requestFDCloseEvent)
	case "CS_Heartbeat":
		util.MyPrint("CS_Heartbeat")
		//err = frameSync.Heartbeat(reqHeartbeat)
	case "FdCreate":
		util.MyPrint("FdCreate")
		//err = frameSync.CreateFD(reqFDCreateEvent)
	default:
		return data, errors.New(gameMatch.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "2"))
	}

	return data, err
}

func (gameMatch *GameMatch) MakeRouterErrNotFound(prefix string, funcName string, index string) string {
	errMsg := prefix + " , FuncName not found-" + index + " :" + funcName
	gameMatch.Option.Log.Error(prefix + " , FuncName not found-" + index + " :" + funcName)
	return errMsg
}
