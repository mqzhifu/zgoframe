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

type SendSms struct {
	Gorm *gorm.DB
}

func NewSendSms(gorm *gorm.DB) *SendSms {
	sendSms := new(SendSms)
	sendSms.Gorm = gorm
	return sendSms
}

func (SendSms *SendSms) Send(projectId int, info request.SendSMS) (err error) {
	util.MyPrint("im in sendsms.send , formInfo:", info)
	if info.RuleId <= 0 || info.Receiver == "" || info.SendIp == "" || info.SendUid <= 0 {
		return errors.New("RuleId || Receiver || SendIp || SendUid is empty")
	}

	checkMobileRs := util.CheckMobileRule(info.Receiver)
	if !checkMobileRs {
		return errors.New("checkMobileRs Receiver err" + info.Receiver)
	}

	var rule model.SmsRule
	err = SendSms.Gorm.Where("id = ? ", info.RuleId).First(&rule).Error
	util.MyPrint(err)
	if err != nil {
		return errors.New("id not in db," + strconv.Itoa(info.RuleId))
	}

	err = SendSms.CheckRule(rule)
	if err != nil {
		return err
	}

	err = SendSms.CheckLimiterPeriod(rule, info.Receiver)
	if err != nil {
		return err
	}
	err = SendSms.CheckLimiterDay(rule, info.Receiver)
	if err != nil {
		return err
	}

	content := SendSms.ReplaceContentTemplate(rule.Content, info.ReplaceVar)
	//ProjectId  int    `json:"app_id" db:"define:tinyint(1);comment:项目ID;defaultValue:0"`           //项目ID

	//util.ExitPrint("content:", content)
	smsLog := model.SmsLog{
		ProjectId: projectId,
		RuleId:    rule.Id,
		Receiver:  info.Receiver,
		Content:   content,
		Title:     rule.Title,
		SendIp:    info.SendIp,
		SendUid:   info.SendUid,
	}

	if rule.Type == model.EMAIL_TYPE_AUTHCODE {
		if rule.ExpireTime <= 0 {
			return errors.New("rule.ExpireTime < 0 ")
		}

		smsLog.ExpireTime = util.GetNowTimeSecondToInt() + rule.ExpireTime
		code := util.GetRandIntNum(9999)
		smsLog.AuthCode = strconv.Itoa(code)
		smsLog.AuthStatus = 1

		content = strings.Replace(smsLog.Content, "{auth_code}", smsLog.AuthCode, -1)
		content = strings.Replace(content, "{auth_expire_time}", strconv.Itoa(rule.ExpireTime), -1)
		smsLog.Content = content
	}

	SendSms.Gorm.Create(&smsLog)
	util.MyPrint("smsLog new id:", smsLog.Id)
	//global.V.Email.SendOneEmailAsync(info.Receiver, rule.Title, content)
	return nil

}

func (SendSms *SendSms) ReplaceContentTemplate(content string, replaceVar map[string]string) string {
	if len(replaceVar) <= 0 {
		return content
	}

	for k, v := range replaceVar {
		content = strings.Replace(content, "{"+k+"}", v, -1)
	}
	return content
}

func (SendSms *SendSms) CheckRule(rule model.SmsRule) error {
	if rule.Period < 30 {
		return errors.New("rule err , Period < 30")
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

func (SendSms *SendSms) CheckLimiterPeriod(rule model.SmsRule, receiver string) error {
	var count int64
	now := util.GetNowTimeSecondToInt()
	nowBefore := now - rule.Period
	where := "created_at >= ?  and  created_at <= ?  and receiver = ? and rule_id = ?"
	SendSms.Gorm.Model(model.SmsLog{}).Where(where, nowBefore, now, receiver, rule.Id).Count(&count)

	util.MyPrint("CheckLimiterPeriod count:", count, " PeriodTimes:", rule.PeriodTimes)
	if count >= int64(rule.PeriodTimes) {
		return errors.New("PeriodTimes err : " + strconv.Itoa(int(count)) + "  > " + strconv.Itoa(rule.PeriodTimes))
	}

	return nil
}

func (SendSms *SendSms) GetNowDayStartTime() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	return t.Unix()
}

func (SendSms *SendSms) CheckLimiterDay(rule model.SmsRule, receiver string) error {
	var count int64

	start := SendSms.GetNowDayStartTime()
	end := start + 24*60*60 - 1
	where := "created_at >= ?  and  created_at <= ?  and receiver = ? and rule_id = ?"
	SendSms.Gorm.Model(model.SmsLog{}).Where(where, start, end, receiver, rule.Id).Count(&count)

	if count > int64(rule.DayTimes) {
		return errors.New("DayTimes err : " + strconv.Itoa(int(count)) + "  > " + strconv.Itoa(rule.DayTimes))
	}

	return nil
}

func (SendSms *SendSms) Verify(ruleId int, mobile string, authCode string) error {
	var smsLog model.SmsLog
	now := util.GetNowTimeSecondToInt()
	util.MyPrint("RegisterSms now:", now)
	err := SendSms.Gorm.Where("receiver = ?   and rule_id = ?  auth_status = 1 ", mobile, ruleId).First(&smsLog).Error
	if err != nil {
		return errors.New("未检查出发送过短信...")
	}

	if smsLog.AuthCode != authCode {
		return errors.New("验证码错误......")
	}

	if smsLog.ExpireTime < now {
		var editSmsLog model.SmsLog
		editSmsLog.Id = smsLog.Id
		editSmsLog.AuthStatus = 3
		SendSms.Gorm.Updates(editSmsLog)

		return errors.New("已失效...(记录已变更状态:已失效)")
	}

	var smsLogEdit model.SmsLog
	smsLogEdit.Id = smsLog.Id
	smsLogEdit.AuthStatus = 2
	SendSms.Gorm.Updates(smsLogEdit)
	return nil
}
