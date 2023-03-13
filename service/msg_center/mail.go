package msg_center

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
	"zgoframe/http/request"
	"zgoframe/model"
	"zgoframe/service"
	"zgoframe/util"
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

func (mail *Mail) Send(info request.SendMail) (recordNewId int, err error) {
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

	//群发全部用户，不需要指定收件人，其它类型都需要指定收件人/组ID/tagId
	if rule.PeopleType != service.MAIL_PEOPLE_ALL {
		if info.Receiver == "" {
			return 0, errors.New("Receiver empty")
		}
	}
	//创建记录之前，先更新一下已失效的记录
	//mail.CheckExpireAndUpStatus()

	//定时：发送时间
	now := util.GetNowTimeSecondToInt()
	if info.SendTime > 0 {
		if info.SendTime <= now {
			return 0, errors.New("info.SendTime <= now")
		}
	}

	switch rule.PeopleType {
	case service.MAIL_PEOPLE_PERSON: //点对点
		recordNewId, err = mail.SendPerson(rule, info, 0)
	case service.MAIL_PEOPLE_ALL: //群发
		newInfo := info
		newInfo.Receiver = "all"
		recordNewId, err = mail.SendGroup(rule, info)
	case service.MAIL_PEOPLE_GROUP: //根据组，群发
		recordNewId, err = mail.SendGroup(rule, info)
	case service.MAIL_PEOPLE_TAG: //根据tag标签，群发
		recordNewId, err = mail.SendGroup(rule, info)
	case service.MAIL_PEOPLE_UIDS: //指定UID，群发
		uids := strings.Split(info.Receiver, ",")
		for _, v := range uids {
			newInfo := info
			newInfo.Receiver = v
			mail.SendPerson(rule, newInfo, 0)
		}
	default:
		return 0, errors.New("PeopleType error")
	}

	return recordNewId, err

}

////检测：已失效(未使用过)的 短信，并更新状态为：已失效
////验证码类型的字段sms_log中没直接存，所以 expire_time > 0 即可.
//func (mail *Mail) CheckExpireAndUpStatus() {
//	var mailLog model.MailLog
//	now := util.GetNowTimeSecondToInt()
//	upRsObj := mail.Gorm.Model(&mailLog).Where("expire_time > 0  and expire_time <=  ? and status = ?  ", now, model.AUTH_CODE_STATUS_NORMAL).Update("status", model.AUTH_CODE_STATUS_EXPIRE)
//	if upRsObj.Error != nil {
//		//if upRsObj.Error == gorm.ErrRecordNotFound {
//		//	util.MyPrint("CheckExpireAndUpStatus not record.")
//		//} else {
//		util.MyPrint("CheckExpireAndUpStatus gorm err:" + upRsObj.Error.Error())
//		//}
//	} else {
//		util.MyPrint("CheckExpireAndUpStatus up record RowsAffected:" + strconv.Itoa(int(upRsObj.RowsAffected)))
//	}
//}

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

	if rule.Content == "" || rule.Title == "" || rule.Type <= 0 {
		return errors.New("rule err , Content || Title   || Type empty")
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
	if rule.PeopleType == service.MAIL_PEOPLE_PERSON || rule.PeopleType == service.MAIL_PEOPLE_UIDS {
		//替换模板动态内容
		content = mail.ReplaceContentTemplate(rule.Content, info.ReplaceVar)
	}
	r, _ := strconv.Atoi(info.Receiver)
	if r == service.MAIL_ADMIN_USER_UID {
		return recordNewId, errors.New("Receiver Uid 不能等于管理员 ID")
	}
	//替换模板动态内容
	mailLog := model.MailLog{
		ProjectId:    rule.ProjectId,
		RuleId:       rule.Id,
		Receiver:     r,
		Content:      content,
		Title:        rule.Title,
		SendIp:       info.SendIp,
		SendUid:      info.SendUid,
		MailGroupId:  MailGroupId,
		SendTime:     info.SendTime,
		ReceiverRead: service.RECEIVER_READ_FALSE,
		ReceiverDel:  service.RECEIVER_DEL_FALSE,
	}
	//失效时间
	if rule.ExpireTime > 0 {
		mailLog.ExpireTime = util.GetNowTimeSecondToInt() + rule.ExpireTime
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
		SendTime: info.SendTime,
	}

	err = mail.Gorm.Create(&mailGroup).Error
	if err != nil {
		return 0, errors.New("gorm err:" + err.Error())
	}

	util.MyPrint("SendGroup new id:", mailGroup.Id, " content:", mailGroup.Content)
	return mailGroup.Id, err
}

