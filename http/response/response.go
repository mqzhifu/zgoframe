//http 响应公共处理
package httpresponse

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"zgoframe/core/global"
	"zgoframe/http/request"
)

//公共HTTP响应结构体
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 4   //公共HTTP响应状态码：失败
	SUCCESS = 200 //公共HTTP响应状态码：成功
)

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

type SysCaptchaResponse struct {
	CaptchaId string `json:"captchaId"`
	PicPath   string `json:"picPath"`
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	myHeader := request.GetMyHeader(c)
	//rid := c.GetHeader("request_id")
	c.Header("X-Request-Id", myHeader.RequestId)
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

//快速响应-无输出数据
func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

//快速响应-有简单的输出信息
func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

//快速响应-有复杂的输出数据
func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", c)
}

//快速响应-即有简单数据，也有复杂数据
func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

//快速响应-失败，无任何输出信息
func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

//快速响应-失败，有些简单的输出信息
func FailWithMessage(message string, c *gin.Context) {
	global.V.Zap.Error("失败", zap.Any("err", message))
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	global.V.Zap.Error("失败", zap.Any("err", message))
	Result(ERROR, data, message, c)
}
