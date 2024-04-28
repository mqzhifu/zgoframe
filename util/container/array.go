package container

import (
	"errors"
	"fmt"
)

type ArrayList struct {
	Pool []*ListNode

	NodeMax int //最多包含多少个节点,0代表不限制
	Len     int //当前总节点数
	Order   int //节点的存储顺序 ：无序/升序/降序
	Debug   int //调试 模式

}

// loop bool
func NewArrayList(order int, nodeMax int, debug int) *ArrayList {
	arrayList := new(ArrayList)
	arrayList.NodeMax = nodeMax
	arrayList.Order = order
	//arrayList.Loop = loop//是否为循环链表，数组做容器用不到该参数
	arrayList.Debug = debug
	//这里为了简化代码，直接开辟了一个最大值的数组大小，实际上应该是：有个算法，定期扩容 或 缩小
	for i := 0; i < nodeMax; i++ {
		arrayList.Pool = append(arrayList.Pool, nil)
	}

	return arrayList
}

func (arrayList *ArrayList) IsEmpty() bool {
	if arrayList.Len <= 0 {
		return true
	}
	return false
}

// 输出信息，用于debug
func (arrayList *ArrayList) Print(a ...interface{}) (n int, err error) {
	if arrayList.Debug > 0 {
		return fmt.Println(a)
	}
	return
}

func (arrayList *ArrayList) CheckMaxNode() error {
	if arrayList.NodeMax > 0 { // <=0  证明没有限制node总数
		if arrayList.Len >= arrayList.NodeMax {
			msg := "linkedList.Len > linkedList.NodeMax"
			return arrayList.makeError(msg)
		}
	}
	return nil
}

func (arrayList *ArrayList) Length() int {
	return arrayList.Len
}

func (arrayList *ArrayList) GetMiddleNode() (empty bool, node *ListNode) {
	GetAllEmpty, nodeList := arrayList.GetAll(ListSearchCondition{})
	if GetAllEmpty {
		return true, node
	}
	if arrayList.Len == 1 {
		return false, arrayList.Pool[0]
	}
	var middle int
	if arrayList.Len%2 == 0 {
		middle = arrayList.Len / 2
	} else {
		middle = arrayList.Len/2 + 1
	}
	//linkedList.Print(" linkedList.Len / 2:", middle)
	for k, v := range nodeList {
		//linkedList.Print("k:",k , " v:",v.Keyword)
		if k == middle-1 {
			return false, v
		}
	}
	return false, node
}

