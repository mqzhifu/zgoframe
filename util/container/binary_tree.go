package container

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"zgoframe/util"
)

const (
	DIRECTION_LEFT  = 1 //方向-左
	DIRECTION_RIGHT = 2 //方向-右

	NODE_KEYWORD_PLACEHOLDER = 3

	TREE_NODE_MAX = 100 //最大节点数
	KEYWORD_NIL   = 999 //有些需求，某些空节点必须得有，但是keyword得给一个int 占位符

	FLAG_NORMAL    = 1 //普通二叉树
	FLAG_BALANCE   = 2 //平衡二叉树
	FLAG_RED_BLACK = 3 //红黑二叉树
)

type TreeNode struct {
	Parent  *TreeNode //父节点-指针
	Left    *TreeNode //左节点-指针
	Right   *TreeNode //右节点-指针
	Keyword int       //关键值
	Data    interface{}
	//无论是深度优先还是广度优先：1. 每次发生变化，连带着都得跟着改
	DeepDesc  int //计算深度，遍历时，递归计算，降序 ，深度优先
	DeepAsc   int //计算深度，遍历时，层级计算，升序 ，广度优先
	LeftDeep  int //左节点深度
	RightDeep int //右节点深度
}

// 二叉树类
type BinaryTree struct {
	RootNode *TreeNode //根节点-指针
	Flag     int
	NodeMax  int //最大节点数
	Len      int
	Debug    int
}

// 创建-二叉树类
// nodeMax:树最多可包含的节点数
// flag:类型  普通 红黑 平衡
func NewBinaryTree(nodeMax int, flag int, debug int) *BinaryTree {
	binaryTree := new(BinaryTree)

	if nodeMax > TREE_NODE_MAX {
		nodeMax = TREE_NODE_MAX
	}

	binaryTree.NodeMax = nodeMax
	binaryTree.Flag = flag
	binaryTree.Debug = debug

	return binaryTree
}

func (binaryTree *BinaryTree) IsEmpty() bool {
	if binaryTree.GetLength() <= 0 {
		return true
	}
	return false
}

// 获取当前树的：节点总数
func (binaryTree *BinaryTree) GetLength() int {
	return binaryTree.Len
}

// 获取根节点元素-指针
func (binaryTree *BinaryTree) GetRootNode() *TreeNode {
	return binaryTree.RootNode
}

// 创建一个新的节点
func (binaryTree *BinaryTree) NewOneNode(keyword int, data interface{}) *TreeNode {
	treeNode := TreeNode{
		Keyword: keyword,
		Data:    data,
		Left:    nil,
		Right:   nil,
		Parent:  nil,
	}
	return &treeNode
}

// 插入一个新节点
// compare:共计比对多少次，才执行了插入操作
func (binaryTree *BinaryTree) InsertOneNode(keyword int, data interface{}) (compare int, err error) {
	binaryTree.Print("InsertOneNode  keyword:", keyword)

	if binaryTree.GetLength() >= binaryTree.NodeMax {
		msg := "NodeMax > " + strconv.Itoa(binaryTree.NodeMax)
		return compare, binaryTree.makeError(msg)
	}

	newNode := binaryTree.NewOneNode(keyword, data)
	if binaryTree.IsEmpty() {
		binaryTree.RootNode = newNode
		binaryTree.Len = 1
		return compare, nil
	}

	node := binaryTree.GetRootNode()
	searchNode, direction, compare, err := binaryTree.InsertOneNodeRecursionCompare(node, newNode, 0, nil, 0)
	if err != nil {
		return compare, err
	}
	//if searchNode == nil{
	//	binaryTree.Print("InsertOneNodeRecursionCompare return: ",searchNode, " dir:" ,direction)
	//}else{
	//	binaryTree.Print("InsertOneNodeRecursionCompare return: ",searchNode.Keyword," dir:",direction)
	//}

	if searchNode == nil { //只有根节点,因为根节点的父节点为空
		searchNode = node
	}
	if direction == DIRECTION_LEFT {
		if searchNode.Left != nil { //表示：新节点应该插入到该节点与该节点的子节点之间
			searchNode.Left.Parent = newNode
		}
		searchNode.Left = newNode

	} else {
		if searchNode.Right != nil { //表示：新节点应该插入到该节点与该节点的子节点之间
			searchNode.Right.Parent = newNode
		}
		searchNode.Right = newNode
	}
	newNode.Parent = searchNode

	binaryTree.Print("compare times:", compare, " downNode:", searchNode.Keyword)

	binaryTree.Len++
	//if binaryTree.Flag == FLAG_BALANCE {
	//	binaryTree.CheckInsertBalance(newNode)
	//}

	return compare, nil
}

