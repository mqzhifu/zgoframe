package grab_order

import (
	"errors"
	"time"
)

type UserBucketElement struct {
	Uid         int   `json:"uid"`
	Weight      int   `json:"weight"`
	SuccessTime int   `json:"success_time"`
	FailedTime  int   `json:"failed_time"`
	CreateTime  int64 `json:"create_time"`
	UpdateTime  int64 `json:"update_time"`
}

// 参与抢单-用户-桶
type UserBucket struct {
	CategoryId    int                       `json:"category_id"`
	MinAmount     int                       `json:"min_amount"`
	MaxAmount     int                       `json:"max_amount"`
	PriorityQueue PriorityQueue             `json:"-"` //有序队列   oid => 金额
	ElementsMap   map[int]UserBucketElement `json:"elements_map"`
}

func NewUserBucket(categoryId int, minAmount int, maxAmount int) *UserBucket {
	userBucket := new(UserBucket)
	userBucket.CategoryId = categoryId
	userBucket.MinAmount = minAmount
	userBucket.MaxAmount = maxAmount
	userBucket.PriorityQueue = PriorityQueue{}
	userBucket.ElementsMap = make(map[int]UserBucketElement)

	return userBucket
}

func (userBucket *UserBucket) PushNewOne(uid int) {
	item := &Item{
		value:    uid,
		priority: DEFAULT_PRIORITY,
	}
	userBucketElement := UserBucketElement{
		Uid:         uid,
		CreateTime:  time.Now().Unix(),
		Weight:      DEFAULT_PRIORITY,
		SuccessTime: 0,
		FailedTime:  0,
	}

	userBucket.ElementsMap[uid] = userBucketElement

	userBucket.PriorityQueue.Push(item)
}

func (userBucket *UserBucket) PopOne() (userBucketElement UserBucketElement, err error) {
	if userBucket.PriorityQueue.Len() <= 0 {
		return userBucketElement, errors.New("len = 0 ")
	}

	priorityQueueElement := userBucket.PriorityQueue.Pop().(Item)
	userBucketElement = userBucket.ElementsMap[priorityQueueElement.value]
	delete(userBucket.ElementsMap, priorityQueueElement.value)

	return userBucketElement, nil
}

func (userBucket *UserBucket) PushOne(userBucketElement UserBucketElement) {
	userBucket.ElementsMap[userBucketElement.Uid] = userBucketElement
	item := &Item{
		value:    userBucketElement.Uid,
		priority: userBucketElement.Weight,
	}
	userBucket.PriorityQueue.Push(item)
}
