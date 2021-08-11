package initialize

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"time"
	"zgoframe/core/global"
	httpmiddleware "zgoframe/http/middleware"
	"zgoframe/http/router"
	"zgoframe/util"
)

func StartHttpGin(){
	//go func(){
	//	err = global.V.Gin.Run(global.C.Gin.Ip + ":"+ global.C.Gin.Port)
	//	global.V.Zap.Error("V.Gin.Run err:" + err.Error())
	//}()

	dns := global.C.Http.Ip + ":"+ global.C.Http.Port
	server := &http.Server{
		Addr:           dns,
		Handler:        global.V.Gin,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	global.V.HttpServer = server
	go func (){
		err := server.ListenAndServe()
		util.MyPrint("server.ListenAndServe() err:",err)
	}()
}

func HttpServerShutdown(){
	cancelCtx , _ := context.WithCancel(context.Background())
	global.V.HttpServer.Shutdown(cancelCtx)
}

//GIN: 监听HTTP   中间件  文件上传
func GetNewHttpGIN()(*gin.Engine,error) {
	ginRouter := gin.Default()
	ginRouter.StaticFS("/static",http.Dir(global.C.Http.StaticPath))


	ginRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))
	
	ginRouter.GET("/metrics/count", func(c *gin.Context) {
		global.V.Metric.CounterInc("paySuccess")
	})

	ginRouter.GET("/metrics/gauge", func(c *gin.Context) {
		global.V.Metric.CounterInc("payUser")
	})


	//var AccessCounter = prometheus.NewCounterVec(
	//	prometheus.CounterOpts{
	//		Name: "grpc_request_count",
	//	},
	//	[]string{"service1","rs"},
	//)
	//
	////var AccessCounter = prometheus.NewCounterVec(
	////	prometheus.CounterOpts{
	////		Name: "api_requests_total",
	////	},
	////	[]string{"method", "path"},
	////)
	//
	//prometheus.MustRegister(AccessCounter)
	//
	//
	//ginRouter.GET("/counter", func(c *gin.Context) {
	//	//purl, _ := url.Parse(c.Request.RequestURI)
	//	//AccessCounter.With(prometheus.Labels{
	//	//	"service1": c.Request.Method,
	//	//	"rs":   purl.Path,
	//	//}).Add(1)
	//
	//	AccessCounter.With(prometheus.Labels{
	//		"service1": "s1",
	//		"rs":   "success",
	//	}).Add(1)
	//
	//
	//	AccessCounter.With(prometheus.Labels{
	//		"service1": "s1",
	//		"rs":   "failed",
	//	}).Add(1)
	//
	//
	//
	//})







	ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ginRouter.Use(httpmiddleware.Cors())

	PublicGroup := ginRouter.Group("")
	{
		router.InitBaseRouter(PublicGroup)

	}

	ginRouter.Use(httpmiddleware.RateMiddleware())
	PrivateGroup := ginRouter.Group("")
	PrivateGroup.Use(httpmiddleware.JWTAuth()).Use(httpmiddleware.CasbinHandler())
	{
		router.InitUserRouter(PrivateGroup)

	}

	return ginRouter,nil


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