// 插入节点时，递归查找新元素应该插入到哪个元素的左右
// node:当前比对节点
// insertNode:要插入的新节点
// direction:算是个上下文保留值，记录 当前节点是从上个节点的 左/右 方向过来的，因为最后一次循环，肯定节点是nil
// parentNode:算是个上下文保留值，记录 当前节点的父节点，因为最后一次循环，肯定节点是nil
// compare:比较次数
// 返回值：
// parentNode:父节点
// nodeDirection:从父节点的哪个方向过来的
// compareTimes:一共经常了多少次比较
func (binaryTree *BinaryTree) InsertOneNodeRecursionCompare(node *TreeNode, insertNode *TreeNode, direction int, parentNode *TreeNode, compare int) (downNode *TreeNode, nodeDirection int, compareTimes int, err error) {
	if node == nil {
		return parentNode, direction, compare, nil
	}
	compare++
	if insertNode.Keyword < node.Keyword {
		return binaryTree.InsertOneNodeRecursionCompare(node.Left, insertNode, DIRECTION_LEFT, node, compare)
	} else if insertNode.Keyword > node.Keyword {
		return binaryTree.InsertOneNodeRecursionCompare(node.Right, insertNode, DIRECTION_RIGHT, node, compare)
	} else {
		msg := "NodeKeyword: not allow repeat ."
		return downNode, direction, compare, binaryTree.makeError(msg)
	}
}

func (binaryTree *BinaryTree) FindOneByKeyword(keyword int) (downNode *TreeNode, empty bool, compare int) {
	binaryTree.Print("FindOneByKeyword:", keyword)
	if binaryTree.IsEmpty() {
		return downNode, true, compare
	}
	node := binaryTree.GetRootNode() //从根节点开始查找
	searchNode, empty, compare := binaryTree.FindOneByKeywordRecursionCompare(keyword, node, 0)
	binaryTree.Print("FindOneByKeyword compare times:", compare)
	return searchNode, empty, compare
}

// 根据 keyword 递归查找一个节点
func (binaryTree *BinaryTree) FindOneByKeywordRecursionCompare(keyword int, node *TreeNode, compare int) (rsNode *TreeNode, empty bool, compareTimes int) {
	if node.Keyword == keyword {
		return node, false, compareTimes
	}

	if node == nil {
		return node, true, compareTimes
	}
	compareTimes++
	if keyword > node.Keyword {
		return binaryTree.FindOneByKeywordRecursionCompare(keyword, node.Right, compareTimes)
	} else if keyword < node.Keyword {
		return binaryTree.FindOneByKeywordRecursionCompare(keyword, node.Left, compareTimes)
	} else { //这里是出现了 相等的情况，按说insert时做了限制，但为防万一，还是加条输出吧
		binaryTree.makeError("FindOneByKeywordRecursionCompare in case :else...")
	}
	return rsNode, true, compareTimes
}

// 获取当前树的深度/高度

func (binaryTree *BinaryTree) GetDeep(flag int) int {
	if binaryTree.IsEmpty() {
		return 0
	}
	var nodeList map[int][]*TreeNode
	if flag == 1 {
		_, nodeList = binaryTree.EachDeepByBreadthFirst(false)
		return len(nodeList)
	} else {
		deep := binaryTree.EachDeepByDeepFirst(binaryTree.GetRootNode())
		return deep
	}
}

