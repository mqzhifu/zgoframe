package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/util"
)

// @Tags persistence
// @Summary 收集日志
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
// @Summary 上传文件
// @Description 目前是本地存一份，同步到OSS一份，目录结构是根据天做hash
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param file formData file true "文件,html中的input file"
// @Param category formData int true "上传的文件类型，1全部2图片3文档4视频,后端会根据类型做验证"
// @Param module formData string false "模块/业务名，可用于给文件名加前缀目录"
// @Accept multipart/form-data
// @Produce  application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /persistence/file/upload [post]
func Upload(c *gin.Context){
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		util.MyPrint("err1:", err.Error())
		return
	}

	category ,_:= strconv.Atoi (c.PostForm("category") )
	module  := c.PostForm("module")

	fileUpload := GetUploadObj(category,module)
	relativeFileName,err := fileUpload.UploadOne(header)

	util.MyPrint("uploadRs:",relativeFileName, " err:",err)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(),c)
	}else{
		httpresponse.OkWithAll(global.C.Oss.Bucket + "." +global.C.Oss.Endpoint + "/" + relativeFileName, "已上传", c)
	}

}

func GetUploadObj(category int,module string)*util.FileUpload{
	//projectId := request.GetProjectId(c)
	fileUploadOption := util.FileUploadOption{
		FilePrefix		: module,
		MaxSize			: 8,
		Category		: category,
		FileHashType	: util.FILE_HASH_DAY,
		Path			: global.C.Upload.Path,
		OssAccessKeyId	: global.C.Oss.AccessKeyId,
		OssEndpoint		: global.C.Oss.Endpoint,
		OssBucketName 	: global.C.Oss.Bucket,
		OssAccessKeySecret: global.C.Oss.AccessKeySecret,
	}

	fileUpload := util.NewFileUpload( fileUploadOption )
	return fileUpload
}



// @Tags persistence
// @Summary 上传多文件
// @Description 目前是本地存一份，同步到OSS一份，目录结构是根据天做hash
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param files formData file true "文件,html中的input files multiple "
// @Param category formData int true "上传的文件类型，1全部2图片3文档4视频,后端会根据类型做验证"
// @Param module formData string false "模块/业务名，可用于给文件名加前缀目录"
// @Accept multipart/form-data
// @Produce  application/json
// @Success 200 {object} []httpresponse.UploadRs "每个图片的上传结果"
// @Router /persistence/file/upload/multi [post]
func UploadMulti(c *gin.Context){
	category ,_:= strconv.Atoi (c.PostForm("category") )
	module  := c.PostForm("module")

	form, err := c.MultipartForm()
	if err != nil {
		httpresponse.FailWithMessage(err.Error(),c)
		return
	}

	fileUpload := GetUploadObj(category,module)
	// 获取所有图片
	files := form.File["files"]
	if len(files) < 1{
		httpresponse.FailWithMessage("请至少上传一个文件.",c)
		return
	}
	util.MyPrint("files len:",len(files))


	errList := []httpresponse.UploadRs{}
	for _, file := range files {
		url,err := fileUpload.UploadOne(file)
		errMsg := ""
		if err != nil{
			errMsg = err.Error()
		}

		errList = append(errList, httpresponse.UploadRs { Url: url,Err: errMsg,FileName :file.Filename }  )
	}

	httpresponse.OkWithAll(errList,"ok",c)
}
