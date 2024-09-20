package seed_business

import (
	"errors"
	"time"
	"zgoframe/protobuf/pb"
	"zgoframe/service/bridge"
	"zgoframe/util"
)

func (twinAgora *TwinAgora) ListeningBridgeMsg() {
	for {
		select {
		case msg := <-twinAgora.Op.ServiceBridge.NativeServiceList.TwinAgora:
			twinAgora.NativeServiceFuncRouter(msg)
		default:
			time.Sleep(time.Millisecond * bridge.BRIDGE_SLEEP_TIME)
		}
	}
}

func (twinAgora *TwinAgora) NativeServiceFuncRouter(msg pb.Msg) (data []byte, err error) {
	prefix := "RouterServiceTwinAgora"

	reqCallPeople := pb.CallPeopleReq{}
	reqFDCloseEvent := pb.FDCloseEvent{}
	reqHeartbeat := pb.Heartbeat{}
	reqFDCreateEvent := pb.FDCreateEvent{}
	reqCallVote := pb.CallVote{}
	reqRoomHeartbeat := pb.RoomHeartbeatReq{}
	reqCancelCallPeople := pb.CancelCallPeopleReq{}
	reqPeopleEntry := pb.PeopleEntry{}
	reqPeopleLeaveRes := pb.PeopleLeaveRes{}

	protoServiceFunc, _ := twinAgora.Op.ProtoMap.GetServiceFuncById(int(msg.SidFid))

	switch protoServiceFunc.FuncName {
	case "CS_CallPeople":
		err = util.ConnParserContentMsg(msg, &reqCallPeople)
	case "CS_PeopleLeave":
		err = util.ConnParserContentMsg(msg, &reqPeopleLeaveRes)
	case "CS_CancelCallPeople":
		err = util.ConnParserContentMsg(msg, &reqCancelCallPeople)
	case "FdClose":
		err = util.ConnParserContentMsg(msg, &reqFDCloseEvent)
	case "CS_PeopleEntry":
		err = util.ConnParserContentMsg(msg, &reqPeopleEntry)
	case "FdCreate":
		err = util.ConnParserContentMsg(msg, &reqFDCreateEvent)
	case "CS_RoomHeartbeat":
		err = util.ConnParserContentMsg(msg, &reqRoomHeartbeat)
	case "CS_Heartbeat":
		err = util.ConnParserContentMsg(msg, &reqHeartbeat)
	case "CS_CallPeopleAccept":
		err = util.ConnParserContentMsg(msg, &reqCallVote)
	case "CS_CallPeopleDeny":
		err = util.ConnParserContentMsg(msg, &reqCallVote)
	default:
		return data, errors.New(twinAgora.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "1"))
	}

	if err != nil {
		twinAgora.Log.Error(prefix + " , ParserContentMsg err:" + err.Error())
		return data, err
	}

	reqCallPeople.SourceUid = msg.SourceUid
	reqFDCloseEvent.SourceUid = msg.SourceUid
	reqHeartbeat.SourceUid = msg.SourceUid
	reqFDCreateEvent.SourceUid = msg.SourceUid
	reqCallVote.SourceUid = msg.SourceUid
	reqRoomHeartbeat.SourceUid = msg.SourceUid
	reqCancelCallPeople.SourceUid = msg.SourceUid
	reqPeopleEntry.SourceUid = msg.SourceUid
	reqPeopleLeaveRes.SourceUid = msg.SourceUid

	switch protoServiceFunc.FuncName {
	case "CS_CallPeople":
		twinAgora.CallPeople(reqCallPeople)
	case "CS_Heartbeat":
		twinAgora.UserHeartbeat(reqHeartbeat)
	case "CS_PeopleLeave":
		twinAgora.PeopleLeave(reqPeopleLeaveRes)
	case "CS_RoomHeartbeat":
		twinAgora.RoomHeartbeat(reqRoomHeartbeat)
	case "FdClose":
		twinAgora.FDCloseEvent(reqFDCloseEvent)
	case "CS_PeopleEntry":
		twinAgora.PeopleEntry(reqPeopleEntry)
	case "CS_CancelCallPeople":
		twinAgora.CancelCallPeople(reqCancelCallPeople)
	case "FdCreate":
		twinAgora.FDCreateEvent(reqFDCreateEvent)
	case "CS_CallPeopleAccept":
		twinAgora.CallPeopleAccept(reqCallVote)
	case "CS_CallPeopleDeny":
		twinAgora.CallPeopleDeny(reqCallVote)
	default:
		return data, errors.New(twinAgora.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "2"))
	}

	return data, err
}

func (twinAgora *TwinAgora) MakeRouterErrNotFound(prefix string, funcName string, index string) string {
	errMsg := prefix + " , FuncName not found-" + index + " :" + funcName
	twinAgora.Log.Error(prefix + " , FuncName not found-" + index + " :" + funcName)
	return errMsg
}
