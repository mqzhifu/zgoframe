package model

//@description 声网回调存储
type AgoraCallbackRecord struct {
	MODEL
	NoticeId    string `json:"notice_id" form:"notice_id" db:"define:varchar(255);comment:通知 ID，标识来自业务服务器的一次事件通知;defaultValue:''"`
	ProductId   int    `json:"product_id" form:"product_id" db:"define:int;comment:业务Id,1rtc2旁路推流CDN3云端录制4Cloud Player5旁路推流- 服务端;defaultValue:0"`
	EventType   int    `json:"event_type" form:"event_type" db:"define:int;comment:事件类型ID;defaultValue:0"`
	ChannelName string `json:"channel_name" form:"channel_name" db:"define:varchar(255);comment:频道名称;defaultValue:''"`
	SessionId   string `json:"session_id" form:"session_id" db:"define:varchar(255);comment:声网返回的sid;defaultValue:''"`
	NotifyMs    string `json:"notify_ms" form:"notify_ms" db:"define:varchar(13);comment:对方推送时间;defaultValue:''"`
	Payload     string `json:"payload" form:"payload" db:"define:text;comment:详细内容"`
}

func (agoraCallbackRecord *AgoraCallbackRecord) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "声网回调记录"

	return m
}
