package util

import (
	"github.com/go-gomail/gomail"
	"go.uber.org/zap"
	"net/smtp"
)

type EmailOption struct {
	Host string
	Port int
	FromEmail string
	Password string
	Log *zap.Logger
}

type MyEmail struct {
	Dialer  *gomail.Dialer
	EmailOption EmailOption
}

func NewMyEmail (emailOption EmailOption)*MyEmail{
	myEmail := new (MyEmail)

	myEmail.Dialer = gomail.NewDialer(emailOption.Host, emailOption.Port,emailOption.FromEmail, emailOption.Password)
	auth := smtp.PlainAuth("", emailOption.FromEmail,  "glnteewafftmcaje", emailOption.Host )
	myEmail.Dialer.Auth = auth

	myEmail.EmailOption = emailOption

	return myEmail
}
//同步 - 发送一封邮件
func(myEmail *MyEmail) SendOneEmailSync(to  string,Subject string,msg string)error{

	m := myEmail.GetInitSendOneEmailInfo(to,Subject,msg)
	err := myEmail.Dialer.DialAndSend(m)
	return err
}
//异步 - 发送一封邮件
func(myEmail *MyEmail) SendOneEmailAsync(to  string,Subject string,msg string)error{
	m := myEmail.GetInitSendOneEmailInfo(to,Subject,msg)
	go  myEmail.Dialer.DialAndSend(m)
	return nil
}

func(myEmail *MyEmail)GetInitSendOneEmailInfo(to  string,Subject string,msg string)*gomail.Message{
	myEmail.EmailOption.Log.Info("myEmail GetInitSendOneEmailInfo : "+to + " subject:" +Subject)

	m := gomail.NewMessage()
	m.SetHeader("From",myEmail.EmailOption.FromEmail)
	m.SetHeader("Subject", Subject)
	m.SetHeader("To", to)

	m.SetBody("text/html", msg)
	return m
}
