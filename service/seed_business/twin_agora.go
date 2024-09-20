package seed_business

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"strconv"
	"zgoframe/model"
	"zgoframe/protobuf/pb"
	"zgoframe/service/bridge"
	"zgoframe/util"
)

// 眼镜端发起呼叫
func (twinAgora *TwinAgora) CallPeople(callPeopleReq pb.CallPeopleReq) {
	twinAgora.Log.Info("in func CallPeople:")
	callPeopleRes := pb.CallPeopleRes{}
	callPeopleRes.AgoraAppId = callPeopleReq.AgoraAppId
	callPeopleRes.AgoraChannel = callPeopleReq.AgoraChannel

	if callPeopleReq.Uid <= 0 {
		callPeopleRes.ErrCode = 400
		callPeopleRes.ErrMsg = twinAgora.Lang.NewString(400)
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(callPeopleReq.SourceUid, "SC_CallPeople", callPeopleRes)
		//data, _ := proto.Marshal(&callPeopleRes)
		//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallPeople", 9999, callPeopleReq.SourceUid, string(data), "", 0)
		callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallPeople", TargetUid: callPeopleReq.Uid, Data: &callPeopleRes}
		twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)
		return
	}

	if callPeopleReq.PeopleType != int32(model.USER_ROLE_DOCTOR) {
		callPeopleRes.ErrCode = 420
		callPeopleRes.ErrMsg = twinAgora.Lang.NewString(420)
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(callPeopleReq.SourceUid, "SC_CallPeople", callPeopleRes)
		//data, _ := proto.Marshal(&callPeopleRes)
		//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallPeople", 9999, callPeopleReq.SourceUid, string(data), "", 0)
		callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallPeople", TargetUid: callPeopleReq.Uid, Data: &callPeopleRes}
		twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)
		return
	}

	if callPeopleReq.TargetUid > 0 {
		callPeopleRes.ErrCode = 421
		callPeopleRes.ErrMsg = twinAgora.Lang.NewString(421)
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(callPeopleReq.SourceUid, "SC_CallPeople", callPeopleRes)
		//data, _ := proto.Marshal(&callPeopleRes)
		//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallPeople", 9999, callPeopleReq.SourceUid, string(data), "", 0)
		callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallPeople", TargetUid: callPeopleReq.Uid, Data: &callPeopleRes}
		twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)
		return
	}

	myRTCUser, ok := twinAgora.GetUserById(int(callPeopleReq.Uid))
	if ok && myRTCUser.RoomId != "" {
		RTCRoomInfo, err := twinAgora.GetRoomById(myRTCUser.RoomId)
		if err != nil {
			callPeopleRes.ErrCode = 501
			callPeopleRes.ErrMsg = twinAgora.Lang.NewReplaceOneString(501, myRTCUser.RoomId)
			twinAgora.MakeError(callPeopleRes.ErrMsg)
			//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(callPeopleReq.SourceUid, "SC_CallPeople", callPeopleRes)
			//data, _ := proto.Marshal(&callPeopleRes)
			//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallPeople", 9999, callPeopleReq.SourceUid, string(data), "", 0)

			callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallPeople", TargetUid: callPeopleReq.Uid, Data: &callPeopleRes}
			twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)

			return
		}

		if RTCRoomInfo.Status == RTC_ROOM_STATUS_CALLING {
			callPeopleRes.ErrCode = 514
			callPeopleRes.ErrMsg = twinAgora.Lang.NewReplaceOneString(514, myRTCUser.RoomId)
			twinAgora.MakeError(callPeopleRes.ErrMsg)
			//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(callPeopleReq.SourceUid, "SC_CallPeople", callPeopleRes)
			//data, _ := proto.Marshal(&callPeopleRes)
			//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallPeople", 9999, callPeopleReq.SourceUid, string(data), "", 0)
			callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallPeople", TargetUid: callPeopleReq.Uid, Data: &callPeopleRes}
			twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)
			return
		}

		if RTCRoomInfo.Status == RTC_ROOM_STATUS_EXECING {
			callPeopleRes.ErrCode = 513
			callPeopleRes.ErrMsg = twinAgora.Lang.NewReplaceOneString(513, myRTCUser.RoomId)
			twinAgora.MakeError(callPeopleRes.ErrMsg)
			//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(callPeopleReq.SourceUid, "SC_CallPeople", callPeopleRes)
			//data, _ := proto.Marshal(&callPeopleRes)
			//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallPeople", 9999, callPeopleReq.SourceUid, string(data), "", 0)
			callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallPeople", TargetUid: callPeopleReq.Uid, Data: &callPeopleRes}
			twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)
			return
		}
		//该房间状态已经结束了，但未做清算处理(持久化)，这里做个容错吧
		if RTCRoomInfo.Status == RTC_ROOM_STATUS_END {
			callPeopleRes.ErrCode = 510
			callPeopleRes.ErrMsg = twinAgora.Lang.NewReplaceOneString(510, RTCRoomInfo.Id)
			twinAgora.MakeError(callPeopleRes.ErrMsg)
			//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(callPeopleReq.SourceUid, "SC_CallPeople", callPeopleRes)
			//data, _ := proto.Marshal(&callPeopleRes)
			//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallPeople", 9999, callPeopleReq.SourceUid, string(data), "", 0)
			callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallPeople", TargetUid: callPeopleReq.Uid, Data: &callPeopleRes}
			twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)
			return
		}
	}

	var userDoctorList []model.User
	err := twinAgora.Gorm.Where(" role =  ?", model.USER_ROLE_DOCTOR).Find(&userDoctorList).Error
	if err != nil {
		callPeopleRes.ErrCode = 402
		callPeopleRes.ErrMsg = twinAgora.Lang.NewString(402)
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(callPeopleReq.SourceUid, "SC_CallPeople", callPeopleRes)
		//data, _ := proto.Marshal(&callPeopleRes)
		//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallPeople", 9999, callPeopleReq.Uid, string(data), "", 0)

		callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallPeople", TargetUid: callPeopleReq.Uid, Data: &callPeopleRes}
		twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)
		return
	}

	//寻找在线的专家
	var onlineUserDoctorList []model.User
	onlineUserDoctorListUids := ""
	for _, userConn := range twinAgora.RTCUserPool {
		for _, user := range userDoctorList {
			//if userConn.Id == user.Id && userConn.Status == util.CONN_STATUS_EXECING {
			if userConn.Id == user.Id {
				onlineUserDoctorList = append(onlineUserDoctorList, user)
				onlineUserDoctorListUids += strconv.Itoa(user.Id) + " , "
			}
		}
	}

	debugLogInfo := "RTCUserPool len:" + strconv.Itoa(len(twinAgora.RTCUserPool)) + "RTCUserPool len:" + strconv.Itoa(len(userDoctorList)) + "RTCUserPool len:" + strconv.Itoa(len(onlineUserDoctorList))
	twinAgora.Log.Debug(debugLogInfo)

	if len(onlineUserDoctorList) <= 0 {
		callPeopleRes.ErrCode = 403
		callPeopleRes.ErrMsg = twinAgora.Lang.NewString(403)
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(callPeopleReq.SourceUid, "SC_CallPeople", callPeopleRes)
		//data, _ := proto.Marshal(&callPeopleRes)
		//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallPeople", 9999, callPeopleReq.SourceUid, string(data), "", 0)
		callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallPeople", TargetUid: callPeopleReq.Uid, Data: &callPeopleRes}
		twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)
		return
	}

	if len(onlineUserDoctorList) > 1 {
		callPeopleRes.ErrCode = 422
		callPeopleRes.ErrMsg = twinAgora.Lang.NewReplaceOneString(422, onlineUserDoctorListUids)
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(callPeopleReq.SourceUid, "SC_CallPeople", callPeopleRes)
		//data, _ := proto.Marshal(&callPeopleRes)
		//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallPeople", 9999, callPeopleReq.SourceUid, string(data), "", 0)
		callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallPeople", TargetUid: callPeopleReq.Uid, Data: &callPeopleRes}
		twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)
		return
	}

	myRTCRoom := twinAgora.CreateRTCRoom(callPeopleReq, onlineUserDoctorList)
	receiveUidsStr := ""
	for _, user := range onlineUserDoctorList {
		receiveUidsStr += strconv.Itoa(user.Id) + "," //专家接收列表
		callReply := pb.CallReply{}
		callReply.RoomId = myRTCRoom.Id
		callReply.AgoraAppId = callPeopleReq.AgoraAppId
		callReply.AgoraChannel = callPeopleReq.AgoraChannel
		callReply.Content = strconv.Itoa(int(callPeopleReq.Uid)) + " calling....... please reply:" + callPeopleReq.Channel
		//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(int32(user.Id), "SC_CallReply", callReply)

		//data, _ := proto.Marshal(&callReply)
		//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallReply", 9999, int32(user.Id), string(data), "", 0)
		callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallReply", TargetUid: int32(user.Id), Data: &callReply}
		twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)
	}

	myRTCUser.RoomId = myRTCRoom.Id
	//先给呼叫者返回消息，告知已成功请等待专家响应
	callPeopleRes.RoomId = myRTCRoom.Id
	callPeopleRes.ErrCode = 200
	callPeopleRes.ErrMsg = "waiting doctor reply"
	callPeopleRes.ReceiveUid = receiveUidsStr[0 : len(receiveUidsStr)-1]

	//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(callPeopleReq.SourceUid, "SC_CallPeople", callPeopleRes)
	//data, _ := proto.Marshal(&callPeopleRes)
	//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallPeople", 9999, callPeopleReq.SourceUid, string(data), "", 0)

	callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallPeople", TargetUid: callPeopleReq.Uid, Data: &callPeopleRes}
	twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)

	return
}

