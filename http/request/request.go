//http 请求公共处理
package request

type SystemConfig struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type RLoginThird struct {
	Register
	ThirdId      string
	PlatformType int
}
