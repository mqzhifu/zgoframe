package initialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"zgoframe/core/global"
	"zgoframe/util"
)

func GetNewGorm() (*gorm.DB,error) {
	switch global.C.System.DbType {
	case "mysql":
		return GormMysql()
	default:
		return GormMysql()
	}
}

func GormMysql() (*gorm.DB,error) {
	m := global.C.Mysql
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Ip + ":" + m.Port + ")/" + m.DbName + "?" + m.Config
	fmt.Println("GormMysql:"+dsn)
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(m.LogMode)); err != nil {
		fmt.Println("MySQL启动异常", err.Error())
		return nil,err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)


		return db,nil
	}
}

func GormShutdown(){
	db , _ := global.V.Gorm.DB()
	db.Close()
}

func gormConfig(mod bool) *gorm.Config {
	var config = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true,NamingStrategy: schema.NamingStrategy{SingularTable: true}}
	config.Logger = util.Default.LogMode(logger.Info)
	//switch global.G.Config.Mysql.LogZap {
	//case "silent", "Silent":
	//	config.Logger = internal.Default.LogMode(logger.Silent)
	//case "error", "Error":
	//	config.Logger = internal.Default.LogMode(logger.Error)
	//case "warn", "Warn":
	//	config.Logger = internal.Default.LogMode(logger.Warn)
	//case "info", "Info":
	//	config.Logger = internal.Default.LogMode(logger.Info)
	//case "zap", "Zap":
	//	config.Logger = internal.Default.LogMode(logger.Info)
	//default:
	//	if mod {
	//		config.Logger = internal.Default.LogMode(logger.Info)
	//		break
	//	}
	//	config.Logger = internal.Default.LogMode(logger.Silent)
	//}
	return config
}

func TestGorm(){
	//db := util.NewDb(global.V.Gorm)
	//userModel := model.User{}
	//_ ,err := db.GetRowById(&userModel,1)
	////user := userInterface.(*model.User)
	//util.MyPrint(userModel.Username,userModel.Id,err)
	//
	//
	//userModel2 := []model.User{}
	//_ ,err = db.GetRowByIds(&userModel2,[]int{1,2,3})
	//util.MyPrint(userModel2,err)
	//
	//
	//userModel3 := model.User{}
	//_ ,err = db.GetRow(&userModel3," username = 'mqzhifu@sina.com' " )
	//util.MyPrint(userModel3,err)
	//
	//
	//userModel4 := []model.User{}
	//query := util.DbQueryListPara{
	//	Where: " id = 1",
	//}
	//_ ,err = db.GetList(&userModel4,query)
	//util.MyPrint(userModel2,err)
	//
	//util.ExitPrint(123123213)


	//user3 ,err := userModel.GetRow(" username = 'mqzhifu@sina.com'")
	//util.MyPrint(user3.Id,user3.Username,err)
	//
	//
	//user4, err  := userModel.GetRowByIds([]int{1,2})
	//util.MyPrint(user4,err)
	//
	//user2 ,err := userModel.GetRowById(程序员代码面试指南 IT名企算法与数据结构题目最优解 ,左程云著 ,P51310000)
	//util.ExitPrint(user2.Id,user2.Username,err)
}