// 这里假设验证都成功了，不做二次验证了
func (twinAgora *TwinAgora) CreateRTCRoom(callPeopleReq pb.CallPeopleReq, onlineUserDoctorList []model.User) RTCRoom {
	twinAgora.Log.Info("CreateRTCRoom ", zap.Int32("uid", callPeopleReq.Uid))
	//给所有在线的专家发送邀请通知
	var receiveUids []int
	for _, user := range onlineUserDoctorList {
		receiveUids = append(receiveUids, user.Id) //专家接收列表
	}
	RTCRoomOne := RTCRoom{
		AddTime:     util.GetNowTimeSecondToInt(),
		CallUid:     int(callPeopleReq.Uid),
		Status:      RTC_ROOM_STATUS_CALLING,
		Uids:        append(receiveUids, int(callPeopleReq.Uid)),
		Id:          uuid.NewV4().String(),
		ReceiveUids: receiveUids,
		//OnlineUids:  []int{int(callPeopleReq.Uid)},
	}
	myRTCUser, ok := twinAgora.GetUserById(int(callPeopleReq.Uid))
	if ok {
		myRTCUser.RoomId = RTCRoomOne.Id
	} else {
		newRTCUser := RTCUser{}
		newRTCUser.RoomId = RTCRoomOne.Id
		newRTCUser.Uptime = util.GetNowTimeSecondToInt()
		twinAgora.RTCUserPool[int(callPeopleReq.Uid)] = &newRTCUser
	}
	twinAgora.Log.Info("CreateRTCRoom id:" + RTCRoomOne.Id)
	twinAgora.RTCRoomPool[RTCRoomOne.Id] = &RTCRoomOne
	return RTCRoomOne
}

