package util

import (
	"errors"
	"github.com/go-gomail/gomail"
	"go.uber.org/zap"
	"net/smtp"
	"strconv"
)

type EmailOption struct {
	Host      string      //smtp.163.com
	Port      string      //25
	FromEmail string      //发件人
	Password  string      //密码
	AuthCode  string      //授权码
	Log       *zap.Logger //日志
}

type MyEmail struct {
	Dialer      *gomail.Dialer //发送邮件的客户端
	EmailOption EmailOption    //邮件配置
}

//创建一个邮件客户端
func NewMyEmail(emailOption EmailOption) (*MyEmail, error) {
	myEmail := new(MyEmail)
	if emailOption.Host == "" || emailOption.Port == "" || emailOption.Password == "" {
		return nil, errors.New("NewMyEmail err: Host or Port or Password is empty")
	}
	port, _ := strconv.Atoi(emailOption.Port)
	myEmail.Dialer = gomail.NewDialer(emailOption.Host, port, emailOption.FromEmail, emailOption.Password)
	auth := smtp.PlainAuth("", emailOption.FromEmail, emailOption.AuthCode, emailOption.Host)
	myEmail.Dialer.Auth = auth

	myEmail.EmailOption = emailOption

	return myEmail, nil
}

//同步 - 发送一封邮件
func (myEmail *MyEmail) SendOneEmailSync(to string, Subject string, msg string) error {

	MyPrint(myEmail.EmailOption.Host, myEmail.EmailOption.Port, myEmail.EmailOption.FromEmail, myEmail.EmailOption.Password, myEmail.EmailOption.AuthCode)

	m := myEmail.GetInitSendOneEmailInfo(to, Subject, msg)
	err := myEmail.Dialer.DialAndSend(m)
	return err
}

//异步 - 发送一封邮件
func (myEmail *MyEmail) SendOneEmailAsync(to string, Subject string, msg string) error {
	m := myEmail.GetInitSendOneEmailInfo(to, Subject, msg)
	go myEmail.Dialer.DialAndSend(m)
	return nil
}

func (myEmail *MyEmail) GetInitSendOneEmailInfo(to string, Subject string, msg string) *gomail.Message {
	myEmail.EmailOption.Log.Info("myEmail GetInitSendOneEmailInfo : " + to + " subject:" + Subject)

	m := gomail.NewMessage()
	m.SetHeader("From", myEmail.EmailOption.FromEmail)
	m.SetHeader("Subject", Subject)
	m.SetHeader("To", to)

	m.SetBody("text/html", msg)
	return m
}
