package util

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"context"
	"strconv"
	"strings"
)

type RedisElement struct {
	KeyTemplate string
	Key 		string
	Expire 		int
	Desc 		string
	Index 		string
	Replace 	[]string
}

type MyRedis struct{
	Redis	*redis.Client
	Option MyRedisOption
}

type MyRedisOption struct{
	KeyPrefix string
	KeySeparate string
	ElementPool map[string]RedisElement

	Ip string
	Port string
	Password string
	DbNumber int
	Log *zap.Logger
}

func NewMyRedis(myRedisOption MyRedisOption)(*MyRedis,error){
	myRedis := new(MyRedis)
	myRedis.Option = myRedisOption

	client := redis.NewClient(&redis.Options{
		Addr:     myRedis.Option.Ip + ":"+ myRedis.Option.Port,
		Password: myRedis.Option.Password, // no password set
		DB:       myRedis.Option.DbNumber,       // use default DB
	})

	pong, err := client.Do(context.Background(),"ping").Result()

	if err != nil {
		myRedis.Option.Log.Error("redis connect ping failed, err:", zap.Any("err", err))
		return nil,err
	}

	myRedis.Option.Log.Info("redis connect ping response:", zap.String("pong",pong.(string)))

	client.AddHook(NewTracingHook())

	myRedis.Redis = client

	return myRedis,nil
}

func (myRedis *MyRedis)Debug(msg string, fields ...zap.Field){
	myRedis.Option.Log.Info(msg , fields ...)
}

//redis key 可能大部分都是动态的
func (myRedis *MyRedis) GetElementByIndex( keyIndex string , values ...string)( redisElement RedisElement, err error){
	one ,ok := myRedis.Option.ElementPool[keyIndex]
	if !ok {
		msg := "GetKey ERR:" + keyIndex + " , not in pool~"
		myRedis.Option.Log.Error(msg)
		return redisElement,errors.New(msg)
	}

	if one.KeyTemplate == ""{
		msg := "GetKey ERR:" + keyIndex + " , KeyTemplate empty~"
		myRedis.Option.Log.Error(msg)
		return redisElement,errors.New(msg)
	}

	one.Key = one.KeyTemplate
	if len(values) > 0 {
		var valuesTmp []string
		for k,v := range values{
			one.Key = strings.Replace(one.Key,"{" + strconv.Itoa(k) + "}",v,-1)
			valuesTmp = append(valuesTmp,v)
		}
		one.Replace = valuesTmp
	}

	one.Key = myRedis.Option.KeyPrefix + one.Key

	myRedis.Debug("GetKey :"+keyIndex  , zap.Any(" element:",one))

	return one,nil
}

//func (myRedis *MyRedis)GetLinkElementByIndex(redisElement *RedisElement )*MyRedis{
//	element ,_ := myRedis.GetElementByIndex(redisElement.Index,redisElement.Replace ... )
//	redisElement = &element
//	return myRedis
//}
//
//
//func (myRedis *MyRedis)Eval(script string,keys []string, args ...interface{}){
//	myRedis.Redis.Eval(script,keys ,args...)
//}
////set 一个永久有效的值
//func (myRedis *MyRedis)Set( element RedisElement , value string )(string,error){
//	//key , _ , err  := myRedis.GetKey(keyIndex,keyReplaceStrArr...)
//	//if err != nil{
//	//	return "",err
//	//}
//	myRedis.Debug(" set "+ element.Key + " val:" + value)
//	cmdVal, cmdErr := myRedis.Redis.Set(element.Key,value,0).Result()
//	myRedis.After(cmdErr)
//	return cmdVal, cmdErr
//}

func (myRedis *MyRedis)After(err error)(){
	if err != nil{
		if err.Error() == "redis: nil"{
			myRedis.Option.Log.Warn("redis result key not exist.")
		}else{
			myRedis.Option.Log.Error("redis result err:"+err.Error())
		}
	}
}

