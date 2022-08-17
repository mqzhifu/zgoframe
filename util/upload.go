package util

import (
	"encoding/base64"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"io"
	"math"
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
)

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
type FileUpload struct {
	FileTypeMap sync.Map
	Option      FileUploadOption
}

type FileUploadOption struct {
	//Path 				string	//文件存储位置
	UploadDir       string //存储上传图片的目录名
	StaticDir       string //存储表态文件的目录名，它是UploadDir的上一级
	ProjectRootPath string //当前项目的绝对路径，它是StaticDir的上一级
	LocalDirPath    string //最终的：文件上传->本地硬盘路径
	Category        int    //上传的文件类型(扩展名分类):1全部2图片3文档,后端会根据类型做验证
	FileHashType    int    //文件存储时，添加前缀目录：hash类型
	MaxSize         int    //文件最大：MB ,默认：nginx是10Mb ,golang是9mb，不建议太大，且修改要与NGINX同步改，不然无效。文件太大建议使用新方法做分片传输
	FilePrefix      string //模块/业务名，可用于给文件名加前缀目录
	//阿里云-OSS相关
	OssAccessKeyId     string
	OssAccessKeySecret string
	OssEndpoint        string
	OssBucketName      string
	OssLocalDomain     string
}

var imgs = []string{"jpg", "jpeg", "png", "gif", "x-png", "png", "bmp", "pjpeg", "x-icon"}
var docs = []string{"txt", "doc", "docx", "dotx", "json", "cvs", "xls", "xlsx", "sql", "msword", "pptx", "pdf", "wps", "vsd"}
var video = []string{"mp3", "mp4", "avi", "rm", "mkv", "wmv", "mov", "flv", "rmvb"}

func NewFileUpload(Option FileUploadOption) *FileUpload {
	fileUpload := new(FileUpload)
	fileUpload.Option = Option

	//fileUpload.InitMap()
	return fileUpload
}
func (fileUpload *FileUpload) GetConstListFileUploadType() map[string]int {
	list := make(map[string]int)
	list["全部"] = FILE_TYPE_ALL
	list["图片"] = FILE_TYPE_IMG
	list["文档"] = FILE_TYPE_DOC
	list["视频"] = FILE_TYEP_VIDEO
	return list
}

//上传一个文件
func (fileUpload *FileUpload) UploadOne(header *multipart.FileHeader, syncOss int) (uploadRs UploadRs, err error) {
	//验证扩展名是否合法
	fileExtName, err := fileUpload.GetExtName(header.Filename)
	if err != nil {
		return uploadRs, err
	}
	//获取当前文件的大小
	fileSizeMB := Round(float64(header.Size)/1024/1024, 4)
	//MyPrint("fileSizeMB:",fileSizeMB)
	if fileUpload.Option.MaxSize > 0 && fileSizeMB > float64(fileUpload.Option.MaxSize) {
		return uploadRs, errors.New("大于限制：" + strconv.Itoa(fileUpload.Option.MaxSize) + " m")
	}

	MyPrint("UploadOne fileExtName:", fileExtName, " header size bytes:", header.Size, " mb:", fileSizeMB)
	////再次检查文件的类型是否正确
	//err = fileUpload.checkFileContentType(header,fileExtName)
	//if err != nil{
	//	return "",err
	//}
	//获取文件存储的绝对路径
	localDiskDir, relativePath, err := fileUpload.checkLocalDiskPath()
	MyPrint("localDiskDir:", localDiskDir)
	if err != nil {
		return uploadRs, err
	}
	//文件名：文件类型_NowUnixStamp_文件扩展名
	fileName := fileUpload.GetNewFileName(fileExtName)
	//hashDir := fileUpload.GetHashDirName()//获取相对路径，只是return时有用
	//if fileUpload.Option.FilePrefix != ""{
	//	hashDir =   fileUpload.Option.FilePrefix + "/" + hashDir
	//}

	newFileName := localDiskDir + "/" + fileName
	MyPrint("uploadOne file:", newFileName)
	//把用户上传的文件(内存中)，转移到本机的硬盘上
	out, err := os.Create(newFileName)
	if err != nil {
		return uploadRs, err
	}
	defer out.Close()
	file, err := header.Open()
	_, err = io.Copy(out, file)
	if err != nil {
		return uploadRs, err
	}
	//同步到阿里云
	if syncOss == 1 {
		err = fileUpload.UploadAliOSS(newFileName, relativePath, fileName)
		if err != nil {
			return uploadRs, err
		}
	}

	uploadRs.RelativePath = relativePath
	uploadRs.StaticDir = fileUpload.Option.StaticDir
	uploadRs.UploadDir = fileUpload.Option.UploadDir
	uploadRs.Filename = fileName
	uploadRs.LocalDiskPath = localDiskDir
	uploadRs.LocalIpUrl = fileUpload.GetLocalIpUrl(uploadRs)
	uploadRs.LocalDomainUrl = fileUpload.GetLocalDomainUrl(uploadRs)
	uploadRs.OssUrl = fileUpload.GetOssUrl(uploadRs)
	uploadRs.OssLocalUrl = fileUpload.GetOssLocalUrl(uploadRs)

	return uploadRs, nil
}
func (fileUpload *FileUpload) GetNewFileName(fileExtName string) string {
	return strconv.Itoa(fileUpload.Option.Category) + "_" + strconv.Itoa(GetNowTimeSecondToInt()) + "." + fileExtName
}

