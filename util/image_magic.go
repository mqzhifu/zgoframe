package util

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
	"strings"
)

type ImageMagic struct {
	AllowExtList []string
	NodeUnit     int //一张图片要切多少块
}

type SliceNode struct {
	Id string
	X1 int
	X2 int
	Y1 int
	Y2 int
}

func NewImageMagic() *ImageMagic {
	image := new(ImageMagic)
	image.AllowExtList = []string{"jpg", "png", "jpeg", "bmp"}
	image.NodeUnit = 3
	return image
}

// 将一张图片 切割(平均)成若干小图片
func (imageMagic *ImageMagic) Slice(oriFilePath string, outPath string) error {
	MyPrint("oriFilePath:" + oriFilePath)
	//从路径中，找到文件名（用于碎片文件名的前缀）
	fileNameArr := strings.Split(oriFilePath, "/")
	fileName := fileNameArr[len(fileNameArr)-1]

	imgObj, imgType, err := imageMagic.LoadFileContentMapImage(oriFilePath)
	if err != nil {
		MyPrint("LoadFileContentMapImage,err:" + err.Error())
		return err
	}
	imgObjWidth := imgObj.Bounds().Size().X
	imgObjHeight := imgObj.Bounds().Size().Y

	MyPrint("LoadFileContentMapImage ,imgType:"+imgType+",x:", imgObjWidth, ",y:", imgObjHeight)

	err = imageMagic.CheckFileType(imgType)
	if err != nil {
		return err
	}
	//横坐标的格子
	everyNodeWidthList := imageMagic.computeNodeSize(imgObjWidth, imageMagic.NodeUnit)
	//纵坐标的格子
	everyNodeHeightList := imageMagic.computeNodeSize(imgObjHeight, imageMagic.NodeUnit)
	//合并格子，计算出最终：每块节点的 起始坐标和结束坐标
	nodeList := imageMagic.Merge(everyNodeWidthList, everyNodeHeightList)
	for _, node := range nodeList {
		imageMagic.DrawImage(imgType, node, imgObj, outPath, fileName)
	}

	return nil
}

// 根据坐标点，画(从原图)出新的图片
func (imageMagic *ImageMagic) DrawImage(imgType string, node SliceNode, img image.Image, outPath string, oriFileName string) error {
	fileNamePrefix := strings.Split(oriFileName, ".")[0]
	target := outPath + "/" + fileNamePrefix + "_" + node.Id + "." + imgType
	MyPrint("target:" + target)

	f, err := os.Create(target) //创建文件
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

// 横纵坐标，画格子
func (imageMagic *ImageMagic) computeNodeSize(total int, unit int) []int {
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

// 合并格子，计算出最终：每块节点的 起始坐标和结束坐标
func (imageMagic *ImageMagic) Merge(everyNodeWidthList []int, everyNodeHeightList []int) []SliceNode {
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

	//ExitPrint(333)
	return nodeList
}

func (imageMagic *ImageMagic) CheckFileType(fileType string) error {
	for _, v := range imageMagic.AllowExtList {
		if fileType == v {
			return nil
		}
	}
	return errors.New("not support file type:" + fileType)
}

// 打开原文件数据流，给到 image.Decode，解析出图片
func (imageMagic *ImageMagic) LoadFileContentMapImage(oriImgPath string) (img image.Image, imgType string, err error) {
	fs, err := os.Stat(oriImgPath)
	if err != nil {
		return nil, "", errors.New("file not exist :" + err.Error())
	}

	if fs.Size() < 10 { //正常一张图片的大小都不会小于10字节(图片头信息也不会小于10字节)
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