//set 一个会失效的值
//func (myRedis *MyRedis)SetEX( element RedisElement , value string ,expireSecond int )(string,error){
//	if expireSecond < 0 {
//		errMsg := "expireSecond < 0"
//		return "",errors.New(errMsg)
//	}
//
//	if element.Expire < 0 {
//		errMsg := "element.Expire < 0"
//		return "",errors.New(errMsg)
//	}
//
//	if expireSecond == 0 {//参数里的失效时间优先级更高，如果没有，再从配置文件里读失效时间
//		expireSecond = element.Expire
//	}
//
//	cmdVal, cmdErr :=  myRedis.Redis.Set(element.Key,value,time.Second * time.Duration(expireSecond)).Result()
//	myRedis.After(cmdErr)
//	return cmdVal, cmdErr
//}
//
//func (myRedis *MyRedis)Get(element RedisElement)(string,error){
//	//key ,_,err  := myRedis.GetKey(keyIndex,keyReplaceStrArr...)
//	//if err != nil{
//	//	return "",err
//	//}
//	cmdVal, cmdErr :=   myRedis.Redis.Get(element.Key).Result()
//	myRedis.After(cmdErr)
//	return cmdVal, cmdErr
//}
//
//func (myRedis *MyRedis)Flush()(string, error){
//	return myRedis.Redis.FlushAll().Result()
//}
//
//func (myRedis *MyRedis)Del(element RedisElement)(int64, error) {
//	cmdVal, cmdErr := myRedis.Redis.Del(element.Key).Result()
//	myRedis.After(cmdErr)
//	return cmdVal, cmdErr
//}
//
//func (myRedis *MyRedis)Keys(match string) ([]string, error){
//	cmdVal, cmdErr := myRedis.Redis.Keys(match).Result()
//	return cmdVal, cmdErr
//}
////设置KEY过期时间
//func (myRedis *MyRedis)Expire(element RedisElement ,expireSecond int)(bool,error){
//	if expireSecond <=0 {
//		err := errors.New("expireSecond <=0")
//		return false,err
//	}
//
//	cmdVal, cmdErr :=  myRedis.Redis.Expire(element.Key, time.Second * time.Duration(expireSecond)).Result()
//	myRedis.After(cmdErr)
//	return cmdVal, cmdErr
//}
////开启一个事务
//func (myRedis *MyRedis)Multi()redis.Pipeliner{
//	return myRedis.Redis.TxPipeline()
//}
//
////执行一个事务
//func (myRedis *MyRedis)Exec(pip redis.Pipeliner)(Cmder []redis.Cmder, err error){
//	Cmder,err = pip.Exec()
//	myRedis.After(err)
//	return Cmder, err
//}
////取消一个事务
//func (myRedis *MyRedis)Discard(pip redis.Pipeliner)error{
//	err :=  pip.Discard()
//	myRedis.After(err)
//	return err
//}
//
//func (myRedis *MyRedis)Watch(){
//
//}
//
//func (myRedis *MyRedis)Unwatch(){
//
//}
//
//func (myRedis *MyRedis)Incr(element RedisElement)(int64, error){
//	cmdVal, cmdErr := myRedis.Redis.Incr(element.Key).Result()
//	myRedis.After(cmdErr)
//	return cmdVal, cmdErr
//}
//
//func (myRedis *MyRedis)IncrBy(element RedisElement,num int64)(int64, error){
//	cmdVal, cmdErr := myRedis.Redis.IncrBy(element.Key,num).Result()
//	myRedis.After(cmdErr)
//	return cmdVal, cmdErr
//}
//
//
//func (myRedis *MyRedis)Decr(element RedisElement)(int64, error){
//	cmdVal, cmdErr := myRedis.Redis.Decr(element.Key).Result()
//	myRedis.After(cmdErr)
//	return cmdVal, cmdErr
//}
//
//
//func (myRedis *MyRedis)DecrBy(element RedisElement,num int64)(int64, error){
//	cmdVal, cmdErr := myRedis.Redis.DecrBy(element.Key,num).Result()
//	myRedis.After(cmdErr)
//	return cmdVal, cmdErr
//}
////列表相关   start
//func (myRedis *MyRedis)LPush(){
//
//}
//
//func (myRedis *MyRedis)RPop(){
//
//}
//
//func (myRedis *MyRedis)LLen(){
//
//}
//
//func (myRedis *MyRedis)LRange(){
//
//}
//
////列表相关   end
//
////hash相关   start
//func (myRedis *MyRedis)HSet(){
//
//}
//
//func (myRedis *MyRedis)HGet(){
//
//}
//
//func (myRedis *MyRedis)HGetall(){
//
//}


//hash相关   end

//队列

//func (myRedis *MyRedis)XInfoGroups(queueName string)(rs map[string]string,err error ){
//	myRedis.Redis.Do("xinfo groups " + queueName)
//	return rs,err
//}
//

//
//func (myRedis *MyRedis)Exists(key string)(rs int64,err error ){
//	rs,err = myRedis.Redis.Exists(key).Result()
//	return rs,err
//}


//队列

//func (myRedis *MyRedis)XAdd(){
//	myRedis.Redis.SetNX()
//}

//func (myRedis *MyRedis)GetExpireTime(paraTime int ,RedisKeyDescTime int){
//
//}

//func (myRedis *MyRedis)CheckLockExpireTime(expireSecond int) int {
//	if expireSecond < 0{
//		return -1
//	}else if expireSecond ==  0{
//		return -2
//	}else if expireSecond > 20 {
//		return -1
//	}else{
//		return 1
//	}
//}
//
//func (myRedis *MyRedis)GetLock(keyIndex string ,expireSecond int,keyReplaceStrArr ...string)(bool,string,error){
//	val := uuid.NewV4()
//	lockKey ,redisKeyDesc, err := myRedis.GetKey(keyIndex,"addgold")
//	if err != nil{
//		return false,"",err
//	}
//
//	if expireSecond < 0 {
//		errMsg := "expireSecond < 0"
//		return false,"",errors.New(errMsg)
//	}
//
//
//	if redisKeyDesc.Expire < 0 {
//		errMsg := "redisKeyDesc.Expire < 0"
//		return false,"",errors.New(errMsg)
//	}
//
//	if expireSecond == 0 {//参数里的失效时间优先级更高，如果没有，再从配置文件里读失效时间
//		expireSecond = redisKeyDesc.Expire
//	}
//
//	expire := time.Second * time.Duration(expireSecond)
//	rs,err := myRedis.Redis.SetNX(lockKey,val.String(),expire).Result()
//
//	return rs,val.String(),err
//}
//
//func (myRedis *MyRedis)DelLock(keyIndex string,val string,keyReplaceStrArr ...string)(int64,error){
//	getValue,err := myRedis.Get(keyIndex,keyReplaceStrArr...)
//	if err != nil{
//		return 0,err
//	}
//
//	if getValue == ""{
//		errMsg := "getValue == '' "
//		return  0,errors.New(errMsg)
//	}
//
//	if getValue != val{
//		errMsg := "  getValue != val "
//		return  0,errors.New(errMsg)
//	}
//
//	delNum ,err := myRedis.Del(keyIndex,keyReplaceStrArr...)
//	return delNum ,err
//
//}