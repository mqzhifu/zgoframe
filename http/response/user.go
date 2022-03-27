package httpresponse

import (
	"zgoframe/model"
)

type SysUserResponse struct {
	User model.User `json:"user"` //用户基础信息
}

//@description 登陆成功返回结果
type LoginResponse struct {
	//登陆成功结果响应
	User       model.User `json:"user"`         //用户基础信息
	Token      string     `json:"token"`        //生成的token
	ExpiresAt  int64      `json:"expires_at"`   //token失效时间
	IsNewToken bool       `json:"is_new_token"` //重复登陆也可以成功，但返回的是旧的TOKEN，非新生成token
	IsNewReg   bool       `json:"is_new_reg"`   //3方登陆时，为了简化操作，如果没注册将自动注册
}
