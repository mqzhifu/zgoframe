package gamematch

import "strconv"

/*待解决问题
3.成功后超时 与  匹配成功重试  冲突
4.玩家status 状态是否要考虑个超时处理  定时清一下
5.优化 日志输出 ，太乱了
*/

/*待解决问题：
1:各协程之间是否需要加锁
(目前是：每个rule，都有一组守护协程，和一组redis key 存数据，所以各个rule是互不影响 的
	所以rule之间不需要加锁了
	但是单个ruleId,不加锁，可能出现的问题：
	1、报名中，有取消指令发出，但是匹配协程依然还在计算，并且计算成功了
	2、匹配中的玩家已超时，但是超时协程未执行，匹配协程先执行了，匹配成功.....
	3、匹配成功已超时，但是超时协程未执行，PUSH协程先执行了....
	4、超时的2个协程挂了，但是 报名 PUSH 匹配 协程均是正常，那匹配依然会成功，PUSH依然还是会推送

	这里有2个维度的问题：
	1、如何保证所有守护协程是正常的？进程的健康可以由外部shell控制，协程呢？
	2、如果保证上面一条是正常的，那核心点就是匹配协程了
)
2:groupId目前是使用外部的ID，是否考虑换成内部
3:查看当时进程有多少个协程，是否有metrics的方式，且最可以UI可视化
4:负载，如何在多台机器，各开一个守护进程~负载请求
(每个ruleId,每台机器负载监听哪几个ruleId，或者，哪个ruleId负载高，单独放在一台机器上
	单个ruleId，每台机器均可以启动，只要redis换个HOST即可
)
5:pnic 异常机制处理
6:log 改成 开协程 异步写
7:压测 redis etcd log 匹配 http 报名
8:
9:
*/

/*
	0. 加锁
	0. 先，从最大的范围搜索，扫整个集合，如果元素过少，直接在这个维度就结束了
	0. 缩小搜索范围，把整个集合，划分成10个维度，每个维度单纯计算，如果成功，那就结束了，如果人数过多，还会再划分一次最细粒度的匹配
	0. 这种是介于上面两者中间~即不是全部集合，也不是单独计算一个维度，而是逐步，放大到:最大90%集合，1-1，1-2....1-9
	0. 解锁

	总结：以上的算法，其实就是不断切换搜索范围（由最大到中小，再到中大），加速匹配时间
*/

//
// 组    group ( 属性A * 百分比 + 属性B * 百分比 ... )  = 最终权重值     ( = ,>= ,<=, > X < , >= X < ,>= X <= X )
// 单用户 normal ( 属性A * 百分比 + 属性B * 百分比 ... )  = 最终权重值
// 团队  team   ( 属性A * 百分比 + 属性B * 百分比 ... )  = 最终权重值

/*
fof :

1元		2张
2元		5张
3元		40张
4元		30张
5元		60张

有序集合_小组人数

	weight groupId
有序集合
	小组权重	小组ID
	weight 	groupId
hash_groupId
	小组详情信息
	groupId
		weight
		person
		timeout
		ATime
hash_groupId_players
	小组成员

集合
	groupPlayers	playerId
*/

/*
这里才是最终的匹配计算公式~
//匹配的最小单们是组，而不是玩家，每个组有几个人~  这<几个人>就是匹配的最小属性
//根据每个组里的人数，加和，等于固定一个正整数

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
	//zlib.ExitPrint(inc)
}
