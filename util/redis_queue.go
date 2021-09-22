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

//type XInfoStream struct {
//	Length          int64		//总消息数
//	RadixTreeKeys   int64		//
//	RadixTreeNodes  int64		//
//	Groups          int64		//有几个消费者组：绑定该队列
//	LastGeneratedID string		//最后一条消费的ID值
//	FirstEntry      XMessage	//第一条消息
//	LastEntry       XMessage	//第二条消息
//}

//type XMessage struct {
//	ID     string
//	Values map[string]interface{}
//}
//
//type XInfoGroup struct {
//	Name            string		//消费者组-名称
//	Consumers       int64		//该组包含了几个消费者
//	Pending         int64		//消费者已拿走了消息，但是未ACK，的消息数量,配合XPENDING 获取列表
//	LastDeliveredID string		//游标，本组最后投递的ID值
//}



//队列结构体
type RedisQueue struct {
	Name 				string
	MaxMsgContentSize 	int	//最大 单条消息体内容
	MaxMsgNum 			int64	//最大 消息数
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
	ConsumerReceiveMsgCallback func(msgList []RedisQueueMsg,err error)//最好不用这个东西
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
	Queue		string

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
	ConsumerGroupPool 	map[string]RedisConsumerGroup	//消费者组
	ConsumerPool 		map[string]RedisConsumer
	QueuePool 			map[string]RedisQueue			//队列
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
	//根据调用者初始化的所有 队列名，进行初始化
	for _,queueName := range redisQueueManager.Option.QueueNameList{
		queue := RedisQueue{Name: queueName}
		//从 redis 获取一个队列的当前：实时详细信息
		XInfoStream , err := redisQueueManager.Redis.Redis.XInfoStream(redisQueueManager.GetContext(),queueName).Result()
		if err != nil{
			redisQueueManager.Log.Error("queueName ("+queueName+") XInfoStream  from redis err : "+err.Error())
			continue
		}
		queue.XInfoStream = XInfoStream
		//存储：消费者从redis里每次读取的消费
		queue.MsgListChan = make(chan RedisQueueMsg , 1000 )
		//将该队列加入到池中
		redisQueueManager.QueuePool[queueName] = queue
		//从 redis 读取一下:该队列的消费者组信息 ps:这里可能为空
		xInfoGroupList , err := redisQueueManager.Redis.Redis.XInfoGroups(redisQueueManager.GetContext(),queueName).Result()
		if err != nil{
			redisQueueManager.Log.Error("queueName ("+queueName+") XInfoGroups  from redis err : "+err.Error())
			continue
		}

		if len(xInfoGroupList) == 0{
			redisQueueManager.Log.Error(queueName  + " , len xInfoGroupList == 0")
			continue
		}
		//遍历一下消费者组，一个队列可能有N个消费者组
		for _,xInfoGroup := range xInfoGroupList{
			consumerGroup := RedisConsumerGroup{
				QueueName: queueName,//将该组绑定到一个队列名上
				XInfoGroup: xInfoGroup,
			}
			//加入到池中
			redisQueueManager.ConsumerGroupPool[xInfoGroup.Name] = consumerGroup
		}
	}
	//ExitPrint(234234234)
	return nil
}

//队列相关=================================
//这个方法有点多余 ，因为在REIDS中：操作一个KYE如果不存在，会直接创建...回头我再想想吧
func (redisQueueManager *RedisQueueManager)CreateQueue(queue RedisQueue)error{
	element := RedisElement{Key:queue.Name }
	exist ,_ := redisQueueManager.Redis.Exist(element)
	if exist > 0 {
		return errors.New("redis keys exist = true")
	}

	return nil
}

func(redisQueueManager *RedisQueueManager) GetQueueAllList()map[string]RedisQueue{
	return redisQueueManager.QueuePool
}

func (redisQueueManager *RedisQueueManager)GetQueueByName(queueName string)(RedisQueue,bool){
	rq ,ok :=  redisQueueManager.QueuePool[queueName]
	return rq,ok
}

//队列=================================


//消费者-组=================================

//获取全部消费者组列表
func (redisQueueManager *RedisQueueManager)GetConsumerGroupAllList()map[string]RedisConsumerGroup{
	return redisQueueManager.ConsumerGroupPool
}
//根据组名 获取一个消费者组的信息
func (redisQueueManager *RedisQueueManager)GetConsumerGroupByName(consumerGroupName string)(consumerGroup RedisConsumerGroup,empty bool){
	if len(redisQueueManager.ConsumerGroupPool) == 0{
		return consumerGroup,true
	}
	for _,v:= range redisQueueManager.ConsumerGroupPool{
		if v.XInfoGroup.Name == consumerGroupName{
			return v,false
		}
	}
	return consumerGroup,true
}
//根据队列名，获取该队列下的所有绑定消费者组列表
func (redisQueueManager *RedisQueueManager)GetConsumerGroupListByQueue(queueName string)(list map[string]RedisConsumerGroup){
	if len(redisQueueManager.ConsumerGroupPool) == 0{
		return nil
	}
	for k,v:= range redisQueueManager.ConsumerGroupPool{
		if v.QueueName == queueName{
			list[k] = v
		}
	}
	return list
}

func (redisQueueManager *RedisQueueManager)CreateConsumerGroup(){

}

func (redisQueueManager *RedisQueueManager)DelConsumerGroup(){

}
//========================================

