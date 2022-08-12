package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	HTTP_DATA_CONTENT_TYPE_JSON   = 1
	HTTP_DATA_CONTENT_TYPE_NOrMAL = 2
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

func (httpCurl *HttpCurl) Get() (res string, err error) {
	return httpCurl.Curl(2, "")
}

func (httpCurl *HttpCurl) Post(data string) (res string, err error) {
	return httpCurl.Curl(1, data)
}

func (httpCurl *HttpCurl) PostJson(data interface{}) (res string, err error) {
	dataBytes, err := json.Marshal(data)
	//MyPrint("dataBytes:", string(dataBytes))
	if err != nil {
		MyPrint(httpCurl.Prefix + "json.Marshal err:")
		return res, err
	}
	dataStr := string(dataBytes)
	return httpCurl.Curl(1, dataStr)
}

func (httpCurl *HttpCurl) Curl(method int, data string) (res string, err error) {
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
		return res, errors.New(httpCurl.Prefix + err.Error())
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	res = string(bodyBytes)
	MyPrint(httpCurl.Prefix+" response read body:", res, " err:", err)

	return res, nil
}
