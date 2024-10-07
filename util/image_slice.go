package util

import (
	"errors"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
)

type ImageSlice struct {
	AllowExtList   []string
	NodeUnit       int    //一张图片要切多少块
	FileMinSize    int64  //文件最小尺寸(字节)
	ImageMinSize   int    //图片最小尺寸(像素)
	OriImagePath   string //原图路径
	ShardImagePath string //切片图路径
	ThumbnailPath  string //缩略图路径
}

type SliceNode struct {
	Id string //列号+行号
	X1 int
	X2 int
	Y1 int
	Y2 int
}

func NewImageSlice(OriImagePath string, ShardImagePath string, ThumbnailPath string) *ImageSlice {
	imageSlice := new(ImageSlice)
	imageSlice.AllowExtList = []string{"jpg", "png", "jpeg", "bmp"}
	imageSlice.NodeUnit = 3       // 共计：3 X 3 = 9块
	imageSlice.FileMinSize = 10   //文件内容小于10个字节这就有问题了，一个图片的头信息都不止10个字节
	imageSlice.ImageMinSize = 100 //100 x 100 ，100 以下的图片，切出来是没有意义 的，太小块了

	imageSlice.OriImagePath = OriImagePath
	imageSlice.ShardImagePath = ShardImagePath
	imageSlice.ThumbnailPath = ThumbnailPath

	return imageSlice
}
func (imageSlice *ImageSlice) GetImgShape(width int, height int) string {
	imageShape := "rectangle"
	if width == height {
		imageShape = "square"
	}
	return imageShape
}

// 将一张图片 切割(平均)成若干小图片
func (imageSlice *ImageSlice) Slice(oriFileName string) error {
	MyPrint("oriFileName:" + oriFileName)

	imgObj, imgType, err := imageSlice.LoadFileContentMapImage(oriFileName)
	if err != nil {
		MyPrint("LoadFileContentMapImage,err:" + err.Error())
		return err
	}
	imgObjWidth := imgObj.Bounds().Size().X
	imgObjHeight := imgObj.Bounds().Size().Y
	imageShape := imageSlice.GetImgShape(imgObjWidth, imgObjHeight)
	MyPrint("LoadFileContentMapImage ,imgType:"+imgType+",x:", imgObjWidth, ",y:", imgObjHeight, ",imageShape:", imageShape)

	if imgObjWidth < imageSlice.ImageMinSize || imgObjHeight < imageSlice.ImageMinSize {
		return errors.New("image size too small:" + strconv.Itoa(imageSlice.ImageMinSize))
	}

	err = imageSlice.CheckFileType(imgType)
	if err != nil {
		return err
	}
	//一个图片切成一块，没有意义。或者当 width < NodeUnit 或 height < NodeUnit 时，也没有意义
	if imgObjWidth/imageSlice.NodeUnit < 2 || imgObjHeight/imageSlice.NodeUnit < 2 {
		return errors.New("图片长宽/切片数(" + strconv.Itoa(imgObjWidth) + "/" + strconv.Itoa(imageSlice.NodeUnit) + ")，太小了")
	}
	//横坐标的格子
	everyNodeWidthList := imageSlice.ComputeNodeSize(imgObjWidth, imageSlice.NodeUnit)
	//纵坐标的格子
	everyNodeHeightList := imageSlice.ComputeNodeSize(imgObjHeight, imageSlice.NodeUnit)
	//合并格子，计算出最终：每块节点的 起始坐标和结束坐标
	nodeList := imageSlice.DrawGrid(everyNodeWidthList, everyNodeHeightList)
	//获取切片图的存放目录
	outDirPrefix := imageSlice.CheckNewShardDir(oriFileName)
	for _, node := range nodeList {
		imageSlice.DrawImage(imgType, node, imgObj, outDirPrefix)
	}

	return nil
}

// 根据坐标点，画(从原图)出新的图片
func (imageSlice *ImageSlice) DrawImage(imgType string, node SliceNode, img image.Image, outDirPrefix string) error {
	shardImgFilePath := outDirPrefix + "/" + node.Id + "." + imgType
	MyPrint("shardImgFilePath:" + shardImgFilePath)

	f, err := os.Create(shardImgFilePath) //创建文件
	if err != nil {
		return err
	}

	if imgType == "png" {
		rgbImg := img.(*image.NRGBA)
		subImg := rgbImg.SubImage(image.Rect(node.X1, node.Y1, node.X2, node.Y2)).(*image.NRGBA)
		png.Encode(f, subImg)
	} else if imgType == "jpeg" || imgType == "jpg" {
		rgbImg := img.(*image.YCbCr)
		subImg := rgbImg.SubImage(image.Rect(node.X1, node.Y1, node.X2, node.Y2)).(*image.YCbCr)
		//imageMagic.saveImage(target, subImg, 100, imgType)
		var opt jpeg.Options
		opt.Quality = 100

		jpeg.Encode(f, subImg, &opt)
	} else if imgType == "bmp" {
		rgbImg := img.(*image.RGBA)
		subImg := rgbImg.SubImage(image.Rect(node.X1, node.Y1, node.X2, node.Y2)).(*image.RGBA)
		//imageMagic.saveImage(target, subImg, 100, imgType)
		png.Encode(f, subImg)
	} else {
		return errors.New("not support image type:" + imgType)
	}
	return nil
}