//检测：是否有群发消息，需要生成到mail_log表中
func (mail *Mail) CheckGroupMsg(uid int) error {
	var mailGroupList []model.MailGroup
	//目前仅支持：全部群发
	where := " people_type = " + strconv.Itoa(service.MAIL_PEOPLE_ALL)
	where += " and " + mail.GetSqlWhereSendTime()
	err := mail.Gorm.Where(where).Find(&mailGroupList).Error
	if err != nil {
		return err
	}

	if len(mailGroupList) <= 0 {
		return errors.New("mailGroupList empty~")
	}

	var mailList []model.MailLog
	where = " receiver = '" + strconv.Itoa(uid) + "' " + " and mail_group_id > 0 "
	err = mail.Gorm.Where(where).Find(&mailList).Error
	if err != nil {
		return err
	}

	for _, group := range mailGroupList {
		exist := false
		for _, v := range mailList {
			if v.MailGroupId == group.Id {
				exist = true
			}
		}

		if !exist {
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

			mail.SendPerson(rule, info, group.Id)
		}
	}
	return nil

}

func (mail *Mail) GetSqlWhereExpire() string {
	now := util.GetNowTimeSecondToInt()
	where := "( expire_time <=0 or ( expire_time > 0  and expire_time <=  " + strconv.Itoa(now) + " ) )"
	return where
}

func (mail *Mail) GetSqlWhereSendTime() string {
	now := util.GetNowTimeSecondToInt()
	where := "( send_time <= 0 or ( send_time > 0  and send_time <=  " + strconv.Itoa(now) + " ) )"
	return where
}

func (mail *Mail) GetUserListByUid(uid int, form request.MailList) (mailList []model.MailLog, err error) {
	mail.CheckGroupMsg(uid)

	//var mailList []model.MailLog
	where := " 1 = 1 and "
	if form.BoxType == service.MAIL_IN_BOX { //收件箱
		where += " receiver = '" + strconv.Itoa(uid) + "' "
	} else if form.BoxType == service.MAIL_OUT_BOX { //发件箱
		where += " sendUid = " + strconv.Itoa(uid)
	} else if form.BoxType == service.MAIL_ALL_BOX { //全部
		where += " receiver = '" + strconv.Itoa(uid) + "' or " + "sendUid = " + strconv.Itoa(uid)
	} else {
		return mailList, errors.New("boxType err")
	}

	if form.ReceiverRead == service.RECEIVER_READ_TRUE { //接收者，已读
		where += " and receiver_read = " + strconv.Itoa(service.RECEIVER_READ_TRUE)
	} else if form.ReceiverRead == service.RECEIVER_READ_FALSE { //接收者，未读
		where += " and receiver_read = " + strconv.Itoa(service.RECEIVER_READ_FALSE)
	} else {
		return mailList, errors.New("ReceiverRead err")
	}

	if form.ReceiverDel == service.RECEIVER_DEL_TRUE { //接收者，已删除
		where += " and receiver_del = " + strconv.Itoa(service.RECEIVER_DEL_TRUE)
	} else if form.ReceiverRead == service.RECEIVER_DEL_FALSE { //接收者，未删除
		where += " and receiver_del = " + strconv.Itoa(service.RECEIVER_DEL_FALSE)
	} else {
		return mailList, errors.New("ReceiverDel err")
	}
	//未失效
	if form.Expire == service.MAIL_EXPIRE_FALSE {
		where += " and " + mail.GetSqlWhereExpire()
	}

	where += " and " + mail.GetSqlWhereSendTime()

	var count int64
	err = mail.Gorm.Model(model.MailLog{}).Where(where).Count(&count).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return mailList, nil
		}
	}

	err = mail.Gorm.Where(where).Find(&mailList).Error
	if err != nil {
		return mailList, err
	}

	return mailList, nil
}

func (mail *Mail) GetOneByUid(uid int, form request.MailInfo) (model.MailLog, error) {
	var mailLog model.MailLog
	//var mailLog model.MailLog
	err := mail.Gorm.First(&mailLog, form.Id).Error
	if err != nil {
		return mailLog, err
	}
	now := util.GetNowTimeSecondToInt()
	if mailLog.ExpireTime > 0 && mailLog.ExpireTime < now {
		return mailLog, errors.New("has expire")
	}

	if mailLog.SendTime > 0 && mailLog.SendTime < now {
		return mailLog, errors.New("未到发送时间")
	}

	if mailLog.Receiver != uid && mailLog.SendUid != uid {
		return mailLog, errors.New("该邮件不属于该用户")
	}

	if form.AutoReceiverRead == service.RECEIVER_READ_TRUE {
		var upDataModel model.MailLog
		upDataModel.Id = mailLog.Id
		mail.Gorm.Model(&upDataModel).Update("receiver_read", service.RECEIVER_READ_TRUE)
	}

	if form.AutoReceiverDel == service.RECEIVER_DEL_TRUE {
		var upDataModel model.MailLog
		upDataModel.Id = mailLog.Id
		mail.Gorm.Model(&upDataModel).Update("receiver_del", service.RECEIVER_DEL_TRUE)
	}

	return mailLog, nil
}

func (mail *Mail) GetUnreadCnt(uid int) int {
	where := "receiver = " + strconv.Itoa(uid) + " and receiver_read =  " + strconv.Itoa(service.RECEIVER_DEL_FALSE)
	where += " and " + mail.GetSqlWhereSendTime() + " and " + mail.GetSqlWhereExpire()
	var count int64
	err := mail.Gorm.Model(model.MailLog{}).Where(where).Count(&count).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0
		}
	}

	return int(count)
}
