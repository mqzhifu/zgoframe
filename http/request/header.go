package request

import (
	"errors"
	"github.com/gin-gonic/gin"
	"zgoframe/model"
)

func GetMyHeader(c *gin.Context) (hr HeaderRequest, err error) {
	myHeaderInterface, exists := c.Get("myHeader")
	if !exists {
		return hr, errors.New("get myHeader is empty~")
	}
	myHeader, ok := myHeaderInterface.(HeaderRequest)
	if !ok {
		return hr, errors.New("assertions failed: HeaderRequest")
	}
	return myHeader, nil
}

func GetMyProject(c *gin.Context) (project model.Project, err error) {
	myHeaderInterface, exists := c.Get("project")
	if !exists {
		return project, errors.New("get GetMyProject is empty~")
	}
	project, ok := myHeaderInterface.(model.Project)
	if !ok {
		return project, errors.New("assertions failed: HeaderRequest")
	}
	return project, nil
}

//func GetParserTokenData(c *gin.Context) (parserTokenData ParserTokenData, err error) {
//	parserTokenDataInter, exists := c.Get("parserTokenData")
//	if !exists {
//		global.V.Zap.Error("parserTokenData empty")
//		return parserTokenData, errors.New("parserTokenData empty")
//	}
//	parserTokenData = parserTokenDataInter.(ParserTokenData)
//	return parserTokenData, nil
//}

// 1. 从token中解出来的值里获取
// 2. 从DB中获取
func GetUid(c *gin.Context) (int, error) {
	user, err := GetUser(c)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

// 有4种方式获取：
// 1. 从token解出来的结构体内获取
// 2. 从token解出来的结构体内，再从DB中获取
// 3. header中也可以取这个值
func GetProjectId(c *gin.Context) (int, error) {

	customClaims, err := GetClaims(c)
	if err != nil {
		return 0, errors.New("Claims key not exist")
	}

	return customClaims.ProjectId, nil
}

func GetProjectIdByHeader(c *gin.Context) int {
	header, _ := GetMyHeader(c)
	return header.ProjectId
}

func GetSourceType(c *gin.Context) (int, error) {
	customClaims, err := GetClaims(c)
	if err != nil {
		return 0, errors.New("Claims key not exist")
	}

	return customClaims.SourceType, nil
}

func GetUser(c *gin.Context) (user model.User, err error) {
	u, exist := c.Get("user")
	if !exist {
		return user, errors.New("not exist")
	}
	user = u.(model.User)
	return user, nil
}

func GetClaims(c *gin.Context) (customClaims CustomClaims, err error) {
	cc, exist := c.Get("customClaims")
	if !exist {
		return customClaims, errors.New("not exist")
	}
	customClaims = cc.(CustomClaims)
	return customClaims, nil
}
