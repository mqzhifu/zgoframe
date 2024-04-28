package container

import (
	"errors"
)

//集合|列表
//ps:这里把有序和无序做了兼容合并，程序上略有点不OOP，主要是为了实验，生产中：还是建议把有序跟无序分开来写

// 一个节点
type ListNode struct {
	Keyword  int         //关键值
	Data     interface{} //附加数据体
	Previous *ListNode   //上一个节点地址
	Next     *ListNode   //下一个节点地址
	Location int         //当前节点在容器中的位置，给数组结构使用
}

// 根据关键字，做一些范围性的搜索
type ListSearchCondition struct {
	IsEmpty   bool
	Offset    int
	Limit     int
	Keyword   int
	Compare   string
	Direction int
}

func ListSearchCompare(node *ListNode, listSearchCondition ListSearchCondition) (bool, error) {
	if listSearchCondition.Compare == "" {
		return true, nil
	}

	//if listSearchCondition.Keyword == 0{
	//	return true,nil
	//}

	rs := false
	if listSearchCondition.Compare == ">" {
		if node.Keyword > listSearchCondition.Keyword {
			rs = true
		}
	} else if listSearchCondition.Compare == ">=" {
		if node.Keyword >= listSearchCondition.Keyword {
			rs = true
		}
	} else if listSearchCondition.Compare == "<" {
		if node.Keyword < listSearchCondition.Keyword {
			rs = true
		}
	} else if listSearchCondition.Compare == "<=" {
		if node.Keyword <= listSearchCondition.Keyword {
			rs = true
		}
	} else {
		return rs, errors.New("listSearchCondition.Compare str err, eg:> >= < <=")
	}
	return rs, nil
}

const (
	//处理方向
	DIRECTION_FIRST = 1 //从首部开始
	DIRECTION_END   = 2 //从尾部开始

	//节点的存储顺序
	ORDER_NONE = 0 //无
	ORDER_ASC  = 1 //升序
	ORDER_DESC = 2 //降序

	NODE_MAX = 100 //一个容器内允许装载的节点最大数
	NODE_MIN = 10

	STACK_FLAG_ARRAY       = 1
	STACK_FLAG_LINKED_LIST = 1
)

type ListInterface interface {
	//判断列表是否为：空
	IsEmpty() bool
	//添加元素时，判断：是否超过当前链表节点最大数
	//不过，构造函数里如果做了容错，此方法就同意义了
	CheckMaxNode() error
	//当前列表中节点总数
	Length() int
	//从首部插入一个节点，有序链表：并不一定是从头插入，可能计算完后，在中间某个节点插入
	//此方法适用无序链表
	InsertNodeByFirst(keyword int, data interface{}) (int, error)
	//从尾部插入一个节点，有序链表：并不一定是从尾插入，可能计算完后，在中间某个节点插入
	//此方法适用无序链表
	InsertNodeByEnd(keyword int, data interface{}) (int, error)
	//从指定位置插入一个节点，后面的节点顺延
	//此方法适用无序链表
	InsertNodeByLocation(location int, keyword int, data interface{}) (int, error)
	//从指定关键字位置：上方，插入一个节点，后面的节点顺延
	//此方法适用无序链表
	InsertNodeByKeyword(searchKeyword, keyword int, data interface{}) (int, error)
	//一次插入多节点,无序的好处理，有序的还牵扯到排序，有点麻烦，先不写了,另外，有单节点插入可以替代
	//此方法适用无序链表
	InsertMultiNode([]ListNode) (int, error)
	//location:在哪个位置点插入 , 当：< 证明为空
	//direction:在位置的：上方/下方插入该结点
	//searchKeyword:先根据关键字找到该节点(暂时仅支持正整数)，不允许为负数，当<=0 时，证明该参数为空
	//ps : location  与 searchKeyword 互斥 , 有序下 两个参数均不允许用
	InsertNode(direction int, location int, searchKeyword int, keyword int, data interface{}) (int, error)
	//删除头部节点
	DelFirstNode() (node *ListNode, err error)
	//删除尾部节点
	DelEndNode() (node *ListNode, err error)
	//根据关键字删除一个节点
	DelNodeByKeyword(keyword int) (node *ListNode, err error)
	//指定位置，删除一个节点
	DelOneNodeByLocation(direction int, location int, limit int) (node []*ListNode, err error)
	//删除一个节点：根据一个结构体
	DelOneNode(listNode *ListNode) error
	//寻找新节点的插入位置：有序列表，插入时，得找到符合条件的位置。如下3种情况：
	//1.最大端：没有搜索到合适的节点，降序时：该值为最大，找不到小于它的值，所以 ，应该插入队首，升序时：该值为最小，也应该插入到队首，所以：该值如果存在，即证明新节点应该插入队首
	//2.最小端：没有搜索到合适的节点，降序时：该值为最小，可以找到节点，但该节点为队尾，所以：该值如果存在，即证明新节点应该插入队尾
	//3.中间：排除掉以上2点特殊情况后，就是：某个节点的上下现位置
	//所以，该函数返回：就得把3种情况都显示
	FindOneNodeInsertLocationByKeyword(keyword int) (max bool, min bool, searchNode *ListNode, err error)
	//根据 关键词 查找一个节点
	FindOneNodeByKeyword(keyword int) (empty bool, searchNode *ListNode)
	//根据 关键词 查找一个节点 ，如果找到返回该节点，同时删除该节点
	//注意：因为keyword是有重复的可能，一次是返回多个节点
	FindOneNodeByKeywordAndDel(keyword int) (empty bool, searchNode []*ListNode)
	//从某个位置 获取一个节点，只支持从上到下的顺序~
	FindOneNodeByLocation(location int) (empty bool, searchNode *ListNode)
	//从某个位置 获取一个节点，只支持从上到下的顺序~如果找到返回该节点，同时删除该节点
	FindOneNodeByLocationAndDel(location int) (empty bool, searchNode *ListNode)
	//中间节点地址，偶数指向两个中位节点：其中一个
	GetMiddleNode() (empty bool, node *ListNode)
	//获取链表所有节点，从头部向下开始
	GetAllByFirst(listSearchCondition ListSearchCondition) (empty bool, nodeList []*ListNode)
	//获取链表所有节点，从尾部向上开始
	GetAllByEnd(listSearchCondition ListSearchCondition) (empty bool, nodeList []*ListNode)
	//获取链表所有节点，根据方向 和 搜索条件
	GetAll(listSearchCondition ListSearchCondition) (empty bool, nodeList []*ListNode)
	//统计链表中重复keyword的情况
	NodeRepeatTotal() (repeatList map[int]int, empty bool)
}

// 创建一个新的节点
func NewListNode(keyword int, data interface{}) *ListNode {
	//linkedList.Print("NewListNode keyword:",keyword , " data:",data)
	ListNode := new(ListNode)
	ListNode.Keyword = keyword
	ListNode.Data = data
	return ListNode
}