// 层级遍历-深度优先
func (binaryTree *BinaryTree) EachDeepByDeepFirst(node *TreeNode) int {
	if node == nil {
		return 0
	}
	left := binaryTree.EachDeepByDeepFirst(node.Left)
	right := binaryTree.EachDeepByDeepFirst(node.Right)
	node.LeftDeep = left
	node.RightDeep = right

	result := 0
	if left > right {
		result = left + 1
	} else {
		result = right + 1
	}

	node.DeepDesc = result

	return result
}

// 层级遍历-广度优先
// nodeNilFill:空节点的元素，是否需要填充，主要是方便打印输出
// 借助队列，压入一个节点，然后弹出，并保存该节点，同时把该节点的左右节点再压入队列，重复此操作
func (binaryTree *BinaryTree) EachDeepByBreadthFirst(nodeNilFill bool) (empty bool, finalNode map[int][]*TreeNode) {
	//保存最终结果   层级/该层级下面的所有节点
	nodeContainer := make(map[int][]*ListNode)
	if binaryTree.IsEmpty() {
		return true, finalNode
	}
	//创建一个 无序 队列(数组类型)
	list := NewQueue(binaryTree.NodeMax, STACK_FLAG_ARRAY, ORDER_NONE, 0)
	//先把首节点压入队列中
	firstNode := binaryTree.GetRootNode()
	list.Push(firstNode.Keyword, firstNode)
	//当前遍历层级
	level := 0
	for {
		level++
		//一次进队列/出队列，遍历出来的一层的所有数据列表
		var nodeList []*ListNode
		//每次弹出一个节点，保存，后面再把该节点的左右节点继续压到队列中。
		//执行一次FOR 证明，遍历完成一层
		for {
			isEmpty, queueNode := list.Pop()
			if isEmpty {
				break
			}
			if queueNode != nil {
				treeNode, ok := queueNode.Data.(*TreeNode)
				if ok {
					treeNode.DeepAsc = level
				}
			}
			nodeContainerOne, ok := nodeContainer[level]
			if ok {
				nodeContainer[level] = append(nodeContainerOne, queueNode)
			} else {
				nodeContainer[level] = []*ListNode{queueNode}
			}
			//保存本次弹出的节点
			nodeList = append(nodeList, queueNode)
		}
		if !nodeNilFill {
			if len(nodeList) <= 0 { //证明没有任何节点了
				break
			}
		} else {
			isEmpty := true
			for _, v := range nodeList {
				if v.Keyword != KEYWORD_NIL {
					isEmpty = false
					break
				}
			}
			if isEmpty {
				break
			}
		}
		if nodeNilFill {
			for _, node := range nodeList {
				//binaryTree.Print("level:",level, " k:", k , " ", node,node.Keyword)
				leftKeyword := KEYWORD_NIL
				rightKeyword := KEYWORD_NIL

				if node == nil {
					list.Push(leftKeyword, nil)
					list.Push(rightKeyword, nil)
					continue
				}

				if node.Keyword == KEYWORD_NIL {
					//list.Push(leftKeyword, nil)
					//list.Push(leftKeyword, nil)
					continue
				}

				treeNode := node.Data.(*TreeNode)
				if treeNode.Left != nil {
					list.Push(treeNode.Left.Keyword, treeNode.Left)
				} else {
					list.Push(leftKeyword, nil)
				}

				if treeNode.Right != nil {
					list.Push(treeNode.Right.Keyword, treeNode.Right)
				} else {
					list.Push(rightKeyword, nil)
				}
			}
		} else {
			//开始将上面弹出的节点的：左右子节点再重新压回到队列中
			for _, node := range nodeList {
				//for k,node:=range nodeList{
				//binaryTree.Print("level:",level, " k:", k , " ", node,node.Keyword)
				if node == nil { //空节点，直接丢弃
					continue
				}
				//这里是对interface 做 nil 判断
				dateValueOf := reflect.ValueOf(node.Data)
				if dateValueOf.IsNil() { //空节点，直接丢弃
					continue
				}
				//断言
				treeNode, ok := node.Data.(*TreeNode)
				if !ok {
					binaryTree.Print("assertions failed.")
					continue
				}
				if treeNode.Left != nil {
					list.Push(treeNode.Left.Keyword, treeNode.Left)
				}
				if treeNode.Right != nil {
					list.Push(treeNode.Right.Keyword, treeNode.Right)
				}
			}
		}
	}
	//
	finalNode = make(map[int][]*TreeNode)
	for k, nodeListRowArr := range nodeContainer {
		var finalNodeListArr []*TreeNode
		for _, nodeList := range nodeListRowArr {
			treeNode, ok := nodeList.Data.(*TreeNode)
			if !ok {
				treeNode = nil
			}
			finalNodeListArr = append(finalNodeListArr, treeNode)
		}
		finalNode[k] = finalNodeListArr
	}

	if nodeNilFill { //占位符模式，有个问题：会多出来一层
		delete(finalNode, level)
	}

	for i, v := range finalNode {
		util.MyPrint(i)
		for _, v2 := range v {
			//fmt.Print(k,v,k2)
			if v2 == nil {
				//binaryTree.Print()
			} else {
				parent := 0
				if v2.Parent != nil {
					parent = v2.Parent.Keyword
				}
				fmt.Println(i, v2.Keyword, " parent:", parent)
			}
		}

	}

	return false, finalNode
}

