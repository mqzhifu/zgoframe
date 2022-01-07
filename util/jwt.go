package util

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

func CreateJwtToken(secretKey string,payload JwtDataPayload)string{
	header := JwtDataHeader{
		Alg: "HS256",
		Typ:"JWT",
	}

	headerJson,_ := json.Marshal(header)
	payloadJson ,_ := json.Marshal(payload)
	//fmt.Println("json : ",string(headerJson),string(payloadJson))

	base64HeaderJson := EncodeSegment(headerJson)
	base64PayloadJson := EncodeSegment(payloadJson)

	base64HeaderPayload := base64HeaderJson + "." + base64PayloadJson

	//fmt.Println("base64HeaderPayload : ",base64HeaderPayload)
	hasher := hmac.New(crypto.SHA256.New , []byte(secretKey))
	hasher.Write([]byte(base64HeaderPayload))

	sign := hasher.Sum(nil)

	base64Sign :=  EncodeSegment(sign)
	//fmt.Println(  " base64Sign : " , base64Sign)
	jwtString := base64HeaderPayload + "." + base64Sign
	//fmt.Println("myself : ",jwtString)


	return jwtString
}

//3方包创建一个TOKEN
func JwtGoCreateToken(secretKey string,payload JwtDataPayload)string{
	type jwtCustomClaims struct {
		jwt.StandardClaims
		//Id 			int
		Uid 		int32
		Expire 		int32
		ATime		int32
		AppId		int32
		Username	string
		// 追加自己需要的信息
		//Uid   uint `json:"uid"`
		//Admin bool `json:"admin"`
	}

	claims := &jwtCustomClaims{
		//StandardClaims: jwt.StandardClaims{
		//	ExpiresAt: int64(time.Now().Add(time.Hour * 72).Unix()),
		//},
		//Id:payload.Id,
		Uid :payload.Uid,
		Expire: payload.Expire,
		ATime: payload.ATime,
		AppId: payload.AppId,
		//Username: payload.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString
}

func ParseJwtToken(secretKey string,tokenStr string)(data JwtData,err error){
	CheckStrEmptyRs := CheckStrEmpty(tokenStr)
	if CheckStrEmptyRs{
		return data,NewCoder(500,"tokenStr is empty")
	}
	tokenStr = strings.Trim(tokenStr," ")
	tokenArr := strings.Split(tokenStr,".")
	if len(tokenArr) != 3{
		return data,NewCoder(501,"tokenStr Split by <.> != 3")
	}

	//fmt.Println(tokenArr)
	headerBase64 := tokenArr[0]
	headerStruct := JwtDataHeader{}
	err = decodeJsonByBase64(headerBase64,&headerStruct)
	if err != nil {
		return data,err
	}
	fmt.Println("headerStruct ",headerStruct)

	payloadBase64 := tokenArr[1]
	payloadStruct := JwtDataPayload{}
	err = decodeJsonByBase64(payloadBase64,&payloadStruct)
	if err != nil {
		return data,err
	}
	fmt.Println("payloadStruct  ",payloadStruct)

	sign := tokenArr[2]

	checkSign := checkJWTSign(headerStruct,payloadStruct,tokenArr,secretKey)
	if !checkSign{
		return data,NewCoder(511,"check sign is err")
	}
	data = JwtData{
		Header: headerStruct,
		Payload: payloadStruct,
		Sign: sign,
	}

	return data,nil
}
//此函数必须依附上面的函数，也就是把JSON都解到struct里，是一个正常的JSON结构体
func checkJWTSign(headerStruct JwtDataHeader,payloadStruct JwtDataPayload,tokenArr []string,secretKey string)bool{
	//MyPrint("checkJWTSign : ",tokenArr)

	hasher := hmac.New(crypto.SHA256.New , []byte(secretKey))
	hasher.Write([]byte(tokenArr[0]+ "." + tokenArr[1]))

	sign := hasher.Sum(nil)
	base64Sign :=  EncodeSegment(sign)
	//MyPrint("checkJWTSign base64Sign : ",base64Sign)
	if base64Sign != tokenArr[2]{
		return false
	}
	//这里还要再严谨一下，把解出来的数据，再重新生成一下jwt token ,看看对不对
	//不过，这里只是把新生成 的payload 传了进入 header 并没有验证
	newToken := CreateJwtToken(secretKey,payloadStruct)
	newTokenArr := strings.Split(newToken,".")
	//fmt.Println(newTokenArr)
	if tokenArr[2] != newTokenArr[2]{
		return false
	}
	return true
}

func decodeJsonByBase64(base64Str string ,dataStruct interface{})error{
	JsonStr,err := DecodeSegment(base64Str)
	if err != nil{
		return NewCoder(502,"decode base64 is err:"+err.Error())
	}

	if  len(JsonStr) == 0{
		return NewCoder(502,"decode base64 len is 0 ")
	}

	err = json.Unmarshal(JsonStr,&dataStruct)
	if err != nil{
		return NewCoder(502,"Unmarshal json err"+err.Error())
	}

	return nil
}

func EncodeSegment(seg []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(seg), "=")
}

func DecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}


func ParseToken(tokenSrt string, SecretKey []byte) (claims jwt.Claims, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenSrt, func(*jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	claims = token.Claims
	return
}

type JwtData struct {
	Header 	JwtDataHeader
	Payload	JwtDataPayload
	Sign 	string
}

type JwtDataHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

//type JwtDataPayload struct {
//	Id 			int
//	Expire 		int
//	ATime		int
//	AppId		int
//	Username	string
//}

type JwtDataPayload struct {
	Uid			int32
	Expire 		int32
	ATime		int32
	AppId		int32
}
