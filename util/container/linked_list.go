package container

import (
	"errors"
	"fmt"
)

/*
	<链表>容器，支持：循环 双向 有序/无序，节点keyword可重复
	注意：目前keyword只支持数字类型的操作，不支持字符串,且只能>=0
*/

// 整个链表
type LinkedList struct {
	First   *ListNode //首节点地址
	End     *ListNode //尾节点地址
	NodeMax int       //最多包含多少个节点,<=0代表不限制
	Len     int       //当前链表总节点数
	Order   int       //节点的存储顺序 ：无序/升序/降序
	Loop    bool      //是否为循环链表
	Debug   int       //调试 模式
}

// 输出信息，用于debug
func (linkedList *LinkedList) Print(a ...interface{}) (n int, err error) {
	if linkedList.Debug > 0 {
		return fmt.Println(a)
	}
	return
}

// 创建一个新的链表
func NewLinkedList(order int, nodeMax int, loop bool, debug int) *LinkedList {
	if order != ORDER_NONE && order != ORDER_DESC && order != ORDER_ASC {
		order = ORDER_NONE
	}

	if nodeMax <= NODE_MIN {
		nodeMax = NODE_MIN
	} else if nodeMax > NODE_MAX {
		nodeMax = NODE_MAX
	}

	linkedList := new(LinkedList)
	linkedList.Order = order
	linkedList.NodeMax = nodeMax
	linkedList.Loop = loop
	linkedList.Debug = debug

	linkedList.Print("NewLinkedList , [options] order :", order, " nodeMax:", nodeMax, " loop:", loop, " debug:", debug)

	return linkedList
}

func (linkedList *LinkedList) IsEmpty() bool {
	if linkedList.Len == 0 || linkedList.First == nil {
		linkedList.Print("[warning] is empty = true")
		return true
	}
	return false
}

func (linkedList *LinkedList) CheckMaxNode() error {
	if linkedList.NodeMax > 0 { // <=0  证明没有限制node总数
		if linkedList.Len >= linkedList.NodeMax {
			msg := "linkedList.Len > linkedList.NodeMax"
			return linkedList.makeError(msg)
		}
	}
	return nil
}

func (linkedList *LinkedList) Length() int {
	return linkedList.Len
}

func (linkedList *LinkedList) InsertNodeByFirst(keyword int, data interface{}) (int, error) {
	linkedList.Print("InsertNodeByFirst keyword :", keyword, "  data:", data)
	return linkedList.InsertNode(DIRECTION_FIRST, 0, -1, keyword, data)
}

func (linkedList *LinkedList) InsertNodeByEnd(keyword int, data interface{}) (int, error) {
	linkedList.Print("InsertNodeByEnd keyword :", keyword, "  data:", data)
	return linkedList.InsertNode(DIRECTION_END, linkedList.Len-1, -1, keyword, data)
}

// 目前只能是正向在某个节点的上方添加
func (linkedList *LinkedList) InsertNodeByLocation(location int, keyword int, data interface{}) (int, error) {
	linkedList.Print("InsertNodeByEnd keyword :", keyword, "  data:", data)
	return linkedList.InsertNode(DIRECTION_FIRST, location, -1, keyword, data)
}

func (linkedList *LinkedList) InsertNodeByKeyword(searchKeyword, keyword int, data interface{}) (int, error) {
	linkedList.Print("InsertNodeByEnd keyword :", keyword, "  data:", data)
	return linkedList.InsertNode(DIRECTION_FIRST, -1, searchKeyword, keyword, data)
}