// 每个原图，会生成一个碎片文件夹，里面存放着碎片文件
func (imageSlice *ImageSlice) CheckNewShardDir(path string) string {
	fullPath := imageSlice.GetShardFullPath(path)
	_, err := PathExists(fullPath)
	if err != nil {
		MyPrint("CheckNewShardDir path:" + path + ",文件夹不存在，创建")
		//文件夹不存在，创建
		os.Mkdir(fullPath, os.ModePerm)
	} else {
		MyPrint("CheckNewShardDir path:" + path + ",文件夹存在，删除")
		//文件夹存在，删除
		os.RemoveAll(path)
	}
	return fullPath
}

// 横纵坐标，画格子
func (imageSlice *ImageSlice) ComputeNodeSize(total int, unit int) []int {
	list := []int{}
	everyNode := total / unit //这里大概率出现小数点，而实际像素必须是整形
	i := 0
	for ; i < unit; i++ {
		list = append(list, everyNode)
	}
	if everyNode*unit == total {
		return list
	}
	//如果不能整除，那么最后一个节点的大小就是 总大小 - 前面节点的大小
	list[unit-1] = total - everyNode*(unit-1)

	return list
}

// 开始画-格子，计算出最终：每块节点的 起始坐标和结束坐标
func (imageSlice *ImageSlice) DrawGrid(everyNodeWidthList []int, everyNodeHeightList []int) []SliceNode {
	startWidth := 0
	nodeList := []SliceNode{}
	for k1, v1 := range everyNodeWidthList {
		startHeight := 0
		for k2, v2 := range everyNodeHeightList {
			node := SliceNode{
				Id: strconv.Itoa(k1) + "_" + strconv.Itoa(k2),
				X1: startWidth,
				Y1: startHeight,
				X2: startWidth + v1,
			}

			if k1 != 0 {
				//node.XStart = startWidth + 1
				//node.YStart = startHeight + 1
			}

			startHeight += v2
			//
			//node.XEnd = startWidth
			//node.YEnd = startHeight
			node.Y2 = startHeight
			nodeList = append(nodeList, node)
		}
		startWidth += everyNodeWidthList[k1]

	}

	for k, v := range nodeList {
		MyPrint(k, "x1:", v.X1, " y1:", v.Y1, "x2:", v.X2, " y2:", v.Y2)
	}

	return nodeList
}

func (imageSlice *ImageSlice) CheckFileType(fileType string) error {
	for _, v := range imageSlice.AllowExtList {
		if fileType == v {
			return nil
		}
	}
	return errors.New("not support file type:" + fileType)
}

// 打开原文件数据流，给到 image.Decode，解析出图片
func (imageSlice *ImageSlice) LoadFileContentMapImage(oriFileName string) (img image.Image, imgType string, err error) {
	oriImgPath := imageSlice.GetOriFullPath(oriFileName)
	fs, err := os.Stat(oriImgPath)
	if err != nil {
		return nil, "", errors.New("file not exist :" + err.Error())
	}

	if fs.Size() < imageSlice.FileMinSize { //正常一张图片的大小都不会小于10字节(图片头信息也不会小于10字节)
		return nil, "", errors.New("file size < 10")
	}
	//打开文件
	file, err := os.Open(oriImgPath)
	if err != nil {
		return
	}
	//defer file.Close()
	img, imgType, err = image.Decode(file)
	if err != nil {
		return nil, "", errors.New("image.Decode err:" + err.Error())
	}

	return img, imgType, err
}

func (imageSlice *ImageSlice) GetOriFullPath(fileName string) string {
	return imageSlice.OriImagePath + fileName
}

func (imageSlice *ImageSlice) GetShardFullPath(fileName string) string {
	return imageSlice.ShardImagePath + fileName + "/"
}

func (imageSlice *ImageSlice) ScaleImg(fileName string, percent int) error {
	imgObj, imgType, err := imageSlice.LoadFileContentMapImage(fileName)
	if err != nil {
		MyPrint("LoadFileContentMapImage,err:" + err.Error())
		return err
	}
	imgObjWidth := imgObj.Bounds().Size().X
	imgObjHeight := imgObj.Bounds().Size().Y
	imageShape := imageSlice.GetImgShape(imgObjWidth, imgObjHeight)
	MyPrint("LoadFileContentMapImage ,imgType:"+imgType+",x:", imgObjWidth, ",y:", imgObjHeight, ",imageShape:", imageShape)

	scaleWidth := imgObjWidth * percent / 100
	scaleHeight := imgObjHeight * percent / 100

	thumbnailImbObj := resize.Thumbnail(uint(scaleWidth), uint(scaleHeight), imgObj, resize.Lanczos3)

	outFilePath := imageSlice.ThumbnailPath + fileName
	outFilePathFd, _ := os.Create(outFilePath)

	switch imgType {
	case "jpeg":
		return jpeg.Encode(outFilePathFd, thumbnailImbObj, &jpeg.Options{100})
	case "png":
		return png.Encode(outFilePathFd, thumbnailImbObj)
	case "gif":
		return gif.Encode(outFilePathFd, thumbnailImbObj, &gif.Options{})
	case "bmp":
		return bmp.Encode(outFilePathFd, thumbnailImbObj)
	default:
		return errors.New("ERROR FORMAT")
	}

	return nil
}

//// 保存图片
//func (imageMagic *ImageMagic) saveImage(path string, subImg image.Image, quality int, tt string) error {
//	f, err := os.Create(path) //创建文件
//	if err != nil {
//		return err
//	}
//	//defer f.Close() //关闭文件
//	var opt jpeg.Options
//	opt.Quality = quality
//	switch tt {
//	case "jpg":
//	case "jpeg":
//		return jpeg.Encode(f, subImg, &opt)
//	case "png":
//		return png.Encode(f, subImg)
//	default:
//	}
//	return nil
//}
