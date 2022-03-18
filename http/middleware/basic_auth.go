package httpmiddleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
)

func BasicAuthHeader() gin.HandlerFunc {
	res := httpresponse.Response{}
	return func(c *gin.Context) {
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
		fmt.Println(project.Access, " - ", header.Access)
		if project.Access != header.Access {
			res.Code = 505
			res.Msg = "ACCESS  error"
			c.AbortWithStatusJSON(500, res)
			return
		}
	}
}
