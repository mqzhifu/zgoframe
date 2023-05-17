package httpmiddleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/http/request"
	"zgoframe/model"
	"zgoframe/util"
)

func Record() gin.HandlerFunc {
	return func(c *gin.Context) {
		prefix := "http middleware <Record>  "
		global.V.Zap.Debug(prefix + "start:")

		var body []byte
		var userId int
		if c.Request.Method != http.MethodGet {
			var err error
			body, err = ioutil.ReadAll(c.Request.Body)
			if err != nil {

				global.V.Zap.Error("read body from request error:", zap.Any("err", err))
			} else {
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
		}
		if claims, ok := c.Get("claims"); ok {
			waitUse := claims.(*request.CustomClaims)
			userId = int(waitUse.Id)
		} else {
			id, err := strconv.Atoi(c.Request.Header.Get("x-user-id"))
			if err != nil {
				userId = 0
			}
			userId = id
		}
		record := model.OperationRecord{
			Ip:     c.ClientIP(),
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   string(body),
			Uid:    userId,
		}
		// 存在某些未知错误 TODO
		//values := c.Request.Header.Values("content-type")
		//if len(values) >0 && strings.Contains(values[0], "boundary") {
		//	record.Body = "file"
		//}
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		startTime := util.GetNowTimeSecondToInt()
		//开始执行用户 业务 函数
		c.Next()
		//用户业务执行完毕后，需要对本次请求做收尾统计，并做持久化
		latency := util.GetNowTimeSecondToInt() - startTime //本次请求的总时长
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Latency = latency
		record.Resp = writer.body.String()

		global.V.Zap.Debug(prefix + "finish , func exec time:" + strconv.Itoa(latency))

		err := global.V.Gorm.Create(&record)
		if err != nil {
			global.V.Zap.Error(prefix+"create record error:", zap.Any("err", err))
		}
		//util.MyPrint("http middleware Create record err", err)
		//fmt.Println("opt final record:", record)
		//if err := service.CreateSysOperationRecord(record); err != nil {
		//	global.V.Zap.Error("create operation record error:", zap.Any("err", err))
		//}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
