//db 模型
package model

import (
	"gorm.io/gorm"
)

type MODEL struct {
	Id        int            `json:"id" gorm:"primarykey" db:"define:int;primarykey:true;unsigned:true;autoIncrement:true;comment:主键自增ID"`
	CreatedAt int64          `json:"created_at" db:"comment:创建时间;define:bigint;defaultValue:0"`
	UpdatedAt int64          `json:"updated_at" db:"comment:最后更新时间;define:bigint;defaultValue:0"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index" json:"-" db:"comment:是否删除;define:bigint;index:true;defaultValue:null"`
}

//type ModelList struct {
//	User User
//	App App
//	SmsRule SmsRule
//	SmsLog SmsLog
//	UserReg UserReg
//	OperationRecord OperationRecord
//}
//
//var V = ModelList{
//	User: User{},
//	App: App{},
//	SmsRule: SmsRule{},
//	SmsLog: SmsLog{},
//	UserReg: UserReg{},
//	OperationRecord : OperationRecord{},
//}

var Db *gorm.DB
