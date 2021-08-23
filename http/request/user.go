package request

import uuid "github.com/satori/go.uuid"

// User register structure
type Register struct {
	AppId 		int `json:"app_id" form:"app_id"`
	Username    string `json:"userName" form:"username"`
	Password    string `json:"passWord" form:"password"`
	NickName    string `json:"nickName" gorm:"default:'QMPlusUser'" form:"nick_name"`
	HeaderImg   string `json:"headerImg" gorm:"default:'http://www.henrongyi.top/avatar/lufu.jpg'" form:"header_img"`
	AuthorityId string `json:"authorityId" gorm:"default:888"`
}

// User login structure
type Login struct {
	AppId 	  int `json:"app_id" form:"app_id"`
	Username  string `json:"username" form:"username"`
	Password  string `json:"password" form:"password"`
	Captcha   string `json:"captcha"  form:"captcha"`
	CaptchaId string `json:"captchaId" form:"captchaId"`
}

type LoginSMS struct {
	AppId 	  string `json:"app_id"`
	Code      string `json:"code"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

type LoginThird struct {
	AppId 	  string `json:"app_id"`
	Code  	  string `json:"Code"`
	Platform  string `json:"platform"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}
// Base sendSMS structure
type SendSMS struct {
	AppId 	  string `json:"app_id"`
	Code  	  string `json:"code"`
	Mobile    string `json:"mobile"`
	RuleId 		int 	`json:"rule_id"`
}

// Modify password structure
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