func (arrayList *ArrayList) InsertNode(direction int, location int, searchKeyword int, keyword int, data interface{}) (int, error) {
	arrayList.Print("InsertNode direction :", direction, " location:", location, " searchKeyword:", searchKeyword, " keyword:", keyword)
	if err := arrayList.CheckMaxNode(); err != nil {
		return 0, err
	}

	if keyword <= 0 {
		msg := "keyword 仅支持 正整形"
		return 0, arrayList.makeError(msg)
	}
	//创建一个新的/空的节点
	newNode := NewListNode(keyword, data)
	if arrayList.IsEmpty() { //当前窗口为空的时候，直接插入就即可
		newNode.Location = 0
		arrayList.Pool[0] = newNode
		arrayList.Len = 1
		return 1, nil
	}

	if arrayList.Order == ORDER_NONE { //无序列表
		if location >= 0 && searchKeyword > 0 {
			msg := "location mutex searchKeyword"
			return 0, arrayList.makeError(msg)
		}

		if location < 0 && searchKeyword <= 0 {
			msg := "location <0 && searchKeyword <=0"
			return 0, arrayList.makeError(msg)
		}

		var searchNode *ListNode
		var empty bool
		if location < 0 && searchKeyword > 0 {
			empty, searchNode = arrayList.FindOneNodeByKeyword(searchKeyword)
			if empty {
				msg := "find keyword is empty!"
				return 0, arrayList.makeError(msg)
			}
			//arrayList.Print(empty,searchNode)
			location = searchNode.Location
		} else {
			if location >= arrayList.Len {
				msg := "location more max len"
				return 0, arrayList.makeError(msg)
			}
		}

		if location == 0 || location == arrayList.Len-1 {
			arrayList.InsertNodeByTop(direction, newNode)
		} else {
			//empty , searchNode = arrayList.FindOneNodeByLocation(location)
			//linkedList.Print(location,empty,searchNode)
			//if empty{
			//	msg := "find keyword is empty!"
			//	return 0,arrayList.makeError(msg)
			//}
			//arrayList.MoveArrayOneElement(searchNode.Location,DIRECTION_END)
			arrayList.MoveArrayOneElement(location, DIRECTION_END)
			newNode.Location = location
			arrayList.Pool[location] = newNode
			arrayList.Len++
		}
	} else {
		if location >= 0 {
			msg := "有序列表不支持：在具体位置插入，因为位置是由程序动态计算出来的"
			return 0, arrayList.makeError(msg)
		}

		if searchKeyword > 0 {
			msg := "有序列表：searchKeyword 参数无效 "
			return 0, arrayList.makeError(msg)
		}
		//有序节点，那就得先找到该节点的应该放在哪个位置上
		max, min, node, err := arrayList.FindOneNodeInsertLocationByKeyword(keyword)
		arrayList.Print("FindOneNodeInsertLocationByKeyword max:", max, " min:", min, " node:", node)
		if err != nil {
			return 0, err
		}
		if max {
			arrayList.InsertNodeByTop(DIRECTION_END, newNode)
		} else if min {
			arrayList.InsertNodeByTop(DIRECTION_FIRST, newNode)
		} else {
			arrayList.Print(" in middle~~~")

			arrayList.MoveArrayOneElement(node.Location, DIRECTION_END)
			newNode.Location = node.Location
			arrayList.Pool[node.Location] = newNode
			arrayList.Len++
		}

	}
	return arrayList.Len, nil
}

func (arrayList *ArrayList) FindOneNodeByLocation(location int) (empty bool, searchNode *ListNode) {
	if location < 0 || location >= arrayList.Len {
		return true, searchNode
	}
	return false, arrayList.Pool[location]
}

func (arrayList *ArrayList) InsertNodeByFirst(keyword int, data interface{}) (int, error) {
	return arrayList.InsertNode(DIRECTION_FIRST, 0, -1, keyword, data)
}

func (arrayList *ArrayList) InsertNodeByLocation(location int, keyword int, data interface{}) (int, error) {
	arrayList.Print("InsertNodeByLocation")
	return arrayList.InsertNode(DIRECTION_FIRST, location, -1, keyword, data)
}

func (arrayList *ArrayList) InsertNodeByKeyword(searchKeyword, keyword int, data interface{}) (int, error) {
	arrayList.Print("InsertNodeByKeyword")
	return arrayList.InsertNode(DIRECTION_FIRST, -1, searchKeyword, keyword, data)
}

func (arrayList *ArrayList) InsertNodeByEnd(keyword int, data interface{}) (int, error) {
	return arrayList.InsertNode(DIRECTION_END, 0, -1, keyword, data)
}

func (arrayList *ArrayList) InsertNodeByTop(direction int, newNode *ListNode) {
	arrayList.Print("InsertNodeByTop ", direction, " newNodeKeyword:", newNode.Keyword)
	if direction == DIRECTION_FIRST { //头部插入一个元素
		newNode.Location = 0
		//将当前所有元素，整体向下移动一个位置
		arrayList.MoveArrayOneElement(0, DIRECTION_END)
		arrayList.Pool[0] = newNode
	} else { //尾部插入一个元素
		newNode.Location = arrayList.Len
		arrayList.Pool[arrayList.Len] = newNode
	}
	arrayList.Len++

}