func (linkedList *LinkedList) InsertNode(direction int, location int, searchKeyword int, keyword int, data interface{}) (int, error) {
	linkedList.Print("InsertNode direction :", direction, " location:", location, " searchKeyword:", searchKeyword, " keyword:", keyword)
	if err := linkedList.CheckMaxNode(); err != nil {
		return 0, err
	}

	if keyword <= 0 {
		msg := "keyword 仅支持 正整形"
		return 0, linkedList.makeError(msg)
	}
	//创建一个新的/空的节点
	newNode := NewListNode(keyword, data)
	if linkedList.IsEmpty() {
		newNode.Location = 0
		//重围节点的上下指针
		if linkedList.Loop { //循环链表
			newNode.Previous = newNode
			newNode.Next = newNode
		} else {
			newNode.Previous = nil
			newNode.Next = nil
		}

		linkedList.First = newNode
		linkedList.End = newNode
		linkedList.Len = 1
		return 1, nil
	}

	if linkedList.Order == ORDER_NONE { //无序列表
		if location >= 0 && searchKeyword > 0 {
			msg := "location mutex searchKeyword"
			return 0, linkedList.makeError(msg)
		}

		if location < 0 && searchKeyword <= 0 {
			msg := "location <0 && searchKeyword <=0"
			return 0, linkedList.makeError(msg)
		}

		var searchNode *ListNode
		var empty bool
		if location < 0 && searchKeyword > 0 {
			empty, searchNode = linkedList.FindOneNodeByKeyword(searchKeyword)
			//linkedList.Print(empty)
			if empty {
				msg := "find keyword is empty!"
				return 0, linkedList.makeError(msg)
			}
			location = searchNode.Location
		} else {
			if location >= linkedList.Len {
				msg := "location more max len"
				return 0, linkedList.makeError(msg)
			}
		}

		if location == 0 || location == linkedList.Len-1 {
			linkedList.InsertNodeByTop(direction, newNode)
		} else {
			empty, searchNode = linkedList.FindOneNodeByLocation(location)
			//linkedList.Print(location,empty,searchNode)
			if empty {
				msg := "find keyword is empty!"
				return 0, linkedList.makeError(msg)
			}
			linkedList.InsertNodeByMiddle(direction, newNode, searchNode)
		}
	} else {
		if location >= 0 {
			msg := "有序列表不支持：在具体位置插入，因为位置是由程序动态计算出来的"
			return 0, linkedList.makeError(msg)
		}

		if searchKeyword > 0 {
			msg := "有序列表：searchKeyword 参数无效 "
			return 0, linkedList.makeError(msg)
		}
		//有序节点，那就得先找到该节点的应该放在哪个位置上
		max, min, node, err := linkedList.FindOneNodeInsertLocationByKeyword(keyword)
		linkedList.Print("FindOneNodeInsertLocationByKeyword max:", max, " min:", min, " node:", node)
		if err != nil {
			return 0, err
		}

		if max {
			linkedList.InsertNodeByTop(DIRECTION_END, newNode)
		} else if min {
			linkedList.InsertNodeByTop(DIRECTION_FIRST, newNode)
		} else {
			linkedList.Print(" in middle~~~keyword:", node.Keyword)
			linkedList.InsertNodeByMiddle(DIRECTION_FIRST, newNode, node)
		}

	}
	return linkedList.Len, nil
}

// 前置条件：做了为空判断，不往头/尾插入，只是往中间位置插入，插入点的上现都存在节点
func (linkedList *LinkedList) InsertNodeByMiddle(direction int, newNode *ListNode, searchNode *ListNode) {
	if direction == DIRECTION_FIRST {
		newNode.Previous = searchNode.Previous
		newNode.Next = searchNode
		searchNode.Previous.Next = newNode
		searchNode.Previous = newNode
	} else {
		newNode.Previous = searchNode
		newNode.Next = searchNode.Next

		searchNode.Next = newNode
		searchNode.Next.Previous = newNode
	}
	linkedList.Len++
}

// 从两端(头|尾)插入一个节点
// 此方法适用无序链表
// 前置条件：做了为空判断
func (linkedList *LinkedList) InsertNodeByTop(direction int, newNode *ListNode) {
	linkedList.Print("InsertNodeByTop")
	//var oldFirstNode *ListNode
	if direction == DIRECTION_FIRST {
		//oldFirstNode := linkedList.First
		//newNode.Next = oldFirstNode
		//oldFirstNode.Previous = newNode
		newNode.Next = linkedList.First
		linkedList.First.Previous = newNode
		if linkedList.Loop {
			newNode.Previous = linkedList.End
			linkedList.End.Next = newNode
			//if linkedList.Len == 1{
			//	linkedList.First.Next = newNode
			//}
		} else {
			newNode.Previous = nil
		}
		linkedList.First = newNode
	} else {
		//oldEndNode := linkedList.End
		//newNode.Previous = oldEndNode
		//oldEndNode.Next = newNode
		newNode.Previous = linkedList.End
		linkedList.End.Next = newNode
		if linkedList.Loop {
			newNode.Next = linkedList.First
			linkedList.First.Previous = newNode
			//if linkedList.Len == 1{
			//	linkedList.First.Previous = newNode
			//}
		} else {
			newNode.Next = nil
		}
		linkedList.End = newNode
	}
	linkedList.Len++
}

