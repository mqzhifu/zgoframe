package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/http/encrypt"
	"zgoframe/http/request"
	"zgoframe/model"
	"zgoframe/util"
)

// 非JWT的接口，公共接口，也是允许访问，但是得从HEADER里提取信用，做基础验证
func HeaderAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		global.V.Zap.Debug("http middleware <HeaderAuth> start:")
		header, err := request.GetMyHeader(c)
		if err != nil {
			util.MyPrint("err:" + err.Error())
			ErrAbortWithResponse(5105, c)
			return
		}
		//验证 SourceType
		if !model.CheckConstInList(model.GetConstListPlatform(), header.SourceType) {
			header.SourceType = model.PLATFORM_UNKNOW
			ErrAbortWithResponse(5100, c)
			return
		}

		if header.ProjectId <= 0 {
			ErrAbortWithResponse(5101, c)
			return
		}

		if header.Access == "" {
			ErrAbortWithResponse(5102, c)
			return
		}

		project, empty := global.V.ProjectMng.GetById(header.ProjectId)
		if empty {
			ErrAbortWithResponse(5103, c)
			return
		}
		if project.Access != header.Access {
			ErrAbortWithResponse(5104, c)
			return
		}
		//把 project 信息放到 context 中，主要是给 响应的时候使用
		c.Set("project", project)
		//CheckSign 得放在 DecodeBody 之前，因为 DecodeBody 会解码，复写 c.Request.Body
		_, err = encrypt.CheckSign(c)
		if err != nil {
			errCode, _ := strconv.Atoi(err.Error())
			ErrAbortWithResponse(errCode, c)
			return
		}

		_, err = encrypt.DecodeBody(c)
		if err != nil {
			util.MyPrint("SetupBody err:" + err.Error())
			ErrAbortWithResponse(5022, c)
			return
		}

		global.V.Zap.Debug("http middleware <HeaderAuth> finish.")
	}
}
