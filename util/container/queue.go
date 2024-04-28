package container

type Queue struct {
	Max   int
	Flag  int //数组|链表
	Order int
	List  ListInterface
}

func NewQueue(max int, flag int, order int, debug int) *Queue {
	queue := new(Queue)
	queue.Max = max
	queue.Flag = flag
	if flag == STACK_FLAG_ARRAY {
		queue.List = NewArrayList(order, max, debug)
	} else if flag == STACK_FLAG_LINKED_LIST {
		queue.List = NewLinkedList(order, max, false, 1)
	}
	return queue
}

// 队首压入
func (queue *Queue) Push(keyword int, data interface{}) (int, error) {
	return queue.List.InsertNodeByFirst(keyword, data)
}

// 队首压入
func (queue *Queue) PushByEnd(keyword int, data interface{}) (int, error) {
	return queue.List.InsertNodeByEnd(keyword, data)
}

// 队尾弹出，push pop 属于快捷方法
func (queue *Queue) Pop() (empty bool, searchNode *ListNode) {
	return queue.List.FindOneNodeByLocationAndDel(queue.List.Length() - 1)
}
func (queue *Queue) PopByFirst() (empty bool, searchNode *ListNode) {
	return queue.List.FindOneNodeByLocationAndDel(0)
}

// 获取尾部一个节点
func (queue *Queue) GetOneByEnd() (empty bool, searchNode *ListNode) {
	return queue.List.FindOneNodeByLocation(queue.List.Length() - 1)
}

func (queue *Queue) IsEmpty() bool {
	return queue.List.IsEmpty()
}
