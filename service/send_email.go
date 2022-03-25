package service

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
	"zgoframe/http/request"
	"zgoframe/model"
	"zgoframe/util"
)

type SendEmail struct {
	Gorm *gorm.DB
}

func NewSendEmail(gorm *gorm.DB) *SendEmail {
	SendEmail := new(SendEmail)
	SendEmail.Gorm = gorm
	return SendEmail
}

func (SendEmail *SendEmail) Send(projectId int, info request.SendEmail) (err error) {
	if info.RuleId <= 0 || info.Receiver == "" || info.SendIp == "" || info.SendUid <= 0 {
		return errors.New("RuleId || Receiver || SendIp || SendUid is empty")
	}

	checkEmailRs := util.CheckEmailRule(info.Receiver)
	if !checkEmailRs {
		return errors.New("CheckEmailRule Receiver err" + info.Receiver)
	}

	var rule model.EmailRule
	err = SendEmail.Gorm.Where("id = ? ", info.RuleId).First(&rule).Error
	if err != nil {
		return errors.New("id not in db," + strconv.Itoa(info.RuleId))
	}

	//if rule.ProjectId != projectId {
	//	return errors.New("projectId != rule.projectId")
	//}

	err = SendEmail.CheckRule(rule)
	if err != nil {
		return err
	}

	err = SendEmail.CheckLimiterPeriod(rule, info.Receiver)
	if err != nil {
		return err
	}
	err = SendEmail.CheckLimiterDay(rule, info.Receiver)
	if err != nil {
		return err
	}

	content := SendEmail.ReplaceContentTemplate(rule.Content, info.ReplaceVar)
	//ProjectId  int    `json:"app_id" db:"define:tinyint(1);comment:项目ID;defaultValue:0"`           //项目ID

	carbonCopyStr := ""
	if len(info.CarbonCopy) > 0 {
		for _, v := range info.CarbonCopy {
			carbonCopyStr += v + " , "
		}
	}

	emailLog := model.EmailLog{
		ProjectId:  projectId,
		RuleId:     rule.Id,
		Receiver:   info.Receiver,
		CarbonCopy: carbonCopyStr,
		Content:    content,
		Title:      rule.Title,
		SendIp:     info.SendIp,
		SendUid:    info.SendUid,
	}

	if rule.Type == model.EMAIL_TYPE_AUTHCODE {
		if rule.ExpireTime <= 0 {

		}
		emailLog.ExpireTime = util.GetNowTimeSecondToInt() + rule.ExpireTime
		code := util.GetRandIntNum(9999)
		emailLog.AuthCode = strconv.Itoa(code)
	}

	SendEmail.Gorm.Create(emailLog)
	//global.V.Email.SendOneEmailAsync(info.Receiver, rule.Title, content)
	return nil

}

func (SendEmail *SendEmail) ReplaceContentTemplate(content string, replaceVar map[string]string) string {
	if len(replaceVar) <= 0 {
		return content
	}

	for k, v := range replaceVar {
		content = strings.Replace(content, "{"+k+"}", v, -1)
	}
	return content
}

func (SendEmail *SendEmail) CheckRule(rule model.EmailRule) error {
	if rule.Period < 60 {
		return errors.New("rule err , Period < 60")
	}

	if rule.DayTimes < 0 {
		return errors.New("rule err , rule.DayTimes < 0")
	}

	if rule.PeriodTimes < 0 {
		return errors.New("rule err , rule.PeriodTimes < 0")
	}

	if rule.Content == "" || rule.Title == "" {
		return errors.New("rule err , rule.Content || rule.Title empty")
	}

	return nil
}

func (SendEmail *SendEmail) CheckLimiterPeriod(rule model.EmailRule, receiver string) error {
	var count int64
	now := util.GetNowTimeSecondToInt()
	nowEnd := now + rule.Period
	where := "created_at >= ?  and  created_at <= ?  and receiver = ? and rule_id = ?"
	SendEmail.Gorm.Model(model.EmailLog{}).Where(where, now, nowEnd, receiver, rule.Id).Count(&count)

	if count > int64(rule.PeriodTimes) {

	}

	return nil
}

func (SendEmail *SendEmail) GetNowDayStartTime() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	return t.Unix()
}

func (SendEmail *SendEmail) CheckLimiterDay(rule model.EmailRule, receiver string) error {
	var count int64

	start := SendEmail.GetNowDayStartTime()
	end := start + 24*60*60 - 1
	where := "created_at >= ?  and  created_at <= ?  and receiver = ? and rule_id = ?"
	SendEmail.Gorm.Model(model.EmailLog{}).Where(where, start, end, receiver, rule.Id).Count(&count)

	if count > int64(rule.DayTimes) {

	}

	return nil
}
