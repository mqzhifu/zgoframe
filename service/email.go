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

type Email struct {
	Gorm    *gorm.DB
	MyEmail *util.MyEmail
}

func NewEmail(gorm *gorm.DB, myEmail *util.MyEmail) *Email {
	sendEmail := new(Email)
	sendEmail.Gorm = gorm
	sendEmail.MyEmail = myEmail
	return sendEmail
}

func (email *Email) Send(projectId int, info request.SendEmail) (recordNewId int, err error) {
	util.MyPrint("im in sendEmail.send , formInfo:", info)
	if info.RuleId <= 0 || info.Receiver == "" || info.SendIp == "" || info.SendUid <= 0 {
		return 0, errors.New("RuleId || Receiver || SendIp || SendUid is empty")
	}

	checkMobileRs := util.CheckEmailRule(info.Receiver)
	if !checkMobileRs {
		return 0, errors.New("Receiver format err： " + info.Receiver)
	}

	var rule model.EmailRule
	err = email.Gorm.Where("id = ? ", info.RuleId).First(&rule).Error
	if err != nil {
		return 0, errors.New("id not in db," + strconv.Itoa(info.RuleId))
	}

	err = email.CheckRule(rule)
	if err != nil {
		return 0, err
	}

	err = email.CheckLimiterPeriod(rule, info.Receiver)
	if err != nil {
		return 0, err
	}
	err = email.CheckLimiterDay(rule, info.Receiver)
	if err != nil {
		return 0, err
	}
	//替换模板动态内容
	content := email.ReplaceContentTemplate(rule.Content, info.ReplaceVar)
	emailLog := model.EmailLog{
		ProjectId: projectId,
		RuleId:    rule.Id,
		Receiver:  info.Receiver,
		Content:   content,
		Title:     rule.Title,
		SendIp:    info.SendIp,
		SendUid:   info.SendUid,
	}
	//如果是验证码类型，要SERVER端生成CODE，并替换到模板中
	if rule.Type == model.RULE_TYPE_AUTH_CODE {
		//验证码必须得有失效时间
		if rule.ExpireTime <= 0 {
			return 0, errors.New("rule.ExpireTime <= 0 ，验证码类型短信，必须得有失效时间")
		}
		//当前时间 + 失效时间
		emailLog.ExpireTime = util.GetNowTimeSecondToInt() + rule.ExpireTime
		//验证码
		code := util.GetRandIntNum(9999)
		emailLog.AuthCode = strconv.Itoa(code)
		//状态
		emailLog.AuthStatus = model.AUTH_CODE_STATUS_NORMAL
		//把刚刚生成的code替换到内容中
		content = strings.Replace(emailLog.Content, "{auth_code}", emailLog.AuthCode, -1)
		content = strings.Replace(content, "{auth_expire_time}", strconv.Itoa(rule.ExpireTime), -1)
		emailLog.Content = content
	}
	//创建记录之前，先更新一下已失效的记录
	email.CheckExpireAndUpStatus()
	err = email.Gorm.Create(&emailLog).Error
	if err != nil {
		return 0, errors.New("gorm err:" + err.Error())
	}

	util.MyPrint("emailLog new id:", emailLog.Id, " content:", emailLog.Content)
	email.MyEmail.SendOneEmailAsync(info.Receiver, emailLog.Title, emailLog.Content)
	return emailLog.Id, nil

}

//检测：已失效(未使用过)的 短信，并更新状态为：已失效
//验证码类型的字段sms_log中没直接存，所以 expire_time > 0 即可.
func (email *Email) CheckExpireAndUpStatus() {
	var emailLog model.EmailLog
	now := util.GetNowTimeSecondToInt()
	upRsObj := email.Gorm.Model(&emailLog).Where("expire_time > 0  and expire_time <=  ? and status = ?  ", now, model.AUTH_CODE_STATUS_NORMAL).Update("status", model.AUTH_CODE_STATUS_EXPIRE)
	if upRsObj.Error != nil {
		//if upRsObj.Error == gorm.ErrRecordNotFound {
		//	util.MyPrint("CheckExpireAndUpStatus not record.")
		//} else {
		util.MyPrint("CheckExpireAndUpStatus gorm err:" + upRsObj.Error.Error())
		//}
	} else {
		util.MyPrint("CheckExpireAndUpStatus up record RowsAffected:" + strconv.Itoa(int(upRsObj.RowsAffected)))
	}
}

func (email *Email) ReplaceContentTemplate(content string, replaceVar map[string]string) string {
	if len(replaceVar) <= 0 {
		return content
	}

	for k, v := range replaceVar {
		content = strings.Replace(content, "{"+k+"}", v, -1)
	}
	return content
}