//消费者
func (redisQueueManager *RedisQueueManager)DelConsumer(){

}

//以消费者组进行消费
func (redisQueueManager *RedisQueueManager) ConsumerByGroup(consumer RedisConsumer)(finalMsgList []RedisQueueMsg,err error){
	blockTimeSecond := time.Second * time.Duration(consumer.BlockTime)

	XReadGroupArgs := redis.XReadGroupArgs{
		Group:consumer.GroupName,
		Streams :[]string{ consumer.Queue , ">"},//>表示未被组内消费的起始消息
		Count   :consumer.MsgCount,
		//Count   :1,//这里先写死，每次读取1条，处理也是一条，简单
		Block   :blockTimeSecond,
		NoAck: consumer.NoAck,
	}
	if consumer.BlockTime <= 0 {
		return redisQueueManager.ConsumerOnceByGroup(XReadGroupArgs)
	}else{
		//firstXReadGroupArgs := XReadGroupArgs
		//firstXReadGroupArgs.Streams = []string{ consumer.Queue ,"$"}
		for{
			select {
				case <- consumer.CancelCtx.Done():
					goto end
				default:
					break
			}

			finalMsgList , err := redisQueueManager.ConsumerOnceByGroup(XReadGroupArgs)
			if err != nil{
				MyPrint("ConsumerByGroup ConsumerOnce err:",err.Error())
				goto end
			}

			if len(finalMsgList) == 0{
				MyPrint("ConsumerByGroup ConsumerOnce len = 0")
				continue
			}
			//var onceLastMsg RedisQueueMsg
			for _,v := range finalMsgList{
				consumer.MsgListChan <- v
				//onceLastMsg = v
			}

			//firstXReadArgs.Streams = []string{ queue.Name ,onceLastMsg.RedisId}
		}
	}

	end:
		return finalMsgList , err
}
//直接消费队列-注：这里可能会阻塞...\
//不建议这么干，因为REDIS内部没有游标控制，得单独再保存，不然每次重启都是重复消费
func  (redisQueueManager *RedisQueueManager)ConsumerByQueue(queue RedisQueue)(finalMsgList []RedisQueueMsg,err error){
	blockTimeSecond := time.Second * time.Duration(queue.BlockTime)

	if queue.MsgCount > queue.MaxMsgNum{
		//一次处理过多条数，有风险
	}

	xReadArgs := redis.XReadArgs{
		Streams :[]string{ queue.Name , queue.Order},
		Count   :queue.MsgCount,
		//Count   :1,//测试：每次读取1条，处理也是一条，简单
		Block   :blockTimeSecond,
	}
	redisQueueManager.Log.Info("start a new consumer by queue:" + queue.Name)

	if queue.BlockTime <= 0 {
		return redisQueueManager.ConsumerOnceByQueue( xReadArgs)
	}else{
		onceLastMsg := xReadArgs
		//firstXReadArgs.Streams = []string{ queue.Name ,"$"}
		for{
			select {
			case <- queue.CancelCtx.Done():
				goto end
			default:
				break
			}

			finalMsgList , err := redisQueueManager.ConsumerOnceByQueue(onceLastMsg)
			if err != nil{
				MyPrint("ConsumerByQueue ConsumerOnce err:",err.Error())
				goto end
			}

			if len(finalMsgList) == 0{
				MyPrint("ConsumerByQueue ConsumerOnce len = 0")
				continue
			}
			//var onceLastMsg RedisQueueMsg
			for _,v := range finalMsgList{
				queue.MsgListChan <- v
				onceLastMsg.Streams = []string{ queue.Name ,v.RedisId}
			}

			time.Sleep(blockTimeSecond)

			//firstXReadArgs.Streams = []string{ queue.Name ,onceLastMsg.RedisId}
		}
	}

end:
	return finalMsgList , err
}
func  (redisQueueManager *RedisQueueManager)ConsumerOnceByQueue(xReadArgs redis.XReadArgs)(finalMsgList []RedisQueueMsg,err error){
	msgList , err := redisQueueManager.Redis.Redis.XRead(redisQueueManager.GetContext(),&xReadArgs).Result()
	return redisQueueManager.ConsumerOnce(msgList,err)
}

func  (redisQueueManager *RedisQueueManager)ConsumerOnceByGroup(XReadGroupArgs redis.XReadGroupArgs)(finalMsgList []RedisQueueMsg,err error){
	msgList , err := redisQueueManager.Redis.Redis.XReadGroup(redisQueueManager.GetContext(),&XReadGroupArgs).Result()
	return redisQueueManager.ConsumerOnce(msgList,err)
}

func  (redisQueueManager *RedisQueueManager)ConsumerOnce( msgList  []redis.XStream,err_in  error)(finalMsgList []RedisQueueMsg,err error){
	//msgList , err := redisQueueManager.Redis.Redis.XRead(redisQueueManager.GetContext(),&xReadArgs).Result()
	if err_in != nil{
		errMsg := "XRead err:"+err_in.Error()
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

//消息

func (redisQueueManager *RedisQueueManager)MsgAdd(queueName string,content string)(string,error){
	_ ,ok := redisQueueManager.GetQueueByName(queueName)
	if !ok {
		msg := " queue name not in pool."
		redisQueueManager.Log.Error(msg)
		return "",errors.New(msg)
	}

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

func GetMsgList(){

}

func ConsumerMsgAck(){

}

//消息