package service

import (
	"gorm.io/gorm"
	"strconv"
	"zgoframe/model"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

type TwinAgora struct {
	Gorm *gorm.DB
}

func NewTwinAgora(Gorm *gorm.DB) *TwinAgora {
	twinAgora := new(TwinAgora)
	twinAgora.Gorm = Gorm
	return twinAgora
}

const (
	USER_ROLE   = 1 //普通用户
	USER_DOCTOR = 2 //专家
	USER_ALL    = 3 //所有人
)

func (twinAgora *TwinAgora) CallPeople(callPeopleReq pb.CallPeopleReq, conn *util.Conn) {
	callPeopleRes := pb.CallPeopleRes{}

	if callPeopleReq.Uid <= 0 {
		callPeopleRes.ErrCode = 1
		callPeopleRes.ErrMsg = ""
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "aaaaaa", callPeopleReq)
		return
	}

	if callPeopleReq.PeopleType != int32(USER_DOCTOR) {
		callPeopleRes.ErrCode = 1
		callPeopleRes.ErrMsg = ""
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "aaaaaa", callPeopleReq)
		return
	}
	var userDoctorList []model.User
	err := twinAgora.Gorm.Where(" role =  ?", USER_DOCTOR).Find(&userDoctorList).Error
	if err != nil {
		callPeopleRes.ErrCode = 1
		callPeopleRes.ErrMsg = ""
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "aaaaaa", callPeopleReq)
		return
	}

	if callPeopleReq.TargetUid > 0 {
		callPeopleRes.ErrCode = 1
		callPeopleRes.ErrMsg = ""
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "aaaaaa", callPeopleReq)
		return
	}

	var onlineUserDoctorList []model.User

	for _, userConn := range conn.ConnManager.Pool {
		for _, user := range userDoctorList {
			if userConn.UserId == int32(user.Id) && userConn.Status == util.CONN_STATUS_EXECING {
				onlineUserDoctorList = append(onlineUserDoctorList, user)
			}
		}
	}

	if len(onlineUserDoctorList) <= 0 {
		callPeopleRes.ErrCode = 1
		callPeopleRes.ErrMsg = ""
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "aaaaaa", callPeopleReq)
		return
	}

	for _, user := range onlineUserDoctorList {
		pushMsgRes := pb.PushMsgRes{}
		pushMsgRes.MsgType = 1
		pushMsgRes.Content = strconv.Itoa(int(callPeopleReq.Uid)) + " 呼叫 视频连接...请进入频道:" + callPeopleReq.Channel
		conn.SendMsgCompressByUid(int32(user.Id))
	}
	return
}
