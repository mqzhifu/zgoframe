package util

import (
	"encoding/base64"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	FILE_TYPE_ALL   = 1
	FILE_TYPE_IMG   = 2
	FILE_TYPE_DOC   = 3
	FILE_TYEP_VIDEO = 4

	UPLOAD_STORE_LOCAL_OFF  = 1
	UPLOAD_STORE_LOCAL_OPEN = 2

	UPLOAD_STORE_OSS_OFF = 0
	UPLOAD_STORE_OSS_ALI = 1
)

type FileDownInfo struct {
	PieceNum         int
	PieceSize        int
	FileSize         int64
	FileRelativePath string
	FileLocalPath    string
}

//上伟文件成功后，返回的数据
type UploadRs struct {
	Filename       string `json:"filename"`         //文件名
	RelativePath   string `json:"relative_path"`    //相对路径：用户自定义的前缀目录 + hash目录名
	UploadDir      string `json:"upload_dir"`       //存储上传图片的目录名
	StaticDir      string `json:"static_dir"`       //存储表态文件的目录名，它是UploadDir的上一级
	LocalDiskPath  string `json:"local_disk_path"`  //本地硬盘存储的路径
	LocalIpUrl     string `json:"local_ip_url"`     //访问本地文件IP-URL地址
	LocalDomainUrl string `json:"local_domain_url"` //访问本地文件DOMAIN-URL地址
	OssLocalUrl    string `json:"oss_local_url"`    //自己的域名绑定在阿里OSS上
	OssUrl         string `json:"oss_url"`          //阿里云的地址
}

//类
type FileManager struct {
	FileTypeMap sync.Map
	Option      FileManagerOption
}

type FileManagerOption struct {
	UploadDir        string //存储上传图片的目录名
	UploadMaxSize    int    //文件最大：MB ,默认：nginx是10Mb ,golang是9mb，不建议太大，且修改要与NGINX同步改，不然无效。文件太大建议使用新方法做分片传输
	UploadStoreLocal int    //上传的文件，是否存储本地
	UploadStoreOSS   int    //上传的文件，是否存储到3方OSS网盘

	DownloadDir     string
	DownloadMaxSize int

	StaticDir       string //存储表态文件的目录名，它是UploadDir的上一级
	ProjectRootPath string //当前项目的绝对路径，它是StaticDir的上一级
	LocalDirPath    string //最终的：文件上传->本地硬盘路径
	Category        int    //上传的文件类型(扩展名分类):1全部2图片3文档,后端会根据类型做验证
	FileHashType    int    //文件存储时，添加前缀目录：hash类型
	FilePrefix      string //模块/业务名，可用于给文件名加前缀目录
	AliOss          *AliOss
	//Path 				string	//文件存储位置
	//阿里云-OSS相关
	//OssAccessKeyId     string
	//OssAccessKeySecret string
	//OssEndpoint        string
	//OssBucketName      string
	//OssLocalDomain     string
}

var imgs = []string{"jpg", "jpeg", "png", "gif", "x-png", "png", "bmp", "pjpeg", "x-icon"}
var docs = []string{"txt", "doc", "docx", "dotx", "json", "cvs", "xls", "xlsx", "sql", "msword", "pptx", "pdf", "wps", "vsd"}
var video = []string{"mp3", "mp4", "avi", "rm", "mkv", "wmv", "mov", "flv", "rmvb"}

func NewFileManagerUpload(Option FileManagerOption) *FileManager {
	fileManager := new(FileManager)
	fileManager.Option = Option

	//fileUpload.InitMap()
	return fileManager
}
func (fileManager *FileManager) GetConstListFileUploadType() map[string]int {
	list := make(map[string]int)
	list["全部"] = FILE_TYPE_ALL
	list["图片"] = FILE_TYPE_IMG
	list["文档"] = FILE_TYPE_DOC
	list["视频"] = FILE_TYEP_VIDEO
	return list
}
func (fileManager *FileManager) GetConstListFileUploadStoreLocal() map[string]int {
	list := make(map[string]int)
	list["关闭"] = UPLOAD_STORE_LOCAL_OFF
	list["打开"] = UPLOAD_STORE_LOCAL_OPEN
	return list
}

