package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
)

//非JWT的接口，公共接口，也是允许访问，但是得从HEADER里提取信用，做基础验证
func HeaderAuth() gin.HandlerFunc {
	res := httpresponse.Response{}
	return func(c *gin.Context) {
		global.V.Zap.Debug("middle HeaderAuth start:")
		header := request.GetMyHeader(c)
		//验证SourceType
		if !request.CheckPlatformExist(header.SourceType) {
			header.SourceType = request.PLATFORM_UNKNOW
			res.Code = 501
			res.Msg = "SourceType unknow"
			c.AbortWithStatusJSON(500, res)
			return
		}

		if header.ProjectId <= 0 {
			res.Code = 502
			res.Msg = "ProjectId <= 0"
			c.AbortWithStatusJSON(500, res)
			return
		}

		if header.Access == "" {
			res.Code = 503
			res.Msg = "ACCESS empty"
			c.AbortWithStatusJSON(500, res)
			return
		}

		project, empty := global.V.ProjectMng.GetById(header.ProjectId)
		if empty {
			res.Code = 504
			res.Msg = "projectId  empty"
			c.AbortWithStatusJSON(500, res)
			return
		}
		//fmt.Println(project.Access, " - ", header.Access)
		if project.Access != header.Access {
			res.Code = 505
			res.Msg = "ACCESS  error"
			c.AbortWithStatusJSON(500, res)
			return
		}
		global.V.Zap.Debug("middle HeaderAuth finish.")
	}
}