// 发起方，取消呼叫
func (twinAgora *TwinAgora) CancelCallPeople(cancelCallPeopleReq pb.CancelCallPeopleReq) error {
	twinAgora.Log.Warn("CancelCallPeople , ", zap.Int32("uid", cancelCallPeopleReq.Uid), zap.String("roomId", cancelCallPeopleReq.RoomId))
	if cancelCallPeopleReq.Uid <= 0 {
		return twinAgora.MakeError(twinAgora.Lang.NewString(400))
	}

	if cancelCallPeopleReq.RoomId == "" {
		return twinAgora.MakeError(twinAgora.Lang.NewString(500))
	}

	myRTCUser, ok := twinAgora.GetUserById(int(cancelCallPeopleReq.Uid))
	if !ok {
		return errors.New(twinAgora.Lang.NewReplaceOneString(401, strconv.Itoa(int(cancelCallPeopleReq.Uid))))
	}

	if myRTCUser.RoomId != cancelCallPeopleReq.RoomId {
		return errors.New(twinAgora.Lang.NewReplaceOneString(412, myRTCUser.RoomId))
	}

	RTCRoomInfo, err := twinAgora.GetRoomById(cancelCallPeopleReq.RoomId)
	if err != nil {
		return err
	}
	//取消呼叫，只能由发起者自己取消
	if RTCRoomInfo.CallUid != int(cancelCallPeopleReq.Uid) {
		return errors.New(twinAgora.Lang.NewReplaceOneString(406, strconv.Itoa(int(cancelCallPeopleReq.Uid))))
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return errors.New(twinAgora.Lang.NewReplaceOneString(512, strconv.Itoa(RTCRoomInfo.Status)))
	}

	//util.MyPrint("RTCRoomInfo.ReceiveUids:", RTCRoomInfo.ReceiveUids, " cancelCallPeopleReq.Uid:", cancelCallPeopleReq.Uid)
	//给所有专家端用户发送取消的消息
	for _, uid := range RTCRoomInfo.ReceiveUids {
		if int(cancelCallPeopleReq.Uid) == uid {
			continue
		}
		//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(int32(uid), "SC_CancelCallPeople", cancelCallPeopleReq)
		//data, _ := proto.Marshal(&cancelCallPeopleReq)
		//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CancelCallPeople", 9999, int32(uid), string(data), "", 0)

		callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CancelCallPeople", TargetUid: int32(uid), Data: &cancelCallPeopleReq}
		twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)

	}
	//清空user的roomId值，RoomEnd只会清空掉已进入房间的用户，而此时房间虽然存在，但没有人进入，用户直接取消了，把给清空一下
	twinAgora.RoomEnd(cancelCallPeopleReq.RoomId, RTC_ROOM_END_STATUS_CANCEL)
	//myRTCUser.RoomId = ""  //这个在roomEnd 里已经做了，不要重复，不然roomEnd 会判断，直接停了
	return nil
}

