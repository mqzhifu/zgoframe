package v1

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
)

// @Tags GameMatch
// @Summary 玩家加入/报名游戏匹配
// @Description  报名是以（组）为单位的，而校验是以 player 为单位的
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Param data body request.HttpReqGameMatchPlayerSign true " "
// @Success 200 {object} gamematch.Group
// @Router /game/match/sign [post]
func GameMatchSign(c *gin.Context) {
	var form request.HttpReqGameMatchPlayerSign
	c.ShouldBind(&form)

	group, err := global.V.MyService.GameMatch.PlayerJoin(form)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	httpresponse.OkWithAll(group, "ok", c)
}

// @Tags GameMatch
// @Summary 取消报名
// @Description  删除已参与匹配的玩家信息，以组为单位，如果组里是多个人，其中一个人取消，组里其它的玩家一并都得跟着取消
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Param data body gamematch.HttpReqGameMatchPlayerCancel true " "
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /game/match/sign/cancel [get]
func GameMatchSignCancel(c *gin.Context) {
	var form request.HttpReqGameMatchPlayerCancel
	c.ShouldBind(&form)

	err := global.V.MyService.GameMatch.Cancel(form)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	httpresponse.OkWithMessage("ok", c)
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
