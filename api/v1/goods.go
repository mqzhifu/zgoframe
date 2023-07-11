package v1

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
)

// @Tags Goods
// @Summary 创建一个商品
// @Description 创建一个商品
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.Goods true "基础信息"
// @Success 200 {int} int "自增ID"
// @Router /goods/create/one [post]
func GoodsCreateOne(c *gin.Context) {
	var L request.Goods
	c.ShouldBind(&L)

	insertId, err := global.V.MyService.User.CreateGoods(L)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		httpresponse.OkWithAll(insertId, "成功", c)
	}

}

// @Tags Orders
// @Summary 创建一个订单
// @Description 创建一个订单
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.Orders true "基础信息"
// @Success 200 {int} int "自增ID"
// @Router /orders/create/one [post]
func OrdersCreateOne(c *gin.Context) {
	var L request.Orders
	c.ShouldBind(&L)

	insertId, err := global.V.MyService.User.CreateOrders(L)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		httpresponse.OkWithAll(insertId, "成功", c)
	}

}

// @Tags Orders
// @Summary 支付一个订单
// @Description 支付一个订单
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.Payment true "基础信息"
// @Success 200 {int} int "自增ID"
// @Router /orders/payment [post]
func OrdersPayment(c *gin.Context) {
	var L request.Payment
	c.ShouldBind(&L)

	err := global.V.MyService.User.OrdersPayment(L)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		httpresponse.OkWithAll("", "成功", c)
	}

}