func (twinAgora *TwinAgora) PeopleEntry(peopleEntry pb.PeopleEntry) error {
	twinAgora.Log.Info("PeopleEntry  ", zap.Int32("uid", peopleEntry.Uid))
	myRTCUser, ok := twinAgora.GetUserById(int(peopleEntry.Uid))
	if !ok {
		return twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(401, strconv.Itoa(int(peopleEntry.Uid))))
	}
	RTCRoomInfo, err := twinAgora.GetRoomById(myRTCUser.RoomId)
	if err != nil {
		return err
	}

	hasSearch := 0
	for _, uid := range RTCRoomInfo.OnlineUids {
		if int(peopleEntry.Uid) == uid {
			hasSearch = 1
			break
		}
	}

	if hasSearch == 1 {
		//您并不在此频道中，请不要乱发消息
		return errors.New(twinAgora.Lang.NewReplaceOneString(407, strconv.Itoa(int(peopleEntry.Uid))))
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_EXECING {
		return twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(511, RTCRoomInfo.Id))
	}
	//util.MyPrint("PeopleEntry:", RTCRoomInfo.Uids, conn.UserId)
	for _, uid := range RTCRoomInfo.Uids {
		if int(peopleEntry.Uid) == uid {
			continue
		}
		//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(int32(uid), "SC_PeopleEntry", peopleEntry)
		//data, _ := proto.Marshal(&peopleEntry)
		callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_PeopleEntry", TargetUid: int32(uid), Data: &peopleEntry}
		//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_PeopleEntry", 9999, int32(uid), string(data), "", 0)
		twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)
	}
	RTCRoomInfo.OnlineUids = append(RTCRoomInfo.OnlineUids, int(peopleEntry.Uid))

	return nil
}

// 某用户离开了房间
func (twinAgora *TwinAgora) PeopleLeave(peopleLeaveRes pb.PeopleLeaveRes) error {
	myRTCUser, ok := twinAgora.GetUserById(int(peopleLeaveRes.Uid))
	if !ok {
		return errors.New(twinAgora.Lang.NewReplaceOneString(401, strconv.Itoa(int(peopleLeaveRes.Uid))))
	}
	RTCRoomInfo, err := twinAgora.GetRoomById(myRTCUser.RoomId)
	if err != nil {
		return errors.New(twinAgora.Lang.NewReplaceOneString(501, myRTCUser.RoomId))
	}

	hasSearch := 0
	for _, uid := range RTCRoomInfo.OnlineUids {
		if int(peopleLeaveRes.Uid) == uid {
			hasSearch = 1
			break
		}
	}

	if hasSearch == 0 {
		//您并不在此频道中，请不要乱发消息
		return errors.New(twinAgora.Lang.NewReplaceOneString(522, myRTCUser.RoomId))
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_EXECING {
		return twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(511, RTCRoomInfo.Id))
	}

	twinAgora.RoomEnd(myRTCUser.RoomId, RTC_ROOM_END_STATUS_USER_LEAVE)

	return nil
}

