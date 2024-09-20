package msg_center

// @parse 站内邮件-是否失效
const (
	MAIL_ADMIN_USER_UID = 9999 //管理员默认UID，主要用于：确定发送者的UID
	
	MAIL_EXPIRE_TRUE  = 1 //是
	MAIL_EXPIRE_FALSE = 1 //否
)

// @parse 站内邮件-发送类型
const (
	MAIL_PEOPLE_PERSON = 1 //单发
	MAIL_PEOPLE_ALL    = 2 //群发
	MAIL_PEOPLE_GROUP  = 3 //指定group
	MAIL_PEOPLE_TAG    = 4 //指定tag
	MAIL_PEOPLE_UIDS   = 5 //指定UIDS
)

// @parse 站内邮件-消息收件箱
const (
	MAIL_IN_BOX     = 1 //收件箱
	MAIL_OUT_BOX    = 2 //发件箱
	MAIL_IN_DEL_BOX = 3 //已删除箱子
	MAIL_ALL_BOX    = 4 //全部
)

// @parse 站内邮件-消息接收是否已读
const (
	RECEIVER_READ_TRUE  = 1 //是
	RECEIVER_READ_FALSE = 2 //否
)

// @parse 站内邮件-消息是否删除
const (
	RECEIVER_DEL_TRUE  = 1 //是
	RECEIVER_DEL_FALSE = 2 //否
)
