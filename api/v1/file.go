package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

func FileUploadReal(c *gin.Context, category int) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		util.MyPrint("err1:", err.Error())
		return
	}
	syncOss, _ := strconv.Atoi(c.PostForm("sync_oss"))
	hashDir, _ := strconv.Atoi(c.PostForm("hash_dir"))
	module := GetFormParaModule(c)

	var uploadRs util.UploadRs
	switch category {
	case util.FILE_TYPE_IMG:
		uploadRs, err = global.V.ImgManager.UploadOne(header, module, hashDir, syncOss)
	case util.FILE_TYPE_DOC:
		uploadRs, err = global.V.DocsManager.UploadOne(header, module, hashDir, syncOss)
	case util.FILE_TYPE_VIDEO:
		uploadRs, err = global.V.VideoManager.UploadOne(header, module, hashDir, syncOss)
	case util.FILE_TYPE_PACKAGES:
		uploadRs, err = global.V.PackagesManager.UploadOne(header, module, hashDir, syncOss)
	}

	util.MyPrint("uploadRs:", uploadRs, " err:", err)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		httpUploadRs := httpresponse.HttpUploadRs{}
		httpUploadRs.UploadRs = uploadRs
		ip, _ := util.GetLocalIp()
		httpUploadRs.FullLocalIpUrl = util.UrlAppendIpHost("http", httpUploadRs.LocalIpUrl, ip, global.C.Http.Port)
		httpUploadRs.FullLocalDomainUrl = util.UrlAppendDomain("http", httpUploadRs.LocalDomainUrl, global.C.Domain.Static, "")
		httpresponse.OkWithAll(httpUploadRs, "已上传", c)
	}
}
func FileUploadRealMulti(c *gin.Context, category int) {
	syncOss, _ := strconv.Atoi(c.PostForm("sync_oss"))
	hashDir, _ := strconv.Atoi(c.PostForm("hash_dir"))
	module := GetFormParaModule(c)
	form, err := c.MultipartForm()
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	//syncOss := c.PostForm("sync_oss")
	//fileUpload := global.GetUploadObj(category, module)
	// 获取所有图片
	files := form.File["files"]
	if len(files) < 1 {
		httpresponse.FailWithMessage("请至少上传一个文件.", c)
		return
	}
	util.MyPrint("files len:", len(files))

	ip, _ := util.GetLocalIp()
	errList := []httpresponse.HttpUploadRs{}
	for _, file := range files {
		httpUploadRs := httpresponse.HttpUploadRs{}

		var uploadRs util.UploadRs
		switch category {
		case util.FILE_TYPE_IMG:
			uploadRs, err = global.V.ImgManager.UploadOne(file, module, hashDir, syncOss)
		case util.FILE_TYPE_DOC:
			uploadRs, err = global.V.DocsManager.UploadOne(file, module, hashDir, syncOss)
		case util.FILE_TYPE_VIDEO:
			uploadRs, err = global.V.VideoManager.UploadOne(file, module, hashDir, syncOss)
		case util.FILE_TYPE_PACKAGES:
			uploadRs, err = global.V.PackagesManager.UploadOne(file, module, hashDir, syncOss)
		}
		//uploadRs, err := global.V.ImgManager.UploadOne(file, module, hashDir, syncOss)
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		httpUploadRs.UploadRs = uploadRs
		httpUploadRs.FullLocalIpUrl = util.UrlAppendIpHost("http", httpUploadRs.LocalIpUrl, ip, global.C.Http.Port)
		httpUploadRs.FullLocalDomainUrl = util.UrlAppendDomain("http", httpUploadRs.LocalDomainUrl, global.C.Domain.Static, "")
		httpUploadRs.Err = errMsg
		errList = append(errList, httpUploadRs)
	}

	httpresponse.OkWithAll(errList, "ok", c)
}

