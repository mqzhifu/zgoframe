package initialize

import (
	"context"
	"time"
	"zgoframe/core/global"
	"zgoframe/util"
)

func GetNewRedis()(*util.MyRedis,error) {
	redisCfg := global.C.Redis
	global.V.Zap.Info("redis conn:"+redisCfg.Ip + ":" +  redisCfg.Port)


	pool := make(map[string]util.RedisElement)
	pool["userInfo"] 	= util.RedisElement{KeyTemplate: "user_info_{0}",Expire: -1,Index: "userInfo"}
	pool["userToken"] 	= util.RedisElement{KeyTemplate: "user_token_{0}_{1}",Expire: -1,Index: "userToken"}
	pool["lock"] 		= util.RedisElement{KeyTemplate: "lock_{0}",Expire: 10,Index: "lock"}
	pool["inc"] 		= util.RedisElement{KeyTemplate: "inc_{0}",Expire: 10,Index: "inc"}

	myRedisKeyOption := util.MyRedisOption{
		Ip: redisCfg.Ip,
		Port: redisCfg.Port,
		Password: redisCfg.Password,
		DbNumber: redisCfg.DbNumber,
		KeyPrefix: global.V.App.Key,
		KeySeparate: "_",
		ElementPool:pool,
		Log: global.V.Zap,
	}

	myRedis ,err := util.NewMyRedis(myRedisKeyOption)
	return myRedis,err
}

func TestRedis(){
	TestQueue()
	//element := util.RedisElement{Index: "userInfo",Replace: []string{"10000"}}
	//
	//r ,err := global.V.Redis.GetLinkElementByIndex(&element).Set(element,"bbbbb")
	//util.MyPrint(r,err)
	//
	//r ,err = global.V.Redis.Get(element)
	//util.MyPrint(r,err)
	//
	//r ,err = global.V.Redis.SetEX(element,"ccccc",10)
	//util.MyPrint(r,err)
	//
	//r ,err = global.V.Redis.Get(element)
	//util.MyPrint(r,err)
	//
	//delRs ,err := global.V.Redis.Del(element)
	//util.MyPrint(delRs,err)
	//
	//
	//element = util.RedisElement{Index: "inc",Replace: []string{"test"}}
	//incrRs ,err := global.V.Redis.GetLinkElementByIndex(&element).Incr(element)
	//util.MyPrint(incrRs,err)
	//
	//r ,err = global.V.Redis.Get(element)
	//util.MyPrint(r,err)
	//
	//incrRs ,err = global.V.Redis.GetLinkElementByIndex(&element).IncrBy(element,11)
	//util.MyPrint(incrRs,err)
	//
	//
	//r ,err = global.V.Redis.Get(element)
	//util.MyPrint(r,err)
	//
	//
	//tx := global.V.Redis.Multi()
	//tx.Set("a","b",0)
	//tx.Exec()
	//tx.Close()
	//tx.Discard()
	//
	//time.Sleep(time.Second * 1)
	//util.ExitPrint("TestRedis exit")
}

func TestQueue(){
	RedisQueueManagerOption := util.RedisQueueManagerOption{
		DeliveryRetry : []int{1,5,10},
		QueueNameList: []string{"myqueue","myqueue2"},
		Redis: global.V.Redis,
		Log: global.V.Zap,
	}

	queueName := "myqueue2"

	redisQueueManager := util.NewRedisQueueManager(RedisQueueManagerOption)
	redisQueueManager.Init()
	redisQueueManager.MsgAdd(queueName,"go_program")

	queue,_ := redisQueueManager.GetQueueByName(queueName)
	queue.Order = "1631941443923-0"
	queue.BlockTime = 1
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	queue.CancelCtx = cancelCtx
	redisQueueManager.ConsumerByQueue(queue)


	time.Sleep(time.Second * 5)
	util.ExitPrint(22222)

	cancelFunc()
}

func RedisShutdown(){
	global.V.Redis.Redis.Close()
}

