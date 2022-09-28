package service

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"zgoframe/model"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

//眼镜端发起呼叫
func (twinAgora *TwinAgora) CallPeople(callPeopleReq pb.CallPeopleReq, conn *util.Conn) {
	util.MyPrint("in func CallPeople:")
	callPeopleRes := pb.CallPeopleRes{}

	if callPeopleReq.Uid <= 0 {
		callPeopleRes.ErrCode = 501
		callPeopleRes.ErrMsg = ERR_UID_ZERO
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(conn.UserId, "SC_CallPeople", callPeopleRes)
		return
	}

	if callPeopleReq.PeopleType != int32(USER_ROLE_DOCTOR) {
		callPeopleRes.ErrCode = 502
		callPeopleRes.ErrMsg = "callPeopleReq.PeopleType !=  " + strconv.Itoa(USER_ROLE_DOCTOR)
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(conn.UserId, "SC_CallPeople", callPeopleRes)
		return
	}

	if callPeopleReq.TargetUid > 0 {
		callPeopleRes.ErrCode = 504
		callPeopleRes.ErrMsg = "暂时不支持 TargetUid > 0 的情况"
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(conn.UserId, "SC_CallPeople", callPeopleRes)
		return
	}

	uid := int(conn.UserId)
	myRTCUser, ok := twinAgora.GetUserById(uid)
	if ok && myRTCUser.RoomId != "" {
		RTCRoomInfo, err := twinAgora.GetRoomById(myRTCUser.RoomId)
		if err != nil {
			callPeopleRes.ErrCode = 511
			callPeopleRes.ErrMsg = ERR_ROOM_ID_NOT_IN_MAP + myRTCUser.RoomId
			twinAgora.MakeError(callPeopleRes.ErrMsg)
			conn.SendMsgCompressByUid(conn.UserId, "SC_CallPeople", callPeopleRes)
			return
		}

		if RTCRoomInfo.Status == RTC_ROOM_STATUS_CALLING {
			callPeopleRes.ErrCode = 520
			callPeopleRes.ErrMsg = "已经存在一条记录：发起呼叫，请不要重复发起，或 等待超时 或 主动挂断"
			twinAgora.MakeError(callPeopleRes.ErrMsg)
			conn.SendMsgCompressByUid(conn.UserId, "SC_CallPeople", callPeopleRes)
			return
		}

		if RTCRoomInfo.Status == RTC_ROOM_STATUS_EXECING {
			callPeopleRes.ErrCode = 521
			callPeopleRes.ErrMsg = "已经存在一条记录：正在与其它人通话中...，不能再发起CALL了"
			twinAgora.MakeError(callPeopleRes.ErrMsg)
			conn.SendMsgCompressByUid(conn.UserId, "SC_CallPeople", callPeopleRes)
			return
		}
		//该房间状态已经结束了，但未做清算处理(持久化)，这里做个容错吧
		if RTCRoomInfo.Status == RTC_ROOM_STATUS_END {
			callPeopleRes.ErrCode = 522
			callPeopleRes.ErrMsg = "房间结束了，但未清楚，请等待后台进程清算"
			twinAgora.MakeError(callPeopleRes.ErrMsg)
			conn.SendMsgCompressByUid(conn.UserId, "SC_CallPeople", callPeopleRes)
			return
		}
	}

	var userDoctorList []model.User
	err := twinAgora.Gorm.Where(" role =  ?", USER_ROLE_DOCTOR).Find(&userDoctorList).Error
	if err != nil {
		callPeopleRes.ErrCode = 503
		callPeopleRes.ErrMsg = "数据库中没有专家用户"
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(conn.UserId, "SC_CallPeople", callPeopleRes)
		return
	}

	//寻找在线的专家
	var onlineUserDoctorList []model.User
	for _, userConn := range conn.ConnManager.Pool {
		for _, user := range userDoctorList {
			if userConn.UserId == int32(user.Id) && userConn.Status == util.CONN_STATUS_EXECING {
				onlineUserDoctorList = append(onlineUserDoctorList, user)
			}
		}
	}

	if len(onlineUserDoctorList) <= 0 {
		callPeopleRes.ErrCode = 510
		callPeopleRes.ErrMsg = "所有专家，均不在线"
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(conn.UserId, "SC_CallPeople", callPeopleRes)
		return
	}
	myRTCRoom := twinAgora.CreateRTCRoom(callPeopleReq, onlineUserDoctorList)
	receiveUidsStr := ""
	for _, user := range onlineUserDoctorList {
		receiveUidsStr += strconv.Itoa(user.Id) + "," //专家接收列表
		callReply := pb.CallReply{}
		callReply.RoomId = myRTCRoom.Id
		callReply.Content = strconv.Itoa(int(callPeopleReq.Uid)) + " 呼叫 视频连接...请进入频道:" + callPeopleReq.Channel
		conn.SendMsgCompressByUid(int32(user.Id), "SC_CallReply", callReply)
	}

	myRTCUser.RoomId = myRTCRoom.Id
	//先给呼叫者返回消息，告知已成功请等待专家响应
	callPeopleRes.RoomId = myRTCRoom.Id
	callPeopleRes.ErrCode = 200
	callPeopleRes.ErrMsg = "请求等待专家响应"
	callPeopleRes.ReceiveUid = receiveUidsStr[0 : len(receiveUidsStr)-1]

	conn.SendMsgCompressByUid(conn.UserId, "SC_CallPeople", callPeopleRes)

	return
}

