package httpmiddleware

// 后台，对用户做角色权限的时候使用
//func CasbinHandler() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		claims, _ := c.Get("claims")
//		waitUse := claims.(*request.CustomClaims)
//		// 获取请求的URI
//		obj := c.Request.URL.RequestURI()
//		// 获取请求方法
//		act := c.Request.Method
//		// 获取用户的角色
//		sub := waitUse.AuthorityId
//		e := service.Casbin()
//		// 判断策略中是否存在
//		success, _ := e.Enforce(sub, obj, act)
//		if global.C.System.ENV == "develop" || success {
//			c.Next()
//		} else {
//			httpresponse.FailWithDetailed(gin.H{}, "权限不足", c)
//			c.Abort()
//			return
//		}
//	}
//}
