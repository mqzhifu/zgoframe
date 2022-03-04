package v1

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"strconv"
	"time"
	"zgoframe/core/global"
	httpmiddleware "zgoframe/http/middleware"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/service"
	"zgoframe/util"
)


// 登录-DB比对成功后，签发jwt
func tokenNext(c *gin.Context, user model.User) {
	util.MyPrint("token next user:",user.Id , "sourceType:",request.GetMyHeader(c).SourceType)
	j := &httpmiddleware.JWT{SigningKey: []byte(global.C.Jwt.Key)} // 唯一签名
	claims := request.CustomClaims{
		ProjectId: user.ProjectId,
		//UUID:        user.Uuid,
		//AuthorityId: user.AuthorityId,
		Id:          user.Id,
		NickName:    user.NickName,
		Username:    user.Username,
		SourceType: request.GetMyHeader(c).SourceType,
		BufferTime:  global.C.Jwt.BufferTime, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,                     // 签名生效时间
			ExpiresAt: time.Now().Unix() + global.C.Jwt.ExpiresTime, // 过期时间 7天  配置文件
			Issuer:    "cocos",                                      // 签名的发行者
		},
	}
	//生成token 串
	token, err := j.CreateToken(claims)
	if err != nil {
		global.V.Zap.Error("获取token失败", zap.Any("err", err))
		httpresponse.FailWithMessage("获取token失败", c)
		return
	}
	//从redis里再取一下：可能有，可能没有
	redisElement ,_:= global.V.Redis.GetElementByIndex("jwt",strconv.Itoa(user.ProjectId),strconv.Itoa(request.GetMyHeader(c).SourceType),strconv.Itoa(user.Id))
	//key := service.GetLoginJwtKey(request.GetMyHeader(c).SourceType,user.AppId,user.Id)
	global.V.Zap.Debug("token key:"+redisElement.Key)
	//util.MyPrint(key)
	jwtStr , err   := global.V.Redis.Get(redisElement)
	util.MyPrint("jwtStr:",jwtStr)

	if  err == redis.Nil {//redis里不存在，要么之前没登陆过，要么失效了...
		_ , err := global.V.Redis.SetEX(redisElement,token,int( global.C.Jwt.ExpiresTime))
		if  err != nil {
		//if err := service.SetRedisJWT(token, key); err != nil {
			global.V.Zap.Error("设置登录状态失败", zap.Any("err", err))
			httpresponse.FailWithMessage("设置登录状态失败", c)
			return
		}
		httpresponse.OkWithDetailed(httpresponse.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	} else if err != nil {
		global.V.Zap.Error("设置登录状态失败", zap.Any("err", err))
		httpresponse.FailWithMessage("设置登录状态失败", c)
	} else {//redis 里已经存在

		//var blackJWT model.JwtBlacklist
		//blackJWT.Jwt = jwtStr
		//if err := service.JsonInBlacklist(blackJWT); err != nil {
		//	httpresponse.FailWithMessage("jwt作废失败", c)
		//	return
		//}

		//写入token到redis,覆盖旧的token
		_ , err := global.V.Redis.SetEX(redisElement,token,int( global.C.Jwt.ExpiresTime))
		if  err != nil {
		//if err := service.SetRedisJWT(token, key); err != nil {
			httpresponse.FailWithMessage("设置登录状态失败", c)
			return
		}
		httpresponse.OkWithDetailed(httpresponse.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	}
}

// @Tags User
// @Summary 用户注册账号
// @Produce  application/json
// @Param data body model.User true "用户名, 昵称, 密码, 角色ID ,AppId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"注册成功"}"
// @Router /base/register [post]
func Register(c *gin.Context) {
	var R request.Register
	_ = c.ShouldBind(&R)
	if err := util.Verify(R, util.RegisterVerify); err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	user := &model.User{Username: R.Username, NickName: R.NickName, Password: R.Password, HeaderImg: R.HeaderImg, AuthorityId: R.AuthorityId ,ProjectId: R.AppId}
	err, userReturn := service.Register(*user,request.GetMyHeader(c))
	if err != nil {
		global.V.Zap.Error("注册失败", zap.Any("err", err))
		httpresponse.FailWithDetailed(httpresponse.SysUserResponse{User: userReturn}, "注册失败", c)
	} else {
		httpresponse.OkWithDetailed(httpresponse.SysUserResponse{User: userReturn}, "注册成功", c)
	}
}

// @Summary 用户退出
// @Description 用户退出
// @Tags User
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"退出成功"}"
// @Router /user/logout [post]
func Logout(c *gin.Context) {
	appId ,_ := request.GetAppId(c)
	uid ,_ := request.GetUid(c)

	redisElement ,_:= global.V.Redis.GetElementByIndex("jwt",strconv.Itoa(appId),strconv.Itoa(request.GetMyHeader(c).SourceType),strconv.Itoa(uid))
	global.V.Redis.Del(redisElement)

	//key := service.GetLoginJwtKey(request.GetMyHeader(c).SourceType,appId,uid)
	//service.DelRedisJWT(key)
}


