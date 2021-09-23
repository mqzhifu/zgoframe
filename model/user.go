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
	Username    string       	`json:"username" db:"define:varchar(50);comment:用户登录名;defaultValue:''"`
	Password    string       	`json:"-" db:"define:varchar(50);comment:用户登录密码;defaultValue:''"`
	NickName    string       	`json:"nick_name" db:"define:varchar(50);comment:用户昵称;defaultValue:''" `
	AuthorityId string       	`json:"authority_id" db:"define:varchar(50);comment:用户角色ID;defaultValue:''"`
	Mobile 		string 			`json:"mobile" db:"define:varchar(50);comment:手机号;defaultValue:''"`
	Email 		string 			`json:"email" db:"define:varchar(50);comment:邮箱;defaultValue:''"`
	Robot 		int				`json:"robot" db:"define:tinyint(1);comment:机器人;defaultValue:0"`
	Status 		int				`json:"status" db:"define:tinyint(1);comment:状态;defaultValue:0"`

	HeaderImg   string       	`json:"headerImg" gorm:"" db:"define:varchar(50);comment:用户头像;defaultValue:''"`
	Authority   SysAuthority 	`json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;" db:"define:varchar(50);comment:用户角色;defaultValue:''"`

	//Type 		int 			`json:"reg_type" db:"define:tinyint(1);comment:类型,1普通2游客;defaultValue:0"`
}

func(user *User)Db(){

}

func(user *User) TableOptions()map[string]string{
	m := make(map[string]string)
	m["comment"] = "用户表"

	return m
}
func(user *User)Count(){

}



////
////根据主键ID查找一条记录
//func (user *User) GetRowById(Id int)(*User, error){
//	userRow := &User{}
//	//这里没必要用first last ，因为它两会排序，既然ID是主键不可能出现重复|多条记录，也就不需要排序
//	txDb := Db.Take( userRow,Id)
//	//gorm.ErrRecordNotFound
//	return userRow, txDb.Error
//	//err := global.V.Gorm.Where("id = ?", Id).Take(userRow).Error
//	//return userRow, err
//}
////根据主键ID查找一组记录
//func (user *User) GetRowByIds(Ids []int)([]*User, error){
//	var users []*User
//	txDb := Db.Find( &users,Ids)
//	return users, txDb.Error
//}
//
//
//func (user *User) GetRow(where string)(*User, error){
//	userRow := &User{}
//	err := Db.Where(where).First(userRow).Error
//	return userRow, err
//}
//
//type GetListPara struct {
//	Where	string
//	Fields 	string
//	Order 	string
//	Group   stringΩ
//	Limit 	int
//	Offset	int
//}
//
//func (user *User) GetList(para GetListPara )([]*User, error){
//	var users []*User
//	db := global.V.Gorm.Find(para.Where)
//	if para.Fields != ""{
//		db.Select(para.Fields)
//	}
//
//	if para.Group != ""{
//		db.Order(para.Group)
//	}
//
//	if para.Order != ""{
//		db.Order(para.Order)
//	}
//
//	if para.Limit >= 0 {
//		db.Limit(para.Offset)
//	}
//
//	if para.Limit  >= 0{
//		db.Limit(para.Limit )
//	}
//
//	err := db.Find(users).Error
//	return users, err
//}
//
//func (user *User) UpRowById(Id int)(*User, error){
//
//}
//
//func (user *User) UpRowByIds(Id int)(*User, error){
//
//}
//
//func (user *User) UpRow(){
//
//}
//
//func (user *User) UpRows(){
//
//}
//
//func (user *User) AddOne(){
//
//}
//
