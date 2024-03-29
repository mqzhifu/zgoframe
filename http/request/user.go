package request

// @description 获取访问权限的 token
type AccessToken struct {
	Grant     string `json:"grant"`     //保留字，先不用
	Timestamp int    `json:"timestamp"` //unix 时间戳 10位的
	Sign      string `json:"sign"`      //签名:md5(SecretKey+Timestamp+Access)
}

// @description 注册信息
type Register struct {
	ProjectId int    `json:"project_id" form:"project_id"` //项目Id
	Username  string `json:"userName" form:"username"`     //用户名
	Password  string `json:"passWord" form:"password"`     //登陆密码 转md5存储
	NickName  string `json:"nickName" form:"nick_name" `   //昵称
	HeaderImg string `json:"headerImg" form:"header_img" ` //头像地址
	Sex       int    `json:"sex" form:"sex"`               //性别
	Channel   int    `json:"channel"`                      //来源渠道
	Birthday  int    `json:"birthday" form:"birthday"`     //生日
	Recommend string `json:"recommend" form:"varchar(50)"` //推荐人
	Guest     int    `json:"guest"  `                      //类型,1普通2游客
	ThirdType int    `json:"third_type" `                  //三方平台类型
	ThirdId   string `json:"third_id"`                     //三方平台ID
	ConfirmPs string `json:"confirm_ps"`                   //确认密码
	Test      int    `json:"test"`                         //是否为测试用户1是2否
	ExtDiy    string `json:"ext_diy"`                      //自定义用户属性，暂未实现
}

// @description 注册信息 - 通过手机号
type RegisterSms struct {
	//ProjectId   int    `json:"project_id"`    //项目Id
	Mobile      string `json:"mobile"`        //手机号
	SmsAuthCode string `json:"sms_auth_code"` //短信验证码
	SmsRuleId   int    `json:"sms_rule_id"`   //短信类型，登陆/注册
	Captcha     string `json:"captcha"`       //图片验证码
	CaptchaId   string `json:"captchaId"`     //图片验证码ID
}

// @descriptionw 绑定手机号
type BindMobile struct {
	ProjectId   int    `json:"project_id"`    //项目Id
	Mobile      string `json:"mobile"`        //手机号
	SmsAuthCode string `json:"sms_auth_code"` //短信验证码
	RuleId      int    `json:"rule_id"`       //短信类型，登陆/注册
}

// @descriptionw 绑定邮箱
type BindEmail struct {
	ProjectId   int    `json:"project_id"`    //项目Id
	Email       string `json:"email"`         //邮箱号
	SmsAuthCode string `json:"sms_auth_code"` //邮件验证码
	RuleId      int    `json:"rule_id"`       //邮件类型，登陆/注册
}

// @description 修改用户基础信息
type SetUserInfo struct {
	NickName  string `json:"nickName" form:"nick_name" `   //昵称
	HeaderImg string `json:"headerImg" form:"header_img" ` //头像地址
	Sex       int    `json:"sex" form:"sex"`               //性别
	Birthday  int    `json:"birthday" form:"birthday"`     //生日
}

// @description 正常登陆，需要用户名密码
type Login struct {
	Username  string `json:"username" form:"username"`   //用户名：普通字符串、手机号、邮箱
	Password  string `json:"password" form:"password"`   //密码
	Captcha   string `json:"captcha"  form:"captcha"`    //验证码
	CaptchaId string `json:"captchaId" form:"captchaId"` //验证码-ID
}

// @description 3方登陆
type RLoginThird struct {
	Register
	ThirdId      string `json:"third_id"`      //3方平台用户ID，此值必填
	PlatformType int    `json:"platform_type"` //3方平台类型，如：微信、QQ、facebook、抖音，此值必填
}

//// @description  3方平台登陆
//type LoginThird struct {
//	Code      string `json:"code"`
//	Platform  string `json:"platform"`
//	Captcha   string `json:"captcha"`
//	CaptchaId string `json:"captchaId"`
//}

// @description 短信登陆
type LoginSMS struct {
	//Code        string `json:"code"`
	Captcha     string `json:"captcha"`       //图片验证码
	CaptchaId   string `json:"captchaId"`     //图片验证码ID
	Mobile      string `json:"mobile"`        //手机号
	SmsAuthCode string `json:"sms_auth_code"` //手机验证码
	SmsRuleId   int    `json:"sms_rule_id"`   //短信类型，登陆/注册
}