// @Tags User
// @Summary 用户修改密码
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.ChangePasswordStruct true "用户名, 原密码, 新密码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/changePassword [put]
func ChangePassword(c *gin.Context) {
	var user request.ChangePasswordStruct
	_ = c.ShouldBindJSON(&user)
	if err := util.Verify(user, util.ChangePasswordVerify); err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	U := &model.User{Username: user.Username, Password: user.Password}
	if err, _ := service.ChangePassword(U, user.NewPassword); err != nil {
		global.V.Zap.Error("修改失败", zap.Any("err", err))
		httpresponse.FailWithMessage("修改失败，原密码与当前账户不符", c)
	} else {
		httpresponse.OkWithMessage("修改成功", c)
	}
}

// @Tags User
// @Summary 分页获取用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/getUserList [post]
func GetUserList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	if err := util.Verify(pageInfo, util.PageInfoVerify); err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	if err, list, total := service.GetUserInfoList(pageInfo); err != nil {
		global.V.Zap.Error("获取失败", zap.Any("err", err))
		httpresponse.FailWithMessage("获取失败", c)
	} else {
		httpresponse.OkWithDetailed(httpresponse.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags User
// @Summary 设置用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.User true "ID, 用户名, 昵称, 头像链接"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /user/setUserInfo [put]
func SetUserInfo(c *gin.Context) {
	var user model.User
	_ = c.ShouldBindJSON(&user)
	if err := util.Verify(user, util.IdVerify); err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	if err, ReqUser := service.SetUserInfo(user); err != nil {
		global.V.Zap.Error("设置失败", zap.Any("err", err))
		httpresponse.FailWithMessage("设置失败", c)
	} else {
		httpresponse.OkWithDetailed(gin.H{"userInfo": ReqUser}, "设置成功", c)
	}
}
// @Tags User
// @Summary 删除用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "用户ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /user/deleteUser [delete]
func DeleteUser(c *gin.Context) {
	httpresponse.OkWithMessage("用户怎么能随便删除呢？不想要鸡腿了？", c)
	return

	//var reqId request.GetById
	//_ = c.ShouldBindJSON(&reqId)
	//if err := util.Verify(reqId, util.IdVerify); err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//	return
	//}
	//jwtId := getUserID(c)
	//if jwtId == int(reqId.Id) {
	//	httpresponse.FailWithMessage("删除失败, 自杀失败", c)
	//	return
	//}
	//if err := service.DeleteUser(reqId.Id); err != nil {
	//	global.V.Zap.Error("删除失败!", zap.Any("err", err))
	//	httpresponse.FailWithMessage("删除失败", c)
	//} else {
	//	httpresponse.OkWithMessage("删除成功", c)
	//}
}

//// 从Gin的Context中获取从jwt解析出来的用户ID
//func GetUserId(c *gin.Context) int {
//	if claims, exists := c.Get("claims"); !exists {
//		global.V.Zap.Error("从Gin的Context中获取从jwt解析出来的用户ID失败, 请检查路由是否使用jwt中间件")
//		return 0
//	} else {
//		waitUse := claims.(*request.CustomClaims)
//		return waitUse.Id
//	}
//}

// 从Gin的Context中获取从jwt解析出来的用户UUID
//func getUserUuid(c *gin.Context) string {
//	if claims, exists := c.Get("claims"); !exists {
//		global.V.Zap.Error("从Gin的Context中获取从jwt解析出来的用户UUID失败, 请检查路由是否使用jwt中间件")
//		return ""
//	} else {
//		waitUse := claims.(*request.CustomClaims)
//		return waitUse.UUID.String()
//	}
//}

// 从Gin的Context中获取从jwt解析出来的用户角色id
//func getUserAuthorityId(c *gin.Context) string {
//	if claims, exists := c.Get("claims"); !exists {
//		global.V.Zap.Error("从Gin的Context中获取从jwt解析出来的用户UUID失败, 请检查路由是否使用jwt中间件")
//		return ""
//	} else {
//		waitUse := claims.(*request.CustomClaims)
//		return waitUse.AuthorityId
//	}
//}
