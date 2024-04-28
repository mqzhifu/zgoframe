package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
	"zgoframe/util/container"
)

// @Tags Test
// @Summary 获取列表
// @Description 先序
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param flag path string true "1先序2中序3后序4全部"
// @Produce  application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /test/binary/tree/list/{flag} [get]
func BinaryTreeListByFlag(c *gin.Context) {

	flagStr := c.Param("flag")
	flag, _ := strconv.Atoi(flagStr)
	if flag <= 0 {
		flag = 4
	}

	if flag == 4 {
		firstList := global.V.BinaryTree.EachByFirst()
		middleList := global.V.BinaryTree.EachByMiddle()
		afterFlist := global.V.BinaryTree.EachByAfter()

		type RsList struct {
			First  []int `json:"first"`
			Middle []int `json:"middle"`
			After  []int `json:"after"`
		}

		rsList := RsList{}

		rsList.First = TreeNodeToArr(firstList)
		rsList.Middle = TreeNodeToArr(middleList)
		rsList.After = TreeNodeToArr(afterFlist)

		httpresponse.OkWithAll(rsList, "获取列表成功1", c)
	} else {
		list := container.GetNewIntArrayAndFillEmpty(global.V.BinaryTree.GetLength())
		node := global.V.BinaryTree.GetRootNode()
		global.V.BinaryTree.EachByOrder(node, &list, flag)

		rsList := TreeNodeToArr(list)
		httpresponse.OkWithAll(rsList, "获取列表成功2", c)
	}

}

func TreeNodeToArr(list []*container.TreeNode) []int {
	var rsList []int
	for _, v := range list {
		util.MyPrint(v.Keyword)
		rsList = append(rsList, v.Keyword)
	}

	return rsList
}

// @Tags Test
// @Summary 添加元素
// @Description 添加一个元素
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param keyword path string true "节点关键值"
// @Produce  application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /test/binary/tree/insert/one/{keyword} [get]
func BinaryTreeInsertOne(c *gin.Context) {

	//binaryTree := util.NewBinaryTree(100, 1, 1)
	//llen := binaryTree.GetLength()
	//util.MyPrint(llen)

	keywordStr := c.Param("keyword")
	keyword, _ := strconv.Atoi(keywordStr)
	//util.MyPrint("keyword:", keyword)
	compare, err := global.V.BinaryTree.InsertOneNode(keyword, "")

	util.MyPrint(compare, err)
	util.MyPrint(global.V.BinaryTree.GetLength())

	httpresponse.OkWithAll(compare, "添加成功", c)
}

// @Tags Test
// @Summary 打印整棵树
// @Description 根据深度
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /test/binary/tree/each/deep [get]
func BinaryTreeEachDeepByBreadthFirst(c *gin.Context) {

	empty, finalNode := global.V.BinaryTree.EachDeepByBreadthFirst(true)
	util.MyPrint(empty, finalNode)
	rs := make(map[int][]int)
	if !empty {
		for level, nodeList := range finalNode {
			var nodeArr []int
			for _, node := range nodeList {
				//util.MyPrint(node)
				if node == nil {
					nodeArr = append(nodeArr, 999)
				} else {
					nodeArr = append(nodeArr, node.Keyword)
				}
			}
			rs[level] = nodeArr
		}
	}
	//util.MyPrint(global.V.BinaryTree.GetLength())
	httpresponse.OkWithAll(rs, "输出树成功 ", c)
}
