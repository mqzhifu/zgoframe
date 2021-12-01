package util

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

//检查一个文件是否已存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename);os.IsNotExist(err){
		exist = false
	}
	return exist
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}


func ReadLine(fileName string) ([]string,error){
	f, err := os.Open(fileName)
	if err != nil {
		return nil,err
	}
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