// @Tags file
// @Summary 上传一张图片( http-form 表单模式 )
// @Security ApiKeyAuth
// @Description 单图片上限2M。支持格式："jpg", "jpeg", "png", "gif", "x-png", "png", "bmp", "pjpeg", "x-icon", "svg", "webp"。
// @Param	X-Source-Type 	header 		string 	true 	"来源" Enums(11,12,21,22)
// @Param	X-Project-Id  	header 		string 	true 	"项目ID" default(6)
// @Param	X-Access      	header 		string 	true 	"访问KEY" default(imzgoframe)
// @Param	file 			formData 	file 	true 	"文件(html中的input的name)"
// @Param	module 			formData 	string 	false 	"模块/业务名，可用于给文件名加前缀目录，注：开头和结尾都不要加反斜杠"
// @Param  	sync_oss 		formData 	int 	false 	"是否同步到云 oss 1 是,2 否" default(2)
// @Param  	hash_dir 		formData 	int 	false 	"自动创建:前缀目录(hash), 0 不使用,1 月, 2 天,3 小时" default(0)
// @Accept 	multipart/form-data
// @Produce	application/json
// @Success 200 {object} httpresponse.HttpUploadRs "上传结果"
// @Router 	/file/upload/img/one [POST]
func FileUploadImgOne(c *gin.Context) {
	FileUploadReal(c, util.FILE_TYPE_IMG)
}

// @Tags file
// @Summary 上传多张图片
// @Security ApiKeyAuth
// @Description 单图片上限2M。支持格式："jpg", "jpeg", "png", "gif", "x-png", "png", "bmp", "pjpeg", "x-icon", "svg", "webp"。
// @Param	X-Source-Type 	header 		string 	true 	"来源" Enums(11,12,21,22)
// @Param	X-Project-Id  	header 		string 	true 	"项目ID" default(6)
// @Param	X-Access      	header 		string 	true 	"访问KEY" default(imzgoframe)
// @Param	files 			formData 	file 	true 	"文件(html中的input的name)"
// @Param	module 			formData 	string 	false 	"模块/业务名，可用于给文件名加前缀目录，注：开头和结尾都不要加反斜杠"
// @Param  	sync_oss 		formData 	int 	false 	"是否同步到云 oss 1 是,2 否" default(2)
// @Param  	hash_dir 		formData 	int 	false 	"自动创建:前缀目录(hash), 0 不使用,1 月, 2 天,3 小时" default(0)
// @Accept multipart/form-data
// @Produce  application/json
// @Success 200 {object} []httpresponse.HttpUploadRs "每个图片的上传结果"
// @Router /file/upload/img/multi [post]
func FileUploadImgMulti(c *gin.Context) {
	FileUploadRealMulti(c, util.FILE_TYPE_IMG)
}

// @Tags file
// @Summary 上传图片 - 流模式 - base64
// @Security ApiKeyAuth
// @Description 有时前端并没有具体文件，而是在与用户交互中：动态产生的文件(图片)流，如：截图(canvas)，这时候直接把文件流传输后端即可
// @Description 单图片上限2M。支持格式："jpg", "jpeg", "png", "gif", "x-png", "png", "bmp", "pjpeg", "x-icon", "svg", "webp"。
// @Param	X-Source-Type 	header 		string 	true 	"来源" Enums(11,12,21,22)
// @Param	X-Project-Id  	header 		string 	true 	"项目ID" default(6)
// @Param	X-Access      	header 		string 	true 	"访问KEY" default(imzgoframe)
// @Param data body request.UploadFile true "基础信息"
// @Accept application/json
// @Produce  application/json
// @Success 200 {object} httpresponse.HttpUploadRs "下载结果"
// @Router /file/upload/img/one/stream/base64 [POST]
func FileUploadImgOneStreamBase64(c *gin.Context) {
	var form request.UploadFile
	c.ShouldBind(&form)

	if form.Stream == "" {
		httpresponse.FailWithMessage("stream empty!!!", c)
		return
	}
	module := GetModule(c, form.Module)
	uploadRs, err := global.V.ImgManager.UploadOneByStream(form.Stream, util.FILE_TYPE_IMG, module, form.HashDir)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		httpUploadRs := httpresponse.HttpUploadRs{}
		httpUploadRs.UploadRs = uploadRs
		ip, _ := util.GetLocalIp()
		httpUploadRs.FullLocalIpUrl = util.UrlAppendIpHost("http", httpUploadRs.LocalIpUrl, ip, global.C.Http.Port)
		httpUploadRs.FullLocalDomainUrl = util.UrlAppendDomain("http", httpUploadRs.LocalDomainUrl, global.C.Domain.Static, "")
		httpresponse.OkWithAll(httpUploadRs, "已上传", c)
	}

}

