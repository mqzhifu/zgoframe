package model

type AgoraCloudRecord struct {
	MODEL
	Uid              int    `json:"uid" form:"uid" db:"define:int;comment:用户ID;defaultValue:0"`
	ListenerAgoraUid int    `json:"listener_agora_uid" form:"listener_agora_uid" db:"define:int;comment:监听者的UID;defaultValue:0"`
	ChannelName      string `json:"channel_name" form:"channel_name" db:"define:varchar(255);comment:频道名称;defaultValue:''"`
	ResourceId       string `json:"resource_id" form:"resource_id" db:"define:text;comment:声网返回的rid"`
	SessionId        string `json:"session_id" form:"session_id" db:"define:varchar(255);comment:声网返回的sid;defaultValue:''"`
	Status           int    `json:"status" form:"status" db:"define:int;comment:0未知1已申请rid2已开始3已结束;defaultValue:0"`
	ServerStatus     int    `json:"server_status" form:"server_status" db:"define:int;comment:后端状态1未处理2已收到声网回调,开始合并视频3视频处理成功4处理异常;defaultValue:0"`
	StartTime        int    `json:"start_time" form:"start_time" db:"define:int;comment:开始录制时间;defaultValue:0"`
	EndTime          int    `json:"end_time" form:"end_time" db:"define:int;comment:结束录制时间时间;defaultValue:0"`
	ConfigInfo       string `json:"config_info" form:"config_info" db:"define:text;comment:请求声网,开始录制时设置的配置信息"`
	AcquireConfig    string `json:"acquire_config" form:"acquire_config" db:"define:varchar(255);comment:获取RID时的配置信息;defaultValue:''"`
	StopResInfo      string `json:"stop_res_info" form:"stop_res_info" db:"define:text;comment:请求声网,停止录制时返回的文件信息"`
	VideoUrl         string `json:"video_url" form:"video_url" db:"define:varchar(255);comment:最终录制好的视频的URL地址;defaultValue:''"`
	ErrLog           string `json:"err_log" form:"err_log" db:"define:text;comment:请求3方返回的错误信息"`
}

func (agoraCloudRecord *AgoraCloudRecord) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "声网录屏"

	return m
}