func (email *Email) CheckRule(rule model.EmailRule) error {
	if rule.Period < model.RULE_PERIORD_MIN {
		return errors.New("rule err : 最小频率周期-时间(秒):" + strconv.Itoa(model.RULE_PERIORD_MIN))
	}

	if rule.DayTimes <= 0 {
		return errors.New("rule err , 每天发送总次数 <= 0")
	}

	if rule.PeriodTimes <= 0 {
		return errors.New("rule err , 最小频率周期-发送次数 < 0")
	}

	if rule.Content == "" || rule.Title == "" || rule.Purpose <= 0 || rule.Type <= 0 {
		return errors.New("rule err , Content || Title || Purpose || Type empty")
	}

	return nil
}

//检查一定周期内，发送次数
func (email *Email) CheckLimiterPeriod(rule model.EmailRule, receiver string) error {
	var count int64
	now := util.GetNowTimeSecondToInt()
	nowBefore := now - rule.Period
	where := "created_at >= ?  and  created_at <= ?  and receiver = ? and rule_id = ?"
	err := email.Gorm.Model(model.EmailLog{}).Where(where, nowBefore, now, receiver, rule.Id).Count(&count).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("gorm err:" + err.Error())
		}
	}

	util.MyPrint("CheckLimiterPeriod count:", count, " PeriodTimes:", rule.PeriodTimes)
	if count >= int64(rule.PeriodTimes) {
		return errors.New("PeriodTimes err : " + strconv.Itoa(int(count)) + "  >= " + strconv.Itoa(rule.PeriodTimes))
	}

	return nil
}

//检查一天内，发送的总次数
func (email *Email) CheckLimiterDay(rule model.EmailRule, receiver string) error {
	var count int64

	start := GetNowDayStartTime()
	end := start + 24*60*60 - 1
	where := "created_at >= ?  and  created_at <= ?  and receiver = ? and rule_id = ?"
	err := email.Gorm.Model(model.EmailLog{}).Where(where, start, end, receiver, rule.Id).Count(&count).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("gorm err:" + err.Error())
		}
	}

	if count > int64(rule.DayTimes) {
		return errors.New("DayTimes err : " + strconv.Itoa(int(count)) + "  >= " + strconv.Itoa(rule.DayTimes))
	}

	return nil
}

//对 短信 验证码 进行校验
func (email *Email) Verify(ruleId int, mobile string, authCode string) error {
	if ruleId <= 0 || mobile == "" || authCode == "" {
		return errors.New("ruleId   || mobile   || authCode is empty")
	}
	var emailLog model.EmailLog
	now := util.GetNowTimeSecondToInt()
	util.MyPrint("SMS-Verify: ruleId", ruleId, "  mobile:", mobile, " auchCode:", authCode)

	checkMobileRs := util.CheckMobileRule(mobile)
	if !checkMobileRs {
		return errors.New("Receiver format err： " + mobile)
	}

	var rule model.EmailRule
	err := email.Gorm.First(&rule, ruleId).Error
	if err != nil {
		return errors.New("check Rule:" + err.Error())
	}

	//err := SendSms.Gorm.Where("receiver = ? and rule_id = ? and  auth_status = ？", mobile, ruleId, AUTH_CODE_STATUS_NORMAL).First(&smsLog).Error
	//err := SendSms.Gorm.First(&smsLog, "receiver = ? and rule_id = ? and  auth_status = ？", mobile, ruleId, 1).Error
	err = email.Gorm.Where("receiver = ?  and rule_id = ? and auth_status = ? ", mobile, ruleId, model.AUTH_CODE_STATUS_NORMAL).First(&emailLog).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未检查出发送过短信...")
		} else {
			return errors.New("Verify gorm err:" + err.Error())
		}
	}

	if emailLog.AuthCode != authCode {
		return errors.New("验证码错误......")
	}

	if emailLog.ExpireTime < now {
		var editEmailLog model.EmailLog
		editEmailLog.Id = emailLog.Id
		editEmailLog.AuthStatus = model.AUTH_CODE_STATUS_EXPIRE
		email.Gorm.Updates(editEmailLog)

		return errors.New("已失效...(记录已变更状态:已失效)")
	}

	var emailLogEdit model.EmailLog
	emailLogEdit.Id = emailLog.Id
	emailLogEdit.AuthStatus = model.AUTH_CODE_STATUS_OK
	upRsObj := email.Gorm.Updates(emailLogEdit)
	if upRsObj.Error != nil {
		//if !errors.Is(upRsObj.Error, gorm.ErrRecordNotFound) {
		return errors.New("Verify gorm err:" + upRsObj.Error.Error())
		//} else {
		//	return errors.New("Verify gorm search not found:" + strconv.Itoa(smsLog.Id))
		//}
	}

	util.MyPrint("RowsAffected", upRsObj.RowsAffected)

	return nil
}

//取出：当前天的起始时间  2022-02-25 00:00:00
func (email *Email) GetNowDayStartTime() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	return t.Unix()
}