//撮当前上传目录的：hash前缀目录
func (fileUpload *FileUpload) GetHashDirName() string {
	dirName := ""
	switch fileUpload.Option.FileHashType {
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

func (fileUpload *FileUpload) OssLs() error {
	AccessKeyId := fileUpload.Option.OssAccessKeyId
	AccessKeySecret := fileUpload.Option.OssAccessKeySecret
	endpoint := fileUpload.Option.OssEndpoint

	client, err := oss.New(endpoint, AccessKeyId, AccessKeySecret)
	//MyPrint("oss New:",client,err)
	if err != nil {
		return err
	}

	bucketName := fileUpload.Option.OssBucketName

	MyPrint("oss endpoint:", endpoint, " AccessKeyId:", AccessKeyId, " AccessKeySecret:", AccessKeySecret, " bucketName:", bucketName)

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	listObjectsResult, err := bucket.ListObjects(oss.Prefix("agoraRecord/ckck/1660644031/"))
	//MyPrint("ListObjectsResult:", listObjectsResult, " err:", err)
	for k, v := range listObjectsResult.Objects {
		MyPrint(k, v)
	}
	ExitPrint(33)
	return nil
}

//将本地文件上传到阿里云-OSS
func (fileUpload *FileUpload) UploadAliOSS(localFilePath string, relativePath, FileName string) error {
	//这里阿里云有个小BUG，所有的路径不能以反斜杠(/)开头
	if relativePath[0:1] == "/" {
		relativePath = relativePath[1:]
	}
	AccessKeyId := fileUpload.Option.OssAccessKeyId
	AccessKeySecret := fileUpload.Option.OssAccessKeySecret
	endpoint := fileUpload.Option.OssEndpoint

	client, err := oss.New(endpoint, AccessKeyId, AccessKeySecret)
	//MyPrint("oss New:",client,err)
	if err != nil {
		return err
	}

	relativePathFile := relativePath + "/" + FileName

	bucketName := fileUpload.Option.OssBucketName

	MyPrint("oss endpoint:", endpoint, " AccessKeyId:", AccessKeyId, " AccessKeySecret:", AccessKeySecret, " bucketName:", bucketName)

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}
	//MyPrint("bucket:",bucket,err)
	MyPrint("oss localFilePath:", localFilePath, " relativePathFile:", relativePathFile)
	err = bucket.PutObjectFromFile(relativePathFile, localFilePath)
	MyPrint("PutObjectFromFile:", err)
	return err

}

//根据文件名(字符串)，取文件的扩展名，同时验证该扩展名是否合法
func (fileUpload *FileUpload) GetExtName(fileName string) (extName string, err error) {
	uploadFileName := strings.TrimSpace(fileName)
	if uploadFileName == "" {
		return "", errors.New("header.Filename  is empty ")
	}

	filenameArr := strings.Split(uploadFileName, ".")
	if len(filenameArr) < 2 {
		return "", errors.New("文件名中未包含: . ")
	}
	//去除扩展名的首尾空格，再全转化成小写
	fileExtName := strings.ToLower(strings.TrimSpace(filenameArr[len(filenameArr)-1]))
	fileExtNameFilter := fileUpload.FilterByExtString(fileUpload.Option.Category, fileExtName)
	if !fileExtNameFilter {
		return "", errors.New("文件扩展名非法")
	}

	return fileExtName, nil
}

//检查本地硬盘文件存储路径
func (fileUpload *FileUpload) checkLocalDiskPath() (localDiskDir string, relativePath string, err error) {
	//硬盘上存储的目录
	localDiskDir = fileUpload.Option.ProjectRootPath + "/" + fileUpload.Option.StaticDir + "/" + fileUpload.Option.UploadDir
	if fileUpload.Option.FilePrefix != "" {
		localDiskDir += "/" + fileUpload.Option.FilePrefix
		relativePath += fileUpload.Option.FilePrefix
	}
	localDiskDir += "/" + fileUpload.GetHashDirName()
	relativePath += fileUpload.GetHashDirName()
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

type FileDownInfo struct {
	PieceNum         int
	PieceSize        int
	FileSize         int64
	FileRelativePath string
	FileLocalPath    string
}

func (fileUpload *FileUpload) DownloadFileInfo(fileRelativePath string) (fileDownInfo FileDownInfo, err error) {
	if fileRelativePath == "" {
		return fileDownInfo, errors.New("fileRelativePath is empty")
	}
	//分成多少片
	//pieceNum := 10
	//触发:分片阀值(文件大小)
	maxMb := fileUpload.Option.MaxSize
	fileSizeSizeLimit := maxMb * 1024 * 1024 //MB 转 bytes
	//获取本地图片存储路径
	localDiskDir, _, _ := fileUpload.checkLocalDiskPath()
	fileDiskDir := localDiskDir + "/" + fileRelativePath
	MyPrint("fileDiskDir:", fileDiskDir)
	//判断文件是否存在
	fileInfo, err := FileExist(fileDiskDir)
	if err != nil {
		return fileDownInfo, err
	}
	//c.String(200,"11")
	//return nil
	//c.Header("Transfer-Encoding", "chunked")
	//c.Header("Content-Type", "image/jpeg")
	pieceNum := 0
	MyPrint("fileInfo.Size:", fileInfo.Size(), " fileSizeSizeLimit:", fileSizeSizeLimit)
	if fileInfo.Size() < int64(fileSizeSizeLimit) {
		return fileDownInfo, errors.New("小于" + strconv.Itoa(maxMb) + " mb ，请走正常接口即可")
	} else if fileInfo.Size() < 100*1024*1024 { //100MB 以内
		pieceNum = 10
	} else { //1Gb 以上的，就直接切分成100份了
		pieceNum = 100
	}
	//每片大小
	perPieceSize := int(math.Ceil(float64(fileInfo.Size() / int64(pieceNum))))

	fileDownInfo.FileSize = fileInfo.Size()
	fileDownInfo.PieceNum = pieceNum
	fileDownInfo.PieceSize = perPieceSize
	fileDownInfo.FileLocalPath = localDiskDir

	return fileDownInfo, nil
}

//文件下载
//正常普通的小文件，直接走nginx+static+cdn 即可，能调用此方法的肯定是文件过大的，得分片处理
func (fileUpload *FileUpload) Download(fileRelativePath string, c *gin.Context) error {
	fileDownInfo, err := fileUpload.DownloadFileInfo(fileRelativePath)
	if err != nil {
		return err
	}
	fileDiskDir := fileDownInfo.FileLocalPath + "/" + fileRelativePath
	fd, err := os.OpenFile(fileDiskDir, os.O_RDONLY, 6)
	if err != nil { //&& err != io.EOF
		return errors.New("file open err:" + err.Error())
	}
	//设置响应头信息
	c.Header("Content-Rang", "bytes")
	fd.Close()

	//headerRange := c.Header("Range")

	//MyPrint("OpenFile enter for earch:")
	//buffer := make([]byte, perPieceSize)
	//for{
	//	readDataLen, err := fd.Read(buffer)// len：读取文件中的数据长度
	//	if err == io.EOF{
	//		go c.String(200,"")
	//		MyPrint("in eof.")
	//		break
	//	}
	//	if err != nil {
	//		MyPrint("err not nil")
	//		MyPrint(err)
	//		break
	//	}
	//	MyPrint("once readDataLen:",readDataLen)
	//	if readDataLen == perPieceSize{
	//		MyPrint("nomal read and push http")
	//		c.String(200,string(buffer))
	//	}else{
	//		MyPrint("last read readDataLen:",readDataLen)
	//		s := buffer[0:readDataLen]
	//		go c.String(200,string(s))
	//	}
	//}
	//
	return nil

}

//流的大小：不能小于100个字节，因为要截取出头部的100个字节，做类型匹配及查错
func (fileUpload *FileUpload) UploadOneByStream(stream string, category int, syncOss int) (uploadRs UploadRs, err error) {
	if category != FILE_TYPE_IMG {
		return uploadRs, errors.New("目前category仅支持：图片流")
	}

	base64TypePrefix := "data:"
	base64TypeImgPrefix := "image/"
	base64TypeImg := imgs

	imgType := ""
	imgTypePrefixStr := ""

	if len(stream) < 100 {
		return uploadRs, errors.New("stream size < 100 bytes")
	}

	for _, v := range base64TypeImg {
		typeStr := base64TypePrefix + base64TypeImgPrefix + v
		//取取stream前100个字节，进行匹配
		if strings.Contains(stream[0:100], typeStr) {
			imgType = v
			imgTypePrefixStr = typeStr
			break
		}
	}

	if imgType == "" {
		return uploadRs, errors.New("no match , img type err.")
	}

	if imgType == "jpeg" {
		imgType = "jpg"
	}

	MyPrint("UploadOneByStream imgTypePrefixStr:", imgTypePrefixStr, " len:", len(imgTypePrefixStr))
	streamData := stream[len(imgTypePrefixStr)+8:]
	MyPrint("streamData:", streamData)
	data, err := base64.StdEncoding.DecodeString(streamData)
	if err != nil {
		return uploadRs, errors.New("base64 DecodeString err:" + err.Error())
	}

	localDiskDir, relativePath, err := fileUpload.checkLocalDiskPath()
	if err != nil {
		return uploadRs, err
	}
	//文件名：文件类型_NowUnixStamp_文件扩展名
	fileName := fileUpload.GetNewFileName(imgType)
	localDiskDirFile := localDiskDir + "/" + fileName
	MyPrint("localDiskDirFile:", localDiskDirFile)
	fd, err := os.OpenFile(localDiskDirFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return uploadRs, errors.New("open file err:" + err.Error())
	}
	defer fd.Close()
	_, err = fd.Write(data)
	if err != nil {
		return uploadRs, errors.New("file write err:" + err.Error())
	}

	if syncOss == 1 {
		//err = fileUpload.UploadAliOSS(newFileName,relativePath,fileName)
		//if err != nil{
		//	return uploadRs,err
		//}
	}

	uploadRs.RelativePath = relativePath
	uploadRs.StaticDir = fileUpload.Option.StaticDir
	uploadRs.UploadDir = fileUpload.Option.UploadDir
	uploadRs.Filename = fileName
	uploadRs.LocalDiskPath = localDiskDir
	uploadRs.LocalIpUrl = fileUpload.GetLocalIpUrl(uploadRs)
	uploadRs.LocalDomainUrl = fileUpload.GetLocalDomainUrl(uploadRs)
	uploadRs.OssUrl = fileUpload.GetOssUrl(uploadRs)
	uploadRs.OssLocalUrl = fileUpload.GetOssLocalUrl(uploadRs)

	return uploadRs, nil
}

func (fileUpload *FileUpload) GetOssLocalUrl(uploadRs UploadRs) string {
	return fileUpload.Option.OssLocalDomain + "/" + uploadRs.RelativePath + "/" + uploadRs.Filename
}

func (fileUpload *FileUpload) GetOssUrl(uploadRs UploadRs) string {
	return fileUpload.Option.OssBucketName + "." + fileUpload.Option.OssEndpoint + "/" + uploadRs.RelativePath + "/" + uploadRs.Filename
}

//ip访问的话，目录前面会多一个 static
func (fileUpload *FileUpload) GetLocalIpUrl(uploadRs UploadRs) string {
	return uploadRs.StaticDir + "/" + uploadRs.UploadDir + "/" + uploadRs.RelativePath + "/" + uploadRs.Filename
}

//域名访问的话，少一个static
func (fileUpload *FileUpload) GetLocalDomainUrl(uploadRs UploadRs) string {
	return uploadRs.UploadDir + "/" + uploadRs.RelativePath + "/" + uploadRs.Filename
}

func (fileUpload *FileUpload) GetAllowFileTypeList(category int) (rs []string, err error) {

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
func (fileTypeFilter *FileUpload) GetAllowFileTypeListToStr(category int) string {
	listStr := ""
	list, _ := fileTypeFilter.GetAllowFileTypeList(category)
	for _, v := range list {
		listStr += v + " "
	}
	return listStr
}

func (fileUpload *FileUpload) FilterByExtString(category int, extName string) bool {
	list, _ := fileUpload.GetAllowFileTypeList(category)
	for _, v := range list {
		if v == extName {
			return true
		}
	}
	return false
}
