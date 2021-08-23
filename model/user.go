package model

import (
	"github.com/satori/go.uuid"
)

//type User struct {
//	global.MODEL
//	UUID        uuid.UUID    	`json:"uuid" 		gorm:"comment:UUID"	db:"comment:uuid"`
//	AppId       int   			`json:"app_id" 		gorm:"comment:app_id"`
//	Sex       	int   			`json:"sex" 		gorm:"comment:性别1男2女"`
//	Birthday    int   			`json:"birthday" 	gorm:"comment:出生日期"`
//	Username    string       	`json:"userName" 	gorm:"comment:用户登录名"`
//	Password    string       	`json:"-"  			gorm:"comment:用户登录密码"`
//	NickName    string       	`json:"nickName" 	gorm:"comment:用户昵称" `
//	AuthorityId string       	`json:"authorityId" gorm:"comment:用户角色ID"`
//	Mobile 		string 			`json:"mobile" 		gorm:"comment:手机号"`
//	Email 		string 			`json:"email" 		gorm:"comment:邮箱"`
//	Type 		int 			`json:"reg_type" 	gorm:"comment:类型,1普通2游客"`
//	Robot 		int				`json:"robot" 		gorm:"comment:机器人"`
//	Status 		int				`json:"status" 		gorm:"comment:状态"`
//
//	HeaderImg   string       	`json:"headerImg" 	gorm:"default:http://qmplusimg.henrongyi.top/head.png;comment:用户头像"`
//	Authority   SysAuthority 	`json:"authority" 	gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
//}



type User struct {
	MODEL
	Uuid        uuid.UUID    	`json:"uuid" db:"define:varchar(50);comment:uuid;unique:uuid;index:uuid;defaultValue:''"`
	AppId       int   			`json:"app_id" db:"define:tinyint(1);comment:app_id;defaultValue:0"`
	Sex       	int   			`json:"sex" db:"define:tinyint(1);comment:性别1男2女;defaultValue:0"`
	Birthday    int   			`json:"birthday" db:"define:int;comment:出生日期;defaultValue:0"`
	Username    string       	`json:"userName" db:"define:varchar(50);comment:用户登录名;defaultValue:''"`
	Password    string       	`json:"-" db:"define:varchar(50);comment:用户登录密码;defaultValue:''"`
	NickName    string       	`json:"nickName" db:"define:varchar(50);comment:用户昵称;defaultValue:''" `
	AuthorityId string       	`json:"authorityId" db:"define:varchar(50);comment:用户角色ID;defaultValue:''"`
	Mobile 		string 			`json:"mobile" db:"define:varchar(50);comment:手机号;defaultValue:''"`
	Email 		string 			`json:"email" db:"define:varchar(50);comment:邮箱;defaultValue:''"`
	Robot 		int				`json:"robot" db:"define:tinyint(1);comment:机器人;defaultValue:0"`
	Status 		int				`json:"status" db:"define:tinyint(1);comment:状态;defaultValue:0"`

	HeaderImg   string       	`json:"headerImg" gorm:"" db:"define:varchar(50);comment:用户头像;defaultValue:''"`
	Authority   SysAuthority 	`json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;" db:"define:varchar(50);comment:用户角色;defaultValue:''"`

	//Type 		int 			`json:"reg_type" db:"define:tinyint(1);comment:类型,1普通2游客;defaultValue:0"`
}


func(user *User) TableOptions()map[string]string{
	m := make(map[string]string)
	m["comment"] = "用户表"

	return m
}
