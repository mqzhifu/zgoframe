package v1

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
	"time"
	"zgoframe/core/global"
	httpmiddleware "zgoframe/http/middleware"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/util"
)

// 登录-DB比对成功后，签发jwt
// 这里分成2个部分，1是JWT字符串2redis，加REDIS是防止恶意攻击
func tokenNext(c *gin.Context, user model.User) (loginResponse httpresponse.LoginResponse, err error) {
	global.V.Zap.Debug("tokenNext id:" + strconv.Itoa(user.Id) + " sourceType:" + strconv.Itoa(request.GetMyHeader(c).SourceType))
	j := httpmiddleware.NewJWT()
	claims := request.CustomClaims{
		//UUID:        user.Uuid,
		//AuthorityId: user.AuthorityId,
		Id:         user.Id,
		ProjectId:  user.ProjectId,
		NickName:   user.NickName,
		Username:   user.Username,
		SourceType: request.GetMyHeader(c).SourceType,
		//BufferTime: global.C.Jwt.BufferTime, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 10,                       // 签名生效时间，这里提前10秒，用于容错
			ExpiresAt: time.Now().Unix() + global.C.Jwt.ExpiresTime, // 过期时间 7天  配置文件
			Issuer:    "cocos",                                      // 签名的发行者
		},
	}
	//生成token 字符串
	token, err := j.CreateToken(claims)
	if err != nil {
		return loginResponse, errors.New("创建token失败:" + err.Error())
	}
	//从redis里再取一下：可能有，可能没有(redis key=sourceType+uid ，因为可能多端同时登陆，所以得有 sourceType)
	redisElement, _ := global.V.Redis.GetElementByIndex("jwt", strconv.Itoa(request.GetMyHeader(c).SourceType), strconv.Itoa(user.Id))
	//key := service.GetLoginJwtKey(request.GetMyHeader(c).SourceType,user.AppId,user.Id)
	global.V.Zap.Debug("token key:" + redisElement.Key)
	jwtStr, err := global.V.Redis.Get(redisElement)
	util.MyPrint("jwtStr:", jwtStr, " err:", err, " ")

	if err == redis.Nil { //redis里不存在，要么之前没登陆过，要么失效了...
		//token 写入redis 并设置失效时间
		_, err = global.V.Redis.SetEX(redisElement, token, int(global.C.Jwt.ExpiresTime))
		if err != nil {
			return loginResponse, errors.New("redis 设置登录状态失败 1" + err.Error())
		}
		//httpresponse.OkWithDetailed(httpresponse.LoginResponse{
		//	User:      user,
		//	Token:     token,
		//	ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		//}, "登录成功", c)
		loginResponse.Token = token
		loginResponse.ExpiresAt = claims.ExpiresAt
		loginResponse.IsNewToken = true
		return loginResponse, nil
	} else if err != nil {
		//util.MyPrint("im in 2")
		//httpresponse.FailWithMessage("redis 设置登录状态失败 2"+err.Error(), c)
		return loginResponse, errors.New("redis 设置登录状态失败 2" + err.Error())
	} else { //redis 里已经存在
		//出现这种情况，就是重复登陆，有两种选择
		//1. 允许重复登陆了，为了兼容性，重新再写入一次
		//重新写入token到redis
		//_, err = global.V.Redis.SetEX(redisElement, token, int(global.C.Jwt.ExpiresTime))
		//if err != nil {
		//	return customClaims, errors.New("redis 设置登录状态失败 3" + err.Error())
		//}
		//2. 不允许重复登陆 ，报个错，返回旧的TOKEN
		loginResponse.Token = jwtStr
		loginResponse.IsNewToken = false
		j := httpmiddleware.NewJWT()
		oldClaims, _ := j.ParseToken(jwtStr)
		loginResponse.ExpiresAt = oldClaims.ExpiresAt

		//var blackJWT model.JwtBlacklist
		//blackJWT.Jwt = jwtStr
		//if err := service.JsonInBlacklist(blackJWT); err != nil {
		//	httpresponse.FailWithMessage("jwt作废失败", c)
		//	return
		//}

		return loginResponse, nil
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
	uid, _ := request.GetUid(c)

	redisElement, _ := global.V.Redis.GetElementByIndex("jwt", strconv.Itoa(request.GetMyHeader(c).SourceType), strconv.Itoa(uid))
	global.V.Redis.Del(redisElement)

	httpresponse.OkWithDetailed("ok", "退出成功", c)
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
	var form request.ChangePasswordStruct
	_ = c.ShouldBindJSON(&form)

	if form.NewPassword == "" || form.NewPasswordConfirm == "" {
		httpresponse.OkWithMessage("NewPassword |NewPasswordConfirm empty", c)
		return
	}

	if form.NewPassword != form.NewPasswordConfirm {
		httpresponse.OkWithMessage("密码与确认密码不一致", c)
		return
	}
	uid, _ := request.GetUid(c)
	err := global.V.MyService.User.ChangePassword(uid, form.NewPassword)
	if err != nil {
		httpresponse.FailWithMessage("修改失败:"+err.Error(), c)
	} else {
		httpresponse.OkWithMessage("修改成功", c)
	}
}

