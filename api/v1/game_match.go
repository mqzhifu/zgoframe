package v1

import (
	"zgoframe/service/gamematch"
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
)

// @Tags GameMatch
// @Summary 玩家报名
// @Description  玩家报名
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Param data body gamematch.HttpReqBusiness true " "
// @Success 200 {object} gamematch.Group
// @Router /game/match/sign [get]
func GameMatchSign(c *gin.Context) {
	var form gamematch.HttpReqBusiness
	c.ShouldBind(&form)

	code,httpReqBusiness  := global.V.MyService.GameMatch.BusinessCheckData(form)
	if code != 0{
		err := global.V.Err.New(code)
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	err := global.V.MyService.GameMatch.CheckHttpSignData(httpReqBusiness)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	signRsData, err := global.V.MyService.GameMatch.Sign(httpReqBusiness)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	httpresponse.OkWithAll(signRsData,"ok",c)
}

// @Tags GameMatch
// @Summary 取消报名 -
// @Description  删除已参与匹配的玩家信息
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Param data body gamematch.HttpReqBusiness true " "
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /game/match/sign/cancel [get]
func GameMatchSignCancel(c *gin.Context) {
	var form gamematch.HttpReqBusiness
	c.ShouldBind(&form)

	code,httpReqBusiness  := global.V.MyService.GameMatch.BusinessCheckData(form)
	if code != 0{
		err := global.V.Err.New(code)
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	err :=  global.V.MyService.GameMatch.CheckHttpSignCancelData(httpReqBusiness)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	signClass :=  global.V.MyService.GameMatch.GetContainerSignByRuleId(httpReqBusiness.RuleId)
	global.V.Zap.Info("del by groupId")
	err = signClass.CancelByGroupId(httpReqBusiness.GroupId)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	httpresponse.OkWithAll("成功","ok",c)
}

//}else if uri == "/success/del"{//匹配成功记录，不想要了，删除一掉
//code,msg = httpd.successDelHandler(postJsonStr)
//}else if uri == "/config"{//
//code,msg = httpd.ConfigHandler(postJsonStr)
//}else if uri == "/rule/add" {//添加一条rule
////code,msg = httpd.ruleAddOne(postDataMap)
//}else if uri == "/tools/getErrorInfo" {//所有错误码列表
//code,msg = httpd.getErrorInfoHandler()
//}else if uri == "/tools/clearRuleByCode"{//清空一条rule的所有数组，用于测试
//code,msg = httpd.clearRuleByCodeHandler(postJsonStr)
//}else if uri == "/tools/getNormalMetrics"{//html api
//code,msg = httpd.normalMetrics()
//}else if uri == "/tools/getRedisMetrics"{//html api
//code,msg = httpd.redisMetrics()
//}else if uri == "/tools/RedisStoreDb"{//html api
//code,msg = httpd.RedisStoreDb()
//}else if uri == "/tools/getHttpReqBusiness"{//html api


