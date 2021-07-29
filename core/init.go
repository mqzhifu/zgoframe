package core

import (
	"github.com/gin-gonic/gin"
	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path"
	"time"
	v1 "zgoframe/api/v1"
	"zgoframe/global"
	httpmiddleware "zgoframe/http/middleware"
	"zgoframe/util"
	"zlib"
	"github.com/go-redis/redis"
	"zgoframe/http/router"
)

const(
	DEFAULT_CONFIT_TYPE  = "toml"
	DEFAULT_CONFIG_FILE_NAME = "config"
)

func Init(ENV string ,configType string , configFileName string){
	myViper,err := GetNewViper(configType,configFileName)
	if err != nil{
		zlib.ExitPrint("GetNewViper err:",err)
	}
	global.V.Vip = myViper
	global.V.Zap = GetNewZapLog()
	global.V.Redis = GetRedis()
	global.V.Gin = GetGIN()
	global.V.Gorm = GetGorm()

	global.C.System.ENV = ENV
}

func GetNewViper(ConfigType string,ConfigName string)(*viper.Viper,error){
	zlib.MyPrint("ConfigType:",ConfigType ," , ConfigName:",ConfigName)
	myViper := viper.New()
	myViper.SetConfigType(ConfigType)
	//myViper.SetConfigName(ConfigName + "." + ConfigType)
	myViper.SetConfigFile(ConfigName + "." + ConfigType)
	err := myViper.ReadInConfig()
	if err != nil{
		zlib.MyPrint("myViper.ReadInConfig() err :",err)
		return myViper,err
	}

	config := global.Config{}
	err = myViper.Unmarshal(&config)
	if err != nil{
		zlib.MyPrint(" myViper.Unmarshal err:",err)
		return myViper,err
	}
	global.C = config
	return myViper,nil
}
//GIN: 监听HTTP   中间件  文件上传
func GetGIN()*gin.Engine {
	ginRouter := gin.Default()
	ginRouter.StaticFS("/static",http.Dir(global.C.Gin.StaticPath))

	//ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ginRouter.Use(httpmiddleware.Cors())

	PublicGroup := ginRouter.Group("")
	{
		 // 注册基础功能路由 不做鉴权
		BaseRouter := PublicGroup.Group("base")
		{
			BaseRouter.POST("login", v1.Login)
			BaseRouter.POST("captcha", v1.Captcha)
		}
	}

	ginRouter.Use(httpmiddleware.RateMiddleware())
	PrivateGroup := ginRouter.Group("")
	PrivateGroup.Use(httpmiddleware.JWTAuth()).Use(httpmiddleware.CasbinHandler())
	{
		router.InitUserRouter(PrivateGroup)

	}


	ginRouter.Run(global.C.Gin.Ip + ":"+global.C.Gin.Port)

	return ginRouter


	//	initialize.InitWkMode()
//	Router := initialize.Routers()
//	Router.Static("/form-generator", "./resource/page")
//
//	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
//	s := initServer(address, Router)
//	// 保证文本顺序输出
//	// In order to ensure that the text order output can be deleted
//	time.Sleep(10 * time.Microsecond)
//	global.GVA_LOG.Info("server run success on ", zap.String("address", address))
//
//	fmt.Printf(`
//	欢迎使用 pg-account
//	当前版本:V2.3.8
//	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
//	默认前端文件运行地址:http://127.0.0.1:8080
//`, address)
//	global.GVA_LOG.Error(s.ListenAndServe().Error())
}

func GetGorm() *gorm.DB {
	switch global.C.System.DbType {
	case "mysql":
		return GormMysql()
	default:
		return GormMysql()
	}
}

func GormMysql() *gorm.DB {
	m := global.C.Mysql
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Ip + ":" + m.Port + ")/" + m.DbName + "?" + m.Config
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(m.LogMode)); err != nil {
		global.V.Zap.Error("MySQL启动异常", zap.Any("err", err))
		os.Exit(0)
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}

func gormConfig(mod bool) *gorm.Config {
	var config = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
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

func GetRedis()*redis.Client {
	redisCfg := global.C.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Ip + ":"+ redisCfg.Port,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DbNumber,       // use default DB
	})
	pong, err := client.Ping().Result()

	if err != nil {
		global.V.Zap.Error("redis connect ping failed, err:", zap.Any("err", err))
	} else {
		global.V.Zap.Info("redis connect ping response:", zap.String("pong",pong))
	}
	return client
}



//以下均是，zap-log 初始化
var level zapcore.Level
func GetNewZapLog() (logger *zap.Logger) {
	if ok, _ := util.PathExists(global.C.Zap.Dir); !ok { // 判断是否有Director文件夹
		zlib.MyPrint("create directory:", global.C.Zap.Dir)
		_ = os.Mkdir(global.C.Zap.Dir, os.ModePerm)
	}

	switch global.C.Zap.Level { // 初始化配置文件的Level
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(getEncoderCore(), zap.AddStacktrace(level))
	} else {
		logger = zap.New(getEncoderCore())
	}
	if global.C.Zap.ShowLine{
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore() (core zapcore.Core) {
	writer, err := GetWriteSyncer() // 使用file-rotatelogs进行日志分割
	if err != nil {
		zlib.MyPrint("Get Write Syncer Failed err:", err.Error())
		return
	}
	return zapcore.NewCore(getEncoder(), writer, level)
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if global.C.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  global.C.Zap.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case global.C.Zap.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case global.C.Zap.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case global.C.Zap.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case global.C.Zap.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(global.C.Zap.Prefix + "2006/01/02 - 15:04:05.000"))
}

func GetWriteSyncer() (zapcore.WriteSyncer, error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(global.C.Zap.Dir, "%Y-%m-%d.log"),
		zaprotatelogs.WithLinkName(global.C.Zap.LinkName),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	if global.C.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}
//以上均是，zap-log 初始化