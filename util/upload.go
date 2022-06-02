package util

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	FILE_TYPE_ALL = 1
	FILE_TYPE_IMG = 2
	FILE_TYPE_DOC = 3
	FILE_TYEP_VIDEO = 4
)

type FileUpload struct {
	FileTypeMap sync.Map
	Option FileUploadOption
}


type FileUploadOption struct {
	Path 				string	//文件存储位置
	Category 			int		//扩展名分类,上传的文件类型，1全部2图片3文档,后端会根据类型做验证
	FileHashType		int		//存储目录hash类型
	OssAccessKeyId		string
	OssAccessKeySecret	string
	OssEndpoint			string
	OssBucketName		string
	MaxSize 			int 	//文件最大：MB
	FilePrefix 			string 	//模块/业务名，可用于给文件名加前缀目录
}

func NewFileUpload(Option FileUploadOption )*FileUpload{
	fileUpload := new(FileUpload)
	fileUpload.Option = Option

	//fileUpload.InitMap()
	return fileUpload
}

func   (fileUpload *FileUpload)  GetHashDirName( )string{
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

	return dirName;
}


func   (fileUpload *FileUpload)UploadAliOSS(localFilePath string , relativePath , FileName string)error{
	AccessKeyId := fileUpload.Option.OssAccessKeyId
	AccessKeySecret := fileUpload.Option.OssAccessKeySecret
	endpoint :=fileUpload.Option.OssEndpoint

	client ,err := oss.New(endpoint,AccessKeyId,AccessKeySecret)
	MyPrint("oss New:",client,err)
	if err != nil{
		return err
	}

	relativePathFile := relativePath + "/" + FileName

	bucketName := fileUpload.Option.OssBucketName
	bucket , err := client.Bucket(bucketName)
	if err != nil{
		return err
	}

	MyPrint("bucket:",bucket,err)
	err = bucket.PutObjectFromFile(relativePathFile,localFilePath)
	MyPrint("PutObjectFromFile:",err)
	return err

}
func   (fileUpload *FileUpload)GetExtName( fileName string)(extName string ,err error){
	uploadFileName := strings.TrimSpace(fileName)
	if uploadFileName == ""{
		return "",errors.New("header.Filename  is empty ")
	}

	filenameArr := strings.Split(uploadFileName, ".")
	if len(filenameArr) < 2{
		return "",errors.New("文件名中未包含: . ")
	}
	//去除扩展名的首尾空格，再全转化成小写
	fileExtName := strings.ToLower( strings.TrimSpace( filenameArr[len(filenameArr)-1] ))
	fileExtNameFilter := fileUpload.FilterByExtString(fileUpload.Option.Category,fileExtName)
	if !fileExtNameFilter{
		return "",errors.New("文件扩展名非法")
	}

	return fileExtName,nil
}

