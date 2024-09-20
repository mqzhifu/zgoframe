package service

// @parse 报警发送渠道类型
const (
	ALERT_LEVEL_ALL      = -1 //全部
	ALERT_LEVEL_SMS      = 1  //短信
	ALERT_LEVEL_EMAIL    = 2  //邮件
	ALERT_LEVEL_FEISHU   = 4  //飞书
	ALERT_LEVEL_WEIXIN   = 8  //微信
	ALERT_LEVEL_DINGDING = 16 //钉钉
)

// @parse 报警发送方式类型
const (
	ALERT_SEND_SYNC  = 1 //同步
	ALERT_SEND_ASYNC = 2 //异步
)

// @parse http状态
const (
	HTTPD_RULE_STATE_INIT  = 1 //初始化中
	HTTPD_RULE_STATE_OK    = 2 //正常运行中
	HTTPD_RULE_STATE_CLOSE = 3 //已关闭
	HTTPD_RULE_STATE_UKNOW = 4 //未知
)

// @parse 呼叫类型
const (
	CALL_USER_PEOPLE_ALL     = 1 //呼叫所有人
	CALL_USER_PEOPLE_GROUP   = 2 //按照<小组>呼叫
	CALL_USER_PEOPLE_PROVIDE = 3 //用户自己指定呼叫的人
)
