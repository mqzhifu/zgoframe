package util

//直接推送报警(非3方)
import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type AlertPush struct {
	Ip   string
	Port string
	Uri  string
	Url  string
}

func NewAlertPush(ip string, port string, uri string, printfPrefix string) *AlertPush {
	alert := new(AlertPush)
	alert.Ip = ip
	alert.Port = port
	alert.Uri = uri
	url := "http://" + ip + ":" + port + "/" + uri
	alert.Url = url

	MyPrint(printfPrefix + "NewAlertPush:" + alert.Url)

	return alert
}

type AlertMsg struct {
	Content string `json:"-"`

	Labels       AlertMsgLabels      `json:"labels"`
	Annotations  AlertMsgAnnotations `json:"annotations"`
	StartsAt     string              `json:"startsAt"`
	EndsAt       string              `json:"endsAt"`
	GeneratorURL string              `json:"generatorURL"`
}

type AlertMsgAnnotations struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type AlertMsgLabels struct {
	Severity    string `json:"severity"`
	TriggerType string `json:"trigger_type"` //1 应用主动 2 被动监控
	ProjectId   string `json:"project_id"`   //3方仅支持string格式

	Alertname string `json:"alertname"`
	JobName   string `json:"job_name"`
	Instance  string `json:"instance"`
}

func (alertPush *AlertPush) Push(projectId int, levelString string, content string) {
	MyPrint("program has error,need push alert....")
	return

	alertMsgAnnotations := AlertMsgAnnotations{
		Summary:     content,
		Description: content,
	}

	alertMsgLabels := AlertMsgLabels{
		Severity:    levelString,
		TriggerType: "initiative",
		ProjectId:   strconv.Itoa(projectId),
		Alertname:   "serviceDIy",
		JobName:     "bbbb",
		Instance:    "127.0.0.1",
	}
	//RFC3339     = "2006-01-02T15:04:05Z07:00"
	//RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	now := time.Now()
	//nowString := now.Format("2006-01-02T15:04:05Z07:00")
	nowString := now.Format(time.RFC3339Nano)

	endTimeNow := now.Add(time.Second * 10)
	endTimeNowStr := endTimeNow.Format(time.RFC3339Nano)
	alertMsg := AlertMsg{
		Annotations:  alertMsgAnnotations,
		Labels:       alertMsgLabels,
		StartsAt:     nowString,
		EndsAt:       endTimeNowStr,
		GeneratorURL: "http://127.0.0.1/service/diy",
	}

	alertMsgArr := []AlertMsg{alertMsg}
	str, err := json.Marshal(alertMsgArr)

	req, err := http.NewRequest("POST", alertPush.Url, bytes.NewReader(str))
	// req.Header.Set("X-Custom-Header", "myvalue")

	MyPrint("alert push ,url:", alertPush.Url)
	MyPrint("alert push ,content:", string(str))
	//MyPrint("json err:",err , " http request err:",err)

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	//req.Header.Set("Content-Length", strconv.Itoa(len(str)))

	client := &http.Client{}
	//resp, err := client.Do(req)
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	//body,err:=ioutil.ReadAll(resp.Body)
	//MyPrint(resp.Status,string(body))

}
