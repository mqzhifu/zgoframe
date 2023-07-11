// http 传输加密/解密
package encrypt

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
	"zgoframe/http/request"
	"zgoframe/model"
	"zgoframe/util"
)

func DecodeBody(c *gin.Context) (projectInfo model.Project, err error) {
	project, err := request.GetMyProject(c)
	if err != nil {
		return project, errors.New(err.Error())
	}

	if project.DataEncrypt <= 0 {
		return project, nil
	}

	bodyByte, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return project, errors.New("ioutil.ReadAll c.Request.Body err:" + err.Error())
	}

	if len(bodyByte) <= 0 {
		util.MyPrint("DecodeBody len(bodyByte) <= 0")
		return project, nil
	}

	bodyStr := string(bodyByte)
	//util.MyPrint("DataEncrypt ReadAll:", bodyStr, err, "project.SecretKey:", project.SecretKey)
	util.MyPrint("DataEncrypt ReadAll:", bodyStr, err, "project.DataEncrypt:", project.DataEncrypt)

	var bodyDataByte []byte
	switch project.DataEncrypt {
	case model.DATA_ENCRYPT_BASE64:
		bodyDataByte, err = base64.StdEncoding.DecodeString(bodyStr)
		if err != nil {
			return project, err
		}

	case model.DATA_ENCRYPT_AES_CBC_BASE64:
		//R := util.AesEncryptCBC([]byte("123456"), []byte(project.SecretKey))
		//mw := base64.StdEncoding.EncodeToString(R)
		//util.ExitPrint(mw)
		dataBase64Byte, err := base64.StdEncoding.DecodeString(string(bodyByte))
		if err != nil {
			return project, err
		}
		bodyDataByte, err = util.AesDecryptCBC(dataBase64Byte, []byte(project.SecretKey))
		//util.MyPrint("AesDecryptCBC rs:", string(rs), err)
		if err != nil {
			return project, err
		}
	default:
		return project, errors.New("project.DataEncrypt value err:" + strconv.Itoa(project.DataEncrypt))
	}
	bodyData := string(bodyDataByte)
	util.MyPrint("DecodeBody final :", bodyData)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyDataByte))
	return project, nil
}

func CheckSign(c *gin.Context) (projectInfo model.Project, err error) {
	myHeader, err := request.GetMyHeader(c)
	if err != nil {
		return projectInfo, errors.New(strconv.Itoa(5019))
	}

	project, err := request.GetMyProject(c)
	if err != nil {
		return projectInfo, errors.New(strconv.Itoa(5020))
	}

	if project.AuthSign <= 0 { //这里是不需要验证签名的接口
		return projectInfo, nil
	}
	if myHeader.Sign == "" {
		return projectInfo, errors.New(strconv.Itoa(5108))
	}
	bodyByte, err := ioutil.ReadAll(c.Request.Body)
	//util.MyPrint("CheckSign ReadAll:", string(bodyByte), err)
	if err != nil {
		util.MyPrint("ioutil.ReadAll c.Request.Body err:" + err.Error())
		return projectInfo, errors.New(strconv.Itoa(5106))
	}
	str := strconv.Itoa(project.Id) + strconv.Itoa(myHeader.ClientReqTime) + project.SecretKey + string(bodyByte)
	md5Str := util.MD5X(str)
	//util.MyPrint("CheckSign:"+strconv.Itoa(project.Id), strconv.Itoa(myHeader.ClientReqTime), project.SecretKey, string(bodyByte))
	util.MyPrint("CheckSign:" + strconv.Itoa(project.Id) + " + " + strconv.Itoa(myHeader.ClientReqTime) + " + " + "******" + " + " + string(bodyByte) + " = " + md5Str)
	if md5Str != myHeader.Sign {
		util.MyPrint("md5Str != header.Sign")
		return projectInfo, errors.New(strconv.Itoa(5107))
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyByte))

	return projectInfo, nil
}

func CreateSign(c *gin.Context, ClientReqTime int, body string) (projectInfo model.Project, rs string, err error) {
	project, err := request.GetMyProject(c)
	if err != nil {
		return projectInfo, "", errors.New(strconv.Itoa(5020))
	}

	if project.AuthSign <= 0 {
		return project, "", nil
	}
	str := strconv.Itoa(project.Id) + strconv.Itoa(ClientReqTime) + project.SecretKey + body
	md5Str := util.MD5X(str)
	return project, md5Str, nil
}

func EncodeBody(c *gin.Context, body string) (projectInfo model.Project, data string, err error) {
	project, err := request.GetMyProject(c)
	if err != nil {
		return projectInfo, data, errors.New(strconv.Itoa(5020))
	}

	if project.DataEncrypt <= 0 {
		return project, data, nil
	}

	if len(body) <= 0 {
		return project, data, nil
	}

	switch project.DataEncrypt {
	case model.DATA_ENCRYPT_BASE64:
		data = base64.StdEncoding.EncodeToString([]byte(body))
	case model.DATA_ENCRYPT_AES_CBC_BASE64:
		//util.ExitPrint(dataBase64)
		decrypted := util.AesEncryptCBC([]byte(body), []byte(project.SecretKey))
		dataBase64 := base64.StdEncoding.EncodeToString(decrypted)
		data = dataBase64
	default:
		return project, data, errors.New("project.DataEncrypt err:" + strconv.Itoa(project.DataEncrypt))
	}

	return project, data, nil
}