func   (fileUpload *FileUpload)checkFileContentType(header *multipart.FileHeader,fileExtName string)error{
	f,_ := header.Open()
	fSrc, _ := ioutil.ReadAll(f)
	realFileType := fileUpload.GetFileType(fSrc[:10])

	contentType := http.DetectContentType(fSrc[:512])
	MyPrint("checkFileContentType realFileType:",realFileType , " http.DetectContentType:",contentType)

	if realFileType != ""{
		if !fileUpload.FilterByExtString(fileUpload.Option.Category,realFileType){
			return errors.New("ext errors:"+fileUpload.GetAllowFileTypeListToStr(fileUpload.Option.Category))
		}

		if realFileType != fileExtName{
			return errors.New("文件名中的扩展名与文件内容的类型不符")
		}
	}else{
		//这里是证明无法从头里识别出具体类型，如:TXT 不同编辑类型，可能头内容不同，ANSI的更是没有头标识符
		if fileExtName != "txt"{
			return errors.New("文件类型非法:未识别出文件类型")
		}
	}

	return nil;
}
func   (fileUpload *FileUpload)checkLocalDiskDir()(localDiskDir string ,err error){
	//硬盘上存储的目录
	localDiskDir = fileUpload.Option.Path
	if fileUpload.Option.FilePrefix != ""{
		localDiskDir += "/" + fileUpload.Option.FilePrefix
	}
	localDiskDir += "/" + fileUpload.GetHashDirName()
	_,err = PathExists(localDiskDir)
	if err != nil {
		if os.IsNotExist(err){
			MyPrint("dir not exist ,mkdir:"+localDiskDir)
			err = os.MkdirAll(localDiskDir,0666)
			if err != nil{
				return "" , errors.New("mkdir err:"+err.Error())
			}
		}else{
			return "",err
		}
	}else{
		MyPrint("baseDir exist:"+localDiskDir)
	}
	return localDiskDir,nil
}
//func   (fileUpload *FileUpload)UploadOne(file multipart.File,header *multipart.FileHeader)(relativePathFileName string ,err error){
func   (fileUpload *FileUpload)UploadOne( header *multipart.FileHeader)(relativePathFileName string ,err error){
	fileExtName ,err := fileUpload.GetExtName(header.Filename)
	if err != nil{
		return "",err
	}
	fileSizeMB := Round(   float64 (header.Size ) / 1024 / 1024 ,4)
	MyPrint("fileSizeMB:",fileSizeMB)
	if fileUpload.Option.MaxSize > 0 && fileSizeMB > float64( fileUpload.Option.MaxSize){
		return  "" ,errors.New("大于限制："+strconv.Itoa(fileUpload.Option.MaxSize) + " m")
	}

	MyPrint("UploadOne fileExtName:",fileExtName  , " header size:",header.Size , " mb:",fileSizeMB)
	////再次检查文件的类型是否正确
	//err = fileUpload.checkFileContentType(header,fileExtName)
	//if err != nil{
	//	return "",err
	//}

	localDiskDir , err := fileUpload.checkLocalDiskDir()
	if err != nil{
		return "",err
	}
	//ExitPrint("localDiskDir:",localDiskDir)
	fileName := strconv.Itoa(fileUpload.Option.Category) + "_" + strconv.Itoa(GetNowTimeSecondToInt()) + "." +fileExtName
	relativePath := fileUpload.GetHashDirName()
	if fileUpload.Option.FilePrefix != ""{
		relativePath =   fileUpload.Option.FilePrefix + "/" + relativePath
	}

	relativePathFileName =  relativePath + "/" + fileName

	newFileName := localDiskDir + "/" + fileName
	MyPrint("uploadOne file:",newFileName)
	//把用户上传的文件(内存中)，转移到本机的硬盘上
	out, err := os.Create(newFileName)
	if err != nil {
		return "",err
	}
	defer out.Close()
	file, err := header.Open()
	_, err = io.Copy(out, file)
	if err != nil {
		return "",err
	}
	//同步到阿里云
	err = fileUpload.UploadAliOSS(newFileName,relativePath,fileName)
	if err != nil{
		return "",err
	}

	return relativePathFileName,nil
}


func   (fileUpload *FileUpload)GetAllowFileTypeList(category int)(rs []string,err error){
	imgs := []string{"jpg","jpeg","png","gif","x-png","png","bmp","pjpeg"}
	docs := []string{"txt","doc","docx","dotx","json","cvs","xls","xlsx","sql","msword","pptx","pdf","wps","vsd"}
	video := []string{"mp3","mp4","avi","rm","mkv","wmv","mov","flv","rmvb"}

	if category == FILE_TYPE_IMG{
		return imgs,nil
	}else if category == FILE_TYPE_DOC{
		return docs,nil
	}else if category == FILE_TYEP_VIDEO{
		return video,nil
	}else if category == FILE_TYPE_ALL{
		all := append(imgs, docs...)
		all = append(all, video...)
		return all,nil
	}else {
		return nil,errors.New("category err.")
	}

}
//主要是给出错信息使用
func   (fileTypeFilter *FileUpload)GetAllowFileTypeListToStr(category int)string{
	listStr := ""
	list ,_ := fileTypeFilter.GetAllowFileTypeList(category)
	for _,v:= range list{
		listStr += v + " "
	}
	return listStr
}

func (fileUpload *FileUpload)FilterByExtString(category int,extName string)bool{
	list ,_ := fileUpload.GetAllowFileTypeList(category)
	for _,v:= range list{
		if v == extName{
			return true
		}
	}
	return false
}

