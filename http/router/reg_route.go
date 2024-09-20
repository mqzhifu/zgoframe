package router

import (
	"zgoframe/core/global"
	httpmiddleware "zgoframe/http/middleware"
	httpresponse "zgoframe/http/response"
)

func RegGinHttpRoute() {
	httpresponse.ErrManager = global.V.Util.Err
	//公共 中间件: 限流 日志 头部解析
	global.V.Base.Gin.Use(httpmiddleware.Limiter()).Use(httpmiddleware.Record()).Use(httpmiddleware.Header())
	global.V.Base.Gin.Use(httpmiddleware.RecordTimeoutReq())

	//设置非登陆可访问API，但是头里要加基础认证的信息
	PublicGroup := global.V.Base.Gin.Group("")
	//开启跨域，NGINX做了配置暂时可以先不用打开
	//PublicGroup.Use(httpmiddleware.Cors())
	PublicGroup.Use(httpmiddleware.HeaderAuth())
	{
		Base(PublicGroup)
		Persistence(PublicGroup)
		Goods(PublicGroup)
		Orders(PublicGroup)
		Test(PublicGroup)
		Pic(PublicGroup)
		GrabOrder(PublicGroup)
	}
	//管理员/开发/运维 使用，头部要验证，还需要二次验证，主要有些危险的操作
	SystemGroup := global.V.Base.Gin.Group("")
	SystemGroup.Use(httpmiddleware.HeaderAuth()).Use(httpmiddleware.SecondAuth())
	{
		Cicd(SystemGroup)
		ConfigCenter(SystemGroup)
		System(SystemGroup)
		Tools(SystemGroup)
	}

	PrivateGroup := global.V.Base.Gin.Group("")
	//设置正常API（需要验证）
	PrivateGroup.Use(httpmiddleware.HeaderAuth()).Use(httpmiddleware.JWTAuth())
	{
		File(PrivateGroup)
		Gateway(PrivateGroup)
		GameMatch(PrivateGroup)
		TwinAgora(PrivateGroup)
		User(PrivateGroup)
		Mail(PrivateGroup)
		FrameSync(PrivateGroup)
	}
	//3方回调的请求
	nobodyGroup := global.V.Base.Gin.Group("")
	nobodyGroup.Use()
	{
		Callback(nobodyGroup)
	}

}
