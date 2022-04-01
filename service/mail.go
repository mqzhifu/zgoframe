package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
	"zgoframe/http/request"
	"zgoframe/model"
	"zgoframe/util"
)

const (
	//1单发2群发3指定group4指定tag5指定UIDS
	MAIL_PEOPLE_PERSON = 1
	MAIL_PEOPLE_ALL    = 2
	MAIL_PEOPLE_GROUP  = 3
	MAIL_PEOPLE_TAG    = 4
	MAIL_PEOPLE_UIDS   = 5
)

type Mail struct {
	Gorm *gorm.DB
	Log  *zap.Logger
}

func NewMail(gorm *gorm.DB, Log *zap.Logger) *Mail {
	mail := new(Mail)
	mail.Gorm = gorm
	mail.Log = Log
	return mail
}

func (mail *Mail) Send(projectId int, info request.SendMail) (recordNewId int, err error) {
	util.MyPrint("im in sendMail.send , formInfo:", info)
	if info.RuleId <= 0 || info.SendIp == "" || info.SendUid <= 0 {
		return 0, errors.New("RuleId  || SendIp || SendUid is empty")
	}

	var rule model.MailRule
	err = mail.Gorm.Where("id = ? ", info.RuleId).First(&rule).Error
	if err != nil {
		return 0, errors.New("id not in db," + strconv.Itoa(info.RuleId))
	}

	err = mail.CheckRule(rule)
	if err != nil {
		return 0, err
	}

	err = mail.CheckLimiterPeriod(rule, info.Receiver)
	if err != nil {
		return 0, err
	}
	err = mail.CheckLimiterDay(rule, info.Receiver)
	if err != nil {
		return 0, err
	}

	if rule.PeopleType != MAIL_PEOPLE_ALL {
		if info.Receiver == "" {
			return 0, errors.New("Receiver empty")
		}
	}

	//创建记录之前，先更新一下已失效的记录
	mail.CheckExpireAndUpStatus()

	switch rule.PeopleType {
	case MAIL_PEOPLE_PERSON:
		recordNewId, err = mail.SendPerson(rule, info, 0)
	case MAIL_PEOPLE_ALL:
		newInfo := info
		newInfo.Receiver = "all"
		recordNewId, err = mail.SendGroup(rule, info)
	case MAIL_PEOPLE_GROUP:
		recordNewId, err = mail.SendGroup(rule, info)
	case MAIL_PEOPLE_TAG:
		recordNewId, err = mail.SendGroup(rule, info)
	case MAIL_PEOPLE_UIDS:
		uids := strings.Split(info.Receiver, ",")
		for _, v := range uids {
			newInfo := info
			newInfo.Receiver = v
			mail.SendPerson(rule, newInfo, projectId)
		}
	default:
		return 0, errors.New("PeopleType error")
	}

	return recordNewId, err

}

