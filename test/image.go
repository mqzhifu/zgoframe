package test

import (
	"strings"
	"zgoframe/core/global"
	"zgoframe/util"
)

func ImageSlice() {
	//dir := global.MainEnv.RootDir + "/" + global.C.Http.StaticPath + "/puzzle"
	extNameList := []string{"jpg", "png", "jpeg", "bmp"}
	dir := global.V.ImageSlice.OriImagePath
	fileList := util.ForeachDir(dir, extNameList)
	//util.MyPrint(fileList)
	if len(fileList) == 0 {
		util.ExitPrint("no file")
	}

	for _, v := range fileList {

		fileNameArr := strings.Split(v.Name, "/")
		fileName := fileNameArr[len(fileNameArr)-1]

		//global.V.ImageSlice.Slice(fileName)
		global.V.ImageSlice.ScaleImg(fileName, 20)
	}

	util.ExitPrint(33)
	//
}
