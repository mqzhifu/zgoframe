package util

import "time"

func GetNowDateMonth()string{
	//cstZone := time.FixedZone("CST", 8*3600)
	//now := time.Now().In(cstZone)
	//now := time.Now()
	//str := strconv.Itoa(now.Year()) +  strconv.Itoa(int(now.Month())) + strconv.Itoa( now.Day())
	str := time.Now().Format("200601")
	return str
}

func GetNowDate()string{
	//cstZone := time.FixedZone("CST", 8*3600)
	//now := time.Now().In(cstZone)
	//now := time.Now()
	//str := strconv.Itoa(now.Year()) +  strconv.Itoa(int(now.Month())) + strconv.Itoa( now.Day())
	str := time.Now().Format("20060102")
	return str
}
func GetNowDateHour()string{
	//ss := time.Now().Format("20060102")
	//ExitPrint(ss)
	//cstZone := time.FixedZone("CST", 8*3600)
	//now := time.Now().In(cstZone)
	//now := time.Now()
	//date := GetNowDate()
	//str := date + strconv.Itoa(now.Hour())
	str := time.Now().Format("2006010215")
	return str
}

func GetDateHour(now int64)string{
	//ss := time.Now().Format("20060102")
	//ExitPrint(ss)
	//cstZone := time.FixedZone("CST", 8*3600)
	//now := time.Now().In(cstZone)
	//now := time.Now()
	//date := GetNowDate()
	//str := date + strconv.Itoa(now.Hour())
	tm := time.Unix(now, 0)
	str := tm.Format("2006010215")
	return str
}


func GetNowTimeSecondToInt()int{
	return int( time.Now().Unix() )
}

func GetNowTimeSecondToInt64()int64{
	return time.Now().Unix()
}

func GetNowMillisecond()int64{
	return time.Now().UnixNano() / 1e6
}
