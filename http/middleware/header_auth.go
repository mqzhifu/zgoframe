package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
	"zgoframe/model"
)

//非JWT的接口，公共接口，也是允许访问，但是得从HEADER里提取信用，做基础验证
func HeaderAuth() gin.HandlerFunc {
	//res := httpresponse.Response{}
	return func(c *gin.Context) {
		global.V.Zap.Debug("middle HeaderAuth start:")
		header, _ := request.GetMyHeader(c)
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
			//res.Code = 504
			//res.Msg = "projectId  empty"
			//c.AbortWithStatusJSON(500, res)
			ErrAbortWithResponse(5103, c)
			return
		}
		//fmt.Println(project.Access, " - ", header.Access)
		if project.Access != header.Access {
			//res.Code = 505
			//res.Msg = "ACCESS  error"
			//c.AbortWithStatusJSON(500, res)
			ErrAbortWithResponse(5104, c)
			return
		}
		global.V.Zap.Debug("middle HeaderAuth finish.")
	}
}
