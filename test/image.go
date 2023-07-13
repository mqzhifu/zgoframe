package test

import (
	"zgoframe/core/global"
	"zgoframe/util"
)

func ImageSlice() {
	dir := global.MainEnv.RootDir + "/" + global.C.Http.StaticPath + "/puzzle"
	extNameList := []string{"jpg", "png", "jpeg", "bmp"}
	fileList := util.ForeachDir(dir+"/ori", extNameList)
	//util.MyPrint(fileList)
	if len(fileList) == 0 {
		util.ExitPrint("no file")
	}

	for _, v := range fileList {
		//if v.Name != "pig.jpg" {
		//	continue
		//}
		oriFileDir := dir + "/ori/" + v.Name
		outDir := dir + "/shard"
		global.V.ImageMagic.Slice(oriFileDir, outDir)
	}

	util.ExitPrint(33)
	//
}
