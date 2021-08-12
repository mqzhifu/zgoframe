package util

import (
	"encoding/json"
	"bytes"
	"net/http"
)

type Alert struct {
	Ip	string
	Port string
	Uri string
	Url string
}

func NewAlert( ip string , port string, uri string)*Alert{
	alert := new(Alert)
	alert.Ip = ip
	alert.Port = port
	alert.Uri = uri
	url := "http://"+ip + ":" + port + "/" + uri
	alert.Url = url

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

func(alert *Alert) Push(alertMsg AlertMsg){
	str ,err := json.Marshal(alertMsg)

	req, err := http.NewRequest("POST", alert.Url, bytes.NewReader(str))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil{
		panic(err)
	}
}