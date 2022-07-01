package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Tags file
// @Summary 上传一张图片
// @Description 目前是:本地存一份，同步到OSS一份，目录结构是根据(天)做hash,注：form要加上属性 enctype=multipart/form-data
// @Param	X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param	X-Project-Id header string true "项目ID" default(6)
// @Param	X-Access header string true "访问KEY" default(imzgoframe)
// @Param	file formData file true "文件(html中的input的name)"
// @Param	module formData string false "模块/业务名，可用于给文件名加前缀目录"
// @Accept 	multipart/form-data
// @Produce	application/json
// @Success 200 {object} httpresponse.HttpUploadRs "上传结果"
// @Router 	/file/upload/img/one [POST]
func FileUploadImgOne(c *gin.Context){
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		util.MyPrint("err1:", err.Error())
		return
	}

	category := util.FILE_TYPE_IMG
	module  := c.PostForm("module")

	fileUpload := global.GetUploadObj(category,module)
	uploadRs,err := fileUpload.UploadOne(header)

	util.MyPrint("uploadRs:",uploadRs, " err:",err)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(),c)
	}else{
		httpUploadRs := httpresponse.HttpUploadRs{}
		httpUploadRs.UploadRs = uploadRs
		ip , _:= util.GetLocalIp()
		httpUploadRs.FullLocalIpUrl = util.UrlAppendIpHost("http",httpUploadRs.LocalIpUrl,ip,global.C.Http.Port)
		httpUploadRs.FullLocalDomainUrl = util.UrlAppendDomain("http",httpUploadRs.LocalDomainUrl,"local.static.com","")
		httpresponse.OkWithAll( httpUploadRs, "已上传", c)
	}

}

// @Tags file
// @Summary 上传多张图片
// @Description 目前是本地存一份，同步到OSS一份，目录结构是根据天做hash，注：form enctype=multipart/form-data
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param files formData file true "文件(html中的input的name),注：multiple 属性"
// @Param module formData string false "模块/业务名，可用于给文件名加前缀目录"
// @Accept multipart/form-data
// @Produce  application/json
// @Success 200 {object} []httpresponse.HttpUploadRs "每个图片的上传结果"
// @Router /file/upload/img/multi [post]
func FileUploadImgMulti(c *gin.Context){
	//category ,_:= strconv.Atoi (c.PostForm("category") )
	category := util.FILE_TYPE_IMG
	module  := c.PostForm("module")

	form, err := c.MultipartForm()
	if err != nil {
		httpresponse.FailWithMessage(err.Error(),c)
		return
	}

	fileUpload := global.GetUploadObj(category,module)
	// 获取所有图片
	files := form.File["files"]
	if len(files) < 1{
		httpresponse.FailWithMessage("请至少上传一个文件.",c)
		return
	}
	util.MyPrint("files len:",len(files))

	ip , _:= util.GetLocalIp()
	errList := []httpresponse.HttpUploadRs{}
	for _, file := range files {
		httpUploadRs := httpresponse.HttpUploadRs{}

		uploadRs,err := fileUpload.UploadOne(file)
		errMsg := ""
		if err != nil{
			errMsg = err.Error()
		}
		httpUploadRs.UploadRs = uploadRs
		httpUploadRs.FullLocalIpUrl = util.UrlAppendIpHost("http",httpUploadRs.LocalIpUrl,ip,global.C.Http.Port)
		httpUploadRs.FullLocalDomainUrl = util.UrlAppendDomain("http",httpUploadRs.LocalDomainUrl,"local.static.com","")
		httpUploadRs.Err = errMsg
		errList = append(errList, httpUploadRs  )
	}

	httpresponse.OkWithAll(errList,"ok",c)
}

// @Tags file
// @Summary 上传图片 - 流模式 - base64
// @Description 有时前端并没有具体文件，而是在与用户交互中：动态产生的文件流，如：截图(canvas)，这时候直接把文件流传输即可
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param stream formData string false "文件流,例：data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAHgAAAB4CAMAAAAOus ......."
// @Param	module formData string false "模块/业务名，可用于给文件名加前缀目录"
// @Accept multipart/form-data
// @Produce  application/json
// @Success 200 {object} httpresponse.HttpUploadRs "下载结果"
// @Router /file/upload/img/one/stream/base64 [POST]
func FileUploadImgOneStreamBase64(c *gin.Context){
	//category ,_:= strconv.Atoi (c.PostForm("category") )
	//if category  <= 0 {
	//	httpresponse.FailWithMessage("category empty!!!",c)
	//	return
	//}
	category := util.FILE_TYPE_IMG
	fileUpload := global.GetUploadObj(category,"")

	stream  := c.PostForm("stream")

	util.MyPrint("stream:",stream)
	if stream == ""{
		httpresponse.FailWithMessage("stream empty!!!",c)
		return
	}

	//stream ,_ = url.QueryUnescape(stream)

	uploadRs ,err := fileUpload.UploadOneByStream(stream,category)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(),c)
	}else{
		httpUploadRs := httpresponse.HttpUploadRs{}
		httpUploadRs.UploadRs = uploadRs
		ip , _:= util.GetLocalIp()
		httpUploadRs.FullLocalIpUrl = util.UrlAppendIpHost("http",httpUploadRs.LocalIpUrl,ip,global.C.Http.Port)
		httpUploadRs.FullLocalDomainUrl = util.UrlAppendDomain("http",httpUploadRs.LocalDomainUrl,"local.static.com","")
		httpresponse.OkWithAll( httpUploadRs, "已上传", c)
	}



}


