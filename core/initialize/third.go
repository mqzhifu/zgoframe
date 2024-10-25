package initialize

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"zgoframe/core"
	"zgoframe/core/global"
	"zgoframe/util"
)

// 初始化时，一些非必要模块
func (initialize *Initialize) InitEtcd(prefix string) {
	var err error
	//etcd
	if global.C.Etcd.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		configZapReturn := global.Zap{
			LevelInt8: 16,
			FileName:  "etcd.zap",
		}
		global.V.Util.Etcd, err = GetNewEtcd(global.MainCmdParameter.Env, configZapReturn, prefix)
		if err != nil {
			global.V.Base.Zap.Error(prefix + "GetNewEtcd err:" + err.Error())
			panic("InitEtcd")
			//return err
		}
	}
}

func (initialize *Initialize) InitEls(prefix string) {
	if global.C.ElasticSearch.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		// ES 配置
		cfg := elasticsearch.Config{
			Addresses: []string{
				"http://" + global.C.ElasticSearch.Dns,
			},
			Username: global.C.ElasticSearch.Username,
			Password: global.C.ElasticSearch.Password,
		}

		// 创建客户端连接
		typedClient, err := elasticsearch.NewTypedClient(cfg)
		if err != nil {
			global.V.Base.Zap.Error("elasticsearch.NewTypedClient failed, err:" + err.Error())
			panic("InitEls")
		}
		global.V.Base.ES8TypedClient = typedClient

		// 创建客户端连接
		client, err := elasticsearch.NewClient(cfg)
		if err != nil {
			fmt.Printf("elasticsearch.NewTypedClient failed, err:%v\n", err)
			panic("InitEls")
		}
		global.V.Base.ES8Client = client
	}
}

func (initialize *Initialize) InitProtobuf(prefix string) {
	var err error
	//初始化-protobuf 映射文件（原 protobuf 目录 改成 static 下面）
	//将rpc service 中的方法，转化成ID（由PHP生成 的ID map）
	if global.C.Protobuf.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		var fileContentArr []string
		protobufStaticDir := global.C.Http.StaticPath + "/proto/"
		fileContentArr, _ = global.V.Util.StaticFileSystem.GetStaticFileContentLine(protobufStaticDir + global.C.Protobuf.IdMapFileName)
		protobufStaticFullDir := global.MainEnv.RootDir + "/" + protobufStaticDir
		global.V.Util.ProtoMap, err = util.NewProtoMap(global.V.Base.Zap, protobufStaticFullDir, global.C.Protobuf.IdMapFileName, global.V.Util.ProjectMng, fileContentArr)
		if err != nil {
			panic("InitProtobuf err:" + err.Error())
		}
	}
}

func (initialize *Initialize) InitAliOss(prefix string) {
	if global.C.AliOss.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		op := util.AliOssOptions{
			AccessKeyId:     global.C.AliOss.AccessKeyId,
			AccessKeySecret: global.C.AliOss.AccessKeySecret,
			Endpoint:        global.C.AliOss.Endpoint,
			BucketName:      global.C.AliOss.Bucket,
			LocalDomain:     global.C.AliOss.SelfDomain,
		}
		global.V.Util.AliOss = util.NewAliOss(op)
	}
}

// 服务发现
func (initialize *Initialize) InitServiceDiscovery(prefix string) {
	var err error
	//服务管理器，这里跟project manager 有点差不多，不同的只是：project是DB中所有记录,service是type=N的情况
	//ps:之所以单独加一个模块，也是因为service有些特殊的结构变量，与project的结构变量不太一样
	global.V.Util.ServiceManager, _ = util.NewServiceManager(global.V.Base.Gorm)
	//service 服务发现，这里有个顺序，必须先实现化完成:serviceManager
	if global.C.ServiceDiscovery.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		if global.C.Etcd.Status != core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
			panic("ServiceDiscovery need Etcd open!")
		}
		global.V.Util.ServiceDiscovery, err = GetNewServiceDiscovery()
		if err != nil {
			panic(err.Error())
		}
	}
}

func (initialize *Initialize) InitMetric(prefix string) {
	//metrics
	if global.C.Metrics.Status == core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
		myPushGateway := util.PushGateway{
			Status:  global.C.PushGateway.Status,
			Ip:      global.C.PushGateway.Ip,
			Port:    global.C.PushGateway.Port,
			JobName: global.V.Util.Project.Name,
		}
		myMetricsOption := util.MyMetricsOption{
			Log:         global.V.Base.Zap,
			NameSpace:   global.V.Util.Project.Name,
			PushGateway: myPushGateway,
			Env:         global.MainCmdParameter.Env,
		}
		global.V.Util.Metric = util.NewMyMetrics(myMetricsOption)

		if global.C.Http.Status != core.GLOBAL_CONFIG_MODEL_STATUS_OPEN {
			panic("InitMetric metrics need gin open!")
		}
		global.V.Base.Gin.GET("/metrics", gin.WrapH(promhttp.Handler()))
		//测试
		//global.V.Base.Gin.GET("/metrics/count", func(c *gin.Context) {
		//	global.V.Base.Metric.CounterInc("paySuccess")
		//})
		//
		//global.V.Base.Gin.GET("/metrics/gauge", func(c *gin.Context) {
		//	global.V.Base.Metric.CounterInc("payUser")
		//})
		//global.V.Base.Metric.Test()
	}
}