// 传递值时，得用到 引用传参，能保存每次计算后的结果
// 又因为是递归，并且是数组，各种地址参数重复改变一个指针的值
// 所以这里，每次对数据的操作都是指针操作，得特殊处理下
func InsertIntArray(list *[]*TreeNode, treeNode *TreeNode) {
	for i := 0; i < len((*list)); i++ {
		if (*list)[i] == nil {
			(*list)[i] = treeNode
			break
		}
	}
}

// 递归根据方向：遍历树
func (binaryTree *BinaryTree) EachByOrder(node *TreeNode, list *[]*TreeNode, order int) {
	if binaryTree.IsEmpty() {
		return
	}

	if node == nil {
		return
	}

	if order == 2 {
		binaryTree.EachByOrder(node.Left, list, order)
		InsertIntArray(list, node)
		//binaryTree.Print(node.Keyword)
		binaryTree.EachByOrder(node.Right, list, order)
	} else if order == 1 {
		InsertIntArray(list, node)
		//binaryTree.Print(node.Keyword)
		binaryTree.EachByOrder(node.Left, list, order)
		binaryTree.EachByOrder(node.Right, list, order)
	} else if order == 3 {
		binaryTree.EachByOrder(node.Left, list, order)
		binaryTree.EachByOrder(node.Right, list, order)
		InsertIntArray(list, node)
		//binaryTree.Print(node.Keyword)
	} else {
		binaryTree.makeError("EachByOrder case ")
	}
}

// 动态创建一个数组，并把所有元素值：初始化为 KEYWORD_NIL
func GetNewIntArrayAndFillEmpty(len int) []*TreeNode {
	list := []*TreeNode{}
	for i := 0; i < len; i++ {
		list = append(list, nil)
	}
	return list
}

// 先序遍历
func (binaryTree *BinaryTree) EachByFirst() []*TreeNode {
	list := GetNewIntArrayAndFillEmpty(binaryTree.GetLength())
	node := binaryTree.GetRootNode()
	binaryTree.EachByOrder(node, &list, 1)

	return list
}

// 中序遍历
func (binaryTree *BinaryTree) EachByMiddle() []*TreeNode {
	list := GetNewIntArrayAndFillEmpty(binaryTree.GetLength())
	node := binaryTree.GetRootNode()
	binaryTree.EachByOrder(node, &list, 2)

	return list
}

// 后序遍历
func (binaryTree *BinaryTree) EachByAfter() []*TreeNode {
	list := GetNewIntArrayAndFillEmpty(binaryTree.GetLength())
	node := binaryTree.GetRootNode()
	binaryTree.EachByOrder(node, &list, 3)

	return list
}

// 输出信息，用于debug
func (binaryTree *BinaryTree) Print(a ...interface{}) (n int, err error) {
	if binaryTree.Debug > 0 {
		return fmt.Println(a)
	}
	return
}

// 创建一个error,统一管理
func (binaryTree *BinaryTree) makeError(msg string) error {
	binaryTree.Print("[errors] " + msg)
	return errors.New(msg)
}
