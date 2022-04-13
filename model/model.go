//db 模型
package model

import (
	"gorm.io/gorm"
)

type MODEL struct {
	Id        int            `json:"id" gorm:"primarykey" db:"define:int;primarykey:true;unsigned:true;autoIncrement:true;comment:主键自增ID"` //自增ID
	CreatedAt int64          `json:"created_at" db:"comment:创建时间;define:bigint;defaultValue:0"`                                            //创建时间
	UpdatedAt int64          `json:"updated_at" db:"comment:最后更新时间;define:bigint;defaultValue:0"`                                          //最后更新时间
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index" json:"-" db:"comment:是否删除;define:bigint;index:true;defaultValue:null"`        //是否删除
}


var Db *gorm.DB