//检测：已失效(未使用过)的 短信，并更新状态为：已失效
//验证码类型的字段sms_log中没直接存，所以 expire_time > 0 即可.
func (mail *Mail) CheckExpireAndUpStatus() {
	var mailLog model.MailLog
	now := util.GetNowTimeSecondToInt()
	upRsObj := mail.Gorm.Model(&mailLog).Where("expire_time > 0  and expire_time <=  ? and status = ?  ", now, model.AUTH_CODE_STATUS_NORMAL).Update("status", model.AUTH_CODE_STATUS_EXPIRE)
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

func (mail *Mail) ReplaceContentTemplate(content string, replaceVar map[string]string) string {
	if len(replaceVar) <= 0 {
		return content
	}

	for k, v := range replaceVar {
		content = strings.Replace(content, "{"+k+"}", v, -1)
	}
	return content
}

func (mail *Mail) CheckRule(rule model.MailRule) error {
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

	if rule.PeopleType <= 0 {
		return errors.New("rule err , PeopleType empty")
	}

	return nil
}

//检查一定周期内，发送次数
func (mail *Mail) CheckLimiterPeriod(rule model.MailRule, receiver string) error {
	var count int64
	now := util.GetNowTimeSecondToInt()
	nowBefore := now - rule.Period
	where := "created_at >= ?  and  created_at <= ?  and receiver = ? and rule_id = ?"
	err := mail.Gorm.Model(model.MailLog{}).Where(where, nowBefore, now, receiver, rule.Id).Count(&count).Error
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
func (mail *Mail) CheckLimiterDay(rule model.MailRule, receiver string) error {
	var count int64

	start := GetNowDayStartTime()
	end := start + 24*60*60 - 1
	where := "created_at >= ?  and  created_at <= ?  and receiver = ? and rule_id = ?"
	err := mail.Gorm.Model(model.MailLog{}).Where(where, start, end, receiver, rule.Id).Count(&count).Error
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

//取出：当前天的起始时间  2022-02-25 00:00:00
func (mail *Mail) GetNowDayStartTime() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	return t.Unix()
}

func (mail *Mail) SendPerson(rule model.MailRule, info request.SendMail, MailGroupId int) (recordNewId int, err error) {
	content := ""
	if rule.PeopleType == MAIL_PEOPLE_PERSON || rule.PeopleType == MAIL_PEOPLE_UIDS {
		//替换模板动态内容
		content = mail.ReplaceContentTemplate(rule.Content, info.ReplaceVar)
	}
	//替换模板动态内容
	mailLog := model.MailLog{
		ProjectId:   rule.ProjectId,
		RuleId:      rule.Id,
		Receiver:    info.Receiver,
		Content:     content,
		Title:       rule.Title,
		SendIp:      info.SendIp,
		SendUid:     info.SendUid,
		MailGroupId: MailGroupId,
	}

	err = mail.Gorm.Create(&mailLog).Error
	if err != nil {
		return 0, errors.New("gorm err:" + err.Error())
	}

	util.MyPrint("mailLog new id:", mailLog.Id, " content:", mailLog.Content)
	return mailLog.Id, err
}

func (mail *Mail) SendGroup(rule model.MailRule, info request.SendMail) (recordNewId int, err error) {
	//替换模板动态内容
	content := mail.ReplaceContentTemplate(rule.Content, info.ReplaceVar)
	mailGroup := model.MailGroup{
		RuleId:   rule.Id,
		Receiver: info.Receiver,
		Content:  content,
		Title:    rule.Title,
		SendIp:   info.SendIp,
		SendUid:  info.SendUid,
	}

	err = mail.Gorm.Create(&mailGroup).Error
	if err != nil {
		return 0, errors.New("gorm err:" + err.Error())
	}

	util.MyPrint("SendGroup new id:", mailGroup.Id, " content:", mailGroup.Content)
	return mailGroup.Id, err
}

func (mail *Mail) CheckGroupMsg(uid int) {
	var mailGroupList []model.MailGroup
	err := mail.Gorm.Where("people_type = ?", MAIL_PEOPLE_ALL).Find(&mailGroupList).Error
	if err != nil {

	}

	//ids := ""
	//for _, v := range mailGroupList {
	//	ids = strconv.Itoa(v.Id) + " , "
	//}

	var mailList []model.MailLog
	where := "receiver = '" + strconv.Itoa(uid) + "' " + " and mail_group_id > 0 "
	err = mail.Gorm.Where(where).Find(&mailList).Error
	if err != nil {

	}

	for _, group := range mailGroupList {
		exist := false
		for _, v := range mailList {
			if v.MailGroupId == group.Id {
				exist = true
			}
		}
		rule := model.MailRule{
			Title:      group.Title,
			Content:    group.Content,
			PeopleType: group.PeopleType,
		}
		info := request.SendMail{
			SendUid:  group.SendUid,
			SendIp:   group.SendIp,
			Receiver: strconv.Itoa(uid),
		}
		if !exist {
			mail.SendPerson(rule, info, group.Id)
		}
	}

}

func (mail *Mail) GetUserListByUid(uid int, actionType int, readType int) error {
	mail.CheckGroupMsg(uid)

	var mailList []model.MailLog
	where := ""
	if actionType == 1 { //收件箱
		where = "receiver = '" + strconv.Itoa(uid) + "' "
	} else if actionType == 2 { //发件箱
		where = "sendUid = " + strconv.Itoa(uid)
	} else if actionType == 3 { //全部
		where = "receiver = '" + strconv.Itoa(uid) + "' or " + "sendUid = " + strconv.Itoa(uid)
	} else {
		return errors.New("actionType err")
	}

	if readType > 0 {
		where += "receiver_read = " + strconv.Itoa(readType)
	}

	where += "receiver_del = 0"

	err := mail.Gorm.Where(where).Find(&mailList).Error
	if err != nil {

	}
	return nil
}

func (mail *Mail) GetOneByUid(uid int, id int, autoRead int) (mailLog model.MailLog) {
	//var mailLog model.MailLog
	err := mail.Gorm.Find(&mailLog, id).Error
	if err != nil {

	}

	if autoRead == 1 {
		mail.Gorm.Update("receiver_read", 1)
	}

	return mailLog
}
