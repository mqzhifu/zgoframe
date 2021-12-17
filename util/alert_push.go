package util

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type AlertPush struct {
	Ip	string
	Port string
	Uri string
	Url string
}

func NewAlertPush( ip string , port string, uri string)*AlertPush{
	alert := new(AlertPush)
	alert.Ip = ip
	alert.Port = port
	alert.Uri = uri
	url := "http://"+ip + ":" + port + "/" + uri
	alert.Url = url


	MyPrint("NewAlertPush:"+ alert.Url)

	return alert
}
type AlertMsg struct {
	Lables 		map[string]string
	Annotations map[string]string
	StartsAt string
	EndsAt string
	GeneratorURL string
}
//[
//	{
//		"labels": {"label": "value", ...},
//		"annotations": {"label": "value", ...},
//		"generatorURL": "string",
//		"startsAt": "2020-01-01T00:00:00.000+08:00", # optional
//		"endsAt": "2020-01-01T01:00:00.000+08:00" # optional
//	},
//	...
//]

func(alertPush *AlertPush) Push(alertMsg AlertMsg){
	//MyPrint("program has error,need push alert....")
	return

	str ,err := json.Marshal(alertMsg)

	req, err := http.NewRequest("POST", alertPush.Url, bytes.NewReader(str))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil{
		panic(err)
	}
}