// 移动数组中的元素，其实是调整位置，当添加/删除一个节点时，得整体向下/向上移动整个数组，保证数组下标能对上
func (arrayList *ArrayList) MoveArrayOneElement(offset int, direction int) {
	arrayList.Print("MoveArrayOneElement offset:", offset, " direction:", direction)
	//因为调用之前做了 长度 判定，后面肯定会有个 空节点
	if direction == DIRECTION_FIRST { //向上移动
		//从数组头部开始遍历，用下面的值，覆盖上面的值
		//offset 对应的下标元素将会被覆盖掉
		for i := offset; i < arrayList.Len-1; i++ {
			arrayList.Pool[i] = arrayList.Pool[i+1]
			arrayList.Pool[i].Location = i
		}
	} else { //向下移动，从尾部开始移动
		//从尾部开始遍历，将当前值覆盖到下一个元素
		for i := arrayList.Len - 1; i >= offset; i-- {
			//arrayList.Print(i,arrayList.Pool[i].Keyword)
			//将当前位置的值 传 给下一个节点
			arrayList.Pool[i+1] = arrayList.Pool[i]
			arrayList.Pool[i+1].Location = i + 1
		}

	}

}

func (arrayList *ArrayList) FindOneNodeInsertLocationByKeyword(keyword int) (max bool, min bool, searchNode *ListNode, err error) {
	if arrayList.Order == ORDER_NONE {
		//此函数，对于：无序链表无用，不允许它使用
		msg := "arrayList.Order == ORDER_NONE"
		arrayList.Print(msg)
		return max, min, searchNode, errors.New(msg)
	}
	eachCompareRs := 0
	locationIndex := -1
	for i := 0; i < arrayList.Len; i++ {
		if arrayList.Order == ORDER_DESC {
			if keyword >= arrayList.Pool[i].Keyword {
				locationIndex = i
				eachCompareRs = 1
				break
			}
		} else {
			if keyword <= arrayList.Pool[i].Keyword {
				locationIndex = i
				eachCompareRs = 1
				break
			}
		}
	}

	arrayList.Print("eachCompareRs : ", eachCompareRs)
	//没找到适合的位置，证明该节点为：极值，即：整个链表遍历到了最后
	if eachCompareRs == 0 {
		max = true
	} else {
		if locationIndex == 0 {
			arrayList.Print("End min = true")
			min = true
		}
	}
	return max, min, searchNode, nil

}

func (arrayList *ArrayList) GetAllByFirst(listSearchCondition ListSearchCondition) (empty bool, nodeList []*ListNode) {
	arrayList.Print("GetAllByFirst ")
	listSearchCondition.Direction = DIRECTION_FIRST
	return arrayList.GetAll(listSearchCondition)
}

func (arrayList *ArrayList) GetAllByEnd(listSearchCondition ListSearchCondition) (empty bool, nodeList []*ListNode) {
	arrayList.Print("GetAllByEnd ")
	listSearchCondition.Direction = DIRECTION_END
	return arrayList.GetAll(listSearchCondition)
}

