package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	"image/jpeg"
	"os"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/util"
)

// @Tags Pic
// @Summary 切割一张图片( http-form 表单模式 )
// @Security ApiKeyAuth
// @Description 一张图片平均分割成若干等分
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
// @Router 	/pic/split [POST]
func Split(c *gin.Context) {
	util.MyPrint(3333)
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		util.MyPrint("err1:", err.Error())
		return
	}

	uploadRs, err := global.V.ImgManager.UploadOne(header, "pic_split", 0, 2)

	fileName := uploadRs.LocalDiskPath + "/" + uploadRs.Filename
	f, err := os.Open(fileName)
	if err != nil {
		util.MyPrint("err1")
		return
	}

	defer f.Close()

	img, extName, err := image.Decode(f)
	if err != nil {
		util.MyPrint("err2:", err)
		return
	}
	util.MyPrint("extName:", extName)
	fmt.Println("图片颜色模型:", img.ColorModel()) // 图片颜色模型

	splitNumber := 10 //平均，切割多少块

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	everyBlockWidth := width / splitNumber
	everyBlockHeight := height / splitNumber

	util.MyPrint("width:", width, "height:", height, " everyBlockWidth:", everyBlockWidth, " everyBlockHeight:", everyBlockHeight)

	//fmt.Println("图片长宽:", img.Bounds().Dx(), "----", img.Bounds().Dy(), "----", img.Bounds().Size().String()) // 图片长宽

	//fmt.Println("该像素点的颜色:", img.At(100, 100))                                                             // 该像素点的颜色

	imgYCbCr := img.(*image.YCbCr)

	everyBlockWidthStepX0 := 0
	everyBlockHeightStepY0 := 0
	everyBlockWidthStepX1 := 0
	everyBlockHeightStepY1 := 0
	for i := 0; i < splitNumber; i++ {
		fmt.Println("=============", i, "======")
		everyBlockHeightStepY0 = i * everyBlockHeight
		everyBlockHeightStepY1 = everyBlockHeightStepY0 + everyBlockHeight
		for j := 0; j < splitNumber; j++ {
			everyBlockWidthStepX0 = j * everyBlockWidth
			everyBlockWidthStepX1 = (j + 1) * everyBlockWidth

			subImage := imgYCbCr.SubImage(image.Rect(everyBlockWidthStepX0, everyBlockHeightStepY0, everyBlockWidthStepX1, everyBlockHeightStepY1)).(*image.YCbCr)
			// 保存图片
			LocalDiskUploadBasePath := global.V.ImgManager.GetLocalDiskUploadBasePath()
			fmt.Println(LocalDiskUploadBasePath)
			fileName2 := "/new" + strconv.Itoa(i) + strconv.Itoa(j) + ".jpg"
			create, _ := os.Create(LocalDiskUploadBasePath + fileName2)
			err = jpeg.Encode(create, subImage, &jpeg.Options{100})

			if err != nil {
				util.MyPrint("err3:", err)
				return
			}
			fmt.Println(everyBlockWidthStepX0, everyBlockHeightStepY0, everyBlockWidthStepX1, everyBlockHeightStepY1, " fileName2:", fileName2)
		}
		//os.Exit(3)

	}

}
