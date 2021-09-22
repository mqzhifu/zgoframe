package util

import (
	"errors"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"context"
	"strconv"
	"strings"
	"time"
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

	pong, err := client.Ping(context.Background()).Result()

	if err != nil {
		myRedis.Option.Log.Error("redis connect ping failed, err:", zap.Any("err", err))
		return nil,err
	}

	myRedis.Option.Log.Info("redis connect ping response:", zap.String("pong",pong))

	client.AddHook(NewTracingHook(myRedis.Option.Log))

	myRedis.Redis = client
	return myRedis,nil
}

//
func  (myRedis *MyRedis)GetContext()context.Context{
	return context.Background()
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

	if one.KeyTemplate == "" {
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
//set 一个永久有效的值
func (myRedis *MyRedis)Set( element RedisElement , value string )(string,error){
	myRedis.Debug(" set "+ element.Key + " val:" + value)
	return myRedis.Redis.Set(myRedis.GetContext(),element.Key,value,0).Result()
}
//set 一个会失效的值
func (myRedis *MyRedis)SetEX( element RedisElement , value string ,expireSecond int )(string,error){
	if expireSecond < 0 {
		errMsg := "expireSecond < 0"
		return "",errors.New(errMsg)
	}

	if element.Expire < 0 {
		errMsg := "element.Expire < 0"
		return "",errors.New(errMsg)
	}

	if expireSecond == 0 {//参数里的失效时间优先级更高，如果没有，再从配置文件里读失效时间
		expireSecond = element.Expire
	}

	return  myRedis.Redis.Set(myRedis.GetContext(),element.Key,value,time.Second * time.Duration(expireSecond)).Result()
}

func (myRedis *MyRedis)Get(element RedisElement)(string,error){
	return    myRedis.Redis.Get(myRedis.GetContext(),element.Key).Result()
}
//删除一个key
func (myRedis *MyRedis)Del(element RedisElement)(int64, error) {
	return myRedis.Redis.Del(myRedis.GetContext(),element.Key).Result()
}
//设置KEY过期时间
func (myRedis *MyRedis)Expire(element RedisElement ,expireSecond int)(bool,error){
	if expireSecond <=0 {
		err := errors.New("expireSecond <=0")
		return false,err
	}

	return  myRedis.Redis.Expire(myRedis.GetContext(),element.Key, time.Second * time.Duration(expireSecond)).Result()
}

func (myRedis *MyRedis)Exist(element RedisElement)(int64, error){
	return myRedis.Redis.Exists(myRedis.GetContext(),element.Key).Result()
}

func (myRedis *MyRedis)Incr(element RedisElement)(int64, error){
	return myRedis.Redis.Incr(myRedis.GetContext(),element.Key).Result()
}

func (myRedis *MyRedis)IncrBy(element RedisElement,num int64)(int64, error){
	return myRedis.Redis.IncrBy(myRedis.GetContext(),element.Key,num).Result()
}


func (myRedis *MyRedis)Decr(element RedisElement)(int64, error){
	return myRedis.Redis.Decr(myRedis.GetContext(),element.Key).Result()
}
//
func (myRedis *MyRedis)DecrBy(element RedisElement,num int64)(int64, error){
	return myRedis.Redis.DecrBy(myRedis.GetContext(),element.Key,num).Result()
}

func (myRedis *MyRedis)GetLock(element RedisElement ,expireSecond int)(bool,string,error){
	val := uuid.NewV4()
	if expireSecond < 0 {
		errMsg := "expireSecond < 0"
		return false,"",errors.New(errMsg)
	}

	if element.Expire < 0 {
		errMsg := "element.Expire < 0"
		return false,"",errors.New(errMsg)
	}

	if expireSecond == 0 {//参数里的失效时间优先级更高，如果没有，再从配置文件里读失效时间
		expireSecond = element.Expire
	}

	expire := time.Second * time.Duration(expireSecond)
	rs,err := myRedis.Redis.SetNX(myRedis.GetContext(),  element.Key,val.String(),expire).Result()

	return rs,val.String(),err
}

func (myRedis *MyRedis)DelLock(element RedisElement,val string)(int64,error){
	getValue,err := myRedis.Get(element)
	if err != nil{
		return 0,err
	}

	if getValue == ""{
		errMsg := "getValue == '' "
		return  0,errors.New(errMsg)
	}

	if getValue != val{
		errMsg := "  getValue != val "
		return  0,errors.New(errMsg)
	}

	delNum ,err := myRedis.Del(element)
	return delNum ,err

}
////列表相关   start
//func (myRedis *MyRedis)LPush(){
//
//}
//
//func (myRedis *MyRedis)RPop(){
//
//}
//
//set 一个永久有效的值
func (myRedis *MyRedis)LLen( element RedisElement )(int64,error){
	return myRedis.Redis.LLen(myRedis.GetContext(),element.Key).Result()
}
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