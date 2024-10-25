package util

import (
	"embed"
	"strings"
)

type StaticFileSystem struct {
	StaticFileSys embed.FS //静态文件
	BuildStatic   string   //是否编译静态文件
}

func NewStaticFileSystem(static embed.FS, buildStatic string) *StaticFileSystem {
	fileSystem := new(StaticFileSystem)
	fileSystem.StaticFileSys = static
	fileSystem.BuildStatic = buildStatic
	return fileSystem
}

// 这个是兼容模式，是正常硬盘上的物理文件，也可以是 EDFS 上的文件
func (staticFileSystem StaticFileSystem) GetStaticFileContent(pathFile string) (a string, err error) {
	var content string
	if staticFileSystem.BuildStatic == "on" {
		contentBytes, err := staticFileSystem.StaticFileSys.ReadFile(pathFile)
		MyPrint(err)
		if err != nil {
			return "", err
		}
		content = string(contentBytes)
	} else {
		content, err = ReadString(pathFile)
		if err != nil {
			return "", err
		}
	}

	return content, err
}

// 这个是兼容模式，是正常硬盘上的物理文件，也可以是EDFS上的文件
func (staticFileSystem StaticFileSystem) GetStaticFileContentLine(pathFile string) (a []string, err error) {
	var content string
	if staticFileSystem.BuildStatic == "on" {
		content, err = staticFileSystem.GetStaticFileContent(pathFile)
		if err != nil {
			return a, err
		}
	} else {
		content, err = ReadString(pathFile)
		if err != nil {
			return a, err
		}
	}

	return strings.Split(content, "\n"), err
}
