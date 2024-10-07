package initialize

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"zgoframe/core/global"
	"zgoframe/util"
)

func GetNewEtcd(env int, configZapReturn global.Zap, prefix string) (myEtcd *util.MyEtcd, err error) {
	//这个是给3方库：clientv3使用的
	//有点操蛋，我回头想想如何优化掉
	zl := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(configZapReturn.LevelInt8)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:      "json",
		EncoderConfig: zap.NewProductionEncoderConfig(),
		//OutputPaths:      []string{"stderr"},
		OutputPaths:      []string{"stdout", configZapReturn.FileName},
		ErrorOutputPaths: []string{"stderr"},
	}

	option := util.EtcdOption{
		ProjectName: global.V.Util.Project.Name,
		ProjectENV:  env,
		//ProjectKey		: global.V.Base.Project.Key,
		FindEtcdUrl: global.C.Etcd.Url,
		Username:    global.C.Etcd.Username,
		Password:    global.C.Etcd.Password,
		Ip:          global.C.Etcd.Ip,
		Port:        global.C.Etcd.Port,
		Log:         global.V.Base.Zap,
		ZapConfig:   zl,
		PrintPrefix: prefix,
	}
	myEtcd, err = util.NewMyEtcdSdk(option)
	return myEtcd, err
}
