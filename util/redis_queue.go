package util

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const (
	READ_QUEUE_ORDER_FONT = "0"
	READ_QUEUE_ORDER_TAIL = "$"
)


//队列结构体
type RedisQueue struct {
	Name 				string
	MaxMsgContentSize 	int	//最大 单条消息体内容
	MaxMsgNum 			int	//最大 消息数
	MaxConsumerNum 		int	//最大 消费者数

	RedisConsumerClientConfig
	//下面是动态从redis读取的
	*redis.XInfoStream
}
//客户端开始消费后的，控制信息配置
type RedisConsumerClientConfig struct {
	BlockTime		int		//直接消费队列时，阻塞时间,空为不阻塞
	MsgCount 		int64		//直接消费队列时，每次读取消息条数，空为一次读取所有消息
	Order			string //直接消费队列时读取方向，队首 队尾
	//ConsumerReceiveMsgCallback func(msgList []RedisQueueMsg,err error)//最好不用这个东西
	MsgListChan		chan RedisQueueMsg
	CancelCtx  		context.Context	//当consumer为阻塞模式时，用于结束阻塞协程
}

//消费者组-结构体
type RedisConsumerGroup struct {
	QueueName string			//队列名
	//下面是动态从redis读取的
	XInfoGroup redis.XInfoGroup
}
//消费者-结构体
type RedisConsumer struct {
	Name		string		//消费者名称
	GroupName	string
	Queue		[]string

	NoAck 		bool		//读取一次消息后，确认机制，1自动 2手动
	RedisConsumerClientConfig
}
//单条消息结构体
type RedisQueueMsg struct {
	Id 			string	`json:"id"`
	RedisId 	string	`json:"redis_id"`
	RetryTimes 	int		`json:"retry_times"`//第几次重试
	CreateTime 	int		`json:"create_time"`
	Content 	string	`json:"content"`
	ContentType string	`json:"content_type"`
	QueueName 	string	`json:"queue_name"`
}

type RedisQueueManagerOption struct {
	DeliveryRetry		[]int	//   重度策略
	QueueNameList 		[]string
	Log 				*zap.Logger
	Redis 				*MyRedis
}


type RedisQueueManager struct {
	ConsumerGroupPool 	map[string]RedisConsumerGroup
	ConsumerPool 		map[string]RedisConsumer
	QueuePool 			map[string]RedisQueue
	Option 				RedisQueueManagerOption
	Log 				*zap.Logger
	Redis 				*MyRedis
}

func NewRedisQueueManager(option RedisQueueManagerOption)*RedisQueueManager{
	redisQueueManager := new(RedisQueueManager)
	redisQueueManager.Option = option
	redisQueueManager.Redis = option.Redis
	redisQueueManager.Log = option.Log

	return redisQueueManager
}

//func (myRedis *MyRedis)XInfoStream(queueName string)(rs map[string]string,err error ){
//	r ,e := myRedis.Redis.XInfoStream(context.Background(),queueName).Result()
//
//	ExitPrint(r,e)
//	return rs,err
//}

func  (redisQueueManager *RedisQueueManager)GetContext()context.Context{
	return context.Background()
}

func (redisQueueManager *RedisQueueManager)Init()error{
	if len(redisQueueManager.Option.QueueNameList) == 0{
		return errors.New("QueueNameList len = 0")
	}

	redisQueueManager.ConsumerGroupPool = make(map[string]RedisConsumerGroup)
	redisQueueManager.QueuePool 		= make(map[string]RedisQueue)
	redisQueueManager.ConsumerPool 		= make(map[string]RedisConsumer)

	for _,queueName := range redisQueueManager.Option.QueueNameList{
		queue := RedisQueue{Name: queueName}
		XInfoStream , err := redisQueueManager.Redis.Redis.XInfoStream(redisQueueManager.GetContext(),queueName).Result()
		if err != nil{
			redisQueueManager.Log.Error(err.Error())
			continue
		}
		queue.XInfoStream = XInfoStream
		queue.MsgListChan = make(chan RedisQueueMsg , 1000 )
		redisQueueManager.QueuePool[queueName] = queue


		xInfoGroupList , err := redisQueueManager.Redis.Redis.XInfoGroups(redisQueueManager.GetContext(),queueName).Result()
		if err != nil{
			redisQueueManager.Log.Error(err.Error())
			continue
		}

		if len(xInfoGroupList) == 0{
			redisQueueManager.Log.Error("len xInfoGroupList == 0")
			continue
		}

		for _,xInfoGroup := range xInfoGroupList{
			consumerGroup := RedisConsumerGroup{
				QueueName: queueName,
				XInfoGroup: xInfoGroup,
			}
			redisQueueManager.ConsumerGroupPool[xInfoGroup.Name] = consumerGroup
		}
	}
	//ExitPrint(234234234)
	return nil
}

//队列
//func (redisQueueManager *RedisQueueManager)CreateQueue(queue RedisQueue){
//	exist ,_ := redisQueueManager.Redis.Exists(queue.Name)
//	if exist > 0 {
//
//	}
//}

func(redisQueueManager *RedisQueueManager) GetQueueList(){

}

func (redisQueueManager *RedisQueueManager)GetQueueByName(queueName string)(RedisQueue,bool){
	rq ,ok :=  redisQueueManager.QueuePool[queueName]
	return rq,ok
}

func (redisQueueManager *RedisQueueManager)InfoQueue(){

}
//队列


//消费者-组
func (redisQueueManager *RedisQueueManager)InfoConsumerGroup(){

}

