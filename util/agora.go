package util

import (
	"encoding/base64"
)

// 基于 Golang 实现的 HTTP 基本认证示例，使用 RTC 的服务端 RESTful API
func GetHTTPBaseAuth(customerKey string, customerSecret string) string {
	// 客户 ID
	//customerKey :=
	// 客户密钥
	//customerSecret := "a8b9fd618edb4061a7d8abd8f734ccaf"

	// 拼接客户 ID 和客户密钥并使用 base64 进行编码
	plainCredentials := customerKey + ":" + customerSecret
	base64Credentials := base64.StdEncoding.EncodeToString([]byte(plainCredentials))

	MyPrint("-------------------base64Credentials:", base64Credentials)
	return base64Credentials
	//
	//url := "https://api.agora.io/dev/v1/projects"
	//method := "GET"
	//
	//payload := strings.NewReader(``)
	//
	//client := &http.Client{}
	//req, err := http.NewRequest(method, url, payload)
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//// 增加 Authorization header
	//req.Header.Add("Authorization", "Basic "+base64Credentials)
	//req.Header.Add("Content-Type", "application/json")
	//
	//// 发送 HTTP 请求
	//res, err := client.Do(req)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//defer res.Body.Close()
	//
	//body, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(string(body))
}
