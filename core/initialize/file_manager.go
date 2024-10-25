package initialize

import (
	"zgoframe/core/global"
	"zgoframe/util"
)

// 文件管理，上传图片 - 公共类
// 公共物理路径，统一放在 static 目录下
func InitFileManager() {
	baseDir := global.MainEnv.RootDir + "/" + global.C.Http.StaticPath + "/puzzle"
	oriImagePath := baseDir + "/ori/"
	shardImagePath := baseDir + "/shard/"
	thumbnailPath := baseDir + "/thumbnail/"
	global.V.Util.ImageSlice = util.NewImageSlice(oriImagePath, shardImagePath, thumbnailPath)

	//文件类型：图片
	fileUploadOption := util.FileManagerOption{
		//FilePrefix:       module,
		//UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OFF,
		Category:         util.FILE_TYPE_IMG,
		UploadDir:        global.C.FileManager.UploadPath,
		UploadMaxSize:    global.C.FileManager.UploadDocImgMaxSize,
		UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OPEN,
		UploadStoreOSS:   util.UPLOAD_STORE_OSS_ALI,
		DownloadDir:      global.C.FileManager.DownloadPath,
		DownloadMaxSize:  global.C.FileManager.DownloadMaxSize,
		FileHashType:     util.FILE_HASH_DAY,
		StaticDir:        global.C.Http.StaticPath,
		ProjectRootPath:  global.MainEnv.RootDir,
		AliOss:           global.V.Util.AliOss,
	}
	global.V.Util.ImgManager = util.NewFileManagerUpload(fileUploadOption)

	//文件类型：文档
	fileUploadOption = util.FileManagerOption{
		//FilePrefix:       module,
		//UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OFF,
		Category:         util.FILE_TYPE_DOC,
		UploadDir:        global.C.FileManager.UploadPath,
		UploadMaxSize:    global.C.FileManager.UploadDocVideoMaxSize,
		UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OPEN,
		UploadStoreOSS:   util.UPLOAD_STORE_OSS_ALI,
		DownloadDir:      global.C.FileManager.DownloadPath,
		DownloadMaxSize:  global.C.FileManager.DownloadMaxSize,
		FileHashType:     util.FILE_HASH_DAY,
		StaticDir:        global.C.Http.StaticPath,
		ProjectRootPath:  global.MainEnv.RootDir,
		AliOss:           global.V.Util.AliOss,
	}
	global.V.Util.DocsManager = util.NewFileManagerUpload(fileUploadOption)

	//文件类型：视频
	fileUploadOption = util.FileManagerOption{
		//FilePrefix:       module,
		//UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OFF,
		Category:         util.FILE_TYPE_VIDEO,
		UploadDir:        global.C.FileManager.UploadPath,
		UploadMaxSize:    global.C.FileManager.UploadDocVideoMaxSize,
		UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OPEN,
		UploadStoreOSS:   util.UPLOAD_STORE_OSS_ALI,
		DownloadDir:      global.C.FileManager.DownloadPath,
		DownloadMaxSize:  global.C.FileManager.DownloadMaxSize,
		FileHashType:     util.FILE_HASH_DAY,
		StaticDir:        global.C.Http.StaticPath,
		ProjectRootPath:  global.MainEnv.RootDir,
		AliOss:           global.V.Util.AliOss,
	}
	global.V.Util.VideoManager = util.NewFileManagerUpload(fileUploadOption)

	//文件类型：安装包
	fileUploadOption = util.FileManagerOption{
		//FilePrefix:       module,
		//UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OFF,
		Category:         util.FILE_TYPE_PACKAGES,
		UploadDir:        global.C.FileManager.UploadPath,
		UploadMaxSize:    global.C.FileManager.UploadDocPackagesMaxSize,
		UploadStoreLocal: util.UPLOAD_STORE_LOCAL_OPEN,
		UploadStoreOSS:   util.UPLOAD_STORE_OSS_ALI,
		DownloadDir:      global.C.FileManager.DownloadPath,
		DownloadMaxSize:  global.C.FileManager.DownloadMaxSize,
		FileHashType:     util.FILE_HASH_DAY,
		StaticDir:        global.C.Http.StaticPath,
		ProjectRootPath:  global.MainEnv.RootDir,
		AliOss:           global.V.Util.AliOss,
	}
	global.V.Util.PackagesManager = util.NewFileManagerUpload(fileUploadOption)
}
