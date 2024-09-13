package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Tags FrameSync
// @Summary 配置信息
// @Description  配置信息
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /frame/sync/config [get]
func FrameSyncConfig(c *gin.Context) {
	op := global.V.Service.FrameSync.Option
	httpresponse.OkWithAll(op, "ok", c)
}

// @Tags FrameSync
// @Summary 一个房间的玩家历史操作记录
// @Description 用于断线重连，数据量可能较大
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.FrameSyncRoomHistory true "基础信息"
// @Produce  application/json
// @Success 200 {object} request.FrameSyncRoomHistory
// @Router /frame/sync/room/history [post]
func FrameSyncRoomHistory(c *gin.Context) {
	bodyByts, _ := ioutil.ReadAll(c.Request.Body)
	form := request.FrameSyncRoomHistory{}
	json.Unmarshal(bodyByts, &form)
	// var form request.FrameSyncRoomHistory
	// c.ShouldBind(&form)

	util.MyPrint("=======------", form)

	room, empty := global.V.Service.FrameSync.RoomManage.GetById(form.RoomId)
	if empty {
		httpresponse.FailWithMessage("roomId Empty id:"+form.RoomId, c)
	} else {
		httpresponse.OkWithAll(room.LogicFrameHistory, "ok", c)
	}
}