//这里假设验证都成功了，不做二次验证了
func (twinAgora *TwinAgora) CreateRTCRoom(callPeopleReq pb.CallPeopleReq, onlineUserDoctorList []model.User) RTCRoom {
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
	}
	myRTCUser, ok := twinAgora.GetUserById(int(callPeopleReq.Uid))
	if ok {
		myRTCUser.RoomId = RTCRoomOne.Id
	} else {
		newRTCUser := RTCUser{}
		newRTCUser.RoomId = RTCRoomOne.Id
		newRTCUser.Uptime = util.GetNowTimeSecondToInt()
		twinAgora.RTCUserPool[myRTCUser.Id] = &newRTCUser
	}
	util.MyPrint("CreateRTCRoom id:", RTCRoomOne.Id, " , uids:", RTCRoomOne.Uids)
	twinAgora.RTCRoomPool[RTCRoomOne.Id] = &RTCRoomOne
	return RTCRoomOne
}

//发起方，取消呼叫
func (twinAgora *TwinAgora) CancelCallPeople(cancelCallPeopleReq pb.CancelCallPeopleReq, conn *util.Conn) error {
	util.MyPrint("CancelCallPeople , uid:", cancelCallPeopleReq.Uid, " roomId:", cancelCallPeopleReq.RoomId)
	if cancelCallPeopleReq.Uid <= 0 {
		return twinAgora.MakeError(ERR_UID_ZERO)
	}

	if cancelCallPeopleReq.RoomId == "" {
		return twinAgora.MakeError(ERR_ROOM_ID_EMPTY)
	}

	myRTCUser, ok := twinAgora.GetUserById(int(cancelCallPeopleReq.Uid))
	if !ok {
		return errors.New(ERR_UID_NOT_IN_MAP)
	}

	if myRTCUser.RoomId != cancelCallPeopleReq.RoomId {
		return errors.New("myRTCUser.RoomId != cancelCallPeopleReq.RoomId")
	}

	RTCRoomInfo, err := twinAgora.GetRoomById(cancelCallPeopleReq.RoomId)
	if err != nil {
		return err
	}
	//取消呼叫，只能由发起者自己取消
	if RTCRoomInfo.CallUid != int(cancelCallPeopleReq.Uid) {
		return errors.New("你不是发起呼叫者，不能取消操作")
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return errors.New(ERR_ROOM_STATUS_NOT_CALL + strconv.Itoa(RTCRoomInfo.Status))
	}
	//给所有专家端用户发送取消的消息
	for _, uid := range RTCRoomInfo.ReceiveUids {
		if int(cancelCallPeopleReq.Uid) == uid {
			continue
		}
		conn.SendMsgCompressByUid(cancelCallPeopleReq.Uid, "SC_CancelCallPeople", cancelCallPeopleReq)
	}
	twinAgora.RoomEnd(cancelCallPeopleReq.RoomId, RTC_ROOM_END_STATUS_CANCEL)
	return nil
}