func (arrayList *ArrayList) GetAll(listSearchCondition ListSearchCondition) (empty bool, nodeList []*ListNode) {
	if arrayList.IsEmpty() {
		return true, nodeList
	}

	if listSearchCondition.Direction == DIRECTION_FIRST {
		for i := 0; i < arrayList.Len; i++ {
			nodeList = append(nodeList, arrayList.Pool[i])
		}
	} else {
		for i := arrayList.Len - 1; i >= 0; i-- {
			nodeList = append(nodeList, arrayList.Pool[i])
		}
	}
	//这里得优化下，不应该把数据都一次全遍历出来
	if !listSearchCondition.IsEmpty {
		offset := 0                         //从第几个节点开始 获取结果集
		limit := len(nodeList)              //获取多少个节点 结果集
		nodeListLen := len(nodeList)        //暂存一下，后面该变量会被覆盖
		if listSearchCondition.Offset > 0 { //证明调用者设置了该值
			if listSearchCondition.Offset <= len(nodeList)-1 {
				offset = listSearchCondition.Offset
			}
		}

		if listSearchCondition.Limit > 0 {
			if listSearchCondition.Limit <= len(nodeList) {
				limit = listSearchCondition.Limit
			}
		}

		//arrayList.Print("offset:",offset, " limit :",limit)
		var newNodeList []*ListNode
		for i := offset; i < offset+limit; i++ {
			if i >= nodeListLen {
				break
			}

			compare, _ := ListSearchCompare(nodeList[i], listSearchCondition)
			if compare {
				if listSearchCondition.Keyword > 0 {
					//arrayList.Print(nodeList[i].Keyword,listSearchCondition.Keyword)
					if nodeList[i].Keyword == listSearchCondition.Keyword {
						//arrayList.Print("222333 ")
						//linkedList.Print("GetAll search keyword: next ",nodeList[i].Next.Keyword , " pre:", nodeList[i].Previous.Keyword)
						newNodeList = append(newNodeList, nodeList[i])
					}
				} else {
					newNodeList = append(newNodeList, nodeList[i])
				}
				//newNodeList = append(newNodeList,nodeList[i])
			}
		}
		//arrayList.Print("newNodeList:",newNodeList)
		nodeList = newNodeList
	}

	return false, nodeList
}

func (arrayList *ArrayList) DelFirstNode() (node *ListNode, err error) {
	return arrayList.Pool[0], arrayList.DelOneNode(arrayList.Pool[0])
}

func (arrayList *ArrayList) DelEndNode() (node *ListNode, err error) {
	return arrayList.Pool[arrayList.Len-1], arrayList.DelOneNode(arrayList.Pool[arrayList.Len-1])
}

func (arrayList *ArrayList) DelOneNodeByLocation(direction int, location int, limit int) (listNodeList []*ListNode, err error) {
	//if  linkedList.IsEmpty(){
	//	return listNodeList,linkedList.makeError("no search , is empty!")
	//}

	if location < 0 {
		msg := "location < 0"
		return listNodeList, arrayList.makeError(msg)
	}

	if location > arrayList.Len-1 {
		msg := "location >  linkedList.Len - 1"
		return listNodeList, arrayList.makeError(msg)
	}

	if location == 0 {
		arrayList.DelFirstNode()
		return append(listNodeList, arrayList.Pool[0]), nil
	}

	if location == arrayList.Len-1 {
		arrayList.DelEndNode()
		return append(listNodeList, arrayList.Pool[arrayList.Len-1]), nil
	}
	var nodeList []*ListNode
	var empty bool
	if direction == DIRECTION_FIRST {
		empty, nodeList = arrayList.GetAll(ListSearchCondition{Direction: direction})
	} else {
		empty, nodeList = arrayList.GetAll(ListSearchCondition{Direction: direction})
	}

	if empty {
		msg := "is empty~"
		return listNodeList, arrayList.makeError(msg)
	}

	arrayList.Print("location:", location, " limit:", limit)
	delLimitCnt := 0
	for k, listNode := range nodeList {
		if k >= location {
			if delLimitCnt >= limit {
				break
			}
			delLimitCnt++
			listNodeList = append(listNodeList, listNode)
			arrayList.DelOneNode(listNode)
		}
	}
	return listNodeList, nil
}

func (arrayList *ArrayList) DelNodeByKeyword(keyword int) (node *ListNode, err error) {
	arrayList.Print("DelNodeByKeyword : ", keyword, " now-len :", arrayList.Len)
	//列表支持keyword 重复，所以可能一次搜索到的结果集是多个
	empty, nodeList := arrayList.GetAll(ListSearchCondition{Keyword: keyword, IsEmpty: false})
	if empty {
		return node, arrayList.makeError("no search")
	}
	//arrayList.Print("DelNodeByKeyword getall :",nodeList)
	for _, v := range nodeList {
		arrayList.DelOneNode(v)
	}
	return node, nil
}

