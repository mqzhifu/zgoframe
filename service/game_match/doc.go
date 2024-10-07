package gamematch

import "strconv"

/*
这里才是最终的匹配计算公式~

匹配的最小单们是组，而不是玩家，每个组有几个人~  这<几个人>就是匹配的最小属性
根据每个组里的人数，加和，等于固定一个正整数

1人组：1个
2人组：4个
3人组：2个
4人组：6个
5人组：3个

假设，满足一局游戏，需要21个人，如下：
1+2+3+4+5 = 21
*/
func (match *Match) calculateNumberTotal(sum int, groupPerson map[int]int) map[int][5]int {

	result := make(map[int][5]int)
	inc := 0
	for a := 0; a <= groupPerson[5]; a++ {
		for b := 0; b <= groupPerson[4]; b++ {
			for c := 0; c <= groupPerson[3]; c++ {
				for d := 0; d <= groupPerson[2]; d++ {
					for e := 0; e <= groupPerson[1]; e++ {
						if 5*a+4*b+3*c+2*d+e == sum {
							//ttt := [5]int{a,b,c,d,e}
							ttt := [5]int{e, d, c, b, a}
							result[inc] = ttt
							match.Log.Debug("5 x " + strconv.Itoa(a) + " + 4 x" + strconv.Itoa(b) + "+ 3 x " + strconv.Itoa(c) + " + 2 x " + strconv.Itoa(d) + " 1 x " + strconv.Itoa(e) + "=" + strconv.Itoa(sum))
							inc++
						}
					}
				}
			}
		}
	}
	return result

	//aCnt := groupPerson[0]
	//bCnt := 2
	//cCnt := 20
	//
	//inc := 0
	//for a:=0;a<=aCnt;a++{
	//	for b:=0;b<=bCnt;b++{
	//		for c:=0;c<=cCnt;c++{
	//			if a + 2 * b + 5 * c == sum {
	//				inc++
	//				zlib.MyPrint("1 x ",a," + 2 x",b,"+ 5 x ",c,"=",sum)
	//			}
	//		}
	//	}
	//}
}
