package initialize

import (
	"context"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
	"zgoframe/core/global"
	httpmiddleware "zgoframe/http/middleware"
	"zgoframe/http/router"
	"zgoframe/util"
)

func StartHttpGin(){
	dns := global.C.Http.Ip + ":"+ global.C.Http.Port
	server := &http.Server{
		Addr:           dns,
		Handler:        global.V.Gin,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	global.V.Zap.Warn("StartHttpGin : "+dns)

	global.V.HttpServer = server
	go func (){
		err := server.ListenAndServe()
		util.MyPrint("server.ListenAndServe() err:",err)
	}()
}

func HttpServerShutdown(){
	cancelCtx , cancel := context.WithTimeout(context.Background(),time.Second * 3)
	global.V.HttpServer.Shutdown(cancelCtx)
	cancel()
}

func HandleNotFound(c *gin.Context){
	handleErr := "404 not found."
	//handleErr.Request = c.Request.Method + " " + c.Request.URL.String()
	c.JSON(404,handleErr)
	return
}
//GIN: 监听HTTP   中间件  文件上传
func GetNewHttpGIN(zapLog *zap.Logger)(*gin.Engine,error) {
	HttpZapLog = zapLog
	ginRouter := gin.Default()
	//单独的日志记录，GIN默认的日志不会持久化的
	ginRouter.Use(ZapLog())
	//加载静态目录
	//	Router.Static("/form-generator", "./resource/page")
	ginRouter.StaticFS("/static",http.Dir(global.C.Http.StaticPath))
	//加载swagger api 工具
	ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//设置跨域
	ginRouter.Use(httpmiddleware.Cors())


	ginRouter.NoMethod(HandleNotFound)

	return ginRouter,nil


}

func RegGinHttpRoute(){
	//设置非登陆可访问API
	PublicGroup := global.V.Gin.Group("")
	PublicGroup.Use(httpmiddleware.OperationRecord()).Use(httpmiddleware.RateMiddleware()).Use(httpmiddleware.ProcessHeader())
	{
		router.InitBaseRouter(PublicGroup)
	}

	//加载限流中间件
	PrivateGroup :=  global.V.Gin.Group("")
	//设置正常API（需要验证）
	//httpmiddleware.CasbinHandler()
	PrivateGroup.Use(httpmiddleware.OperationRecord()).Use(httpmiddleware.RateMiddleware()).Use(httpmiddleware.ProcessHeader(),httpmiddleware.JWTAuth())
	{
		router.InitUserRouter(PrivateGroup)
		router.InitLogslaveRouter(PrivateGroup)
		router.InitSysRouter(PrivateGroup)
	}

}


var HttpZapLog *zap.Logger
func ZapLog()gin.HandlerFunc {
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
		context :=  strconv.Itoa(c.Writer.Status()) + s + c.Request.Method + s + path + s + query + c.ClientIP()
		// + s + c.Request.UserAgent() + c.Errors.ByType(gin.ErrorTypePrivate).String()

		HttpZapLog.Info(context)
		//global.V.Zap.Info("eeeeee", zap.String("time", `http://foo.com`))

		c.Next()
	}
}
