package initialize

import (
	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"strings"
	"time"
	"zgoframe/core/global"
	"zgoframe/util"
)

//以下均是，zap-log 初始化
var level zapcore.Level
var zapDir string
var zapFileName string
var zapInConsole int
func GetNewZapLog(alert *util.AlertPush,moduleName string,FileName string,InConsole int) (logger *zap.Logger,err error) {
	zapDir = global.C.Zap.Dir + "/" + moduleName
	zapFileName = global.C.Zap.LinkName +  "_" + FileName


	util.MyPrint("GetNewZapLog:",moduleName ,FileName ,InConsole )

	zapInConsole = InConsole
	if ok, _ := util.PathExists(zapDir); !ok { // 判断是否有Director文件夹
		util.MyPrint("create directory:", zapDir)
		err = os.Mkdir(zapDir, os.ModePerm)
		if err != nil{
			return nil,err
		}
	}

	util.MyPrint("zap.dir:"+zapDir + " "+zapFileName)

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

	hook := zap.Hooks(func(entry zapcore.Entry) error {
		if !global.C.Zap.AutoAlert{
			//alert.Push()
			return nil
		}
		num := zap.ErrorLevel | zap.PanicLevel |  zap.FatalLevel |  zap.DPanicLevel
		if entry.Level & num == 0{
			global.V.AlertPush.Push(util.AlertMsg{})
		}
		return nil
	})

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(getEncoderCore(), zap.AddStacktrace(level),hook)
	} else {
		logger = zap.New(getEncoderCore(),hook)
	}
	if global.C.Zap.ShowLine{
		logger = logger.WithOptions(zap.AddCaller())
	}
	logger = logger.With(zap.Int("appId", global.V.App.Id))
	//logger.With(zap.String("appId","5"))
	return logger,nil
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore() (core zapcore.Core) {
	writer, err := GetWriteSyncer() // 使用file-rotatelogs进行日志分割
	if err != nil {
		util.MyPrint("Get Write Syncer Failed err:", err.Error())
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
		//EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeCaller: diy,
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

func GetWriteSyncer() (zapcore.WriteSyncer, error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(zapDir, "%Y-%m-%d.log"),
		zaprotatelogs.WithLinkName(zapFileName),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	if zapInConsole == 1 {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}else{
		return zapcore.NewMultiWriteSyncer( zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}
//以上均是，zap-log 初始化