func (fileManager *FileManager) GetConstListFileUploadStoreOSS() map[string]int {
	list := make(map[string]int)
	list["关闭"] = UPLOAD_STORE_OSS_OFF
	list["阿里"] = UPLOAD_STORE_OSS_ALI
	return list
}

//上传一个文件
func (fileManager *FileManager) UploadOne(header *multipart.FileHeader) (uploadRs UploadRs, err error) {
	//验证扩展名是否合法
	fileExtName, err := fileManager.GetExtName(header.Filename)
	if err != nil {
		return uploadRs, err
	}
	//获取当前文件的大小
	fileSizeMB := Round(float64(header.Size)/1024/1024, 4)
	//MyPrint("fileSizeMB:",fileSizeMB)
	if fileManager.Option.UploadMaxSize > 0 && fileSizeMB > float64(fileManager.Option.UploadMaxSize) {
		return uploadRs, errors.New("大于限制：" + strconv.Itoa(fileManager.Option.UploadMaxSize) + "(mb)")
	}

	MyPrint("UploadOne fileExtName:", fileExtName, " header size bytes:", header.Size, " mb:", fileSizeMB)
	//获取文件存储的绝对路径
	localDiskDir, relativePath, err := fileManager.checkLocalDiskPath()
	MyPrint("localDiskDir:", localDiskDir)
	if err != nil {
		return uploadRs, err
	}
	//文件名：文件类型_NowUnixStamp_文件扩展名
	fileName := fileManager.GetNewFileName(fileExtName)
	newFileName := localDiskDir + "/" + fileName
	MyPrint("uploadOne file:", newFileName)
	if fileManager.Option.UploadStoreLocal == UPLOAD_STORE_LOCAL_OPEN {
		//把用户上传的文件(内存中)，转移到本机的硬盘上
		out, err := os.Create(newFileName)
		defer out.Close()
		if err != nil {
			return uploadRs, errors.New("本地存储文件失败1:" + err.Error())
		}
		file, err := header.Open()
		_, err = io.Copy(out, file)
		if err != nil {
			return uploadRs, errors.New("本地存储文件失败2:" + err.Error())
		}
	}

	//同步到阿里云
	if fileManager.Option.UploadStoreOSS == UPLOAD_STORE_OSS_ALI {
		if fileManager.Option.UploadStoreLocal == UPLOAD_STORE_LOCAL_OPEN {
			//如果本地存储打开了，流里的数据已经读完了，不能重复读，那就用本地已保存的文件传到OSS上
			err = fileManager.Option.AliOss.UploadOneByFile(newFileName, relativePath, fileName)
		} else {
			fileStream, err := header.Open()
			if err != nil {
				return uploadRs, errors.New(" header.Open()读取失败:" + err.Error())
			}

			err = fileManager.Option.AliOss.UploadOneByStream(fileStream, relativePath, fileName)
			if err != nil {
				return uploadRs, errors.New("上传阿里云OSS失败:" + err.Error())
			}
		}
	}

	uploadRs.RelativePath = relativePath
	uploadRs.StaticDir = fileManager.Option.StaticDir
	uploadRs.UploadDir = fileManager.Option.UploadDir
	uploadRs.Filename = fileName
	uploadRs.LocalDiskPath = localDiskDir
	uploadRs.LocalIpUrl = fileManager.GetLocalIpUrl(uploadRs)
	uploadRs.LocalDomainUrl = fileManager.GetLocalDomainUrl(uploadRs)
	uploadRs.OssUrl = fileManager.GetOssUrl(uploadRs)
	uploadRs.OssLocalUrl = fileManager.GetOssLocalUrl(uploadRs)

	return uploadRs, nil
}
func (fileManager *FileManager) RealUploadOne() {

}

