package service

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetRedisJWT
//@description: 从redis取jwt
//@param: userName string
//@return: err error, redisJWT string

func GetRedisJWT(userName string) (err error, redisJWT string) {
	//redisJWT, err = global.V.Redis.Get(userName)
	//return err, redisJWT
	return err,redisJWT
}

func DelRedisJWT(userName string) int64 {
	//IntCmd , _ := global.V.Redis.Del(userName)
	//return 	IntCmd
	return 0
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetRedisJWT
//@description: jwt存入redis并设置过期时间
//@param: userName string
//@return: err error, redisJWT string

func SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	//timer := time.Duration(global.C.Jwt.ExpiresTime) * time.Second
	//_, err = global.V.Redis.SetEX(userName, jwt, 0)
	return err
}

func GetLoginJwtKey(sourceType int ,appId int ,uid int)string{
	//key := "jwt:login:"+ strconv.Itoa(sourceType) + ":"+ strconv.Itoa(appId) + ":" + strconv.Itoa(int(uid))
	//return key
	return ""
}
