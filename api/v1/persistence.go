package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Tags persistence
// @Summary 收集日志
// @Description 用于后台统计
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.LogData true "用户信息"
// @Produce  application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /persistence/log/push [post]
func LogPush(c *gin.Context) {
	var L request.LogData
	_ = c.ShouldBind(&L)
	//if err := util.Verify(L, util.LogReceiveVerify); err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//	return
	//}
	//L.Uid, _ = request.GetUid(c)
	//L.ProjectId, _ = request.GetProjectId(c)

	str, _ := json.Marshal(L)
	global.V.Zap.Info(string(str))

	httpresponse.OkWithAll("", "已收录", c)
}

// @Tags persistence
// @Summary 上传文件
// @Description 上传文件
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param file formData file true "文件"
// @Param category formData int true "上传的文件类型，1全部2图片3文档"
// @Accept multipart/form-data
// @Produce  application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /persistence/file/upload [post]
func Upload(c *gin.Context){
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		util.MyPrint("err1:", err.Error())
		return
	}

	category ,_:= strconv.Atoi (c.PostForm("category") )

	ossBucketName := "xiaoztest"
	fileUploadOption := util.FileUploadOption{
		Category: category,
		Path: global.C.Upload.Path,
		OssAccessKeyId: global.C.Oss.AccessKeyId,
		OssAccessKeySecret: global.C.Oss.AccessKeySecret,
		OssEndpoint: global.C.Oss.Endpoint,
		OssBucketName : ossBucketName,
		FileHashType: util.FILE_HASH_DAY,
	}

	fileUpload := util.NewFileUpload( fileUploadOption )
	relativeFileName,err := fileUpload.UploadOne(file,header)

	util.MyPrint("uploadRs:",relativeFileName, " err:",err)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(),c)
	}else{
		httpresponse.OkWithAll(ossBucketName + "." +global.C.Oss.Endpoint + "/" + relativeFileName, "已上传", c)
	}

}
