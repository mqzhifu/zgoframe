package util

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
)

type AliOssOptions struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	BucketName      string
	LocalDomain     string
}

type AliOss struct {
	Op AliOssOptions
}

func NewAliOss(op AliOssOptions) *AliOss {
	aliOss := new(AliOss)
	aliOss.Op = op
	return aliOss
}

// 将本地文件上传到阿里云-OSS
func (aliOss *AliOss) UploadOneByFile(localFilePath string, relativePath, FileName string) error {
	//这里阿里云有个小BUG，所有的路径不能以反斜杠(/)开头
	if relativePath[0:1] == "/" {
		relativePath = relativePath[1:]
	}
	AccessKeyId := aliOss.Op.AccessKeyId
	AccessKeySecret := aliOss.Op.AccessKeySecret
	endpoint := aliOss.Op.Endpoint

	client, err := oss.New(endpoint, AccessKeyId, AccessKeySecret)
	//MyPrint("oss New:",client,err)
	if err != nil {
		return err
	}

	relativePathFile := relativePath + "/" + FileName

	bucketName := aliOss.Op.BucketName

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
func (aliOss *AliOss) DelOne(relativePath string) error {
	MyPrint("aliOss DelOne relativePath:", relativePath)
	_, bucket, err := aliOss.GetClientBucket()
	if err != nil {
		return err
	}
	err = bucket.DeleteObject(relativePath)
	return err
}

// 将本地文件上传到阿里云-OSS
func (aliOss *AliOss) UploadOneByStream(reader io.Reader, relativePath, FileName string) error {
	//这里阿里云有个小BUG，所有的路径不能以反斜杠(/)开头
	if relativePath[0:1] == "/" {
		relativePath = relativePath[1:]
	}

	_, bucket, err := aliOss.GetClientBucket()
	if err != nil {
		return err
	}
	//AccessKeyId := aliOss.Op.AccessKeyId
	//AccessKeySecret := aliOss.Op.AccessKeySecret
	//endpoint := aliOss.Op.Endpoint
	//client, err := oss.New(endpoint, AccessKeyId, AccessKeySecret)
	//MyPrint("oss New:",client,err)
	//if err != nil {
	//	return err
	//}
	relativePathFile := relativePath + "/" + FileName
	//bucketName := aliOss.Op.BucketName
	//MyPrint("oss endpoint:", endpoint, " AccessKeyId:", AccessKeyId, " AccessKeySecret:", AccessKeySecret, " bucketName:", bucketName)
	//bucket, err := client.Bucket(bucketName)
	//if err != nil {
	//	return err
	//}
	//MyPrint("bucket:",bucket,err)
	MyPrint("oss localFilePath:", reader, " relativePathFile:", relativePathFile)
	err = bucket.PutObject(relativePathFile, reader)
	MyPrint("PutObjectFromFile:", err)
	return err

}
func (aliOss *AliOss) GetClientBucket() (client *oss.Client, bucket *oss.Bucket, err error) {
	AccessKeyId := aliOss.Op.AccessKeyId
	AccessKeySecret := aliOss.Op.AccessKeySecret
	endpoint := aliOss.Op.Endpoint

	client, err = oss.New(endpoint, AccessKeyId, AccessKeySecret)
	//MyPrint("oss New:",client,err)
	if err != nil {
		return client, bucket, err
	}

	bucketName := aliOss.Op.BucketName

	MyPrint("oss endpoint:", endpoint, " AccessKeyId:", AccessKeyId, " AccessKeySecret:", AccessKeySecret, " bucketName:", bucketName)

	bucket, err = client.Bucket(bucketName)
	if err != nil {
		return client, bucket, err
	}
	return client, bucket, err
}

func (aliOss *AliOss) OssLs(dirPrefix string) (listObjectsResult oss.ListObjectsResult, err error) {
	//这里阿里云有个小BUG，所有的路径不能以反斜杠(/)开头
	if dirPrefix[0:1] == "/" {
		dirPrefix = dirPrefix[1:]
	}

	_, bucket, err := aliOss.GetClientBucket()
	if err != nil {
		return listObjectsResult, err
	}
	listObjectsResult, err = bucket.ListObjects(oss.Prefix(dirPrefix))
	//MyPrint("ListObjectsResult:", listObjectsResult, " err:", err)
	return listObjectsResult, err
}

func (aliOss *AliOss) DownloadFile(ossPathFile string, localPathFile string) error {
	MyPrint("ossPathFile:", ossPathFile, " localPathFile:", localPathFile)
	//这里阿里云有个小BUG，所有的路径不能以反斜杠(/)开头
	if ossPathFile[0:1] == "/" {
		ossPathFile = ossPathFile[1:]
	}
	_, bucket, err := aliOss.GetClientBucket()
	if err != nil {
		return err
	}

	err = bucket.GetObjectToFile(ossPathFile, localPathFile)
	MyPrint("DownloadFile err:", err)
	if err != nil {
		return err
	}
	return nil
}