//func(arrayList *ArrayList) DelOneNodeByKeyword(keyword int)error{
//	arrayList.Print("DelOneNode keyword : ",keyword," now-len :",arrayList.Len)
//	empty,node := arrayList.FindOneNodeByKeyword(keyword)
//	if empty{
//		return errors.New("no search")
//	}
//	return arrayList.DelOneNode(node)
//}

func (arrayList *ArrayList) DelOneNode(node *ListNode) error {
	arrayList.Print("DelOneNode:", node.Keyword)
	if arrayList.IsEmpty() {
		return arrayList.makeError("no search , is empty!")
	}
	//如果当前链表仅剩下一个节点，直接清空即可
	if arrayList.Len == 1 {
		arrayList.Len = 0
		return nil
	}
	if node.Location == 0 {
		//直接将数组整体向上移动一格，0元素会被第一个元素覆盖
		arrayList.MoveArrayOneElement(0, DIRECTION_FIRST)
	} else if node.Location == arrayList.Len-1 {
		//从尾部删除
	} else {
		//从中间删除一个
		arrayList.MoveArrayOneElement(node.Location, DIRECTION_FIRST)
	}

	arrayList.Len--
	return nil
}
func (arrayList *ArrayList) FindOneNodeByKeywordAndDel(keyword int) (empty bool, searchNode []*ListNode) {
	return empty, searchNode
}

func (arrayList *ArrayList) FindOneNodeByLocationAndDel(location int) (empty bool, searchNode *ListNode) {
	getAllEmpty, nodeList := arrayList.GetAllByFirst(ListSearchCondition{})
	if getAllEmpty {
		return getAllEmpty, searchNode
	}
	//非法值，默认为：最后一个节点
	if location >= arrayList.Len || location < 0 {
		location = arrayList.Len
	}

	arrayList.Print("FindOneNodeByLocation  location:", location)
	for k, v := range nodeList {
		if k == location {
			arrayList.DelOneNodeByLocation(DIRECTION_FIRST, k, 1)
			return false, v
		}
	}
	return true, searchNode
}

// 根据 关键词 查找一个节点
func (arrayList *ArrayList) FindOneNodeByKeyword(keyword int) (empty bool, searchNode *ListNode) {
	if arrayList.IsEmpty() {
		return true, searchNode
	}

	for i := 0; i < arrayList.Len; i++ {
		if arrayList.Pool[i].Keyword == keyword {
			return false, arrayList.Pool[i]
		}
	}
	return true, searchNode
}

func (arrayList *ArrayList) makeError(msg string) error {
	arrayList.Print("[errors] " + msg)
	return errors.New(msg)
}

func (arrayList *ArrayList) NodeRepeatTotal() (repeatList map[int]int, empty bool) {
	empty, list := arrayList.GetAllByFirst(ListSearchCondition{})
	if empty {
		return
	}
	var keywordList []int

	if len(list) == 1 {
		return
	}

	for _, v := range list {
		keywordList = append(keywordList, v.Keyword)
	}

	repeatListCnt := make(map[int]int)
	for i := 0; i < len(keywordList)-1; i++ {
		_, ok := repeatListCnt[keywordList[i]]
		if !ok {
			repeatListCnt[keywordList[i]] = 0
		}

		repeatListCnt[keywordList[i]] = repeatListCnt[keywordList[i]] + 1
	}

	repeatList = make(map[int]int)
	for k, v := range repeatListCnt {
		if v > 1 {
			repeatList[k] = v
		}
	}

	return repeatList, false
}

func (arrayList *ArrayList) InsertMultiNode([]ListNode) (int, error) {
	return 0, nil
}