func (linkedList *LinkedList) GetMiddleNode() (empty bool, node *ListNode) {
	GetAllEmpty, nodeList := linkedList.GetAll(ListSearchCondition{})
	if GetAllEmpty {
		return true, node
	}
	if linkedList.Len == 1 {
		return false, linkedList.First
	}
	var middle int
	if linkedList.Len%2 == 0 {
		middle = linkedList.Len / 2
	} else {
		middle = linkedList.Len/2 + 1
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

func (linkedList *LinkedList) GetAllByFirst(listSearchCondition ListSearchCondition) (empty bool, nodeList []*ListNode) {
	linkedList.Print("GetAllByFirst start :", listSearchCondition.IsEmpty, listSearchCondition.Direction)
	listSearchCondition.Direction = DIRECTION_FIRST
	return linkedList.GetAll(listSearchCondition)
}

// 获取链表所有节点，从悎部开始
func (linkedList *LinkedList) GetAllByEnd(listSearchCondition ListSearchCondition) (empty bool, nodeList []*ListNode) {
	linkedList.Print("GetAllByEnd :", listSearchCondition)
	listSearchCondition.Direction = DIRECTION_END
	return linkedList.GetAll(listSearchCondition)
}

func (linkedList *LinkedList) GetAll(listSearchCondition ListSearchCondition) (empty bool, nodeList []*ListNode) {
	if linkedList.IsEmpty() {
		return true, nodeList
	}

	var node *ListNode
	//linkedList.Print(" listSearchCondition.Direction:", listSearchCondition.Direction)
	direction := listSearchCondition.Direction
	if direction <= 0 {
		direction = DIRECTION_FIRST
	}
	if direction == DIRECTION_FIRST {
		node = linkedList.First
	} else {
		node = linkedList.End
	}
	//linkedList.Print("GetAll 第一个元素:",node.Keyword,node.Next)
	cnt := 0
	for {
		cnt++
		//linkedList.Print("GetAll cng:",cnt, " node:",node.Next)
		if direction == DIRECTION_FIRST {
			if node.Next == nil || node == linkedList.End {
				//linkedList.Print("dead loop break in 1 :",node.Keyword)
				nodeList = append(nodeList, node)
				break
			}
		} else {
			if node.Previous == nil || node == linkedList.First {
				//linkedList.Print("dead loop break in 2 :",node.Keyword)
				nodeList = append(nodeList, node)
				break
			}
		}

		nodeList = append(nodeList, node)
		if direction == DIRECTION_FIRST {
			node = node.Next //从头部，向下遍历
		} else {
			node = node.Previous //从尾部，向上遍历
		}
		//linkedList.Print("ddd:",cnt, " " ,node.Keyword)
	}
	linkedList.Print("cnt:", cnt)
	//这里得优化下，不应该把数据都一次全遍历出来

	if !listSearchCondition.IsEmpty {
		offset := 0                  //从第几个节点开始 获取结果集
		limit := len(nodeList)       //获取多少个节点 结果集
		nodeListLen := len(nodeList) //暂存一下，后面该变量会被覆盖
		//linkedList.Print("offset:",offset,",limit:",limit ," ",listSearchCondition.Offset , " ",listSearchCondition.Limit)
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

		linkedList.Print("listSearchCondition offset:", offset, " limit :", limit, " no_linked_len:", linkedList.Len, " keyword:", listSearchCondition.Keyword)
		var newNodeList []*ListNode
		for i := offset; i < offset+limit; i++ {
			//linkedList.Print(nodeList[i].Keyword," next:",nodeList[i].Next.Keyword , " pre:", nodeList[i].Previous.Keyword)
			if i >= nodeListLen {
				linkedList.Print("break in i >= nodeListLen")
				break
			}

			//linkedList.Print(listSearchCondition.Keyword)
			compare, _ := ListSearchCompare(nodeList[i], listSearchCondition)
			if compare {
				if listSearchCondition.Keyword > 0 {
					if nodeList[i].Keyword == listSearchCondition.Keyword {
						//linkedList.Print("GetAll search keyword: next ",nodeList[i].Next.Keyword , " pre:", nodeList[i].Previous.Keyword)
						newNodeList = append(newNodeList, nodeList[i])
					}
				} else {
					newNodeList = append(newNodeList, nodeList[i])
				}

			}
		}
		nodeList = newNodeList
	}

	return false, nodeList
}

func (linkedList *LinkedList) DelFirstNode() (node *ListNode, err error) {
	if linkedList.IsEmpty() {
		return node, linkedList.makeError("no search , is empty!")
	}
	firstNode := linkedList.First
	return firstNode, linkedList.DelOneNode(firstNode)
}

func (linkedList *LinkedList) DelEndNode() (node *ListNode, err error) {
	if linkedList.IsEmpty() {
		return node, linkedList.makeError("no search , is empty!")
	}
	endNode := linkedList.End
	return endNode, linkedList.DelOneNode(endNode)
}

func (linkedList *LinkedList) DelOneNodeByLocation(direction int, location int, limit int) (listNodeList []*ListNode, err error) {
	if linkedList.IsEmpty() {
		return listNodeList, linkedList.makeError("no search , is empty!")
	}

	if location < 0 {
		msg := "location < 0"
		return listNodeList, linkedList.makeError(msg)
	}

	if location > linkedList.Len-1 {
		msg := "location >  linkedList.Len - 1"
		return listNodeList, linkedList.makeError(msg)
	}

	if location == 0 {
		linkedList.DelFirstNode()
		return append(listNodeList, linkedList.First), nil
	}

	if location == linkedList.Len-1 {
		linkedList.DelEndNode()
		return append(listNodeList, linkedList.End), nil
	}
	var nodeList []*ListNode
	var empty bool
	if direction == DIRECTION_FIRST {
		empty, nodeList = linkedList.GetAll(ListSearchCondition{Direction: direction})
	} else {
		empty, nodeList = linkedList.GetAll(ListSearchCondition{Direction: direction})
	}

	if empty {
		msg := "is empty~"
		return listNodeList, linkedList.makeError(msg)
	}

	linkedList.Print("location:", location, " limit:", limit)
	delLimitCnt := 0
	for k, listNode := range nodeList {
		if k >= location {
			if delLimitCnt >= limit {
				break
			}
			delLimitCnt++
			listNodeList = append(listNodeList, listNode)
			linkedList.DelOneNode(listNode)
		}
	}
	return listNodeList, nil
}

func (linkedList *LinkedList) makeError(msg string) error {
	linkedList.Print("[errors] " + msg)
	return errors.New(msg)
}

func (linkedList *LinkedList) DelNodeByKeyword(keyword int) (node *ListNode, err error) {
	linkedList.Print("DelNodeByKeyword : ", keyword, " now-len :", linkedList.Len)
	//列表支持keyword 重复，所以可能一次搜索到的结果集是多个
	empty, nodeList := linkedList.GetAll(ListSearchCondition{Keyword: keyword, IsEmpty: false})
	if empty {
		return node, linkedList.makeError("no search")
	}
	linkedList.Print(nodeList)
	for _, v := range nodeList {
		linkedList.DelOneNode(v)
	}
	return node, nil
}

func (linkedList *LinkedList) DelOneNode(node *ListNode) error {
	//如果当前链表仅剩下一个节点，直接清空即可
	if linkedList.Len == 1 {
		linkedList.Len = 0
		linkedList.First = nil
		linkedList.End = nil
		return nil
	}
	linkedList.Print("del one node : ", node.Keyword)
	//当前节点大于2个，正常操作
	if linkedList.Len > 2 {
		if node == linkedList.First {
			linkedList.Print("del in first")
			linkedList.First = node.Next
			if linkedList.Loop {
				node.Next.Previous = linkedList.End
				linkedList.End.Next = node.Next
			} else {
				node.Next.Previous = nil
			}
		} else if node == linkedList.End {
			linkedList.Print("del in end")
			linkedList.End = node.Previous
			if linkedList.Loop {
				node.Previous.Next = linkedList.First
				linkedList.First.Previous = linkedList.End
			} else {
				node.Previous.Next = nil
			}
		} else {
			linkedList.Print("del in middle", " node.next:", node.Next.Keyword, " node.pre :", node.Previous.Keyword)
			node.Previous.Next = node.Next
			node.Next.Previous = node.Previous
		}
	} else {
		//只有2个节点的时候，删除任意一个，得把两
		if node == linkedList.First {
			linkedList.First = node.Next
			linkedList.End = node.Next
			if linkedList.Loop {
				node.Next.Previous = node.Next
				node.Next.Next = node.Next
			} else {
				node.Next.Previous = nil
				node.Next.Next = nil
			}
		} else if node == linkedList.End {
			linkedList.First = node.Previous
			linkedList.End = node.Previous

			if linkedList.Loop {
				node.Next.Previous = node.Previous
				node.Next.Next = node.Previous
			} else {
				node.Next.Previous = nil
				node.Next.Next = nil
			}
		} else {
			node.Previous.Next = node.Next
			node.Next.Previous = node.Previous
		}
	}
	fmt.Println("len:", linkedList.Len, " linkedList.Len--")
	linkedList.Len--

	return nil
}

func (linkedList *LinkedList) FindOneNodeInsertLocationByKeyword(keyword int) (max bool, min bool, searchNode *ListNode, err error) {
	if linkedList.Order == ORDER_NONE {
		//此函数，对于：无序链表无用，不允许它使用
		msg := "linkedList.Order == ORDER_NONE"
		return max, min, searchNode, linkedList.makeError(msg)
	}
	linkedList.Print(":FindOneNodeInsertLocationByKeyword :", keyword)
	//不管是升序还是降序，都从头部获取整个链表
	getAllEmpty, nodeList := linkedList.GetAllByFirst(ListSearchCondition{})
	if getAllEmpty { //获取列表为空，不应该为空，调用者应该提前做了判断
		msg := "linkedList.GetAll is empty"
		return max, min, searchNode, linkedList.makeError(msg)
	}
	//linkedList.Print("getAllEmpty:",getAllEmpty)
	//遍历搜索对比：找到新插件节点应该放在哪个位置 ，最差的情况是:找不见
	eachCompareRs := 0
	for _, v := range nodeList {
		//降序
		if linkedList.Order == ORDER_DESC {
			if keyword >= v.Keyword {
				searchNode = v
				eachCompareRs = 1
				break
			}
		} else { //升序
			if keyword <= v.Keyword {
				searchNode = v
				eachCompareRs = 1
				break
			}
		}
	}
	linkedList.Print("eachCompareRs : ", eachCompareRs)
	//没找到适合的位置，证明该节点为：极值，即：整个链表遍历到了最后
	if eachCompareRs == 0 {
		max = true
	} else {
		//linkedList.Print("linkedList.First:",linkedList.First , " searchNode:",searchNode , "endNode :",linkedList.End)
		//这里有2种情况，1就是在中间某个位置正常插入即可 2 是在队首上面
		if searchNode == linkedList.First {
			linkedList.Print("End min = true")
			min = true
		}
	}

	return max, min, searchNode, nil

}

// 根据 关键词 查找一个节点
func (linkedList *LinkedList) FindOneNodeByKeyword(keyword int) (empty bool, searchNode *ListNode) {
	linkedList.Print("FindOneNodeByKeyword keyword:", keyword)
	getAllEmpty, nodeList := linkedList.GetAllByFirst(ListSearchCondition{})
	if getAllEmpty {
		return empty, searchNode
	}

	for k, v := range nodeList {
		if v.Keyword == keyword {
			v.Location = k
			return false, v
		}
	}
	return true, searchNode
}

func (linkedList *LinkedList) InsertMultiNode([]ListNode) (int, error) {
	return 0, nil
}

func (linkedList *LinkedList) FindOneNodeByLocation(location int) (empty bool, searchNode *ListNode) {
	//getAllEmpty,nodeList := linkedList.GetAllByFirst(ListSearchCondition{Offset: location,Limit: 1,IsEmpty: false})
	getAllEmpty, nodeList := linkedList.GetAllByFirst(ListSearchCondition{})
	//linkedList.Print(getAllEmpty,nodeList)
	if getAllEmpty {
		return getAllEmpty, searchNode
	}
	//非法值，默认为：最后一个节点
	if location > linkedList.Len || location < 0 {
		location = linkedList.Len
	}

	linkedList.Print("FindOneNodeByLocation  location:", location)
	for k, v := range nodeList {
		if k == location {
			return false, v
		}
	}
	return true, searchNode
}

func (linkedList *LinkedList) FindOneNodeByKeywordAndDel(keyword int) (empty bool, searchNode []*ListNode) {
	return empty, searchNode
}

func (linkedList *LinkedList) FindOneNodeByLocationAndDel(location int) (empty bool, searchNode *ListNode) {
	//getAllEmpty,nodeList := linkedList.GetAllByFirst(ListSearchCondition{Offset: location,Limit: 1,IsEmpty: false})
	getAllEmpty, nodeList := linkedList.GetAllByFirst(ListSearchCondition{})
	//linkedList.Print(getAllEmpty,nodeList)
	if getAllEmpty {
		return getAllEmpty, searchNode
	}
	//非法值，默认为：最后一个节点
	if location >= linkedList.Len || location < 0 {
		location = linkedList.Len
	}

	linkedList.Print("FindOneNodeByLocation  location:", location)
	for k, v := range nodeList {
		if k == location {
			linkedList.DelOneNodeByLocation(DIRECTION_FIRST, k, 1)
			return false, v
		}
	}
	return true, searchNode
}

func (linkedList *LinkedList) NodeRepeatTotal() (repeatList map[int]int, empty bool) {
	empty, list := linkedList.GetAllByFirst(ListSearchCondition{})
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
