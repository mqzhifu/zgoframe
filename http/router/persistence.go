package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func Persistence(Router *gin.RouterGroup) {
	persistenceRouter := Router.Group("persistence")
	{
		persistenceRouter.POST("log/push", v1.LogPush)
		persistenceRouter.POST("log/push/file", v1.LogPushFile)
		persistenceRouter.POST("log/push/file/json", v1.LogPushFileJson)
	}
}

func File(Router *gin.RouterGroup) {
	persistenceRouter := Router.Group("file")
	{
		//persistenceRouter.POST("upload/img", v1.Upload)
		//persistenceRouter.POST("file/upload/multi", v1.UploadMulti)
		////persistenceRouter.POST("file/big/download", v1.FileBigDownload)
		//persistenceRouter.GET("file/big/download", v1.FileBigDownload)
		//persistenceRouter.POST("file/upload/stream", v1.FileUploadStream)

		persistenceRouter.POST("upload/img/one", v1.FileUploadImgOne)
		persistenceRouter.POST("upload/img/one/stream/base64", v1.FileUploadImgOneStreamBase64)
		persistenceRouter.POST("upload/img/multi", v1.FileUploadImgMulti)

		persistenceRouter.POST("upload/doc/one", v1.FileUploadDocOne)
		//persistenceRouter.POST("upload/doc/one/stream/base64", v1.Upload)
		persistenceRouter.POST("upload/doc/multi", v1.FileUploadDocMulti)
		persistenceRouter.POST("upload/packages/one", v1.FileUploadPackagesOne)
		persistenceRouter.POST("upload/video/one", v1.FileUploadVideoOne)
		persistenceRouter.POST("delete/one", v1.FileDeleteOne)
		persistenceRouter.POST("copy/one", v1.FileCopyOne)
		persistenceRouter.POST("move/one", v1.FileMoveOne)

		//
		//persistenceRouter.POST("upload/video/one", v1.Upload)
		//persistenceRouter.POST("upload/video/one/stream/base64", v1.Upload)
		//persistenceRouter.POST("upload/video/multi", v1.Upload)
	}
}
