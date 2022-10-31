package util

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type MyRedisGo struct {
	option   RedisGoOption
	Op       RedisGoOption
	Conn     redis.Conn
	connPool *redis.Pool
}

type RedisGoOption struct {
	Host string
	Port string
	Ps   string
	Log  *zap.Logger
}

func NewRedisConn(redisOption RedisGoOption) (*MyRedisGo, error) {
	myRedis := new(MyRedisGo)
	redisOption.Log.Info("NewRedisConn : " + redisOption.Host + " " + redisOption.Port)
	conn, error := redis.Dial("tcp", redisOption.Host+":"+redisOption.Port)
	if error != nil {
		return nil, error
	}
	myRedis.option = redisOption
	myRedis.Op = redisOption
	myRedis.Conn = conn
	return myRedis, nil
}

func NewRedisConnPool(redisOption RedisGoOption) (*MyRedisGo, error) {
	myRedis := new(MyRedisGo)
	redisOption.Log.Info("NewRedisConn ,host : " + redisOption.Host + " port : " + redisOption.Port + " ps:" + redisOption.Ps)
	myRedisPool := &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     180,
		MaxActive:   200,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			var err error
			var c redis.Conn
			if redisOption.Ps != "" {
				c, err = redis.Dial("tcp", redisOption.Host+":"+redisOption.Port, redis.DialPassword(redisOption.Ps))
			} else {
				c, err = redis.Dial("tcp", redisOption.Host+":"+redisOption.Port)
			}
			if err != nil {
				return nil, err
			}
			// 选择db
			c.Do("SELECT", 1)
			return c, nil
		},
		Wait: true, //如果获取不到，即阻塞
	}
	myRedis.option = redisOption
	myRedis.Op = redisOption
	myRedis.connPool = myRedisPool
	redisOption.Log.Info("test redis conn fd : ping ")

	if redisOption.Ps != "" {
		_, err := myRedis.RedisDo("AUTH", redisOption.Ps)
		if err != nil {
			return nil, errors.New("redis AUTH err:" + redisOption.Ps)
		}
	}
	_, err := myRedis.RedisDo("ping")
	return myRedis, err
}
func (myRedis *MyRedisGo) Shutdown() {
	//myRedis.connPool.Close()
	myRedis.option.Log.Warn("redis shutdown.")
}
func (myRedis *MyRedisGo) GetNewConnFromPool() redis.Conn {
	//myRedis.option.Log.Debug("redis :get new conn FD from pool.")
	conn := myRedis.connPool.Get()
	return conn
}

//指定一个 sock fd
func (myRedis *MyRedisGo) ConnDo(conn redis.Conn, commandName string, args ...interface{}) (reply interface{}, error error) {
	//myRedis.option.Log.Debug("[redis]connDo  :",commandName,args)
	res, error := conn.Do(commandName, args...)
	if error != nil {
		myRedis.option.Log.Error("redis ConnDo err :" + error.Error())
		return nil, error
	}
	return res, error
}
func (myRedis *MyRedisGo) Exec(conn redis.Conn) (reply interface{}, error error) {
	rs, err := myRedis.ConnDo(conn, "exec")
	//myRedis.option.Log.Info("redis : exec , rs : ",rs,"err:",err)
	if err != nil {
		myRedis.option.Log.Error("transaction failed : " + err.Error())
	}
	return rs, err
}
func (myRedis *MyRedisGo) Multi(conn redis.Conn) (reply interface{}, error error) {
	//myRedis.option.Log.Debug("[redis]Multi  ")
	return myRedis.Send(conn, "Multi")
}

func (myRedis *MyRedisGo) Send(conn redis.Conn, commandName string, args ...interface{}) (reply interface{}, error error) {
	err := conn.Send(commandName, args...)
	MyPrint("[redis]Send : ", commandName, args, " err : ", err)
	return reply, err
}

//func  (myRedis *MyRedis)Exec(conn redis.Conn)(reply interface{}, error error){
//	myRedis.option.Log.Debug("[redis]Exec  ")
//	return myRedis.ConnDo(conn,"EXEC")
//}

func (myRedis *MyRedisGo) RedisDo(commandName string, args ...interface{}) (reply interface{}, error error) {
	//myRedis.option.Log.Debug("[redis]redisDo init:",commandName,args)
	conn := myRedis.GetNewConnFromPool()
	defer conn.Close()
	res, error := conn.Do(commandName, args...)
	if error != nil {
		myRedis.option.Log.Warn("redis err :" + error.Error())
		return nil, error
	}
	//MyPrint("RedisDo:", res, error)
	//reflect.ValueOf(res).IsNil(),reflect.ValueOf(res).Kind(),reflect.TypeOf(res)
	//zlib.MyPrint("redisDo exec ,res : ",res," err :",err)
	return res, error
}

func (myRedis *MyRedisGo) RedisDelAllByPrefix(prefix string) {
	myRedis.option.Log.Warn(" action redisDelAllByPrefix : " + prefix)
	res, err := redis.Strings(myRedis.RedisDo("keys", prefix))
	if err != nil {
		ExitPrint("redis keys err :", err.Error())
	}
	myRedis.option.Log.Debug("del element will num :" + strconv.Itoa(len(res)))
	if len(res) <= 0 {
		myRedis.option.Log.Warn(" keys is null,no need del...")
		return
	}
	for _, v := range res {
		res, _ := myRedis.RedisDo("del", v)
		MyPrint("del key ", v, " ,  rs : ", res)
	}
}

func (myRedis *MyRedisGo) redisDelAll(redisPrefix string) {
	myRedis.RedisDelAllByPrefix(redisPrefix)
}
