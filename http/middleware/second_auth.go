package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
)

type SecondAuthUser struct {
	Name string
	Ps   string
}

func GetSecondAuthUserList() []SecondAuthUser {
	userList := []SecondAuthUser{}
	userList = append(userList, SecondAuthUser{Name: "xiaoz", Ps: "qwerASDFzxcv"}, SecondAuthUser{Name: "test", Ps: "qweASD1234560"})
	return userList
}

func SecondAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		global.V.Base.Zap.Debug("http middleware <SecondAuth>  start:")
		myHeader, exist := c.Get("myHeader")
		if !exist {
			ErrAbortWithResponse(5105, c)
			return
		}
		myHeaderSt := myHeader.(request.HeaderRequest)
		if myHeaderSt.SecondAuthUname == "" || myHeaderSt.SecondAuthPs == "" {
			ErrAbortWithResponse(5290, c)
			return
		}
		rs := SecondAuthing(myHeaderSt.SecondAuthUname, myHeaderSt.SecondAuthPs)
		if !rs {
			ErrAbortWithResponse(5291, c)
			return
		}
		global.V.Base.Zap.Debug("http middleware <SecondAuth>  finish.")
		c.Next()
	}
}

func SecondAuthing(name string, ps string) bool {
	userList := GetSecondAuthUserList()
	for _, v := range userList {
		if v.Name == name && v.Ps == ps {
			return true
		}
	}
	return false
}
