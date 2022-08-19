package httpmiddleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
)

type JWT struct {
	SigningKey []byte
}

//创建一个JWT结构，自带密钥
func NewJWT() *JWT {
	return &JWT{
		[]byte(global.C.Jwt.Key),
	}
}

// 根据 JWT，创建一个token ，HS256(SHA-256 + HMAC ,共享一个密钥)
func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	global.V.Zap.Debug("CreateToken")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//快捷函数，方便 回调
func JWTAuth() gin.HandlerFunc {
	global.V.Zap.Debug("im in jwtauth:")
	return RealJWTAuth
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (customClaims request.CustomClaims, err error) {
	global.V.Zap.Debug("ParseToken:" + tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	//util.MyPrint(token.Header, " ", token.Valid, "  ", token.Signature, " ", token.Method.Alg(), " ", err)
	if err != nil { //发生错误
		global.V.Zap.Debug("jwt.ParseWithClaims err:" + err.Error())

		if ve, ok := err.(*jwt.ValidationError); ok { //
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				replaceMap := global.V.Err.MakeOneStringReplace(err.Error())
				err = global.V.Err.NewReplace(5201, replaceMap)
				return customClaims, err
			}
		}

		replaceMap := global.V.Err.MakeOneStringReplace(err.Error())
		err = global.V.Err.NewReplace(5202, replaceMap)
		return customClaims, err

	}
	//if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
	//	return claims, nil
	//}
	claims, ok := token.Claims.(*request.CustomClaims)
	if ok && token.Valid {
		global.V.Zap.Debug("ParseToken success , id: " + strconv.Itoa(claims.Id) + " username:" + claims.Username + " sourceType" + strconv.Itoa(claims.SourceType))
		return *claims, nil
	} else {
		err := global.V.Err.New(5203)
		//global.V.Zap.Debug("ParseToken failed ,err: 断言失败，request.CustomClaims")
		return customClaims, err
	}

}

//给中间件使用
func RealJWTAuth(c *gin.Context) {
	header, _ := request.GetMyHeader(c)
	user, customClaims, err := CheckToken(header)
	if err != nil {
		code, msg, _ := global.V.Err.SplitMsg(err.Error())
		httpresponse.Result(code, nil, msg, c)
		//ErrAbortWithResponse()
		//httpresponse.FailWithAll(gin.H{"reload": true}, err.Error(), c)
		c.Abort()
		return
	}
	//if parserTokenData.NewToken != "" {
	//	c.Header("new-token", parserTokenData.NewToken)
	//	c.Header("new-expires-at", strconv.FormatInt(parserTokenData.Claims.ExpiresAt, 10))
	//}

	c.Set("user", user)
	c.Set("customClaims", customClaims)
	//c.Set("user",ParserTokenData.User)

	c.Next()

}

//检查一个token (解析token)
func CheckToken(myHeader request.HeaderRequest) (u model.User, customClaims request.CustomClaims, err error) {
	//登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
	j := NewJWT()
	// parseToken 解析token包含的信息
	claims, err := j.ParseToken(myHeader.Token)
	if err != nil {
		//if err == TokenExpired {
		//	return u, customClaims, errors.New("授权已过期")
		//}
		//return u, customClaims, errors.New(err.Error())
		return u, customClaims, err
	}

	if claims.ProjectId <= 0 || claims.Id <= 0 || claims.SourceType <= 0 {
		//return u, customClaims, errors.New("ProjectId or claims.Id or claims.SourceType : is null")
		return u, customClaims, global.V.Err.New(5204)
	}
	//请求头里的来源类型要与jwt里的对上
	//if claims.SourceType != parserTokenData.SourceType {
	//	return parserTokenData, errors.New("SourceType错误")
	//}
	//util.MyPrint(claims.AppId,claims.ID,claims.Username)
	//err, user := service.FindUserById(claims.Id)
	//parserTokenData.User = user
	//if err != nil {
	//	//_ = service.JsonInBlacklist(model.JwtBlacklist{Jwt: token})
	//	return parserTokenData, errors.New("id not in db")
	//}
	redisElement, _ := global.V.Redis.GetElementByIndex("jwt", strconv.Itoa(claims.SourceType), strconv.Itoa(claims.Id))
	global.V.Zap.Debug("user token key:" + redisElement.Key)
	jwtStr, err := global.V.Redis.Get(redisElement)
	//if eee == redis.Nil {
	//	util.MyPrint("jwt hit hit okokok")
	//} else {
	//	util.MyPrint("jwt not hit nil no no no no ")
	//}
	if err == redis.Nil {
		//return u, customClaims, errors.New("token 不在redis 中，也可能已失效")
		return u, customClaims, global.V.Err.New(5205)
	}

	if err != nil || jwtStr == "" || err == redis.Nil {
		//return u, customClaims, errors.New("redis 读取token 为空 , 失败:" + err.Error())
		return u, customClaims, global.V.Err.New(5206)
	}

	//if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
	//	claims.ExpiresAt = time.Now().Unix() + global.C.Jwt.ExpiresTime
	//	newToken, _ := j.CreateToken(*claims)
	//	//CustomClaims, _ = j.ParseToken(newToken)
	//	claims, _ = j.ParseToken(newToken)
	//	parserTokenData.NewToken = newToken
	//}
	//parserTokenData.Claims = claims
	customClaims = claims
	var user model.User
	err = global.V.Gorm.Where("id = ? ", claims.Id).First(&user).Error
	if err != nil {
		//return u, customClaims, errors.New("uid not in db :" + strconv.Itoa(claims.Id))
		replaceMap := global.V.Err.MakeOneStringReplace(err.Error() + " " + strconv.Itoa(claims.Id))
		return u, customClaims, global.V.Err.NewReplace(5207, replaceMap)
	}
	//if errors.Is(global.V.Gorm.Where("id = ? ", claims.Id).First(&user).Error, gorm.ErrRecordNotFound) {
	//	return u, customClaims, errors.New("uid not in db :" + strconv.Itoa(claims.Id))
	//}

	if user.Status == model.USER_STATUS_DENY {
		return u, customClaims, errors.New("USER STATUS err")
	}
	//parserTokenData.User = user
	return user, customClaims, nil
}
