package grab_order

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"zgoframe/util"
)

// 参与抢单-用户-桶
type UserBucket struct {
	CategoryId int           `json:"category_id"`
	MinAmount  int           `json:"min_amount"`
	MaxAmount  int           `json:"max_amount"`
	QueueRedis *QueueRedis   `json:"queue_redis"` //有序队列
	Redis      *util.MyRedis `json:"-"`
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

func NewUserBucket(redis *util.MyRedis, categoryId int, minAmount int, maxAmount int) *UserBucket {
	userBucket := new(UserBucket)
	userBucket.CategoryId = categoryId
	userBucket.MinAmount = minAmount
	userBucket.MaxAmount = maxAmount

	redisKey := "grab_order_queue_" + strconv.Itoa(categoryId) + "_" + GetRangeKey(minAmount, maxAmount)
	userBucket.QueueRedis = NewQueueRedis(redis, redisKey)
	return userBucket
}

func NewQueueRedis(r *util.MyRedis, key string) *QueueRedis {
	q := new(QueueRedis)
	q.Key = key

	//fmt.Println("---=====", r)
	q.Redis = r
	return q
}

func (queue *QueueRedis) Len() int {
	redisRs := queue.Redis.Redis.ZCard(context.Background(), queue.Key)
	return int(redisRs.Val())
}

type SetRs struct {
	Score float64 `json:"score"`
	Uid   int     `json:"uid"`
}

func (queue *QueueRedis) GetAll() []SetRs {

	rs := []SetRs{}

	res := queue.Redis.Redis.ZRangeWithScores(context.Background(), queue.Key, 0, -1)
	for _, v := range res.Val() {
		setRs := SetRs{
			Score: v.Score,
			Uid:   v.Member.(int),
		}
		rs = append(rs, setRs)
	}

	return rs
}

// 有序队列中 添加一个UID，可能用户池的UID会重复，但用序队列会自动覆盖
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

func (queue *QueueRedis) DelOneByUid(uid int) {
	queue.Redis.Redis.ZRem(context.Background(), queue.Key, uid)
}
