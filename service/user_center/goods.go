package user_center

import (
	"errors"
	"strconv"
	"strings"
	"zgoframe/http/request"
	"zgoframe/model"
	"zgoframe/util"
)

func (user *User) CreateGoods(goods request.Goods) (insertId int, err error) {
	title := strings.TrimSpace(goods.Title)
	payAllowType := strings.TrimSpace(goods.PayAllowType)
	if goods.Price <= 0 {
		return 0, errors.New("price <= 0 ")
	}

	if title == "" {
		return 0, errors.New("title empty ")
	}

	if goods.Stock <= 0 {
		return 0, errors.New("Stock <= 0 ")
	}

	if payAllowType == "" {
		return 0, errors.New("payAllowType empty ")
	}

	constHandle := util.NewConstHandle()
	enumConstPool := constHandle.EnumConstPool
	payTypeConstList := enumConstPool["PAY_TYPE_"].ConstList
	goodsStatusConstList := enumConstPool["GOODS_STATUS_"].ConstList

	searchStatus := 0
	for _, v := range goodsStatusConstList {
		goodsStatus := v.Value.(int)
		util.MyPrint(goodsStatus)
		if goodsStatus == goods.Status {
			searchStatus = 1
			break
		}
	}

	if searchStatus == 0 {
		return 0, errors.New("status err : " + strconv.Itoa(goods.Status))
	}

	PayAllowTypeArr := strings.Split(payAllowType, ",")
	for _, payTypeStr := range PayAllowTypeArr {
		searchPayType := 0
		payType := 0
		for _, c := range payTypeConstList {
			constValue := c.Value.(int)
			payType, _ = strconv.Atoi(payTypeStr)
			if payType == constValue {
				searchPayType = 1
				break
			}
		}

		if searchPayType == 0 {
			return 0, errors.New("searchPayType err : " + strconv.Itoa(payType))
		}
	}

	g := model.Goods{
		Title:         title,
		Desc:          goods.Desc,
		Status:        goods.Status,
		Type:          goods.Type,
		Price:         goods.Price,
		Stock:         goods.Stock,
		AllowCoupon:   goods.AllowCoupon,
		AllowGoldCoin: goods.AllowGoldCoin,
		PayAllowType:  payAllowType,
		AdminId:       goods.AdminId,
		Memo:          goods.Memo,
	}

	err = user.Gorm.Create(&g).Error
	if err != nil {
		return 0, errors.New("insert db err: " + err.Error())
	}
	return g.Id, nil
}

func (user *User) CreateOrders(orders request.Orders) (insertId int, err error) {

	couponPrice := 0
	goldCoin := 0

	if orders.GoodsId <= 0 {
		return 0, errors.New("goodsId <= 0 ")
	}

	if orders.CouponId > 0 {
		//待补充
	}

	if orders.GoldCoin > 0 {
		goldCoin = orders.GoldCoin //这里假设 1金币=1元
	}

	if orders.Amount <= 0 {
		return 0, errors.New("购买数量 <= 0 ")
	}

	g := model.Goods{}
	err = user.Gorm.Where("id = ?", orders.GoodsId).First(&g).Error
	if err != nil {
		return 0, errors.New("goods not found")
	}

	if g.Stock <= 0 {
		return 0, errors.New("goods stock not enough")
	}

	realPrice := orders.Amount*g.Price - couponPrice - goldCoin

	newOrders := model.Orders{
		Status:    model.ORDERS_STATUS_NORMAL,
		Type:      orders.Type,
		GoodsId:   orders.GoodsId,
		Amount:    orders.Amount,
		Uid:       orders.Uid,
		Source:    orders.Source,
		UserMemo:  orders.UserMemo,
		AdminMemo: orders.AdminMemo,
		Price:     orders.Amount * g.Price,
		RealPrice: realPrice,
	}

	err = user.Gorm.Create(&newOrders).Error
	if err != nil {
		return 0, errors.New("insert db err: " + err.Error())
	}

	return newOrders.Id, nil
}

func (user *User) OrdersPayment(payment request.Payment) (err error) {
	if payment.OrdersId <= 0 {
		return errors.New("OrdersId <= 0")
	}

	o := model.Orders{}
	err = user.Gorm.Where("id = ?", o.Id).First(&o).Error
	if err != nil {
		return errors.New("orders not found")
	}

	if o.Status != model.ORDERS_STATUS_NORMAL {
		return errors.New("orders status err")
	}

	g := model.Goods{}
	err = user.Gorm.Where("id = ?", o.GoodsId).First(&g).Error
	if err != nil {
		return errors.New("goods not found")
	}

	constHandle := util.NewConstHandle()
	enumConstPool := constHandle.EnumConstPool
	payTypeConstList := enumConstPool["PAY_TYPE_"].ConstList
	searchPayType := 0
	for _, v := range payTypeConstList {
		payType := v.Value.(int)
		if payType == payment.Type {
			searchPayType = 1
		}
	}

	if searchPayType == 0 {
		return errors.New("支付类型不被允许")
	}

	return nil
}