// @Tags User
// @Summary 获取当前登陆用户的基础信息
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/getUserInfo [get]
func GetUserInfo(c *gin.Context) {
	user, _ := request.GetUser(c)
	httpresponse.OkWithDetailed(user, "ok", c)
}

// @Tags User
// @Summary 分页获取用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/getUserList [post]
func GetUserInfoList(c *gin.Context) {
	//var pageInfo request.PageInfo
	//_ = c.ShouldBindJSON(&pageInfo)
	//if err := util.Verify(pageInfo, util.PageInfoVerify); err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//	return
	//}
	//if err, list, total := service.GetUserInfoList(pageInfo); err != nil {
	//	global.V.Zap.Error("获取失败", zap.Any("err", err))
	//	httpresponse.FailWithMessage("获取失败", c)
	//} else {
	//	httpresponse.OkWithDetailed(httpresponse.PageResult{
	//		List:     list,
	//		Total:    total,
	//		Page:     pageInfo.Page,
	//		PageSize: pageInfo.PageSize,
	//	}, "获取成功", c)
	//}
}

// @Tags User
// @Summary 绑定手机号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.BindMobile true " "
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /user/bindMobile [put]
func BindMobile(c *gin.Context) {
	var form request.BindMobile
	_ = c.ShouldBind(&form)
	if err := util.Verify(form, util.RegisterVerify); err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	_, exist, err := global.V.MyService.User.FindUserByMobile(form.Mobile)
	if exist {
		httpresponse.FailWithMessage("手机号已绑定过了，请不要重复操作", c)
		return
	}

	if err != nil {
		httpresponse.FailWithMessage("db err:"+err.Error(), c)
		return
	}

	uid, _ := request.GetUid(c)
	err = global.V.MyService.User.BindMobile(uid, form.Mobile)
	if err != nil {
		httpresponse.FailWithMessage("db err:"+err.Error(), c)
		return
	}
	httpresponse.OkWithMessage("绑定成功", c)

}

// @Tags User
// @Summary 设置用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SetUserInfo true "ID, 用户名, 昵称, 头像链接"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /user/setUserInfo [put]
func SetUserInfo(c *gin.Context) {
	var editInfoData request.SetUserInfo
	_ = c.ShouldBindJSON(&editInfoData)
	if err := util.Verify(editInfoData, util.IdVerify); err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	uid, _ := request.GetUid(c)

	user := model.User{
		NickName:  editInfoData.NickName,
		HeaderImg: editInfoData.HeaderImg,
		Sex:       editInfoData.Sex,
		Birthday:  editInfoData.Birthday,
	}
	user.Id = uid

	if err, ReqUser := global.V.MyService.User.SetUserInfo(user); err != nil {
		global.V.Zap.Error("设置失败", zap.Any("err", err))
		httpresponse.FailWithMessage("设置失败", c)
	} else {
		httpresponse.OkWithDetailed(gin.H{"userInfo": ReqUser}, "设置成功", c)
	}
}

// @Tags User
// @Summary 删除用户
// @Description 欧美国家要求比较严，必须得有这功能，国内现在也有但不多，主要是用来删除测试的（危险甚用）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
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
