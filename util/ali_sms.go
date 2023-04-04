package util

import (
	"encoding/json"
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

type AliSmsOp struct {
	AccessKeyId     string //阿里云短信服务的AccessKeyId
	AccessKeySecret string //阿里云短信服务的AccessKeySecret
	Endpoint        string //阿里云短信服务的Endpoint
}

type AliSms struct {
	Op           AliSmsOp                 //阿里云短信服务的配置
	AliSmsClient *dysmsapi20170525.Client //阿里云短信服务的客户端
}

//创建一个阿里云短信服务客户端
func NewAliSms(aliSmsOp AliSmsOp) (*AliSms, error) {
	aliSms := new(AliSms)
	aliSms.Op = aliSmsOp

	//判断参数是否为空
	if aliSmsOp.AccessKeyId == "" || aliSmsOp.AccessKeySecret == "" || aliSmsOp.Endpoint == "" {
		return nil, errors.New("NewAliSms err: AccessKeyId or AccessKeySecret or Endpoint is empty")
	}

	config := &openapi.Config{
		AccessKeyId:     &aliSmsOp.AccessKeyId,
		AccessKeySecret: &aliSmsOp.AccessKeySecret,
		Endpoint:        tea.String(aliSmsOp.Endpoint),
	}
	// 访问的域名
	aliSmsClient, err := dysmsapi20170525.NewClient(config)
	if err != nil {
		return aliSms, err
	}
	aliSms.AliSmsClient = aliSmsClient
	return aliSms, nil
}

func (aliSms *AliSms) Send(Receiver string, templateCode string, signName string, ReplaceVar string) (string, error) {
	if Receiver == "" || templateCode == "" || signName == "" {
		return "", errors.New("Receiver | templateCode | signName is empty ")
	}
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  &Receiver,
		SignName:      &signName,
		TemplateCode:  &templateCode,
		TemplateParam: &ReplaceVar,
	}
	//MyPrint("========")
	responseInfo, err := aliSms.AliSmsClient.SendSms(sendSmsRequest)
	//MyPrint("========")
	//MyPrint(responseInfo, err)
	//MyPrint("========")
	if err != nil {
		MyPrint(err)
		return "", err
	}

	if *responseInfo.StatusCode != 200 {
		return "", errors.New("responseInfo.StatusCode != 200")
	}
	responseInfoBodyBytes, err := json.Marshal(responseInfo.Body)
	if err != nil {
		MyPrint(err)
	}
	//MyPrint("========")
	backInfo := string(responseInfoBodyBytes)
	//MyPrint(backInfo)
	return backInfo, nil
}