// @Tags file
// @Summary 上传一个文档
// @Description 目前是:本地存一份，同步到OSS一份，目录结构是根据(天)做hash,注：form要加上属性 enctype=multipart/form-data
// @Param	X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param	X-Project-Id header string true "项目ID" default(6)
// @Param	X-Access header string true "访问KEY" default(imzgoframe)
// @Param	file formData file true "文件(html中的input的name)"
// @Param	module formData string false "模块/业务名，可用于给文件名加前缀目录"
// @Accept 	multipart/form-data
// @Produce	application/json
// @Success 200 {object} httpresponse.HttpUploadRs "上传结果"
// @Router 	/file/upload/doc/one [POST]
func FileUploadDocOne(c *gin.Context){
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		util.MyPrint("err1:", err.Error())
		return
	}

	category := util.FILE_TYPE_DOC
	module  := c.PostForm("module")

	fileUpload := global.GetUploadObj(category,module)
	uploadRs,err := fileUpload.UploadOne(header)

	util.MyPrint("uploadRs:",uploadRs, " err:",err)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(),c)
	}else{
		httpUploadRs := httpresponse.HttpUploadRs{}
		httpUploadRs.UploadRs = uploadRs
		ip , _:= util.GetLocalIp()
		httpUploadRs.FullLocalIpUrl = util.UrlAppendIpHost("http",httpUploadRs.LocalIpUrl,ip,global.C.Http.Port)
		httpUploadRs.FullLocalDomainUrl = util.UrlAppendDomain("http",httpUploadRs.LocalDomainUrl,"local.static.com","")
		httpresponse.OkWithAll( httpUploadRs, "已上传", c)
	}

}

// @Tags file
// @Summary 上传多个文档
// @Description 目前是本地存一份，同步到OSS一份，目录结构是根据天做hash，注：form enctype=multipart/form-data
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param files formData file true "文件(html中的input的name),注：multiple 属性"
// @Param module formData string false "模块/业务名，可用于给文件名加前缀目录"
// @Accept multipart/form-data
// @Produce  application/json
// @Success 200 {object} []httpresponse.HttpUploadRs "每个图片的上传结果"
// @Router /file/upload/doc/multi [post]
func FileUploadDocMulti(c *gin.Context){
	//category ,_:= strconv.Atoi (c.PostForm("category") )
	category := util.FILE_TYPE_IMG
	module  := c.PostForm("module")

	form, err := c.MultipartForm()
	if err != nil {
		httpresponse.FailWithMessage(err.Error(),c)
		return
	}

	fileUpload := global.GetUploadObj(category,module)
	// 获取所有图片
	files := form.File["files"]
	if len(files) < 1{
		httpresponse.FailWithMessage("请至少上传一个文件.",c)
		return
	}
	util.MyPrint("files len:",len(files))

	ip , _:= util.GetLocalIp()
	errList := []httpresponse.HttpUploadRs{}
	for _, file := range files {
		httpUploadRs := httpresponse.HttpUploadRs{}

		uploadRs,err := fileUpload.UploadOne(file)
		errMsg := ""
		if err != nil{
			errMsg = err.Error()
		}
		httpUploadRs.UploadRs = uploadRs
		httpUploadRs.FullLocalIpUrl = util.UrlAppendIpHost("http",httpUploadRs.LocalIpUrl,ip,global.C.Http.Port)
		httpUploadRs.FullLocalDomainUrl = util.UrlAppendDomain("http",httpUploadRs.LocalDomainUrl,"local.static.com","")
		httpUploadRs.Err = errMsg
		errList = append(errList, httpUploadRs  )
	}

	httpresponse.OkWithAll(errList,"ok",c)
}




// @Tags file
// @Summary 大文件下载
// @Description 大文件走NGINX不现实，而且，中间断了后，无法续传
// @Param 	X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param 	X-Project-Id header string true "项目ID" default(6)
// @Param 	X-Access header string true "访问KEY" default(imzgoframe)
// @Param 	path formData string false "文件相对路径"
// @Accept 	multipart/form-data
// @Produce application/json
// @Success 200 {object} httpresponse.HttpUploadRs "下载结果"
// @Router 	/file/big/download [post]
func FileBigDownload(c *gin.Context){

	fileUpload := global.GetUploadObj(1,"")
	//filePath := c.PostForm("path")
	////filePath := c.Query("path")
	filePath := "1.jpg"
	util.MyPrint("FileBigDownload filePath:",filePath)


	err := fileUpload.Download(filePath,c)
	util.MyPrint(" fileUpload.DownloadBig return err:",err)

	//headerRange := c.Request.Header.Get("Range")
	//c.Header("Content-Ranges","bytes 0-1023/1024")

}

func FileDownloadInfo(c *gin.Context){
	filePath := c.PostForm("path")
	fileUpload := global.GetUploadObj(1,"")
	fileDownInfo ,err := fileUpload.DownloadFileInfo(filePath)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(),c)
		return
	}

	c.Header("Accept-Ranges","bytes")
	c.Header("Content-Length",strconv.Itoa(int(fileDownInfo.FileSize)))
}


