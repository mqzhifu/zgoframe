package v1

import (
	"github.com/gin-gonic/gin"
)

// @Tags FrameSync
// @Summary 配置信息
// @Description  配置信息
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Produce application/json
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /frame/sync/config [get]
func FrameSyncConfig(c *gin.Context) {
	//op := global.V.MyService.RoomManage.GetById()
	//httpresponse.OkWithAll(op, "ok", c)
}