// @Tags file
// @Summary 上传一个文档
// @Security ApiKeyAuth
// @Description 单文件上限20M。支持格式："txt", "doc", "docx", "dotx", "json", "cvs", "xls", "xlsx", "sql", "msword", "ppt", "pptx", "pdf", "wps", "vsd"
// @Param	X-Source-Type 	header 		string 	true 	"来源" Enums(11,12,21,22)
// @Param	X-Project-Id  	header 		string 	true 	"项目ID" default(6)
// @Param	X-Access      	header 		string 	true 	"访问KEY" default(imzgoframe)
// @Param	file 			formData 	file 	true 	"文件(html中的input的name)"
// @Param	module 			formData 	string 	false 	"模块/业务名，可用于给文件名加前缀目录，注：开头和结尾都不要加反斜杠"
// @Param  	sync_oss 		formData 	int 	false 	"是否同步到云 oss 1 是,2 否" default(2)
// @Param  	hash_dir 		formData 	int 	false 	"自动创建:前缀目录(hash), 0 不使用,1 月, 2 天,3 小时" default(0)
// @Accept 	multipart/form-data
// @Produce	application/json
// @Success 200 {object} httpresponse.HttpUploadRs "上传结果"
// @Router 	/file/upload/doc/one [POST]
func FileUploadDocOne(c *gin.Context) {
	FileUploadReal(c, util.FILE_TYPE_DOC)
}

// @Tags file
// @Summary 上传多个文档
// @Description 单文件上限20M。支持格式："txt", "doc", "docx", "dotx", "json", "cvs", "xls", "xlsx", "sql", "msword", "ppt", "pptx", "pdf", "wps", "vsd"
// @Param	X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param	X-Project-Id  header string true "项目ID" default(6)
// @Param	X-Access      header string true "访问KEY" default(imzgoframe)
// @Param	files 		formData file 	true 	"文件(html中的input的name)"
// @Param	module 		formData string false 	"模块/业务名，可用于给文件名加前缀目录，注：开头和结尾都不要加反斜杠"
// @Param  	sync_oss 	formData int 	false 	"是否同步到云oss 1是2否" default(2)
// @Param  	hash_dir 	formData int 	false 	"自动创建前缀目录 0不使用1月2天3小时" default(0)
// @Accept multipart/form-data
// @Produce  application/json
// @Success 200 {object} []httpresponse.HttpUploadRs "每个图片的上传结果"
// @Router /file/upload/doc/multi [post]
func FileUploadDocMulti(c *gin.Context) {
	FileUploadRealMulti(c, util.FILE_TYPE_DOC)
}

// @Tags file
// @Summary 上传一个视频文件
// @Security ApiKeyAuth
// @Description 单文件上限20M。支持格式："mp4", "avi", "rm", "mkv", "wmv", "mov", "flv", "fla", "rmvb", "m3u8", "webm", "ts", "wav"
// @Param	X-Source-Type 	header 		string 	true 	"来源" Enums(11,12,21,22)
// @Param	X-Project-Id  	header 		string 	true 	"项目ID" default(6)
// @Param	X-Access      	header 		string 	true 	"访问KEY" default(imzgoframe)
// @Param	file 			formData 	file 	true 	"文件(html中的input的name)"
// @Param	module 			formData 	string 	false 	"模块/业务名，可用于给文件名加前缀目录，注：开头和结尾都不要加反斜杠"
// @Param  	sync_oss 		formData 	int 	false 	"是否同步到云 oss 1 是,2 否" default(2)
// @Param  	hash_dir 		formData 	int 	false 	"自动创建:前缀目录(hash), 0 不使用,1 月, 2 天,3 小时" default(0)
// @Accept 	multipart/form-data
// @Produce	application/json
// @Success 200 {object} httpresponse.HttpUploadRs "上传结果"
// @Router 	/file/upload/video/one [POST]
func FileUploadVideoOne(c *gin.Context) {
	FileUploadReal(c, util.FILE_TYPE_VIDEO)
}

