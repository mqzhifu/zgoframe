package initialize

import (
	"context"
	"github.com/gin-gonic/gin"
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

func HandleNotFound(c *gin.Context){
	handleErr := "404 not found."
	//handleErr.Request = c.Request.Method + " " + c.Request.URL.String()
	c.JSON(404,handleErr)
	return
}


//GIN: 监听HTTP   中间件  文件上传
func GetNewHttpGIN()(*gin.Engine,error) {
	ginRouter := gin.Default()
	//获取目录加载
	ginRouter.StaticFS("/static",http.Dir(global.C.Http.StaticPath))
	//加载swagger api 工具
	ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//设置跨域
	ginRouter.Use(httpmiddleware.Cors())
	//设置非登陆可访问API
	PublicGroup := ginRouter.Group("")
	{
		router.InitBaseRouter(PublicGroup)
	}
	PublicGroup.Use(httpmiddleware.ProcessHeader())


	//加载限流中间件
	//ginRouter.Use(httpmiddleware.RateMiddleware())
	PrivateGroup := ginRouter.Group("")
	//设置正常API（需要验证）
	PrivateGroup.Use(httpmiddleware.JWTAuth()).Use(httpmiddleware.CasbinHandler())
	{
		router.InitUserRouter(PrivateGroup)
	}

	//global.V.Gin.GET("/sys/quit",  HttpQuit)
	//global.V.Gin.GET("/sys/config", GetConfig)

	//ginRouter.NoMethod(HandleNotFound)

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
