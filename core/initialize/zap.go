package initialize

import (
	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"errors"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"strings"
	"time"
	"zgoframe/core/global"
	"zgoframe/util"
)


func GetNewZapLog(alert *util.AlertPush , configZap global.Zap) (logger *zap.Logger,configZapReturn global.Zap,err error) {
	if configZap.Level == ""{
		return  nil,configZapReturn,errors.New("Level is empty")
	}

	var level zapcore.Level
	switch configZap.Level { // 初始化配置文件的Level
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

	configZap.LevelInt8 = int8(level)
	//每次输出日志后，回调钩子，主要用来报警
	hook := zap.Hooks(func(entry zapcore.Entry) error {
		if !configZap.AutoAlert{//未开始自动报警
			return nil
		}
		//以下级别日志，均要报警
		num := zap.ErrorLevel | zap.PanicLevel |  zap.FatalLevel |  zap.DPanicLevel
		if entry.Level & num == 0{
			alert.Push(int(entry.Level),entry.Message)
		}
		return nil
	})

	writeSyncer ,err := GetWriteSyncer(&configZap)
	if err != nil{
		return nil,configZapReturn,err
	}
	util.MyPrint("zap.dir:"+configZap.BaseDir + " "+configZap.FileName)

	var encoder zapcore.Encoder
	if configZap.Format == "json" {
		encoder =  zapcore.NewJSONEncoder(getEncoderConfig(configZap))
	}else{
		encoder =  zapcore.NewConsoleEncoder(getEncoderConfig(configZap))
	}

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(configZap.LevelInt8))
	//如果是非正常日志，需要加入调用栈的详细信息，方便查错
	//if level == zap.InfoLevel {
		logger = zap.New(core,hook)
	//} else {
	//	logger = zap.New(code, zap.AddStacktrace(level),hook)
	//}
	//每行日志，都添加上：最后调用的文件，方便定位
	if configZap.ShowLine{
		logger = logger.WithOptions(zap.AddCaller())
	}

	configZapReturn = configZap
	return logger,configZapReturn,nil
}
//所有的日志都给加一个公共的项：projectId，方便给日志分类规档
//
func LoggerWithProject(logger *zap.Logger,projectId int )*zap.Logger{
	logger = logger.With(zap.Int("projectId", projectId))
	return logger
}

func getEncoderConfig(configZap global.Zap) (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey		: "message",
		LevelKey		: "level",
		TimeKey			: "time",
		NameKey			: "logger",
		CallerKey		: "caller",
		StacktraceKey	: configZap.StacktraceKey,
		LineEnding		: zapcore.DefaultLineEnding,
		EncodeLevel		: zapcore.LowercaseLevelEncoder,
		EncodeTime		: CustomTimeEncoder,
		EncodeDuration	: zapcore.SecondsDurationEncoder,
		//EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeCaller	: diy,
		ConsoleSeparator: " | ",
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

func diy(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if strings.Contains(caller.String(), "http.go:70")  {
		return
	}
	enc.AppendString(caller.String())
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(global.C.Zap.Prefix + "2006-01-02 - 15:04:05.000"))
}

func GetWriteSyncer(configZap *global.Zap) (zapcore.WriteSyncer, error) {
	if configZap.ModuleName != ""{
		//将ModuleName单独再创建一个子文件夹子
		//configZap.BaseDir += "/" + configZap.ModuleName
	}

	if configZap.SoftLinkFileName == ""{
		return  nil,errors.New("linkName is empty")
	}

	configZap.FileName = configZap.SoftLinkFileName +  "_" + configZap.FileName
	util.MyPrint("GetNewZapLog:",configZap.ModuleName ,configZap.FileName  ,configZap.LogInConsole )

	_,err  := util.PathExists(configZap.BaseDir)
	if err != nil { // 判断是否有Director文件夹
		util.MyPrint("create directory:", configZap.BaseDir)
		err = os.Mkdir(configZap.BaseDir, os.ModePerm)
		if err != nil{
			return nil,err
		}
	}
	//创建一个：输出器(文件|屏幕)，带翻滚功能
	fileWriter, err := zaprotatelogs.New(
		path.Join(configZap.BaseDir, configZap.ModuleName + "_%Y-%m-%d.log"),//日志文件名格式
		zaprotatelogs.WithLinkName(configZap.FileName),//文件软连接名
		zaprotatelogs.WithMaxAge(7*24*time.Hour),//单个文件，失效时间 7天
		zaprotatelogs.WithRotationTime(24*time.Hour),//每1天 做一次日志文件滚动处理
	)
	var zc zapcore.WriteSyncer
	if configZap.LogInConsole {//日志同时输出到屏幕
		zc ,err = zapcore.NewMultiWriteSyncer( zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}else{
		zc, err = zapcore.NewMultiWriteSyncer( zapcore.AddSync(fileWriter)), err
	}

	return zc,err
	//return zapcore.AddSync(fileWriter), err
}

//// getEncoderCore 获取Encoder的zapcore.Core
//func getEncoderCore(configZap global.Zap) (core zapcore.Core) {
//	writer, err := GetWriteSyncer(configZap) // 使用file-rotatelogs进行日志分割
//	if err != nil {
//		util.MyPrint("Get Write Syncer Failed err:", err.Error())
//		return
//	}
//	return zapcore.NewCore(getEncoder(configZap), writer, zapcore.Level(configZap.LevelInt8))
//}

//// getEncoder 获取zapcore.Encoder
//func getEncoder(configZap global.Zap) zapcore.Encoder {
//	if global.C.Zap.Format == "json" {
//		return zapcore.NewJSONEncoder(getEncoderConfig(configZap))
//	}else{
//		return zapcore.NewConsoleEncoder(getEncoderConfig(configZap))
//	}
//}
