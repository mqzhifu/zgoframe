package initialize

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"zgoframe/core/global"
	"zgoframe/util"
)

func GetNewGorm(printPrefix string) ([]*gorm.DB, error) {
	//global.V.Zap.Info(printPrefix + "GetNewGorm , DBType:" + global.C.System.DbType)
	switch global.C.System.DbType {
	//目前仅支持MYSQL ，后期考虑是否加入其它DB
	case "mysql":
		return GormMysql(printPrefix)
	default:
		return GormMysql(printPrefix)
	}
}

func GormMysql(printPrefix string) ([]*gorm.DB, error) {
	var list []*gorm.DB
	for _,m:=range global.C.Mysql{
		if m.Status != "open"{
			continue
		}
		//m :=
		dns := m.Username + ":" + m.Password + "@tcp(" + m.Ip + ":" + m.Port + ")/" + m.DbName + "?" + m.Config
		global.V.Zap.Info(printPrefix + " gorm mysql ," + m.Username + ":" + "****" + "@tcp(" + m.Ip + ":" + m.Port + ")/" + m.DbName + "?" + m.Config)
		mysqlConfig := mysql.Config{
			DSN:                       dns,   // DSN data source name
			DefaultStringSize:         191,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据版本自动配置
		}
		db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(m.LogMode))
		if err != nil {
			global.V.Zap.Error("MySQL启动异常:" + err.Error())
			return nil, err
		}
		//db = db.Debug()
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)

		list = append(list,db)
	}

	return list, nil
}

func GormShutdown() {
	db, _ := global.V.Gorm.DB()
	db.Close()
}

func gormConfig(mod bool) *gorm.Config {
	//DisableForeignKeyConstraintWhenMigrating:当执行DB迁移时，禁用 外键约束
	//NamingStrategy：表名的一些配置，禁用 表名黑夜为复杂的情况，也就是使用单数表名，这里也可以配置统一表名前缀
	var config = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, NamingStrategy: schema.NamingStrategy{SingularTable: true}}

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
