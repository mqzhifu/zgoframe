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
)

func StartHttpGin() {
	dns := global.C.Http.Ip + ":" + global.C.Http.Port
	global.V.Base.Zap.Debug("http gin dns:" + dns)
	server := &http.Server{
		Addr:    dns,
		Handler: global.V.Base.Gin,
		//ReadTimeout:    10 * time.Second,//这里先注释掉，上传大文件的时候，这里可能超时造成 NGINX 502
		//WriteTimeout:   10 * time.Second,
		//MaxHeaderBytes: 1 << 20,
	}

	global.V.Base.Zap.Warn("StartHttpGin : " + dns)

	global.V.Base.HttpServer = server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if strings.Contains(err.Error(), "bind: address already in use") {
				global.V.Base.Zap.Error("server.ListenAndServe() err: bind port failed , " + err.Error())
				global.MainEnv.RootQuitFunc(-5)
				global.MainEnv.RootCancelFunc()
			}
		}
		global.V.Base.Zap.Error("server.ListenAndServe() err:" + err.Error())
	}()
}

func HttpServerShutdown() {
	cancelCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	global.V.Base.HttpServer.Shutdown(cancelCtx)
	cancel()
}

func HandleNotFound(c *gin.Context) {
	handleErr := "404 not found......"
	//handleErr.Request = c.Request.Method + " " + c.Request.URL.String()
	c.JSON(404, handleErr)
	return
}

// GIN: 监听HTTP   中间件  文件上传
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
	if global.MainCmdParameter.BuildStatic == "on" {
		ginRouter.StaticFS(staticFSUriName, http.FS(global.V.Base.StaticFileSys))
	} else {
		ginRouter.StaticFS(staticFSUriName, http.Dir(staticPath))
	}

	//favicon.ico
	ginRouter.StaticFile("/favicon.ico", "./static/favicon.ico")
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
		context := strconv.Itoa(c.Writer.Status()) + s + c.Request.Method + s + path + s + query + c.ClientIP() + s + c.Request.Host
		// + s + c.Request.UserAgent() + c.Errors.ByType(gin.ErrorTypePrivate).String()

		HttpZapLog.Info(context)
		//global.V.Base.Zap.Info("eeeeee", zap.String("time", `http://foo.com`))

		c.Next()
	}
}

//func RegGinHttpRoute() {
//	httpresponse.ErrManager = global.V.Util.Err
//	//公共 中间件: 限流 日志 头部解析
//	global.V.Base.Gin.Use(httpmiddleware.Limiter()).Use(httpmiddleware.Record()).Use(httpmiddleware.Header())
//	global.V.Base.Gin.Use(httpmiddleware.RecordTimeoutReq())
//
//	//设置非登陆可访问API，但是头里要加基础认证的信息
//	PublicGroup := global.V.Base.Gin.Group("")
//	//开启跨域，NGINX做了配置暂时可以先不用打开
//	//PublicGroup.Use(httpmiddleware.Cors())
//	PublicGroup.Use(httpmiddleware.HeaderAuth())
//	{
//		router.Base(PublicGroup)
//		router.Persistence(PublicGroup)
//		router.Goods(PublicGroup)
//		router.Orders(PublicGroup)
//		router.Test(PublicGroup)
//		router.Pic(PublicGroup)
//		router.GrabOrder(PublicGroup)
//	}
//	//管理员/开发/运维 使用，头部要验证，还需要二次验证，主要有些危险的操作
//	SystemGroup := global.V.Base.Gin.Group("")
//	SystemGroup.Use(httpmiddleware.HeaderAuth()).Use(httpmiddleware.SecondAuth())
//	{
//		router.Cicd(SystemGroup)
//		router.ConfigCenter(SystemGroup)
//		router.System(SystemGroup)
//		router.Tools(SystemGroup)
//	}
//
//	PrivateGroup := global.V.Base.Gin.Group("")
//	//设置正常API（需要验证）
//	PrivateGroup.Use(httpmiddleware.HeaderAuth()).Use(httpmiddleware.JWTAuth())
//	{
//		router.File(PrivateGroup)
//		router.Gateway(PrivateGroup)
//		router.GameMatch(PrivateGroup)
//		router.TwinAgora(PrivateGroup)
//		router.User(PrivateGroup)
//		router.Mail(PrivateGroup)
//		router.FrameSync(PrivateGroup)
//	}
//	//3方回调的请求
//	nobodyGroup := global.V.Base.Gin.Group("")
//	nobodyGroup.Use()
//	{
//		router.Callback(nobodyGroup)
//	}
//
//}
