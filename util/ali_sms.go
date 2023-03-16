package util

import (
	"encoding/json"
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

type AliSmsOp struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
}

type AliSms struct {
	Op           AliSmsOp
	AliSmsClient *dysmsapi20170525.Client
}

func NewAliSms(aliSmsOp AliSmsOp) (*AliSms, error) {
	aliSms := new(AliSms)
	aliSms.Op = aliSmsOp

	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: &aliSmsOp.AccessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: &aliSmsOp.AccessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
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
