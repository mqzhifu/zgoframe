package initialize

import (
	"context"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
	"zgoframe/core/global"
	httpmiddleware "zgoframe/http/middleware"
	httpresponse "zgoframe/http/response"
	"zgoframe/http/router"
	"zgoframe/util"
)

func StartHttpGin() {
	dns := global.C.Http.Ip + ":" + global.C.Http.Port
	global.V.Zap.Debug("http gin dns:" + dns)
	server := &http.Server{
		Addr:           dns,
		Handler:        global.V.Gin,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	global.V.Zap.Warn("StartHttpGin : " + dns)

	global.V.HttpServer = server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if strings.Contains(err.Error(), "bind: address already in use") {
				util.MyPrint("server.ListenAndServe() err: bind port failed , ", err.Error())
				global.MainEnv.RootQuitFunc(-5)
				global.MainEnv.RootCancelFunc()
			}
		}
		util.MyPrint("server.ListenAndServe() err:", err)
	}()
}

func HttpServerShutdown() {
	cancelCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	global.V.HttpServer.Shutdown(cancelCtx)
	cancel()
}

func HandleNotFound(c *gin.Context) {
	handleErr := "404 not found......"
	//handleErr.Request = c.Request.Method + " " + c.Request.URL.String()
	c.JSON(404, handleErr)
	return
}

//GIN: 监听HTTP   中间件  文件上传
func GetNewHttpGIN(zapLog *zap.Logger, prefix string) (*gin.Engine, error) {
	staticFSUriName := "/static"
	swaggerUri := "/swagger/*any"

	staticPath := global.MainEnv.RootDir + "/" + global.C.Http.StaticPath
	//保存一下，给外部使用，主要是给HTTP获取配置信息时，使用
	global.C.Http.DiskStaticPath = staticPath
	zapLog.Info(prefix + "GetNewHttpGIN static config , uri: " + staticFSUriName + " , diskPath: " + staticPath)
	zapLog.Info(prefix + "GetNewHttpGIN swagger uri:" + swaggerUri)
	//这里用到了两个log ，一个是gin 自己的LOG，它不会持久化，只输出到屏幕，另一个是zap自建的LOG，用于持久化，但不输出到屏幕
	HttpZapLog = zapLog
	//设置开发模式，日志输出会变少
	//gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.Default()
	//单独的日志记录，GIN默认的日志不会持久化的
	ginRouter.Use(ZapLog())
	//设置静态目录，等待请求
	ginRouter.StaticFS(staticFSUriName, http.Dir(staticPath))
	//加载swagger api 工具
	ginRouter.GET(swaggerUri, ginSwagger.WrapHandler(swaggerFiles.Handler))
	//设置跨域
	ginRouter.Use(httpmiddleware.Cors())
	//404
	ginRouter.NoMethod(HandleNotFound)
	////8<<20 即 8*2^20=8M
	//ginRouter.MaxMultipartMemory=8<<20
	return ginRouter, nil

}

func RegGinHttpRoute() {
	httpresponse.ErrManager = global.V.Err

	global.V.Gin.Use(httpmiddleware.Limiter()).Use(httpmiddleware.Record()).Use(httpmiddleware.Header())

	//设置非登陆可访问API，但是头里要加基础认证的信息
	PublicGroup := global.V.Gin.Group("")
	//开启跨域，NGINX做了配置暂时可以先不用打开
	//PublicGroup.Use(httpmiddleware.Cors())
	PublicGroup.Use(httpmiddleware.HeaderAuth())
	{
		router.InitBaseRouter(PublicGroup)
		router.InitConfigCenterRouter(PublicGroup)
		router.InitGameMatchRouter(PublicGroup)
		router.InitPersistenceRouter(PublicGroup)
		router.InitFileRouter(PublicGroup)

	}
	//给 管理员/开发/运维 使用，正常需要登陆一次并获取TOKEN，还需要二次验证
	SystemGroup := global.V.Gin.Group("")
	SystemGroup.Use(httpmiddleware.JWTAuth()).Use(httpmiddleware.SecondAuth())
	{
		router.InitToolsRouter(SystemGroup)
		router.InitCicdRouter(SystemGroup)
		router.InitSysRouter(SystemGroup)
	}

	PrivateGroup := global.V.Gin.Group("")
	//设置正常API（需要验证）
	PrivateGroup.Use(httpmiddleware.JWTAuth())
	{
		router.InitTwinAgoraRouter(PrivateGroup)
		router.InitUserRouter(PrivateGroup)
		router.InitMailRouter(PrivateGroup)
	}

	nobodyGroup := global.V.Gin.Group("")
	nobodyGroup.Use()
	{
		router.CallbackRouter(nobodyGroup)
		router.InitGatewayRouter(nobodyGroup)
	}

}

var HttpZapLog *zap.Logger

func ZapLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		//start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		//zap.Int("status", c.Writer.Status()),
		//	zap.String("method", c.Request.Method),
		//	zap.String("path", path),
		//	zap.String("query", query),
		//	zap.String("ip", c.ClientIP()),
		//	zap.String("user-agent", c.Request.UserAgent()),
		//	zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		//	zap.Duration("cost", cost),
		s := " "
		context := strconv.Itoa(c.Writer.Status()) + s + c.Request.Method + s + path + s + query + c.ClientIP()
		// + s + c.Request.UserAgent() + c.Errors.ByType(gin.ErrorTypePrivate).String()

		HttpZapLog.Info(context)
		//global.V.Zap.Info("eeeeee", zap.String("time", `http://foo.com`))

		c.Next()
	}
}
