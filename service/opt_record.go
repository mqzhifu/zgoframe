package service

import (
	"zgoframe/core/global"
	"zgoframe/model"
)

//@author: [granty1](https://github.com/granty1)
//@function: CreateSysOperationRecord
//@description: 创建记录
//@param: sysOperationRecord model.SysOperationRecord
//@return: err error

func CreateSysOperationRecord(sysOperationRecord model.OperationRecord) (err error) {
	err = global.V.Gorm.Create(&sysOperationRecord).Error
	return err
}