//流的大小：不能小于100个字节，因为要截取出头部的100个字节，做类型匹配及校验
func (fileManager *FileManager) UploadOneByStream(stream string, category int) (uploadRs UploadRs, err error) {
	if category != FILE_TYPE_IMG {
		return uploadRs, errors.New("目前category仅支持：图片流")
	}

	base64StreamPrefix := "data:"
	base64StreamImgPrefix := "image/"
	imgAllowExtType := imgs

	imgExtType := ""
	imgTypePrefixStr := ""

	if len(stream) < 100 {
		return uploadRs, errors.New("stream size < 100 bytes")
	}

	for _, v := range imgAllowExtType {
		typeStr := base64StreamPrefix + base64StreamImgPrefix + v
		//取取stream前100个字节，进行匹配
		if strings.Contains(stream[0:100], typeStr) {
			imgExtType = v
			imgTypePrefixStr = typeStr
			break
		}
	}

	if imgExtType == "" {
		return uploadRs, errors.New("no match , img type err.")
	}

	if imgExtType == "jpeg" {
		imgExtType = "jpg"
	}

	MyPrint("UploadOneByStream imgTypePrefixStr:", imgTypePrefixStr, " len:", len(imgTypePrefixStr))
	streamData := stream[len(imgTypePrefixStr)+8:]
	MyPrint("streamData:", streamData)
	data, err := base64.StdEncoding.DecodeString(streamData)
	if err != nil {
		return uploadRs, errors.New("base64 DecodeString err:" + err.Error())
	}
	//====

	//获取当前文件的大小
	size := len(streamData) / 1024 / 1024
	fileSizeMB := Round(float64(size), 4)
	//MyPrint("fileSizeMB:",fileSizeMB)
	if fileManager.Option.UploadMaxSize > 0 && fileSizeMB > float64(fileManager.Option.UploadMaxSize) {
		return uploadRs, errors.New("大于限制：" + strconv.Itoa(fileManager.Option.UploadMaxSize) + "(mb)")
	}
	//获取文件存储的绝对路径
	localDiskDir, relativePath, err := fileManager.checkLocalDiskPath()
	MyPrint("localDiskDir:", localDiskDir)
	if err != nil {
		return uploadRs, err
	}
	//文件名：文件类型_NowUnixStamp_文件扩展名
	fileName := fileManager.GetNewFileName(imgExtType)
	newFileName := localDiskDir + "/" + fileName
	MyPrint("uploadOne file:", newFileName)
	if fileManager.Option.UploadStoreLocal == UPLOAD_STORE_LOCAL_OPEN {
		//把用户上传的文件(内存中)，转移到本机的硬盘上
		//out, err := os.Create(newFileName)
		//defer out.Close()
		//if err != nil {
		//	return uploadRs, errors.New("本地存储文件失败1:" + err.Error())
		//}
		//file, err := header.Open()
		//_, err = io.Copy(out, file)
		//if err != nil {
		//	return uploadRs, errors.New("本地存储文件失败2:" + err.Error())
		//}

		fd, err := os.OpenFile(newFileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			return uploadRs, errors.New("open file err:" + err.Error())
		}
		defer fd.Close()
		_, err = fd.Write(data)
		if err != nil {
			return uploadRs, errors.New("file write err:" + err.Error())
		}
	}

	//同步到阿里云
	if fileManager.Option.UploadStoreOSS == UPLOAD_STORE_OSS_ALI {
		rr := strings.NewReader(streamData)
		err = fileManager.Option.AliOss.UploadOneByStream(rr, relativePath, fileName)
		if err != nil {
			return uploadRs, errors.New("上传阿里云OSS失败:" + err.Error())
		}
	}

	uploadRs.RelativePath = relativePath
	uploadRs.StaticDir = fileManager.Option.StaticDir
	uploadRs.UploadDir = fileManager.Option.UploadDir
	uploadRs.Filename = fileName
	uploadRs.LocalDiskPath = localDiskDir
	uploadRs.LocalIpUrl = fileManager.GetLocalIpUrl(uploadRs)
	uploadRs.LocalDomainUrl = fileManager.GetLocalDomainUrl(uploadRs)
	uploadRs.OssUrl = fileManager.GetOssUrl(uploadRs)
	uploadRs.OssLocalUrl = fileManager.GetOssLocalUrl(uploadRs)

	return uploadRs, nil
}

