package model

// 日志记录，主要是HTTP请求，正常有日志，且用MYSQL存也不太合适，所以正常不用
type OperationRecord struct {
	MODEL
	Ip           string `json:"ip" form:"ip" db:"define:varchar(50);comment:ip;defaultValue:''"`
	Method       string `json:"method" form:"method" db:"define:varchar(50);comment:get|post|put|delete;defaultValue:''"`
	Path         string `json:"path" form:"path" db:"define:varchar(50);comment:uri请求路径;defaultValue:''"`
	Status       int    `json:"status" form:"status" db:"define:int;comment:请求状态;defaultValue:0"`
	Latency      int    `json:"latency" form:"latency" db:"define:int;comment:延迟;defaultValue:0"`
	Agent        string `json:"agent" form:"agent" db:"define:text;comment:useragent"`
	ErrorMessage string `json:"error_message" form:"error_message" db:"define:varchar(255);comment:错误信息;defaultValue:''"`
	Body         string `json:"body" form:"body" db:"define:text;comment:请求内容"`
	Resp         string `json:"resp" form:"resp" db:"define:text;comment:返回结果"`
	Uid          int    `json:"uid" form:"uid" db:"define:int;comment:用户Id;defaultValue:0"`
}

func (operationRecord *OperationRecord) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "请求日志"

	return m
}
