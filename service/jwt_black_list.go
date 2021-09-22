package service

//func GetRedisJWT(userName string) (err error, redisJWT string) {
//	//redisJWT, err = global.V.Redis.GetElementByIndex("userInfo")
//	//return err, redisJWT
//	return err,redisJWT
//}

//func DelRedisJWT(userName string) int64 {
//	//IntCmd , _ := global.V.Redis.Del(userName)
//	//return 	IntCmd
//	return 0
//}

//func SetRedisJWT(jwt string, userName string) (err error) {
//	// 此处过期时间等于jwt过期时间
//	//timer := time.Duration(global.C.Jwt.ExpiresTime) * time.Second
//	//_, err = global.V.Redis.SetEX(userName, jwt, 0)
//	return err
//}

//func GetLoginJwtKey(sourceType int ,appId int ,uid int)string{
//	//key := "jwt:login:"+ strconv.Itoa(sourceType) + ":"+ strconv.Itoa(appId) + ":" + strconv.Itoa(int(uid))
//	//return key
//	return ""
//}
