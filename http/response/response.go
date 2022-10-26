//http 响应公共处理
package httpresponse

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"zgoframe/core"
	"zgoframe/core/global"
	"zgoframe/http/request"
	"zgoframe/util"
)

//@description 公共HTTP响应结构体
type Response struct {
	Code int         `json:"code"` //状态码，200是OK，4代表有发生错误
	Data interface{} `json:"data"` //请求时有数据返回，会在此字段中
	Msg  string      `json:"msg"`  //如果有错误会写在此，如果有些提示信息也会放在这里
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// @Description 获取图片验证码
type Captcha struct {
	Id            string `json:"id"`             //使用(验证)时要带上这个ID，后端才能判断
	PicContent    string `json:"pic_content"`    //图片内容,base64
	ContentLength int    `json:"content_length"` //数字的总长度
}

////上传多张图片，处理结果
//type UploadRs struct{
//	Err 		string `json:"err"`			//上传图片是否发生错误，如果为空证明没有，此字段给上传多张图片使用
//	Url 		string	`json:"url"`		//上传成功的图片：URL新地址
//	FileName 	string `json:"file_name"`	//原图片文件名
//}

// @Description 上传图片结果
type HttpUploadRs struct {
	util.UploadRs
	FullLocalIpUrl     string `json:"full_local_ip_url"` //本地存储的访问地址
	FullLocalDomainUrl string `json:"full_local_domain_url"`
	Err                string `json:"err"` //上传图片是否发生错误，如果为空证明没有，此字段给上传多张图片使用
}

//@description 常量信息
type ConstInfo struct {
	List map[string]int
	Key  string
	Name string
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	myHeader, _ := request.GetMyHeader(c)
	//rid := c.GetHeader("request_id")
	headerResponse := request.HeaderResponse{}

	headerResponse.ProjectId = myHeader.ProjectId
	headerResponse.SourceType = myHeader.SourceType
	headerResponse.RequestId = myHeader.RequestId
	headerResponse.TraceId = myHeader.TraceId
	headerResponse.AutoIp = myHeader.AutoIp
	headerResponse.ClientReqTime = myHeader.ClientReqTime
	headerResponse.ReceiveTime = myHeader.ServerReceiveTime
	headerResponse.ResponseTime = util.GetNowTimeSecondToInt()

	httpResponse := util.HttpHeaderSureStructCovertSureMap(headerResponse)
	for k, v := range httpResponse {
		c.Header(k, v)
	}

	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

//快速响应-无输出数据
func Ok(c *gin.Context) {
	Result(core.HTTP_RES_COMM_SUCCESS, map[string]interface{}{}, "操作成功", c)
}

//快速响应-有简单类型(一个字符串)的输出信息
func OkWithMessage(message string, c *gin.Context) {
	Result(core.HTTP_RES_COMM_SUCCESS, map[string]interface{}{}, message, c)
}

//快速响应-有复杂的输出数据
func OkWithData(data interface{}, c *gin.Context) {
	Result(core.HTTP_RES_COMM_SUCCESS, data, "操作成功", c)
}

//快速响应-即有简单数据，也有复杂数据
func OkWithAll(data interface{}, message string, c *gin.Context) {
	Result(core.HTTP_RES_COMM_SUCCESS, data, message, c)
}

//快速响应-失败，无任何输出信息
func Fail(c *gin.Context) {
	Result(core.HTTP_RES_COMM_ERROR, map[string]interface{}{}, "操作失败", c)
}

//快速响应-失败，有些简单的输出信息
func FailWithMessage(message string, c *gin.Context) {
	global.V.Zap.Error("失败", zap.Any("err", message))
	Result(core.HTTP_RES_COMM_ERROR, map[string]interface{}{}, message, c)
}

func FailWithAll(data interface{}, message string, c *gin.Context) {
	global.V.Zap.Error("失败", zap.Any("err", message))
	Result(core.HTTP_RES_COMM_ERROR, data, message, c)
}

var ErrManager = &util.ErrMsg{}

// 一次请求，发生了一些错误，统一输出，但不停止，依然返回
func ErrWithAllByCode(code int, c *gin.Context) error {
	errInfo := ErrManager.New(code)
	//util.MyPrint("ErrWithAllByCode:",code,errInfo)
	Result(code, errInfo.Error(), errInfo.Error(), c)
	return errInfo
}