func (redisQueueManager *RedisQueueManager)GetConsumerGroupList(){

}

func (redisQueueManager *RedisQueueManager)CreateConsumerGroup(){

}

func (redisQueueManager *RedisQueueManager)DelConsumerGroup(){

}
//消费者
func (redisQueueManager *RedisQueueManager)DelConsumer(){

}

func Consumer(redisConsumer RedisConsumer){

}

//消费者



func  (redisQueueManager *RedisQueueManager)ConsumerByQueue(queue RedisQueue)(finalMsgList []RedisQueueMsg,err error){
	blockTimeSecond := time.Second * time.Duration(queue.BlockTime)

	xReadArgs := redis.XReadArgs{
		Streams :[]string{ queue.Name , queue.Order},
		Count   :queue.MsgCount,//这里先写死，每条回传比较简单
		Block   :blockTimeSecond,
	}
	if queue.BlockTime > 0 {
		return redisQueueManager.ConsumerOne(queue,xReadArgs)
	}else{
		//if queue.ConsumerReceiveMsgCallback == nil{
		//
		//}
		for{

			select {
				case <- queue.CancelCtx.Done():
					goto end
				default:
					break
			}

			finalMsgList , err := redisQueueManager.ConsumerOne(queue,xReadArgs)
			if err != nil{
				goto end
			}

			if len(finalMsgList) == 0{
				continue
			}
			for _,v := range finalMsgList{
				queue.MsgListChan <- v
			}
		}
	}

	end:
		return finalMsgList , err
}

func  (redisQueueManager *RedisQueueManager)ConsumerOne(queue RedisQueue , xReadArgs redis.XReadArgs)(finalMsgList []RedisQueueMsg,err error){
	msgList , err := redisQueueManager.Redis.Redis.XRead(redisQueueManager.GetContext(),&xReadArgs).Result()
	if err != nil{
		errMsg := "XRead err:"+err.Error()
		return finalMsgList,errors.New(errMsg)
	}

	if len(msgList) == 0{
		errMsg := "redis XRead queue msg list = 0"
		redisQueueManager.Log.Warn(errMsg)
		return finalMsgList,nil
	}

	for _,XStream := range msgList {
		for _,XMessage := range XStream.Messages{
			createTime ,_:= strconv.Atoi(XMessage.Values["create_time"].(string))
			retryTime ,_:= strconv.Atoi(XMessage.Values["retry_times"].(string))
			msg := RedisQueueMsg{
				CreateTime: createTime,
				RetryTimes:  retryTime,
				ContentType:  XMessage.Values["content_type"].(string),
				Content: XMessage.Values["content"].(string),

				RedisId: XMessage.ID,
				QueueName: XStream.Stream,
			}

			finalMsgList = append(finalMsgList,msg)
		}
	}
	return finalMsgList,nil
}


func (redisQueueManager *RedisQueueManager) ConsumerByGroup(consumer RedisConsumer){
	XReadGroupArgs :=  redis.XReadGroupArgs{
		Group 	:   consumer.GroupName,
		Consumer:consumer.Name,
		Streams :[]string {consumer.Queue[0] , ">" },
		Count   :consumer.MsgCount,
		Block   :time.Second * time.Duration(consumer.BlockTime),
		NoAck   :consumer.NoAck,
	}

	redisQueueManager.Redis.Redis.XReadGroup(redisQueueManager.GetContext(),&XReadGroupArgs)
}

//消息

func (redisQueueManager *RedisQueueManager)MsgAdd(queueName string,content string)(string,error){
	//_ ,ok := redisQueueManager.GetQueueByName(queueName)
	//if !ok {
	//	msg := " queue name not in pool."
	//	redisQueueManager.Log.Error(msg)
	//	return "",errors.New(msg)
	//}

	msg := RedisQueueMsg{
		Id 			:uuid.NewV4().String(),
		RedisId 	:"",
		RetryTimes 	:333,
		CreateTime 	:GetNowTimeSecondToInt(),
		Content 	: content ,
		ContentType :"json",
	}
	//MyPrint(msg.CreateTime)
	//str ,_ := json.Marshal(msg)
	//msgStr := "Id " + msg.Id + " RetryTimes " + strconv.Itoa(msg.RetryTimes) + " CreateTime " + strconv.Itoa(msg.CreateTime) + " ContentType " + msg.ContentType + " Content " + msg.Content
	msgMap := StructCovertMap(msg)
	//ExitPrint(msgMap)

	XAddArgs := redis.XAddArgs{
		Stream : queueName,
		Values: msgMap,
	}

	//NoMkStream bool
	//MaxLen     int64 // MAXLEN N
	//
	//// Deprecated: use MaxLen+Approx, remove in v9.
	//MaxLenApprox int64 // MAXLEN ~ N
	//
	//MinID string
	//// Approx causes MaxLen and MinID to use "~" matcher (instead of "=").
	//Approx bool
	//Limit  int64
	//ID     string
	//Values interface{}

	redisBackMsgId ,e := redisQueueManager.Redis.Redis.XAdd(redisQueueManager.GetContext(),&XAddArgs).Result()
	if e != nil{
		errMsg := " queue XAdd err:" + e.Error()
		redisQueueManager.Log.Error(errMsg)
		return "",errors.New(errMsg)
	}

	redisQueueManager.Log.Info("queue xadd msg success ,id:"+redisBackMsgId)
	return redisBackMsgId,nil


}

func MsgDel(){

}

func MsgLen(){

}

func GetMsgList(){

}

func MsgAck(){

}

//消息