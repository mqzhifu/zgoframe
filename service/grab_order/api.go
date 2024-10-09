package grab_order

import "zgoframe/model"

// 给前端API
func (grabOrder *GrabOrder) GetPayCategory() ([]model.PayCategory, error) {
	list := []model.PayCategory{}
	grabOrder.Gorm.Find(&list)
	return list, nil
}

type BaseData struct {
	Range    AmountRange `json:"range"`
	Settings Settings    `json:"settings"`
}

// 给前端API
func (grabOrder *GrabOrder) GetBaseData() (BaseData, error) {

	baseData := BaseData{
		Range:    grabOrder.AmountRange,
		Settings: grabOrder.Settings,
	}

	return baseData, nil
}

// 给前端API
func (grabOrder *GrabOrder) GetBucketList() (map[int]*OrderBucket, error) {
	return grabOrder.OrderBucketList, nil
}

// 给前端API
func (grabOrder *GrabOrder) GetUserTotal() (map[int]UserElement, error) {
	return grabOrder.UserTotal.UserElementList, nil
}

// 给前端API
func (grabOrder *GrabOrder) GetUserBucketAmountList() (map[int]map[string][]SetRs, error) {
	userBucketAmountListRs := make(map[int]map[string][]SetRs)

	for categoryId, userBucketList := range grabOrder.UserBucketAmountRangeList {
		userBucketAmountListRs[categoryId] = make(map[string][]SetRs)
		for _, userBucket := range userBucketList {
			//fmt.Println("=====", userBucket.QueueRedis.Key)
			rr := userBucket.QueueRedis.GetAll()
			//_, ok := userBucketAmountListRs[categoryId][userBucket.QueueRedis.Key]
			userBucketAmountListRs[categoryId][userBucket.QueueRedis.Key] = rr
		}
	}
	return userBucketAmountListRs, nil
}