// 被呼叫者，接收/拒绝 公共验证
func (twinAgora *TwinAgora) PeopleVote(callVote pb.CallVote) (error, *RTCRoom) {
	twinAgora.Log.Info("PeopleVote :", zap.Int32("uid", callVote.Uid), zap.String("roomId", callVote.RoomId))
	RTCRoomInfo, err := twinAgora.GetRoomById(callVote.RoomId)
	if err != nil {
		return err, RTCRoomInfo
	}
	hasSearch := 0
	for _, uid := range RTCRoomInfo.ReceiveUids {
		if int(callVote.Uid) == uid {
			hasSearch = 1
			break
		}
	}
	//并没有发消息让你判定是否接收
	if hasSearch == 0 {
		return twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(409, strconv.Itoa(int(callVote.Uid)))), RTCRoomInfo
	}

	hasVote := 0
	for _, uid := range RTCRoomInfo.ReceiveUidsAccept {
		if int(callVote.Uid) == uid {
			hasVote = 1
			break
		}
	}
	if hasVote == 1 {
		//您已经投过票了，不要重复操作
		return twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(410, strconv.Itoa(int(callVote.Uid)))), RTCRoomInfo
	}

	hasVote = 0
	for _, uid := range RTCRoomInfo.ReceiveUidsDeny {
		if int(callVote.Uid) == uid {
			hasVote = 1
			break
		}
	}
	if hasVote == 1 {
		//您已经投过票了，不要重复操作
		return twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(411, strconv.Itoa(int(callVote.Uid)))), RTCRoomInfo
	}

	return nil, RTCRoomInfo
}

// 被呼叫者，接收呼叫
func (twinAgora *TwinAgora) CallPeopleAccept(callVote pb.CallVote) error {
	twinAgora.Log.Info("in func :CallPeopleAccept , " + "callVoteRoomId:" + callVote.RoomId)
	err, RTCRoomInfo := twinAgora.PeopleVote(callVote)
	if err != nil {
		return err
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(512, callVote.RoomId))
	}

	rtcUser, exist := twinAgora.GetUserById(int(callVote.Uid))
	if !exist {
		return twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(401, strconv.Itoa(int(callVote.Uid))))
	}
	rtcUser.RoomId = RTCRoomInfo.Id

	RTCRoomInfo.ReceiveUidsAccept = append(RTCRoomInfo.ReceiveUidsAccept, int(callVote.Uid))
	//util.MyPrint("RTCRoomInfo.ReceiveUidsAccept:", RTCRoomInfo.ReceiveUidsAccept)
	//RTCRoomInfo.Uids = append(RTCRoomInfo.Uids, int(callVote.Uid))
	//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(int32(RTCRoomInfo.CallUid), "SC_CallPeopleAccept", callVote)
	//data, _ := proto.Marshal(&callVote)
	callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallPeopleAccept", TargetUid: int32(RTCRoomInfo.CallUid), Data: &callVote}
	//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallPeopleAccept", 9999, int32(RTCRoomInfo.CallUid), string(data), "", 0)
	twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)
	//目前是1v1视频，只要有一个接收，即把房间状态标识为运行中，这里后期优化一下吧
	RTCRoomInfo.Status = RTC_ROOM_STATUS_EXECING
	return nil

}

// 被呼叫者，拒绝呼叫
func (twinAgora *TwinAgora) CallPeopleDeny(callVote pb.CallVote) error {
	twinAgora.Log.Debug("in func :CallPeopleDeny")
	err, RTCRoomInfo := twinAgora.PeopleVote(callVote)
	if err != nil {
		return err
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(512, callVote.RoomId))
	}
	RTCRoomInfo.ReceiveUidsDeny = append(RTCRoomInfo.ReceiveUidsDeny, int(callVote.Uid))
	//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(int32(RTCRoomInfo.CallUid), "SC_CallPeopleDeny", callVote)
	//data, _ := proto.Marshal(&callVote)
	callGatewayMsg := bridge.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_CallPeopleDeny", TargetUid: int32(RTCRoomInfo.CallUid), Data: &callVote}
	//twinAgora.Op.ServiceBridge.CallGateway("TwinAgora", "SC_CallPeopleDeny", 9999, int32(RTCRoomInfo.CallUid), string(data), "", 0)
	twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)

	twinAgora.RoomEnd(callVote.RoomId, RTC_ROOM_END_STATUS_DENY)
	return err

}
