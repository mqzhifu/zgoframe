package request

import uuid "github.com/satori/go.uuid"

// 用户注册
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

	ExtDiy string `json:"ext_diy"` //自定义用户属性，暂未实现
}

// 正常登陆，需要用户名密码
type Login struct {
	Username  string `json:"username" form:"username"`   //用户名：普通字符串、手机号、邮箱
	Password  string `json:"password" form:"password"`   //密码
	Captcha   string `json:"captcha"  form:"captcha"`    //验证码
	CaptchaId string `json:"captchaId" form:"captchaId"` //验证码-ID
}

//短信登陆
type LoginSMS struct {
	Code      string `json:"code"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

//3方平台登陆
type LoginThird struct {
	Code      string `json:"Code"`
	Platform  string `json:"platform"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

// 发送验证码
type SendSMS struct {
	AppId  string `json:"app_id"`
	Code   string `json:"code"`
	Mobile string `json:"mobile"`
	RuleId int    `json:"rule_id"`
}

// 修改密码
type ChangePasswordStruct struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

// Modify  user's auth structure
type SetUserAuth struct {
	UUID        uuid.UUID `json:"uuid"`
	AuthorityId string    `json:"authorityId"`
}
