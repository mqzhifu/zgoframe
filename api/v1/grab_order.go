package v1

import (
	"fmt"
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
	data, err := apiServices().GrabOrder.GetPayCategory()
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
// @Router /grab/order/get/data [GET]
func GrabOrderGetData(c *gin.Context) {
	//bodyByts, _ := ioutil.ReadAll(c.Request.Body)
	//form := request.FrameSyncRoomHistory{}
	//json.Unmarshal(bodyByts, &form)
	// var form request.FrameSyncRoomHistory
	// c.ShouldBind(&form)
	data, err := apiServices().GrabOrder.GetData()
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

	err, uid := apiServices().GrabOrder.CreateOrder(form)
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

	//req := []grab_order.UserGrabInfo{}
	//for _, v := range form {
	//	o := grab_order.UserGrabInfo{
	//		PayCategoryId: v.PayCategoryId,
	//		AmountMin:     v.AmountMin,
	//		AmountMax:     v.AmountMax,
	//	}
	//	req = append(req, o)
	//}
	//
	err = apiServices().GrabOrder.UserOpenGrab(4, form)
	fmt.Println("33==========")
	if err != nil {
		httpresponse.FailWithMessage("err:"+err.Error(), c)
	} else {
		httpresponse.OkWithAll("可以的哟~", "ok", c)
	}
}
