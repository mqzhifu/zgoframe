package model

// @description 电子邮件日志
type ConnRecord struct {
	MODEL
	Uid            int32  `json:"uid" db:"define:int;comment:用户ID;defaultValue:0"`                    //
	Status         int    `json:"status" db:"define:tinyint(1);comment:状态;defaultValue:0"`            //
	ProtocolType   int32  `json:"protocol_type" db:"define:tinyint(1);comment:传输协议类型;defaultValue:0"` //
	ContentType    int32  `json:"content_type" db:"define:tinyint(1);comment:传输内容类型;defaultValue:0"`  //
	SessionId      string `json:"session_id" db:"define:varchar(50);comment:会话ID;defaultValue:''"`    //
	CloseType      int    `json:"close_type" db:"define:tinyint(1);comment:关闭类型;defaultValue:0"`      //
	CloseTime      int64  `json:"close_time" db:"define:int;comment:关闭时间;defaultValue:0"`             //
	TotalOutputNum int    `json:"total_output_num" db:"define:int;comment:发送消息次数;defaultValue:0"`     //
	TotalInputNum  int    `json:"total_input_num" db:"define:int;comment:接收消息次数;defaultValue:0"`      //
	Rtt            int    `json:"rtt" db:"define:int;comment:延迟时间;defaultValue:0"`                    //
}

func (connRecord *ConnRecord) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "长连接的记录"

	return m
}
