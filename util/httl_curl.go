package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type HttpCurl struct {
	Header              map[string]string
	RequestContentType  int
	ResponseContentType int
	Url                 string
	Prefix              string
}

func NewHttpCurl(url string, header map[string]string) *HttpCurl {
	httpCurl := new(HttpCurl)
	httpCurl.Url = url
	httpCurl.Header = header
	httpCurl.Prefix = "httpCurl "
	return httpCurl
}

func (httpCurl *HttpCurl) Get() (httpCode int, body string, err error) {
	return httpCurl.Curl(2, "")
}

func (httpCurl *HttpCurl) Post(data string) (httpCode int, body string, err error) {
	return httpCurl.Curl(1, data)
}

func (httpCurl *HttpCurl) PostJson(data interface{}) (httpCode int, body string, err error) {
	dataBytes, err := json.Marshal(data)
	//MyPrint("dataBytes:", string(dataBytes))
	if err != nil {
		MyPrint(httpCurl.Prefix + "json.Marshal err:")
		return httpCode, body, err
	}
	dataStr := string(dataBytes)
	return httpCurl.Curl(1, dataStr)
}

func (httpCurl *HttpCurl) Curl(method int, data string) (httpCode int, body string, err error) {
	MyPrint(httpCurl.Prefix+" url:", httpCurl.Url, " data:", data)
	client := &http.Client{}
	var request *http.Request
	if method == 1 {
		request, _ = http.NewRequest("POST", httpCurl.Url, strings.NewReader(data))
	} else {
		request, _ = http.NewRequest("GET", httpCurl.Url, nil)
	}

	if len(httpCurl.Header) > 0 {
		for k, v := range httpCurl.Header {
			request.Header.Add(k, v)
		}
	}

	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		return httpCode, body, errors.New(httpCurl.Prefix + err.Error())
	}
	MyPrint(httpCurl.Prefix + " res status code:" + strconv.Itoa(response.StatusCode))
	//if response.StatusCode != 200 {
	//	return res, errors.New("status != 200")
	//}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, body, errors.New("ioutil.ReadAll(response.Body) err:" + err.Error())
	}

	body = string(bodyBytes)
	MyPrint(httpCurl.Prefix+" response read body:", body, " err:", err)

	return response.StatusCode, body, nil
}
