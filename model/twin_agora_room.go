package model

//@description 声网回调存储
type TwinAgoraRoom struct {
	MODEL
	Channel           string `json:"channel" form:"channel" db:"define:varchar(255);comment:频道;defaultValue:''"`
	Status            int    `json:"status" form:"status" db:"define:int;comment:状态,1发起呼叫，2正常通话中，3已结束;defaultValue:0"`
	EndStatus         int    `json:"end_status" form:"end_status" db:"define:int;comment:结束的状态：(1)超时，(2)某一方退出,(3)某一方拒绝(4)发起方主动取消呼叫;defaultValue:0"`
	CallUid           int    `json:"call_uid" form:"call_uid" db:"define:varchar(255);comment:发起呼叫者;defaultValue:''"`
	ReceiveUids       string `json:"receive_uids" form:"receive_uids" db:"define:varchar(255);comment:接收呼叫者消息;defaultValue:''"`
	ReceiveUidsAccept string `json:"receive_uids_accept" form:"receive_uids_accept" db:"define:varchar(13);comment:被呼叫的用户IDS，接收了此次呼叫;defaultValue:''"`
	ReceiveUidsDeny   string `json:"receive_uids_deny" form:"receive_uids_deny" db:"define:text;comment:被呼叫的用户IDS，拒绝了此次呼叫"`
	Uids              string `json:"uids" form:"uids" db:"define:text;comment:ReceiveUidsAccept+CallUid"`
}

func (twinAgoraRoom *TwinAgoraRoom) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "AR远程呼叫,房间记录"

	return m
}
