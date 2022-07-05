package v1

import (
	"bufio"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
	"strings"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/util"
)

// @Tags persistence
// @Summary 收集日志,目前是存于MYSQL中，后期可以优化成文件或ES中
// @Description 用于后台统计
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.StatisticsLogData true "用户信息"
// @Produce  application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /persistence/log/push [post]
func LogPush(c *gin.Context) {
	var form request.StatisticsLogData
	_ = c.ShouldBind(&form)
	//if err := util.Verify(L, util.LogReceiveVerify); err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//	return
	//}
	//L.Uid, _ = request.GetUid(c)
	//L.ProjectId, _ = request.GetProjectId(c)

	header ,err  := request.GetMyHeader(c)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(),c)
		return
	}

	if form.Action == "" {
		httpresponse.FailWithMessage("Action is empty",c)
		return
	}

	headerJsonStr ,_ := json.Marshal(header)
	msgJsonStr , _ := json.Marshal(form.Msg)
	//header.BaseInfo = nil
	//headerCommon := header
	statisticsLog := model.StatisticsLog{
		HeaderCommon:string(headerJsonStr),
		HeaderBase: "",
		Uid: form.Uid,
		ProjectId: form.ProjectId,
		Action:form.Action,
		Category: form.Category,
		Msg: string(msgJsonStr),
	}
	//util.ExitPrint(statisticsLog)
	err = global.V.Gorm.Create(&statisticsLog).Error
	if err != nil{
		httpresponse.FailWithMessage("db insert failed err:"+err.Error(),c)
		return
	}

	//str, _ := json.Marshal(L)
	//global.V.Zap.Info(string(str))

	httpresponse.OkWithAll("", "已收录", c)
}

// @Tags persistence
// @Summary 收集日志(文件)-目前是存于MYSQL中，后期可以优化成文件或ES中
// @Description 用于后台统计
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param file formData file true "文件(html中的input的name)"
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /persistence/log/push/file [post]
func LogPushFile(c *gin.Context) {
	//var form request.StatisticsLogData
	//_ = c.ShouldBind(&form)

	_, header, err := c.Request.FormFile("file")
	if err != nil {
		util.MyPrint("c.Request.FormFile err1:", err.Error())
		return
	}

	category := util.FILE_TYPE_DOC
	module  := "log"

	fileUpload := global.GetUploadObj(category,module)
	uploadRs,err := fileUpload.UploadOne(header,2)
	if err != nil{
		util.MyPrint("fileUpload.UploadOne err:", err.Error())
		return
	}

	localFilePath := uploadRs.LocalDiskPath + "/" + uploadRs.Filename
	fd, err := os.Open(localFilePath)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(fd)
	headerJsonStr ,_ := json.Marshal(header)
	success := 0
	failed := 0
	for {
		//func (b *Reader) ReadString(delim byte) (string, error) {}
		line, err := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		util.MyPrint(line)
		arr := strings.Split(line, ",_,")
		if len(arr) <  4{
			failed++
			continue
		}
		projectId ,_:= strconv.Atoi(arr[0])
		uid ,_:= strconv.Atoi(arr[1])
		logCategory ,_:= strconv.Atoi(arr[2])
		statisticsLog := model.StatisticsLog{
			HeaderCommon:string(headerJsonStr),
			HeaderBase: "",
			ProjectId: projectId ,
			Uid : uid,
			Category :logCategory,
			Action: arr[3],
		}

		//util.ExitPrint(statisticsLog)
		err = global.V.Gorm.Create(&statisticsLog).Error
		if err != nil{
			failed++
			httpresponse.FailWithMessage("db insert failed err:"+err.Error(),c)
			return
		}
		success++

	}

	httpresponse.OkWithAll("", "已收录", c)
}
