package model

type UserThird struct {
	MODEL
	Uid          int    `json:"uid" db:"define:int;comment:用户ID;defaultValue:0"`
	ThirdId      string `json:"third_id" db:"define:varchar(50);comment:三方平台(登陆)用户ID;defaultValue:''"`
	PlatformType int    `json:"sex" db:"define:tinyint(1);comment:3方平台类型;defaultValue:0"`
}

func (userThird *UserThird) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "用户3方注册平台"

	return m
}
