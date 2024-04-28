package container

import (
	"fmt"
)

// 前缀树/字典树
// 暂不支持中文
type TrieTree struct {
	RootNode TrieNode
	Debug    int
}

type TrieNode struct {
	End   bool
	Value int32 //golang 没有字符类型
}

func NewTrieTree() *TrieTree {
	trieTree := new(TrieTree)
	trieTree.Debug = 1

	trieTree.InsertOne("apple")
	return trieTree
}

func (trieTree *TrieTree) InsertOne(word string) {
	if trieTree.IsEmpty() {
		for k, v := range word {
			trieTree.Print(k, v)
			trieNode := TrieNode{
				End:   false,
				Value: v,
			}
			if k == len(word) {
				trieNode.End = true
			}

			if k == 0 {
				trieTree.RootNode = trieNode
			}
		}
	}

}

func (trieTree *TrieTree) IsEmpty() bool {
	return true
}

// 输出信息，用于debug
func (trieTree *TrieTree) Print(a ...interface{}) (n int, err error) {
	if trieTree.Debug > 0 {
		return fmt.Println(a)
	}
	return
}
