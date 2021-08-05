package main

type Trie struct {
	children map[rune]*Trie
	// value    string
	code   string
	isWord bool
}

func Constructor() Trie {
	root := new(Trie)
	root.children = make(map[rune]*Trie)
	root.isWord = false
	return *root
}

func (trie *Trie) Insert(word, code string) {
	for _, v := range word {
		if trie.children[v] == nil {
			node := new(Trie)
			//子节点
			node.children = make(map[rune]*Trie)
			//初始化节点单词标志为假
			node.isWord = false
			trie.children[v] = node
		}
		trie = trie.children[v]
	}
	trie.code = code
	trie.isWord = true
}

// func (trie *Trie) Search(word string) (bool, string) {
// 	for _, v := range word {
// 		if trie.children[v] == nil {
// 			return false, ""
// 		}
// 		trie = trie.children[v]
// 	}
// 	return trie.isWord, trie.code
// }