func (fileUpload *FileUpload) InitMap() {
	var fileTypeMap sync.Map

	fileTypeMap.Store("ffd8ffe000104a464946", "jpg")  //JPEG (jpg)
	fileTypeMap.Store("89504e470d0a1a0a0000", "png")  //PNG (png)
	fileTypeMap.Store("47494638396126026f01", "gif")  //GIF (gif)
	fileTypeMap.Store("49492a00227105008037", "tif")  //TIFF (tif)
	fileTypeMap.Store("424d228c010000000000", "bmp")  //16色位图(bmp)
	fileTypeMap.Store("424d8240090000000000", "bmp")  //24位位图(bmp)
	fileTypeMap.Store("424d8e1b030000000000", "bmp")  //256色位图(bmp)
	fileTypeMap.Store("41433130313500000000", "dwg")  //CAD (dwg)
	fileTypeMap.Store("3c21444f435459504520", "html") //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c68746d6c3e0", "html")        //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c21646f637479706520", "htm")  //HTM (htm)
	fileTypeMap.Store("48544d4c207b0d0a0942", "css")  //css
	fileTypeMap.Store("696b2e71623d696b2e71", "js")   //js
	fileTypeMap.Store("7b5c727466315c616e73", "rtf")  //Rich Text Format (rtf)
	fileTypeMap.Store("38425053000100000000", "psd")  //Photoshop (psd)
	fileTypeMap.Store("46726f6d3a203d3f6762", "eml")  //Email [Outlook Express 6] (eml)
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "doc")  //MS Excel 注意：word、msi 和 excel的文件头一样
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "vsd")  //Visio 绘图
	fileTypeMap.Store("5374616E64617264204A", "mdb")  //MS Access (mdb)
	fileTypeMap.Store("252150532D41646F6265", "ps")
	fileTypeMap.Store("255044462d312e350d0a", "pdf")  //Adobe Acrobat (pdf)
	fileTypeMap.Store("2e524d46000000120001", "rmvb") //rmvb/rm相同
	fileTypeMap.Store("464c5601050000000900", "flv")  //flv与f4v相同
	fileTypeMap.Store("00000020667479706d70", "mp4")
	fileTypeMap.Store("49443303000000002176", "mp3")
	fileTypeMap.Store("000001ba210001000180", "mpg") //
	fileTypeMap.Store("3026b2758e66cf11a6d9", "wmv") //wmv与asf相同
	fileTypeMap.Store("52494646e27807005741", "wav") //Wave (wav)
	fileTypeMap.Store("52494646d07d60074156", "avi")
	fileTypeMap.Store("4d546864000000060001", "mid") //MIDI (mid)
	fileTypeMap.Store("504b0304140000000800", "zip")
	fileTypeMap.Store("526172211a0700cf9073", "rar")
	fileTypeMap.Store("235468697320636f6e66", "ini")
	fileTypeMap.Store("504b03040a0000000000", "jar")
	fileTypeMap.Store("4d5a9000030000000400", "exe")        //可执行文件
	fileTypeMap.Store("3c25402070616765206c", "jsp")        //jsp文件
	fileTypeMap.Store("4d616e69666573742d56", "mf")         //MF文件
	fileTypeMap.Store("3c3f786d6c2076657273", "xml")        //xml文件
	fileTypeMap.Store("494e5345525420494e54", "sql")        //xml文件
	fileTypeMap.Store("7061636b616765207765", "java")       //java文件
	fileTypeMap.Store("406563686f206f66660d", "bat")        //bat文件
	fileTypeMap.Store("1f8b0800000000000000", "gz")         //gz文件
	fileTypeMap.Store("6c6f67346a2e726f6f74", "properties") //bat文件
	fileTypeMap.Store("cafebabe0000002e0041", "class")      //bat文件
	fileTypeMap.Store("49545346030000006000", "chm")        //bat文件
	fileTypeMap.Store("04000000010000001300", "mxp")        //bat文件
	fileTypeMap.Store("504b0304140006000800", "docx")       //docx文件
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "wps")        //WPS文字wps、表格et、演示dps都是一样的
	fileTypeMap.Store("6431303a637265617465", "torrent")
	fileTypeMap.Store("6D6F6F76", "mov")         //Quicktime (mov)
	fileTypeMap.Store("FF575043", "wpd")         //WordPerfect (wpd)
	fileTypeMap.Store("CFAD12FEC5FD746F", "dbx") //Outlook Express (dbx)
	fileTypeMap.Store("2142444E", "pst")         //Outlook (pst)
	fileTypeMap.Store("AC9EBD8F", "qdf")         //Quicken (qdf)
	fileTypeMap.Store("E3828596", "pwl")         //Windows Password (pwl)
	fileTypeMap.Store("2E7261FD", "ram")         //Real Audio (ram)

	fileUpload.FileTypeMap = fileTypeMap
}

// 获取前面结果字节的二进制
func  (fileUpload *FileUpload) bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// 用文件前面几个字节来判断
// fSrc: 文件字节流（就用前面几个字节）
func  (fileUpload *FileUpload) GetFileType(fSrc []byte) string {
	var fileType string
	fileCode := fileUpload.bytesToHexString(fSrc)
	MyPrint("fileCode:",fileCode)
	fileUpload.FileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}

