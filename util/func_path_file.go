package util

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

//检查一个文件是否已存在
func FileExist(filename string) (os.FileInfo,error) {
	if filename == ""{
		msg := "filename empty"
		return nil,errors.New(msg)
	}
	fd, err := os.Stat(filename)
	if err != nil{
		if os.IsNotExist(err){
			return nil,errors.New(filename+":IsNotExist")
		}
		return nil,err
	}
	if fd.IsDir(){
		return fd,errors.New(filename+":is dir")
	}
	return fd,nil
}

func PathExists(path string) (os.FileInfo, error) {
	if path == ""{
		msg := "path empty"
		return nil,errors.New(msg)
	}
	fd, err := os.Stat(path)
	if err != nil{
		if os.IsNotExist(err){//这个不能动，因为外层要使用os.IsNotExist 继续判断
			return nil,err
		}
		return nil,err
	}
	if !fd.IsDir(){
		return fd,errors.New(path+":is not dir")
	}
	return fd,nil
}

//打开一个文件，并按照换行符 读取到一个数组中
func ReadLine(fileName string) ([]string,error){
	f, err := os.Open(fileName)
	if err != nil {
		return nil,err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				return result,nil
			}
			return nil,err
		}
		result = append(result,line)
	}
	return result,nil
}

//打开一个文件，并按照换行符 读取到一个数组中
func ReadString(fileName string) (string,error){
	f, err := os.Open(fileName)
	if err != nil {
		return "",err
	}
	defer f.Close()

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
	return strings,nil
}

//遍历一个目录的所有文件列表，但 子目录不处理
func GetFileListByDir(path string)[]string {
	var fileList []string
	fs,err := ioutil.ReadDir(path)
	if err != nil{
		MyPrint("GetFileListByDir err:",err.Error())
		return fileList
	}
	for _,file:=range fs{
		if file.IsDir(){
			//fmt.Println(path+file.Name())
			//GetFileListByDir(path+file.Name()+"/")
		}else{
			//fmt.Println(path+file.Name())
			fileList = append(fileList,file.Name())
		}
	}
	return fileList
}

func ForeachDir(path string)(dirList []string,err error) {
	//var fileList []string
	fs,err := ioutil.ReadDir(path)
	if err != nil{
		return dirList,errors.New("GetFileListByDir err:" + err.Error())
		//return fileList
	}
	for _,file:=range fs{
		if file.IsDir(){
			//fmt.Println(path+file.Name())
			dirList = append(dirList,file.Name())
		}
	}
	return dirList,nil
}
