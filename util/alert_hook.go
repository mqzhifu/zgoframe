package util

//
import (
	"github.com/go-gomail/gomail"
	"go.uber.org/zap"
	"net/smtp"
)

type AlertHook struct {
	Email *MyEmail
	EmailOption EmailOption
	Log *zap.Logger
}

type EmailOption struct {
	Host string
	Port int
	FromEmail string
	Password string
}

func NewAlertHook(log *zap.Logger)*AlertHook{
	alertHook := new (AlertHook)
	alertHook.Log = log

	emailOption :=EmailOption{
		Host: "smtp.qq.com",
		Port: 465,
		//Port: 587,
		FromEmail: "78878296@qq.com",
		Password: "mM123456",
	}

	alertHook.Email = NewMyEmail(emailOption)
	//myEmail.SendOneEmail()

	log.Info("NewAlertHook")

	return alertHook
}

func SendSMS(){

}

func GetEmailInc(){

}

//func SendEmail(){
//
//}

func SendFeishu(){

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
	return myEmail
}

func(myEmail *MyEmail) SendOneEmail()error{
	addr := "78878296@qq.com"

	m := gomail.NewMessage()
	//m.SetHeader("From",addr + "<" + myEmail.EmailOption.FromEmail + ">")
	m.SetHeader("From",addr)
	m.SetHeader("Subject", "testGoLib")
	m.SetHeader("To", "mqzhifu@sina.com")
	m.SetBody("text/html", "rt")
	err := myEmail.Dialer.DialAndSend(m)
	return err
}