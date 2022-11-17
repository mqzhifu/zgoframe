package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
	"zgoframe/model"
	"zgoframe/util"
)

//非JWT的接口，公共接口，也是允许访问，但是得从HEADER里提取信用，做基础验证
func HeaderAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		global.V.Zap.Debug("http middleware <HeaderAuth> start:")
		header, err := request.GetMyHeader(c)
		if err != nil {
			util.MyPrint("err:" + err.Error())
			ErrAbortWithResponse(5105, c)
			return
		}
		//验证SourceType
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
		global.V.Zap.Debug("http middleware <HeaderAuth> finish.")
	}
}
