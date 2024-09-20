package seed_business

// @parse 房间状态2
const (
	RTC_ROOM_STATUS_CALLING = 1 //房间状态：呼叫中
	RTC_ROOM_STATUS_EXECING = 2 //房间状态：运行中
	RTC_ROOM_STATUS_END     = 3 //房间状态：已结束
)

// @parse 房间结束类型
const (
	RTC_ROOM_END_STATUS_TIMEOUT_CALLING = 10 //房间结束状态标识：呼叫超时(也可能是连接断了)
	RTC_ROOM_END_STATUS_TIMEOUT_EXEC    = 11 //房间结束状态标识：运行超时(也可能是连接断了)
	RTC_ROOM_END_STATUS_CONN_CLOSE      = 2  //房间结束状态标识：用户退出
	RTC_ROOM_END_STATUS_DENY            = 3  //房间结束状态标识：被呼叫人拒绝
	RTC_ROOM_END_STATUS_CANCEL          = 4  //房间结束状态标识：呼叫者取消
	RTC_ROOM_END_STATUS_USER_LEAVE      = 4  //房间结束状态标识：用户离开
)

// @parse RTC推送消息类型
const (
	RTC_PUSH_MSG_EVENT_FD_CREATE_REPEAT = 400 //FD 重复
	RTC_PUSH_MSG_EVENT_UID_NOT_IN_MAP   = 401 // uid 不在MAP中
)
