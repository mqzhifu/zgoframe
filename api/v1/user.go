package v1

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
	"zgoframe/core/global"
	httpmiddleware "zgoframe/http/middleware"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/util"
)

// 登录成功(DB比对成功)后，签发jwt
// 这里分成2个部分，1是JWT字符串2redis，加REDIS是防止恶意攻击
func tokenNext(c *gin.Context, user model.User, loginType int) (loginResponse httpresponse.LoginResponse, err error) {
	//ExpiresAt :=  time.Now().Unix() + global.C.Jwt.ExpiresTime // 过期时间 7天  配置文件
	ExpiresAt := time.Now().Unix() + 60							//测试使用
	haeder , _ := request.GetMyHeader(c)

	util.MyPrint("tokenNext uid:" + strconv.Itoa(user.Id) + " sourceType:" + strconv.Itoa(haeder.SourceType)+ " ExpiresAt:",ExpiresAt)
	j := httpmiddleware.NewJWT()

	claims := request.CustomClaims{
		Id:         user.Id,
		ProjectId:  user.ProjectId,
		NickName:   user.NickName,
		Username:   user.Username,
		SourceType: haeder.SourceType,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 10,                       // 签名生效时间，这里提前10秒，用于容错
			ExpiresAt: ExpiresAt,									//失效时间
			Issuer:    "ck-ar",                                      // 签名的发行者
		},
	}
	//生成token 字符串
	token, err := j.CreateToken(claims)
	if err != nil {
		return loginResponse, errors.New("创建token失败:" + err.Error())
	}
	//从redis里再取一下：可能有，可能没有(redis key=sourceType+uid ，因为可能多端同时登陆，所以得有 sourceType)
	redisElement, _ := global.V.Redis.GetElementByIndex("jwt", strconv.Itoa(haeder.SourceType), strconv.Itoa(user.Id))
	jwtStr, err := global.V.Redis.Get(redisElement)
	util.MyPrint("token key:" + redisElement.Key," RedisJwtStr:", jwtStr, " err:", err, " ")

	if err == redis.Nil { //redis里不存在，要么之前没登陆过，要么失效了...
		//token 写入redis 并设置失效时间
		_, err = global.V.Redis.SetEX(redisElement, token, int(global.C.Jwt.ExpiresTime))
		if err != nil {
			return loginResponse, errors.New("redis 设置登录状态失败 1" + err.Error())
		}
		LoginRecord(c, user.Id, loginType)
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

		LoginRecord(c, user.Id, loginType)

		return loginResponse, nil
	}
}

//用户每次登陆后，日志
//1. 日志，供查询分析
//2. 减少请求方每次头里加上一些重复的统计信息，有session的功能，但不推荐这么用
//3. 长连接不可能每次请求都带头信息的
func LoginRecord(c *gin.Context, uid int, loginType int) {
	header,_ := request.GetMyHeader(c)

	userLogin := model.UserLogin{
		ProjectId:  header.ProjectId,
		SourceType: header.SourceType,
		Uid:        uid,
		Type:       loginType,
		AutoIp:     header.AutoIp,
		Ip:         header.BaseInfo.Ip,

		AppVersion:    header.BaseInfo.AppVersion,
		Os:            header.BaseInfo.OS,
		OsVersion:     header.BaseInfo.OSVersion,
		Device:        header.BaseInfo.Device,
		DeviceVersion: header.BaseInfo.DeviceVersion,
		Lat:           header.BaseInfo.Lat,
		Lon:           header.BaseInfo.Lon,
		DeviceId:      header.BaseInfo.DeviceId,
		Dpi:           header.BaseInfo.DPI,
	}

	global.V.Gorm.Create(&userLogin)
	util.MyPrint("model.UserLogin new Id:", userLogin.Id)
}

// @Summary 用户退出
// @Description 删除 jwt，记录日志。不过只是删除一端的JWT，不同端(source_type)登陆都会生成一个jwt
// @Tags User
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "成功"
// @Router /user/logout [post]
func Logout(c *gin.Context) {
	uid, _ := request.GetUid(c)
	header ,_ := request.GetMyHeader(c)
	redisElement, _ := global.V.Redis.GetElementByIndex("jwt", strconv.Itoa(header.SourceType), strconv.Itoa(uid))
	global.V.Redis.Del(redisElement)

	httpresponse.OkWithAll("ok", "退出成功", c)
	//key := service.GetLoginJwtKey(request.GetMyHeader(c).SourceType,appId,uid)
	//service.DelRedisJWT(key)
}

