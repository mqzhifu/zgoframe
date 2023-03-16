package util

type AliSmsOp struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
}

type AliSms struct {
	Op AliSmsOp
}

func NewAliSms(aliSmsOp AliSmsOp) (*AliSms, error) {
	aliSms := new(AliSms)
	aliSms.Op = aliSmsOp
	return aliSms, nil
}

func (aliSms *AliSms) Send(Receiver string, ThirdTemplateId string, signName string, ReplaceVar string) (string, error) {
	return "", nil

}
