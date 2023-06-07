package util

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// 检查一个文件是否已存在
func FileExist(filename string) (os.FileInfo, error) {
	if filename == "" {
		msg := "filename empty"
		return nil, errors.New(msg)
	}
	fd, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New(filename + ":IsNotExist")
		}
		return nil, err
	}
	if fd.IsDir() {
		return fd, errors.New(filename + ":is dir")
	}
	return fd, nil
}

func PathExists(path string) (os.FileInfo, error) {
	if path == "" {
		msg := "path empty"
		return nil, errors.New(msg)
	}
	fd, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) { //这个不能动，因为外层要使用os.IsNotExist 继续判断
			return nil, err
		}
		return nil, err
	}
	if !fd.IsDir() {
		return fd, errors.New(path + ":is not dir")
	}
	return fd, nil
}

// 打开一个文件，并按照换行符 读取到一个数组中
func ReadLine(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	//defer f.Close()
	buf := bufio.NewReader(f)
	var result []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				return result, nil
			}
			return nil, err
		}
		result = append(result, line)
	}
	return result, nil
}
func UrlAppendIpHost(protocol string, url string, ip string, port string) string {
	fullUrl := protocol + "://" + ip
	if port != "" {
		fullUrl += ":" + port
	}
	fullUrl += "/" + url
	return fullUrl
}

func UrlAppendDomain(protocol string, url string, domain string, port string) string {
	fullUrl := protocol + "://" + domain
	if port != "" {
		fullUrl += ":" + port
	}
	fullUrl += "/" + url
	return fullUrl
}

// 打开一个文件，并按照换行符 读取到一个字符串中
func ReadString(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	//defer f.Close()

	//buf := bufio.NewReader(f)
	var strings string
	for {
		buf := make([]byte, 1024)
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		strings += string(buf)
	}
	return strings, nil
}

type ForeachDirInfo struct {
	Cate string
	Name string
}

// 遍历一个目录的所有文件/目录列表，但 不递归，也就是子目录不处理
func ForeachDir(path string) []ForeachDirInfo {
	var list []ForeachDirInfo
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		MyPrint("ForeachDir err:", err.Error())
		return list
	}
	for _, file := range fs {
		if file.IsDir() {
			foreachDirInfo := ForeachDirInfo{
				Name: file.Name(),
				Cate: "dir",
			}
			list = append(list, foreachDirInfo)
		} else {
			foreachDirInfo := ForeachDirInfo{
				Name: file.Name(),
				Cate: "file",
			}
			list = append(list, foreachDirInfo)
		}
	}
	return list
}

func FileCopy(srcPath string, tarPath string) error {
	err := CheckFileModify(srcPath, tarPath)
	if err != nil {
		return err
	}
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(tarPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func FileMove(srcPath string, tarPath string) error {
	err := CheckFileModify(srcPath, tarPath)
	if err != nil {
		return err
	}

	err = os.Rename(srcPath, tarPath)
	if err != nil {
		return err
	}

	return nil
}

func CheckFileModify(srcPath string, tarPath string) error {
	if srcPath == "" || tarPath == "" {
		return errors.New("relativePath is empty")
	}

	if srcPath == tarPath {
		return errors.New("srcRelativePath == tarRelativePath")
	}

	_, err := FileExist(srcPath)
	if err != nil { //源文件已存在
		return err
	}

	_, err = FileExist(tarPath)
	if err == nil { //目标文件已存在的
		return err
	}

	return nil
}
