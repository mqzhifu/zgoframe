package container

import (
	"math"
	"strconv"
	"zgoframe/util"
)

type GeoHash struct {
	LongitudeMax float64 //经度 最大
	LongitudeMin float64 //经度 最大

	LatitudeMax float64 //纬度 最大
	LatitudeMin float64 //纬度 最大
	//第1位 10000米
	//第2位 1000米
	//第3位 100米
	//第4位 10米
	PointExactly int //小数点精确到第几位
}

func NewGeoHash() *GeoHash {
	geoHash := new(GeoHash)

	geoHash.LatitudeMax = 90
	geoHash.LatitudeMin = -90

	geoHash.LongitudeMax = 180
	geoHash.LongitudeMin = -180

	geoHash.PointExactly = 3

	return geoHash
}

func (geoHash *GeoHash) Calc(latitude float64, longitude float64) string {
	util.MyPrint("latitude:", latitude, " longitude:", longitude)
	if latitude > geoHash.LatitudeMax || latitude < geoHash.LatitudeMin {
		util.MyPrint("geoHash.LatitudeMax:", geoHash.LatitudeMax, " geoHash.LatitudeMin:", geoHash.LatitudeMin)
	}

	if longitude > geoHash.LongitudeMax || longitude < geoHash.LongitudeMin {
		util.MyPrint("geoHash.LongitudeMax:", geoHash.LongitudeMax, " geoHash.LongitudeMin:", geoHash.LongitudeMin)
	}

	LatitudeMaxBinary := geoHash.CalcRecursion(latitude, geoHash.LatitudeMax, geoHash.LatitudeMin, 1)
	LongitudeMinBinary := geoHash.CalcRecursion(longitude, geoHash.LongitudeMax, geoHash.LongitudeMin, 1)
	geoHash.MergeLongitude(LongitudeMinBinary, LatitudeMaxBinary)
	return "00000"
}

func (geoHash *GeoHash) CalcRecursion(number float64, locationOne float64, locationTwo float64, inc int) string {
	util.MyPrint("inc:", inc, " number:", number, " locationOne:", locationOne, " locationTwo:", locationTwo)
	inc++

	if inc > 20 {
		return " err:inc out "
	}

	numberStr := geoHash.Float64ToString(number)
	locationOneStr := geoHash.Float64ToString(locationOne)
	locationTwoStr := geoHash.Float64ToString(locationTwo)
	util.MyPrint("numberStr:", numberStr, " locationOneStr:", locationOneStr, " locationTwo_str:", locationTwoStr)
	if numberStr == locationOneStr || numberStr == locationTwoStr {
		util.MyPrint("success")
		return ""
	}

	//只有这种情况会出现 一个数为：正，一个数为：负，其余情况均是：正正 负负  0 除外
	if (locationOne == geoHash.LatitudeMax && locationTwo == geoHash.LatitudeMin) || (locationOne == geoHash.LongitudeMax && locationTwo == geoHash.LongitudeMin) {

		if number > 0 && number < locationOne {
			return "1" + geoHash.CalcRecursion(number, 0, locationOne, inc)
		} else if number > locationTwo && number < 0 {
			return "0" + geoHash.CalcRecursion(number, 0, locationTwo, inc)
		} else {
			util.ExitPrint("err2")
		}
	}
	//distance ：使用了绝对值， 一定是一个正数
	distance := (math.Abs(locationTwo) - math.Abs(locationOne)) / 2
	util.MyPrint("distance:", distance)

	var locationOneEnd, locationTwoStart float64

	if locationTwo > 0 { //两个数均为正数
		locationOneEnd = locationOne + distance
		locationTwoStart = locationTwo - distance
	} else { //负数
		if locationOne == 0 { //一个为0 一个为负数
			locationOneEnd = locationOne + distance
			locationTwoStart = locationTwoStart + distance
		} else { //两个数均为负数
			locationOneEnd = locationOne - distance
			locationTwoStart = locationTwo + distance
		}
	}

	util.MyPrint("locationOne:", locationOne, " locationOneEnd:", locationOneEnd, " locationTwoStart:", locationTwoStart, " locationTwo:", locationTwo)

	if number > locationOne && number < locationOneEnd {
		return "0" + geoHash.CalcRecursion(number, locationOne, locationOneEnd, inc)
	} else if number > locationTwoStart && number < locationTwo {
		return "1" + geoHash.CalcRecursion(number, locationTwoStart, locationTwo, inc)
	}

	return "err4"
}

// 0 1 2 3 4 5 6 7 8 9 10 11
// 0 0 1 1 2 2 3 3 4 4 5  5
// 合并 二进制 经纬度
// 从0开始
// 偶数为：经度
// 奇数为：纬度
func (geoHash *GeoHash) MergeLongitude(LongitudeMinBinary string, LatitudeMaxBinary string) {
	length := len(LongitudeMinBinary) + len(LatitudeMaxBinary)
	util.MyPrint("MergeLongitude LongitudeMinBinary:", LongitudeMinBinary, " LatitudeMaxBinary:", LatitudeMaxBinary, " totalLength:", length)

	str := make([]byte, length)
	for i := 0; i < length; i++ {
		location := i / 2
		util.MyPrint("=====================i:", i, " location:", location)
		if i%2 == 0 {
			str[i] = LongitudeMinBinary[location]
		} else {
			str[i] = LatitudeMaxBinary[location]
		}
	}

	util.ExitPrint(string(str))
}
func (geoHash *GeoHash) Float64ToString(number float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(number, 'f', geoHash.PointExactly, 64)
}
