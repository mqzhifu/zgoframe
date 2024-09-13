package core

import (
	"reflect"
	"zgoframe/core/global"
	"zgoframe/http/request"
	"zgoframe/test"
	"zgoframe/util"
)

func DoMySelf() {

}

// test_command()
// Gateway()
func DoTestAction(flag string) {
	//util.GetHTTPBaseAuth()
	util.MyPrint("=======================", flag)
	switch flag {
	case "alert":
		global.V.Service.Alert.Send(6, "商品库存不足，请及时补充货源", "warning")
	case "db_table":
		sqlList := global.AutoCreateUpDbTable()
		sqlStrings := ""
		for _, v := range sqlList {
			sqlStrings += v
		}
		util.MyPrint(sqlStrings)
		util.ExitPrint("i want exit 2.")
	case "ali_sms":
		//AliSms := util.NewAliSms()
		//AliSms.Send()
	case "email":
		test.Email()
	case "sms":
		test.Sms()
	case "service_sms":
		SendSMS := request.SendSMS{}
		recordNewId, err := global.V.Service.Sms.Send(6, SendSMS)
		util.MyPrint(recordNewId, err)
	case "cicd":
		test.Cicd()
	case "ProjectAutoCreateUserDbRecord":
		test.ProjectAutoCreateUserDbRecord()
	case "alert_push":
		//global.V.AlertPush.Push(1, "error", "test push alert info.")
	case "grpc":
		test.Grpc()
	case "super_visor":
		test.Test_command()
		test.Test_supervisor()
	case "gateway":
		test.Gateway()
	case "image_slice":
	case "calc":
		dd := "aaaaa"
		ttDD := reflect.TypeOf(dd)
		vvDD := reflect.ValueOf(dd)

		util.MyPrint(ttDD, vvDD)

		//slice1 := []int{0, 1, 2, 3}
		//testSlice(slice1)
		//util.MyPrint(slice1)

		//slice := []int{0, 1, 2, 3}
		//fmt.Printf("slice: %v \n", slice)
		//
		//changeSlice(slice)
		//fmt.Printf("slice: %v\n", slice)

		//a := []int{1, 2, 3, 4, 5}
		//changeSlice(a)
		//fmt.Println(a)

		//time.NewTimer(2 * time.Second)
		//
		//var slice []int
		//var a [3]int
		//util.MyPrint((slice == nil))

		//var b [3]int
		//a := [3]int{}
		//util.MyPrint(a, b)
		//a := [4]int{1, 2, 3, 4}
		//b := []int{1, 2, 3, 4}
		//
		//fmt.Printf("%p %T %v %d %d \n", &a, a, a, len(a), cap(a))
		//fmt.Printf("%p %T %v %d %d \n", &b, b, b, len(b), cap(b))
		//util.ExitPrint(a, b)
		util.ExitPrint("im die")
		//geo := container.NewGeoHash()
		//geo.Calc(39.923201, 116.390705)
	default:
		util.ExitPrint("DoTestAction flag no hit , flag:" + flag)
	}

}

//func testSlice(a []int) {
//	a = append(a, 1)
//	util.MyPrint(a)
//}
//
//func changeSlice(slice []int) {
//	slice[2] = 333
//}
