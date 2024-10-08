package initialize

import (
	"context"
	"time"
	"zgoframe/core/global"
	"zgoframe/util"
)

func GetNewRedis(prefix string) (*util.MyRedis, error) {
	redisCfg := global.C.Redis
	global.V.Base.Zap.Info(prefix + "redis conn:" + redisCfg.Ip + ":" + redisCfg.Port + " ps:****")

	pool := make(map[string]util.RedisElement)
	//pool["userInfo"] 	= util.RedisElement{KeyTemplate: "user_info_{0}",Expire: -1,Index: "userInfo"}
	//pool["userToken"] 	= util.RedisElement{KeyTemplate: "user_token_{0}_{1}",Expire: -1,Index: "userToken"}
	pool["lock"] = util.RedisElement{KeyTemplate: "lock_{0}", Expire: 10, Index: "lock", Desc: "公共锁"}
	pool["inc"] = util.RedisElement{KeyTemplate: "inc_{0}", Expire: 10, Index: "inc", Desc: "公共记数器(自增)"}
	pool["limiter"] = util.RedisElement{KeyTemplate: "limiter_{0}", Expire: 10, Index: "limiter", Desc: "http 每秒限流"}
	pool["jwt"] = util.RedisElement{KeyTemplate: "jwt_{0}_{1}", Expire: 10, Index: "jwt", Desc: "jwt_{sourceType}_{uid}，用户登陆凭证"}
	pool["rtc_token"] = util.RedisElement{KeyTemplate: "rtc_token_{0}_{1}", Expire: 24 * 60 * 60, Index: "rtc_token", Desc: "rtc_token_{username}_{channel}，声网RTC-登陆凭证"}
	pool["rtm_token"] = util.RedisElement{KeyTemplate: "rtm_token_{0}", Expire: 24 * 60 * 60, Index: "rtm_token", Desc: "rtm_token_{username}，声网RTm-登陆凭证"}
	pool["login_failed_limiter_ip"] = util.RedisElement{KeyTemplate: "login_failed_limiter_ip_{0}_{1}", Expire: 24 * 60 * 60, Index: "login_failed_limiter_ip", Desc: "login_failed_limiter_ip_{0}_{1}，登陆失败N次后，要限制登陆IP"}

	myRedisKeyOption := util.MyRedisOption{
		Ip:          redisCfg.Ip,
		Port:        redisCfg.Port,
		Password:    redisCfg.Password,
		DbNumber:    redisCfg.DbNumber,
		KeyPrefix:   global.V.Util.Project.Name,
		KeySeparate: "_",
		ElementPool: pool,
		Log:         global.V.Base.Zap,
	}

	myRedis, err := util.NewMyRedis(myRedisKeyOption)
	return myRedis, err
}

func TestQueue() {
	//先创建一个队列管理器
	RedisQueueManagerOption := util.RedisQueueManagerOption{
		DeliveryRetry: []int{1, 5, 10},
		//测试的所有队列名称
		QueueNameList: []string{"myqueue", "myqueue2"},
		Redis:         global.V.Base.Redis,
		Log:           global.V.Base.Zap,
	}

	redisQueueManager := util.NewRedisQueueManager(RedisQueueManagerOption)
	//初始化
	redisQueueManager.Init()

	//本次操作的队列：名称
	queueName := "myqueue2"
	//发送一条消息
	redisQueueManager.MsgAdd(queueName, "go_program")
	//获取队列的实时信息及结构体
	queue, _ := redisQueueManager.GetQueueByName(queueName)
	//queue.Order = "1631941443923-0"
	queue.Order = "0"   //从第一条消息开始读取
	queue.BlockTime = 2 //每次阻塞时长
	queue.MsgCount = 2  //一次读取条数
	//创建一个上下文，用于结束
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	queue.CancelCtx = cancelCtx
	//开启消费...
	go redisQueueManager.ConsumerByQueue(queue)
	finishForTimes := 1
	for {
		if finishForTimes > 10 {
			cancelFunc()
			break
		}
		select {
		case oneMsg := <-queue.MsgListChan:
			util.MyPrint("client consumer read on msg :", oneMsg)
		default:
			break
		}
		time.Sleep(time.Second * 1)
		util.MyPrint("client consumer sleep 1")
		finishForTimes++
	}
	util.MyPrint("end end end ...")
}

func RedisShutdown() {
	global.V.Base.Redis.Redis.Close()
}
