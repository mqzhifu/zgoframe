package grab_order

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"zgoframe/util"
)

//type UserBucketElement struct {
//	Uid         int   `json:"uid"`
//	Weight      int   `json:"weight"`
//	SuccessTime int   `json:"success_time"`
//	FailedTime  int   `json:"failed_time"`
//	CreateTime  int64 `json:"create_time"`
//	UpdateTime  int64 `json:"update_time"`
//}

// 参与抢单-用户-桶
type UserBucket struct {
	CategoryId int           `json:"category_id"`
	MinAmount  int           `json:"min_amount"`
	MaxAmount  int           `json:"max_amount"`
	QueueRedis *QueueRedis   `json:"queue_redis"` //有序队列
	Redis      *util.MyRedis `json:"-"`
	//PriorityQueue PriorityQueue             `json:"-"` //有序队列   oid => 金额
	//ElementsMap map[int]UserBucketElement `json:"elements_map"`
}

func NewUserBucket(r *util.MyRedis, categoryId int, minAmount int, maxAmount int) *UserBucket {
	userBucket := new(UserBucket)
	userBucket.CategoryId = categoryId
	userBucket.MinAmount = minAmount
	userBucket.MaxAmount = maxAmount
	//userBucket.ElementsMap = make(map[int]UserBucketElement)

	redisKey := "grab_order_queue_" + GetRangeKey(minAmount, maxAmount)
	userBucket.QueueRedis = NewQueueRedis(r, redisKey)
	return userBucket
}

func (userBucket *UserBucket) PopOne() (queueItem QueueItem, err error) {
	if userBucket.QueueRedis.Len() <= 0 {
		return queueItem, errors.New("len = 0 ")
	}
	item, _ := userBucket.QueueRedis.Pop()
	//userBucketElement = userBucket.ElementsMap[userBucketElement.Uid]
	//delete(userBucket.ElementsMap, priorityQueueElement.value)
	//
	//return userBucketElement, nil
	return item, err
}

func (userBucket *UserBucket) PushOne(queueItem QueueItem) {
	userBucket.QueueRedis.Push(queueItem)
}

type QueueRedis struct {
	Key   string        `json:"key"`
	Redis *util.MyRedis `json:"-"`
	//Cnt   int           `json:"cnt"`
}

type QueueItem struct {
	Uid   int `json:"uid"`
	Score int `json:"score"`
}

func NewQueueRedis(r *util.MyRedis, key string) *QueueRedis {
	q := new(QueueRedis)
	q.Key = key

	fmt.Println("---=====", r)
	q.Redis = r
	return q
}

func (queue *QueueRedis) Len() int {
	redisRs := queue.Redis.Redis.ZCard(context.Background(), queue.Key)
	return int(redisRs.Val())
}

func (queue *QueueRedis) Push(item QueueItem) {
	res := queue.Redis.Redis.ZAdd(context.Background(), queue.Key, &redis.Z{Score: float64(item.Score), Member: item.Uid})
	fmt.Println("Redis.ZAdd :", res)
}

func (queue *QueueRedis) Pop() (item QueueItem, err error) {
	list := queue.Redis.Redis.ZRevRangeWithScores(context.Background(), queue.Key, 0, 0).Val()
	if len(list) <= 0 {
		fmt.Println("list is empty")
		return item, errors.New("list is empty ")
	}
	item.Score = int(list[0].Score)
	item.Uid, _ = strconv.Atoi(list[0].Member.(string))

	queue.Redis.Redis.ZRem(context.Background(), queue.Key, list[0].Member)
	return item, err
}