func (twinAgora *TwinAgora) PeopleEntry(peopleEntry pb.PeopleEntry, conn *util.Conn) error {
	util.MyPrint("PeopleEntry  uid:", peopleEntry.Uid)
	myRTCUser, ok := twinAgora.GetUserById(int(peopleEntry.Uid))
	if !ok {
		return twinAgora.MakeError(ERR_UID_NOT_IN_MAP)
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
		return errors.New("您已经此频道中，请不要乱发消息")
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_EXECING {
		return twinAgora.MakeError(ERR_ROOM_STATUS_NOT_EXEC + RTCRoomInfo.Id)
	}

	for _, uid := range RTCRoomInfo.OnlineUids {
		if int(peopleEntry.Uid) == uid {
			continue
		}
		conn.SendMsgCompressByUid(peopleEntry.Uid, "SC_PeopleEntry", peopleEntry)
	}
	RTCRoomInfo.OnlineUids = append(RTCRoomInfo.OnlineUids, int(peopleEntry.Uid))

	return nil
}

//某用户离开了房间
func (twinAgora *TwinAgora) PeopleLeave(peopleLeaveRes pb.PeopleLeaveRes, conn *util.Conn) error {
	myRTCUser, ok := twinAgora.GetUserById(int(peopleLeaveRes.Uid))
	if !ok {
		return errors.New("GetUserById empty ")
	}
	RTCRoomInfo, err := twinAgora.GetRoomById(myRTCUser.RoomId)
	if err != nil {
		return errors.New("roomId NOT in map")
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
		return errors.New("您并不在此频道中，请不要乱发消息")
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_EXECING {
		return twinAgora.MakeError(ERR_ROOM_STATUS_NOT_EXEC + RTCRoomInfo.Id)
	}

	twinAgora.RoomEnd(myRTCUser.RoomId, RTC_ROOM_END_STATUS_USER_LEAVE)

	return nil
}

//被呼叫者，接收/拒绝 公共验证
func (twinAgora *TwinAgora) PeopleVote(callVote pb.CallVote) (error, *RTCRoom) {
	util.MyPrint("PeopleVote :", callVote)
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
		return twinAgora.MakeError("未发消息给你，请不要捣乱"), RTCRoomInfo
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
		return twinAgora.MakeError("您已经投过票了1，请不要捣乱"), RTCRoomInfo
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
		return twinAgora.MakeError("您已经投过票了2，请不要捣乱"), RTCRoomInfo
	}

	return nil, RTCRoomInfo
}

//被呼叫者，接收呼叫
func (twinAgora *TwinAgora) CallPeopleAccept(callVote pb.CallVote, conn *util.Conn) error {
	err, RTCRoomInfo := twinAgora.PeopleVote(callVote)
	if err != nil {
		return err
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return twinAgora.MakeError(ERR_ROOM_STATUS_NOT_CALL + callVote.RoomId)
	}

	rtcUser, _ := twinAgora.GetUserById(int(conn.UserId))
	rtcUser.RoomId = RTCRoomInfo.Id

	RTCRoomInfo.ReceiveUidsAccept = append(RTCRoomInfo.ReceiveUidsAccept, int(callVote.Uid))
	util.MyPrint("RTCRoomInfo.ReceiveUidsAccept:", RTCRoomInfo.ReceiveUidsAccept)
	//RTCRoomInfo.Uids = append(RTCRoomInfo.Uids, int(callVote.Uid))
	conn.SendMsgCompressByUid(int32(RTCRoomInfo.CallUid), "SC_CallPeopleAccept", callVote)
	//目前是1v1视频，只要有一个接收，即把房间状态标识为运行中，这里后期优化一下吧
	RTCRoomInfo.Status = RTC_ROOM_STATUS_EXECING
	return nil

}

//被呼叫者，拒绝呼叫
func (twinAgora *TwinAgora) CallPeopleDeny(callVote pb.CallVote, conn *util.Conn) error {
	err, RTCRoomInfo := twinAgora.PeopleVote(callVote)
	if err != nil {
		return err
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return twinAgora.MakeError(ERR_ROOM_STATUS_NOT_CALL + callVote.RoomId)
	}
	RTCRoomInfo.ReceiveUidsDeny = append(RTCRoomInfo.ReceiveUidsDeny, int(callVote.Uid))
	conn.SendMsgCompressByUid(int32(RTCRoomInfo.CallUid), "SC_CallPeopleDeny", callVote)
	twinAgora.RoomEnd(callVote.RoomId, RTC_ROOM_END_STATUS_DENY)
	return err

}