//func (fileManager *FileManager) DownloadFileInfo(fileRelativePath string) (fileDownInfo FileDownInfo, err error) {
//	if fileRelativePath == "" {
//		return fileDownInfo, errors.New("fileRelativePath is empty")
//	}
//	//分成多少片
//	//pieceNum := 10
//	//触发:分片阀值(文件大小)
//	maxMb := fileUpload.Option.MaxSize
//	fileSizeSizeLimit := maxMb * 1024 * 1024 //MB 转 bytes
//	//获取本地图片存储路径
//	localDiskDir, _, _ := fileUpload.checkLocalDiskPath()
//	fileDiskDir := localDiskDir + "/" + fileRelativePath
//	MyPrint("fileDiskDir:", fileDiskDir)
//	//判断文件是否存在
//	fileInfo, err := FileExist(fileDiskDir)
//	if err != nil {
//		return fileDownInfo, err
//	}
//	//c.String(200,"11")
//	//return nil
//	//c.Header("Transfer-Encoding", "chunked")
//	//c.Header("Content-Type", "image/jpeg")
//	pieceNum := 0
//	MyPrint("fileInfo.Size:", fileInfo.Size(), " fileSizeSizeLimit:", fileSizeSizeLimit)
//	if fileInfo.Size() < int64(fileSizeSizeLimit) {
//		return fileDownInfo, errors.New("小于" + strconv.Itoa(maxMb) + " mb ，请走正常接口即可")
//	} else if fileInfo.Size() < 100*1024*1024 { //100MB 以内
//		pieceNum = 10
//	} else { //1Gb 以上的，就直接切分成100份了
//		pieceNum = 100
//	}
//	//每片大小
//	perPieceSize := int(math.Ceil(float64(fileInfo.Size() / int64(pieceNum))))
//
//	fileDownInfo.FileSize = fileInfo.Size()
//	fileDownInfo.PieceNum = pieceNum
//	fileDownInfo.PieceSize = perPieceSize
//	fileDownInfo.FileLocalPath = localDiskDir
//
//	return fileDownInfo, nil
//}
////文件下载
////正常普通的小文件，直接走nginx+static+cdn 即可，能调用此方法的肯定是文件过大的，得分片处理
//func (fileManager *FileManager) Download(fileRelativePath string, c *gin.Context) error {
//	fileDownInfo, err := fileUpload.DownloadFileInfo(fileRelativePath)
//	if err != nil {
//		return err
//	}
//	fileDiskDir := fileDownInfo.FileLocalPath + "/" + fileRelativePath
//	fd, err := os.OpenFile(fileDiskDir, os.O_RDONLY, 6)
//	if err != nil { //&& err != io.EOF
//		return errors.New("file open err:" + err.Error())
//	}
//	//设置响应头信息
//	c.Header("Content-Rang", "bytes")
//	fd.Close()
//
//	//headerRange := c.Header("Range")
//
//	//MyPrint("OpenFile enter for earch:")
//	//buffer := make([]byte, perPieceSize)
//	//for{
//	//	readDataLen, err := fd.Read(buffer)// len：读取文件中的数据长度
//	//	if err == io.EOF{
//	//		go c.String(200,"")
//	//		MyPrint("in eof.")
//	//		break
//	//	}
//	//	if err != nil {
//	//		MyPrint("err not nil")
//	//		MyPrint(err)
//	//		break
//	//	}
//	//	MyPrint("once readDataLen:",readDataLen)
//	//	if readDataLen == perPieceSize{
//	//		MyPrint("nomal read and push http")
//	//		c.String(200,string(buffer))
//	//	}else{
//	//		MyPrint("last read readDataLen:",readDataLen)
//	//		s := buffer[0:readDataLen]
//	//		go c.String(200,string(s))
//	//	}
//	//}
//	//
//	return nil
//
//}

func (fileManager *FileManager) GetLocalDiskBasePath() string {
	return fileManager.Option.ProjectRootPath + "/" + fileManager.Option.StaticDir + "/" + fileManager.Option.UploadDir
}

func (fileManager *FileManager) GetOssLocalUrl(uploadRs UploadRs) string {
	return fileManager.Option.AliOss.Op.LocalDomain + "/" + uploadRs.RelativePath + "/" + uploadRs.Filename
}

func (fileManager *FileManager) GetOssUrl(uploadRs UploadRs) string {
	return fileManager.Option.AliOss.Op.BucketName + "." + fileManager.Option.AliOss.Op.Endpoint + "/" + uploadRs.RelativePath + "/" + uploadRs.Filename
}

