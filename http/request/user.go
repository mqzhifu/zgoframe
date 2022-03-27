package request

//@description 注册信息
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

	ExtDiy string `json:"ext_diy"` //自定义用户属性，暂未实现
}

//@description 注册信息 - 通过手机号
type RegisterSms struct {
	ProjectId   int    `json:"project_id"`    //项目Id
	Mobile      string `json:"mobile"`        //手机号
	SmsAuthCode string `json:"sms_auth_code"` //短信验证码
	SmsRuleId   int    `json:"sms_rule_id"`   //短信类型，登陆/注册
}

//@descriptionw 绑定手机号
type BindMobile struct {
	ProjectId   int    `json:"project_id"`    //项目Id
	Mobile      string `json:"mobile"`        //手机号
	SmsAuthCode string `json:"sms_auth_code"` //短信验证码
	RuleId      int    `json:"rule_id"`       //短信类型，登陆/注册
}

//@descriptionw 绑定邮箱
type BindEmail struct {
	ProjectId   int    `json:"project_id"`    //项目Id
	Email       string `json:"email"`         //邮箱号
	SmsAuthCode string `json:"sms_auth_code"` //短信验证码
	RuleId      int    `json:"rule_id"`       //短信类型，登陆/注册
}

//@description 修改用户基础信息
type SetUserInfo struct {
	NickName  string `json:"nickName" form:"nick_name" `   //昵称
	HeaderImg string `json:"headerImg" form:"header_img" ` //头像地址
	Sex       int    `json:"sex" form:"sex"`               //性别
	Birthday  int    `json:"birthday" form:"birthday"`     //生日
}

//@description 正常登陆，需要用户名密码
type Login struct {
	Username  string `json:"username" form:"username"`   //用户名：普通字符串、手机号、邮箱
	Password  string `json:"password" form:"password"`   //密码
	Captcha   string `json:"captcha"  form:"captcha"`    //验证码
	CaptchaId string `json:"captchaId" form:"captchaId"` //验证码-ID
}

//@description 短信登陆
type LoginSMS struct {
	Code        string `json:"code"`
	Captcha     string `json:"captcha"`
	CaptchaId   string `json:"captchaId"`
	Mobile      string `json:"mobile"`
	SmsAuthCode string `json:"sms_auth_code"`
	SmsRuleId   int    `json:"sms_rule_id"` //短信类型，登陆/注册
}

//@description  3方平台登陆
type LoginThird struct {
	Code      string `json:"Code"`
	Platform  string `json:"platform"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

//@description 发送验证码
type SendSMS struct {
	RuleId     int               `json:"rule_id"`     //配置规则的ID
	ReplaceVar map[string]string `json:"replace_var"` //邮件内容模块中变量替换
	Receiver   string            `json:"receiver"`    //接收者，email格式
	SendUid    int               `json:"send_uid"`    //发送者ID，管理员是9999，未知8888
	SendIp     string            `json:"send_ip"`     //发送者IP，如为空系统默认取：请求方的IP,最好给真实的，一但被刷，会使用此值
	Captcha    string            `json:"captcha"`     //验证码
	CaptchaId  string            `json:"captchaId"`   //获取验证码时拿到的Id
}

//@description 发送邮件
type SendEmail struct {
	RuleId     int               `json:"rule_id"`     //配置规则的ID
	ReplaceVar map[string]string `json:"replaceVar"`  //邮件内容模块中变量替换
	Receiver   string            `json:"receiver"`    //接收者，email格式
	CarbonCopy []string          `json:"carbon_copy"` //抄送，，email格式
	SendUid    int               `json:"send_uid"`    //发送者ID，管理员是9999，未知8888
	SendIp     string            `json:"send_ip"`     //发送者IP，如为空系统默认取：请求方的IP,最好给真实的，一但被刷，会使用此值
}

//@description 设置/修改密码
type SetPassword struct {
	Password           string `json:"password"`             //旧密码
	NewPassword        string `json:"newPassword"`          //新密码
	NewPasswordConfirm string `json:"new_password_confirm"` //新密码确认
}

//@description 修改密码
type RestPasswordSms struct {
	Mobile             string `json:"mobile"`               //手机号
	SmsAuthCode        string `json:"sms_auth_code"`        //短信验证码
	SmsRuleId          int    `json:"sms_rule_id"`          //短信类型，登陆/注册
	Password           string `json:"password"`             //旧密码
	NewPassword        string `json:"newPassword"`          //新密码
	NewPasswordConfirm string `json:"new_password_confirm"` //新密码确认
}

type CheckMobileExist struct {
	Mobile  string `json:"mobile"`  //手机号
	Purpose int    `json:"purpose"` //用途,1注册2找回密码3修改密码4登陆
}

type CheckUsernameExist struct {
	Username string `json:"username"` //用户名
	Purpose  int    `json:"purpose"`  //用途,1注册2找回密码3修改密码4登陆
}

type CheckEmailExist struct {
	Email   string `json:"email"`   //邮箱
	Purpose int    `json:"purpose"` //用途,1注册2找回密码3修改密码4登陆
}
