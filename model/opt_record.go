package model

// 如果含有time.Time 请自行import time包
type OperationRecord struct {
	MODEL
	Ip           string `json:"ip" form:"ip" db:"define:varchar(50);comment:ip;defaultValue:''"`
	Method       string `json:"method" form:"method" db:"define:varchar(50);comment:get|post;defaultValue:''"`
	Path         string `json:"path" form:"path" db:"define:varchar(50);comment:请求路径;defaultValue:''"`
	Status       int    `json:"status" form:"status" db:"define:int;comment:请求状态;defaultValue:0"`
	Latency      int    `json:"latency" form:"latency" db:"define:int;comment:延迟;defaultValue:0"`
	Agent        string `json:"agent" form:"agent" db:"define:text;comment:useragent"`
	ErrorMessage string `json:"error_message" form:"error_message" db:"define:varchar(255);comment:错误信息;defaultValue:''"`
	Body         string `json:"body" form:"body" db:"define:text;comment:请求内容"`
	Resp         string `json:"resp" form:"resp" db:"define:text;comment:返回结果"`
	UserID       int    `json:"user_id" form:"user_id" db:"define:int;comment:userId;defaultValue:0"`
}

func (operationRecord *OperationRecord) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "请求日志"

	return m
}
