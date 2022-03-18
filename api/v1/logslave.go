package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

type LogData struct {
	ProjectId int    `json:"project_id"`
	Uid       int    `json:"uid" form:"uid"`
	Level     int    `json:"level" form:"level"`
	Msg       string `json:"msg" form:"msg"`
}

func Push(c *gin.Context) {
	var L LogData
	_ = c.ShouldBind(&L)
	if err := util.Verify(L, util.LogReceiveVerify); err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	L.Uid, _ = request.GetUid(c)
	L.ProjectId, _ = request.GetAppId(c)

	str, _ := json.Marshal(L)
	global.V.Zap.Info(string(str))

	httpresponse.OkWithDetailed("", "已收录", c)
}