//ip访问的话，目录前面会多一个 static
func (fileManager *FileManager) GetLocalIpUrl(uploadRs UploadRs) string {
	return uploadRs.StaticDir + "/" + uploadRs.UploadDir + "/" + uploadRs.RelativePath + "/" + uploadRs.Filename
}

//域名访问的话，少一个static
func (fileManager *FileManager) GetLocalDomainUrl(uploadRs UploadRs) string {
	return uploadRs.UploadDir + "/" + uploadRs.RelativePath + "/" + uploadRs.Filename
}

func (fileManager *FileManager) GetAllowFileTypeList(category int) (rs []string, err error) {

	if category == FILE_TYPE_IMG {
		return imgs, nil
	} else if category == FILE_TYPE_DOC {
		return docs, nil
	} else if category == FILE_TYEP_VIDEO {
		return video, nil
	} else if category == FILE_TYPE_ALL {
		all := append(imgs, docs...)
		all = append(all, video...)
		return all, nil
	} else {
		return nil, errors.New("category err.")
	}

}

//主要是给出错信息使用
func (fileManager *FileManager) GetAllowFileTypeListToStr(category int) string {
	listStr := ""
	list, _ := fileManager.GetAllowFileTypeList(category)
	for _, v := range list {
		listStr += v + " "
	}
	return listStr
}

func (fileManager *FileManager) FilterByExtString(category int, extName string) bool {
	list, _ := fileManager.GetAllowFileTypeList(category)
	for _, v := range list {
		if v == extName {
			return true
		}
	}
	return false
}

//把用户上传的文件名，转换成自己想要的文件名：类型ID_当时时间.扩展名
func (fileManager *FileManager) GetNewFileName(fileExtName string) string {
	return strconv.Itoa(fileManager.Option.Category) + "_" + strconv.Itoa(GetNowTimeSecondToInt()) + "." + fileExtName
}

//撮当前上传目录的：hash前缀目录
func (fileManager *FileManager) GetHashDirName() string {
	dirName := ""
	switch fileManager.Option.FileHashType {
	case FILE_HASH_NONE:
		break
	case FILE_HASH_HOUR:
		dirName = GetNowDateHour()
	case FILE_HASH_DAY:
		dirName = GetNowDate()
	case FILE_HASH_MONTH:
		dirName = GetNowDateMonth()
	}

	return dirName
}

//根据文件名(字符串)，取文件的扩展名，同时验证该扩展名是否合法
func (fileManager *FileManager) GetExtName(fileName string) (extName string, err error) {
	if !CheckFileName(fileName) {
		return "", errors.New("文件名不合法：只允许大小写字母+(-_),且必须且只能出现一个:符号(.),最小3，最长111 ")
	}
	//根据.切割文件名字符串
	extName = strings.Split(fileName, ".")[1]
	//判断下：扩展名类型是否合合法
	fileExtNameFilter := fileManager.FilterByExtString(fileManager.Option.Category, extName)
	if !fileExtNameFilter {
		return "", errors.New("文件扩展名非法")
	}

	return extName, nil
}

//检查本地硬盘文件存储路径
func (fileManager *FileManager) checkLocalDiskPath() (localDiskDir string, relativePath string, err error) {
	//硬盘上存储的目录
	localDiskDir = fileManager.GetLocalDiskBasePath()
	if fileManager.Option.FilePrefix != "" {
		localDiskDir += "/" + fileManager.Option.FilePrefix
		relativePath += fileManager.Option.FilePrefix
	}
	localDiskDir += "/" + fileManager.GetHashDirName()
	relativePath += fileManager.GetHashDirName()
	_, err = PathExists(localDiskDir)
	if err != nil {
		if os.IsNotExist(err) {
			MyPrint("dir not exist ,mkdir:" + localDiskDir)
			err = os.MkdirAll(localDiskDir, 0666)
			if err != nil {
				return "", "", errors.New("mkdir err:" + err.Error())
			}
		} else {
			return "", "", err
		}
	} else {
		MyPrint("baseDir exist:" + localDiskDir)
	}

	return localDiskDir, relativePath, nil
}