// @Tags file
// @Summary 上传一个压缩包
// @Security ApiKeyAuth
// @Description 单文件上限 50 M。支持格式："zip", "rar", "apk", "tar", "jar", "7z", "gz", "rz"
// @Param	X-Source-Type 	header 		string 	true 	"来源" Enums(11,12,21,22)
// @Param	X-Project-Id  	header 		string 	true 	"项目ID" default(6)
// @Param	X-Access      	header 		string 	true 	"访问KEY" default(imzgoframe)
// @Param	file 			formData 	file 	true 	"文件(html中的input的name)"
// @Param	module 			formData 	string 	false 	"模块/业务名，可用于给文件名加前缀目录，注：开头和结尾都不要加反斜杠"
// @Param  	sync_oss 		formData 	int 	false 	"是否同步到云 oss 1 是,2 否" default(2)
// @Param  	hash_dir 		formData 	int 	false 	"自动创建:前缀目录(hash), 0 不使用,1 月, 2 天,3 小时" default(0)
// @Accept 	multipart/form-data
// @Produce	application/json
// @Success 200 {object} httpresponse.HttpUploadRs "上传结果"
// @Router 	/file/upload/packages/one [POST]
func FileUploadPackagesOne(c *gin.Context) {
	FileUploadReal(c, util.FILE_TYPE_PACKAGES)
}

// @Tags file
// @Summary 大文件下载(暂未实现，后续补充)
// @Security ApiKeyAuth
// @Description 大文件走NGINX不现实，而且，中间断了后，无法续传
// @Param	X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param	X-Project-Id  header string true "项目ID" default(6)
// @Param	X-Access      header string true "访问KEY" default(imzgoframe)
// @Param 	path formData string false "文件相对路径"
// @Accept 	multipart/form-data
// @Produce application/json
// @Success 200 {object} httpresponse.HttpUploadRs "下载结果"
// @Router 	/file/big/download [post]
func FileBigDownload(c *gin.Context) {

	//fileUpload := global.GetUploadObj(1, "")
	////filePath := c.PostForm("path")
	//////filePath := c.Query("path")
	//filePath := "1.jpg"
	//util.MyPrint("FileBigDownload filePath:", filePath)
	//
	//err := fileUpload.Download(filePath, c)
	//util.MyPrint(" fileUpload.DownloadBig return err:", err)
	//
	////headerRange := c.Request.Header.Get("Range")
	////c.Header("Content-Ranges","bytes 0-1023/1024")

}

//分断续传，使用http header:ranges ，C端首次请求：获取文件基础信息
func FileDownloadInfo(c *gin.Context) {
	//filePath := c.PostForm("path")
	//fileUpload := global.GetUploadObj(1, "")
	//fileDownInfo, err := fileUpload.DownloadFileInfo(filePath)
	//if err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//	return
	//}
	//
	//c.Header("Accept-Ranges", "bytes")
	//c.Header("Content-Length", strconv.Itoa(int(fileDownInfo.FileSize)))
}

func GetFormParaModule(c *gin.Context) string {

	module := c.PostForm("module")
	projectId := request.GetProjectIdByHeader(c)
	if projectId > 0 {
		projectInfo, empty := global.V.ProjectMng.GetById(projectId)
		//util.MyPrint("projectInfo:=====", projectInfo)
		if !empty {
			if module != "" {
				module = projectInfo.Name + "/" + module
			} else {
				module = projectInfo.Name
			}
		}
	}

	util.MyPrint("GetFormParaModule projectId:", strconv.Itoa(projectId), "module:"+module)
	return module
}

func GetModule(c *gin.Context, module string) string {
	projectId := request.GetProjectIdByHeader(c)
	if projectId > 0 {
		projectInfo, empty := global.V.ProjectMng.GetById(projectId)
		//util.MyPrint("projectInfo:=====", projectInfo)
		if !empty {
			if module != "" {
				module = projectInfo.Name + "/" + module
			} else {
				module = projectInfo.Name
			}
		}
	}

	util.MyPrint("GetModule projectId:", strconv.Itoa(projectId), "module:"+module)
	return module
}
