package initialize

import "C"
import (
	"zgoframe/core/global"
	"zgoframe/util"
)

func InitFileManager() {
	baseDir := global.MainEnv.RootDir + "/" + global.C.Http.StaticPath + "/puzzle"
	oriImagePath := baseDir + "/ori/"
	shardImagePath := baseDir + "/shard/"
	thumbnailPath := baseDir + "/thumbnail/"
	global.V.ImageSlice = util.NewImageSlice(oriImagePath, shardImagePath, thumbnailPath)

	fileUploadOption := util.FileManagerOption{
		//FilePrefix:       module,
		Category:         util.FILE_TYPE_IMG,
		UploadDir:        global.C.FileManager.UploadPath,
		UploadMaxSize:    global.C.FileManager.UploadDocImgMaxSize,
		UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OPEN,
		//UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OFF,
		UploadStoreOSS:  util.UPLOAD_STORE_OSS_ALI,
		DownloadDir:     global.C.FileManager.DownloadPath,
		DownloadMaxSize: global.C.FileManager.DownloadMaxSize,
		FileHashType:    util.FILE_HASH_DAY,
		StaticDir:       global.C.Http.StaticPath,
		ProjectRootPath: global.MainEnv.RootDir,
		AliOss:          global.V.AliOss,
	}
	global.V.ImgManager = util.NewFileManagerUpload(fileUploadOption)

	fileUploadOption = util.FileManagerOption{
		//FilePrefix:       module,
		Category:         util.FILE_TYPE_DOC,
		UploadDir:        global.C.FileManager.UploadPath,
		UploadMaxSize:    global.C.FileManager.UploadDocVideoMaxSize,
		UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OPEN,
		//UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OFF,
		UploadStoreOSS:  util.UPLOAD_STORE_OSS_ALI,
		DownloadDir:     global.C.FileManager.DownloadPath,
		DownloadMaxSize: global.C.FileManager.DownloadMaxSize,
		FileHashType:    util.FILE_HASH_DAY,
		StaticDir:       global.C.Http.StaticPath,
		ProjectRootPath: global.MainEnv.RootDir,
		AliOss:          global.V.AliOss,
	}
	global.V.DocsManager = util.NewFileManagerUpload(fileUploadOption)

	fileUploadOption = util.FileManagerOption{
		//FilePrefix:       module,
		Category:         util.FILE_TYPE_VIDEO,
		UploadDir:        global.C.FileManager.UploadPath,
		UploadMaxSize:    global.C.FileManager.UploadDocVideoMaxSize,
		UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OPEN,
		//UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OFF,
		UploadStoreOSS:  util.UPLOAD_STORE_OSS_ALI,
		DownloadDir:     global.C.FileManager.DownloadPath,
		DownloadMaxSize: global.C.FileManager.DownloadMaxSize,
		FileHashType:    util.FILE_HASH_DAY,
		StaticDir:       global.C.Http.StaticPath,
		ProjectRootPath: global.MainEnv.RootDir,
		AliOss:          global.V.AliOss,
	}
	global.V.VideoManager = util.NewFileManagerUpload(fileUploadOption)

	fileUploadOption = util.FileManagerOption{
		//FilePrefix:       module,
		Category:         util.FILE_TYPE_PACKAGES,
		UploadDir:        global.C.FileManager.UploadPath,
		UploadMaxSize:    global.C.FileManager.UploadDocPackagesMaxSize,
		UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OPEN,
		//UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OFF,
		UploadStoreOSS:  util.UPLOAD_STORE_OSS_ALI,
		DownloadDir:     global.C.FileManager.DownloadPath,
		DownloadMaxSize: global.C.FileManager.DownloadMaxSize,
		FileHashType:    util.FILE_HASH_DAY,
		StaticDir:       global.C.Http.StaticPath,
		ProjectRootPath: global.MainEnv.RootDir,
		AliOss:          global.V.AliOss,
	}
	global.V.PackagesManager = util.NewFileManagerUpload(fileUploadOption)
}
