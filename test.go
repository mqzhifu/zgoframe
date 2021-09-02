package main

import (
	"zgoframe/core/global"
	"zgoframe/model"
	"zgoframe/util"
)

func createDbTable(){
	mydb := util.NewDb(global.V.Gorm)
	mydb.CreateTable(&model.User{},&model.SmsLog{},&model.SmsRule{},&model.App{},&model.UserReg{} , &model.OperationRecord{})
	util.ExitPrint("init done.")
}