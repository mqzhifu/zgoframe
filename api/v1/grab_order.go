package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Tags GrabOrder
// @Summary 获取支付分类的列表
// @Description  获取支付分类的列表
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /grab/order/get/pay/category [GET]
func GetPayCategory(c *gin.Context) {
	data, err := ApiServices().GrabOrder.GetPayCategory()
	if err != nil {
		httpresponse.FailWithMessage("err:"+err.Error(), c)
	} else {
		httpresponse.OkWithAll(data, "ok", c)
	}
}

// @Tags GrabOrder
// @Summary 抢单-获取数据
// @Description 获取所有，汇总数据
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} request.FrameSyncRoomHistory
// @Router /grab/order/get/order/bucket/list [GET]
func GrabOrderBucketList(c *gin.Context) {
	data, err := ApiServices().GrabOrder.GetBucketList()
	if err != nil {
		httpresponse.FailWithMessage("err:"+err.Error(), c)
	} else {
		httpresponse.OkWithAll(data, "ok", c)
	}
}

// @Tags GrabOrder
// @Summary 抢单-订单列表
// @Description 获取所有，订单列表
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} request.FrameSyncRoomHistory
// @Router /grab/order/get/base/data [GET]
func GrabOrderGetBaseData(c *gin.Context) {
	data, err := ApiServices().GrabOrder.GetBaseData()
	if err != nil {
		httpresponse.FailWithMessage("err:"+err.Error(), c)
	} else {
		httpresponse.OkWithAll(data, "ok", c)
	}
}

// @Tags GrabOrder
// @Summary 抢单-用户汇总列表
// @Description 获取所有，用户汇总列表
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} request.FrameSyncRoomHistory
// @Router /grab/order/get/user/total [GET]
func GrabOrderGetUserTotal(c *gin.Context) {
	data, err := ApiServices().GrabOrder.GetUserTotal()
	if err != nil {
		httpresponse.FailWithMessage("err:"+err.Error(), c)
	} else {
		httpresponse.OkWithAll(data, "ok", c)
	}
}

// @Tags GrabOrder
// @Summary 抢单-用户汇总列表
// @Description 获取所有，用户汇总列表
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} request.FrameSyncRoomHistory
// @Router /grab/order/get/user/total [GET]
func GrabOrderGetUserBucketAmountList(c *gin.Context) {
	data, err := ApiServices().GrabOrder.GetUserBucketAmountList()
	if err != nil {
		httpresponse.FailWithMessage("err:"+err.Error(), c)
	} else {
		httpresponse.OkWithAll(data, "ok", c)
	}
}

// @Tags GrabOrder
// @Summary 创建一个订单
// @Description 创建一个订单，匹配一个用户来支付
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.GrabOrder true "用户信息"
// @Produce  application/json
// @Success 200 {object} request.FrameSyncRoomHistory
// @Router /grab/order/create [POST]
func GrabOrderCreate(c *gin.Context) {
	var form request.GrabOrder
	err := c.ShouldBind(&form)
	util.MyPrint("form:", form, " err:", err)

	//o := grab_order.Order{
	//	Id:         form.Id,
	//	Uid:        form.Uid,
	//	CategoryId: form.CategoryId,
	//	Timeout:    form.Timeout,
	//}

	err, uid := ApiServices().GrabOrder.CreateOrder(form)
	if err != nil {
		httpresponse.FailWithMessage("err:"+err.Error(), c)
	} else {
		httpresponse.OkWithAll("可以的哟~uid:"+strconv.Itoa(uid), "ok", c)
	}
}

// @Tags GrabOrder
// @Summary 用户打开抢单功能
// @Description 用户打开抢单功能，开始抢单匹配用户
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body []request.GrabOrderUserOpen true "用户信息"
// @Produce  application/json
// @Success 200 {object} request.FrameSyncRoomHistory
// @Router /grab/order/user/open [POST]
func GrabOrderUserOpen(c *gin.Context) {
	var form []request.GrabOrderUserOpen
	err := c.ShouldBind(&form)
	util.MyPrint("GrabOrderUserOpen form:", form, " err:", err)

	err = ApiServices().GrabOrder.UserOpenGrab(form[0].Uid, form)
	if err != nil {
		httpresponse.FailWithMessage("err:"+err.Error(), c)
	} else {
		httpresponse.OkWithAll("可以的哟~", "ok", c)
	}
}
