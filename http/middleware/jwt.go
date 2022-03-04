package httpmiddleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"strconv"
	"time"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/service"
)

//"strconv"
//"time"
//"zgoframe/model"

type JWT struct {
	SigningKey []byte
}


var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.C.Jwt.Key),
	}
}

func JWTAuth() gin.HandlerFunc {
	return RealJWTAuth
}

// 创建一个token
func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}

}

func RealJWTAuth(c *gin.Context) {
	parserTokenData,err := CheckToken(request.GetMyHeader(c))
	if err != nil{
		httpresponse.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
		c.Abort()
		return
	}
	if parserTokenData.NewToken != ""{
		c.Header("new-token", parserTokenData.NewToken)
		c.Header("new-expires-at", strconv.FormatInt(parserTokenData.Claims.ExpiresAt, 10))
	}

	c.Set("parserTokenData", parserTokenData)
	//c.Set("userId", ParserTokenData.Claims.Id)
	//c.Set("user",ParserTokenData.User)

	c.Next()

}

func CheckToken(myHeader request.Header )(parserTokenData request.ParserTokenData,err error){
	parserTokenData.Token = myHeader.Token
	parserTokenData.SourceType = myHeader.SourceType

	if parserTokenData.Token == "" || parserTokenData.SourceType <= 0{
		return parserTokenData,errors.New("SourceType错误，未登录或非法访问")
	}
	//登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
	j := NewJWT()
	// parseToken 解析token包含的信息
	claims, err := j.ParseToken(parserTokenData.Token)
	if err != nil {
		if err == TokenExpired {
			return parserTokenData,errors.New("授权已过期")
		}
		return parserTokenData,errors.New(err.Error())
	}

	if claims.ProjectId <0 || claims.Id < 0 {
		return parserTokenData,errors.New("AppId or Id is null")
	}

	if claims.SourceType != parserTokenData.SourceType{
		return parserTokenData,errors.New("SourceType错误")
	}
	//util.MyPrint(claims.AppId,claims.ID,claims.Username)
	err, user := service.FindUserById(claims.Id)
	parserTokenData.User = user
	if err != nil {
		//_ = service.JsonInBlacklist(model.JwtBlacklist{Jwt: token})
		return parserTokenData,errors.New("id not in db")
	}
	redisElement ,_:= global.V.Redis.GetElementByIndex("jwt",strconv.Itoa(claims.ProjectId),strconv.Itoa(parserTokenData.SourceType),strconv.Itoa(claims.Id))
	//redisLoginJwtKey := service.GetLoginJwtKey(parserTokenData.SourceType,claims.AppId,claims.Id)
	global.V.Zap.Debug("user token key:"+redisElement.Key)
	//err, jwtStr := service.GetRedisJWT(redisLoginJwtKey)
	jwtStr , err  := global.V.Redis.Get(redisElement)
	if err != nil || jwtStr == "" || err == redis.Nil {
		return parserTokenData,errors.New("token 不在redis 中")
	}
	if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
		claims.ExpiresAt = time.Now().Unix() + global.C.Jwt.ExpiresTime
		newToken, _ := j.CreateToken(*claims)
		//CustomClaims, _ = j.ParseToken(newToken)
		claims, _ = j.ParseToken(newToken)
		parserTokenData.NewToken = newToken
	}
	parserTokenData.Claims = claims

	return parserTokenData,nil
}

// 更新token
//func (j *JWT) RefreshToken(tokenString string) (string, error) {
//	jwt.TimeFunc = func() time.Time {
//		return time.Unix(0, 0)
//	}
//	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return j.SigningKey, nil
//	})
//	if err != nil {
//		return "", err
//	}
//	if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
//		jwt.TimeFunc = time.Now
//		claims.StandardClaims.ExpiresAt = time.Now().Unix() + 60*60*24*7
//		return j.CreateToken(*claims)
//	}
//	return "", TokenInvalid
//}

