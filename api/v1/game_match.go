package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

// @Tags GameMatch
// @Summary 玩家加入/报名游戏匹配
// @Description  报名是以（组）为单位的，而校验是以 player 为单位的
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Param data body pb.GameMatchSign true " "
// @Success 200 {object} gamematch.Group
// @Router /game/match/sign [post]
func GameMatchSign(c *gin.Context) {
	var form pb.GameMatchSign
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
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Param data body pb.GameMatchPlayerCancel true " "
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /game/match/sign/cancel [get]
func GameMatchSignCancel(c *gin.Context) {
	var form pb.GameMatchPlayerCancel
	c.ShouldBind(&form)

	err := global.V.MyService.GameMatch.Cancel(form)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	httpresponse.OkWithMessage("ok", c)
}

// @Tags GameMatch
// @Summary 获取一个 RULE 的基础信息
// @Description  RULE是后台录入的，一次匹配的大部分的配置信息
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param id path string true "query rule id"
// @Produce application/json
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /game/match/rule/{id} [get]
func GameMatchGetOneRule(c *gin.Context) {
	ridStr := c.Param("id")
	rid, _ := strconv.Atoi(ridStr)
	// rule, err := global.V.MyService.GameMatch.RuleManager.GetById(rid)
	util.MyPrint("rid:", rid)
	rule := model.GameMatchRule{}
	err := global.V.Gorm.Where("id = ? ", rid).First(&rule).Error
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	httpresponse.OkWithAll(rule, "ok", c)
}

// @Tags GameMatch
// @Summary 获取语言包
// @Description  用于日常调试
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /game/match/lang [get]
func GameMatchGetLang(c *gin.Context) {
	// util.ErrInfo
	lang := global.V.MyService.GameMatch.GetLang()
	httpresponse.OkWithAll(lang, "ok", c)
}

// @Tags GameMatch
// @Summary 配置信息
// @Description  配置信息
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /game/match/config [get]
func GameMatchConfig(c *gin.Context) {
	op := global.V.MyService.GameMatch.GetOption()
	httpresponse.OkWithAll(op, "ok", c)
}

// }else if uri == "/rule/add" {//添加一条rule
// //code,msg = httpd.ruleAddOne(postDataMap)
// }else if uri == "/tools/getErrorInfo" {//所有错误码列表
// code,msg = httpd.getErrorInfoHandler()
// }else if uri == "/tools/clearRuleByCode"{//清空一条rule的所有数组，用于测试
// code,msg = httpd.clearRuleByCodeHandler(postJsonStr)
// }else if uri == "/tools/getNormalMetrics"{//html api
// code,msg = httpd.normalMetrics()
// }else if uri == "/tools/getRedisMetrics"{//html api
// code,msg = httpd.redisMetrics()
// }else if uri == "/tools/RedisStoreDb"{//html api
// code,msg = httpd.RedisStoreDb()
// }else if uri == "/tools/getHttpReqBusiness"{//html api