// @Tags User
// @Summary 设置/修改 密码
// @Description 首次设置 与 修改两个动作可以合成一个，因为没有唯一性验证
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.SetPassword true "用户名, 原密码, 新密码"
// @Success 200 {string} string "ok"
// @Router /user/set/password [put]
func SetPassword(c *gin.Context) {
	var form request.SetPassword
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
// @Summary 获取当前登陆用户的基础信息(使用头里的token解析)
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {object} model.User "用户结构体"
// @Router /user/info [get]
func GetUserInfo(c *gin.Context) {
	user, _ := request.GetUser(c)
	httpresponse.OkWithAll(user, "ok", c)
}

// @Tags User
// @Summary 分页获取用户列表,目前并没有加筛选条件
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "基础信息
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/list [post]
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
// @Summary 设定|修改 - 手机号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.BindMobile true "基础信息"
// @Success 200 {string} string "ok"
// @Router /user/set/mobile [put]
func SetMobile(c *gin.Context) {
	var form request.BindMobile
	_ = c.ShouldBind(&form)

	if form.Mobile == "" || form.SmsAuthCode == "" || form.RuleId <= 0 {
		httpresponse.FailWithMessage("Mobile || SmsAuthCode || SmsRuleId is empty", c)
		return
	}
	err := global.V.MyService.Sms.Verify(form.RuleId, form.Mobile, form.SmsAuthCode)
	if err != nil {
		httpresponse.FailWithMessage("SendSms.Verify err:"+err.Error(), c)
		return
	}

	uid, _ := request.GetUid(c)
	err = global.V.MyService.User.BindMobile(uid, form.Mobile)
	if err != nil {
		httpresponse.FailWithMessage("User.BindMobile err:"+err.Error(), c)
		return
	}

	httpresponse.OkWithMessage("绑定成功", c)

}

// @Tags User
// @Summary 设定|修改 - 邮箱
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.BindEmail true "基础信息"
// @Success 200 {string} string "设置成功"
// @Router /user/set/email [put]
func SetEmail(c *gin.Context) {
	var form request.BindEmail
	_ = c.ShouldBind(&form)

	if form.Email == "" || form.SmsAuthCode == "" || form.RuleId <= 0 {
		httpresponse.FailWithMessage("email || SmsAuthCode || SmsRuleId is empty", c)
		return
	}
	err := global.V.MyService.Email.Verify(form.RuleId, form.Email, form.SmsAuthCode)
	if err != nil {
		httpresponse.FailWithMessage("SendSms.Verify err:"+err.Error(), c)
		return
	}

	uid, _ := request.GetUid(c)
	err = global.V.MyService.User.BindMobile(uid, form.Email)
	if err != nil {
		httpresponse.FailWithMessage("User.BindMobile err:"+err.Error(), c)
		return
	}

	httpresponse.OkWithMessage("绑定成功", c)

}

// @Tags User
// @Summary 设置/修改 用户基础信息
// @Description ""
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SetUserInfo true "ID, 用户名, 昵称, 头像链接"
// @Success 200 {string} string "成功"
// @Router /user/set/info [post]
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
		httpresponse.OkWithAll(gin.H{"userInfo": ReqUser}, "设置成功", c)
	}
}

// @Tags User
// @Summary 删除用户
// @Description 欧美国家要求比较严，必须得有这功能，国内现在也有但不多，目前是用来方便开发/测试的，像脚本做自动化测试生成的用户(需要删除)，以及测试员线上测试时产生的用户数据需要删除（危险甚用）
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Param uids formData string true "用户IDs，多用户时用逗号分割"
// @Produce application/json
// @Success 200 {string} string "ok"
// @Router /user/delete [delete]
func DeleteUser(c *gin.Context) {
	uids := c.PostForm("uids")
	if uids == "" {
		httpresponse.FailWithMessage("uids empty", c)
		return
	}

	uidsArr := strings.Split(uids, ",")
	for _, v := range uidsArr {
		if v == "" {
			continue
		}
		uid, _ := strconv.Atoi(v)
		global.V.MyService.User.Delete(uid)
	}
	//httpresponse.OkWithMessage("用户怎么能随便删除呢？不想要鸡腿了？", c)

	httpresponse.OkWithMessage("ok", c)

}