// @description 发送验证码
type SendSMS struct {
	RuleId     int               `json:"rule_id"`     //配置规则的ID
	ReplaceVar map[string]string `json:"replace_var"` //邮件内容模块中变量替换
	Receiver   string            `json:"receiver"`    //接收者-手机号
	SendUid    int               `json:"send_uid"`    //发送者ID，管理员是9999，未知8888
	SendIp     string            `json:"send_ip"`     //发送者IP，如为空系统默认取：请求方的IP,最好给真实的，一但被刷，会使用此值
	Captcha    string            `json:"captcha"`     //验证码
	CaptchaId  string            `json:"captchaId"`   //获取验证码时拿到的Id
}

// @description 图片验证码
type Captcha struct {
	Width  int `json:"width"`  //图片宽度，默认：240，最大：1000
	Height int `json:"height"` //图片高度，默认：80，最大：1000
}

// @description 发送邮件
type SendEmail struct {
	RuleId     int               `json:"rule_id"`     //配置规则的ID
	ReplaceVar map[string]string `json:"replaceVar"`  //邮件内容模块中变量替换
	Receiver   string            `json:"receiver"`    //接收者（email格式）
	CarbonCopy []string          `json:"carbon_copy"` //抄送（email格式），可以是多人
	SendUid    int               `json:"send_uid"`    //发送者ID，管理员:9999，未知:8888
	SendIp     string            `json:"send_ip"`     //发送者IP，如为空系统默认取：请求方的IP,最好给真实的，一但被刷，会使用此值
}

// @description 发送站内信
type SendMail struct {
	RuleId     int               `json:"rule_id"`    //配置规则的ID
	ReplaceVar map[string]string `json:"replaceVar"` //邮件内容模块中变量替换
	Receiver   string            `json:"receiver"`   //接收者: uid or grpuId or tagId or uids
	SendUid    int               `json:"send_uid"`   //发送者ID，管理员是9999，未知8888
	SendIp     string            `json:"send_ip"`    //发送者IP，如为空系统默认取：请求方的IP,最好给真实的，一但被刷，会使用此值
	SendTime   int               `json:"send_time"`  //定时发送，unixStamp 必须大于当前时间
}

// @description 站内信列表
type MailList struct {
	BoxType      int `json:"box_type"`      //1收件箱2发件箱4全部
	ReceiverRead int `json:"receiver_read"` //1接收者已读2接收者未读
	ReceiverDel  int `json:"receiver_del"`  //1接收者已删除2接收者未删除
	Expire       int `json:"expire"`        //1消息已过期2消息未过期
	PageInfo         //分页
}

// @description 站内信一条消息详情
type MailInfo struct {
	Id               int `json:"id"`
	AutoReceiverRead int `json:"auto_receiver_read"` //自动更新为：接收者已读
	AutoReceiverDel  int `json:"auto_receiver_del"`  //自动更新为：接收者已删除
}

// @description 设置/修改密码
type SetPassword struct {
	Password           string `json:"password"`             //旧密码
	NewPassword        string `json:"newPassword"`          //新密码
	NewPasswordConfirm string `json:"new_password_confirm"` //新密码确认
}

// @description 修改密码
type RestPasswordSms struct {
	Mobile             string `json:"mobile"`               //手机号
	SmsAuthCode        string `json:"sms_auth_code"`        //短信验证码
	SmsRuleId          int    `json:"sms_rule_id"`          //短信类型，登陆/注册
	Password           string `json:"password"`             //旧密码
	NewPassword        string `json:"newPassword"`          //新密码
	NewPasswordConfirm string `json:"new_password_confirm"` //新密码确认
}

// @description 检测手机号是否已存在
type CheckMobileExist struct {
	Mobile  string `json:"mobile"`  //手机号
	Purpose int    `json:"purpose"` //用途,1注册2找回密码3修改密码4登陆
}

// @description 检测用户名是否已存在
type CheckUsernameExist struct {
	Username string `json:"username"` //用户名
	Purpose  int    `json:"purpose"`  //用途,1注册2找回密码3修改密码4登陆
}

// @description 检测邮件是否已存在
type CheckEmailExist struct {
	Email   string `json:"email"`   //邮箱
	Purpose int    `json:"purpose"` //用途,1注册2找回密码3修改密码4登陆
}
