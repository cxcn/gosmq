package main

type Trie struct {
	children map[rune]*Trie
	code     string
}

func NewTrie() *Trie {
	root := new(Trie)
	root.children = make(map[rune]*Trie)
	return root
}

func (trie *Trie) Insert(word, code string) {
	for _, v := range word {
		if trie.children[v] == nil {
			//子节点
			node := new(Trie)
			node.children = make(map[rune]*Trie)
			trie.children[v] = node
		}
		trie = trie.children[v]
	}
	trie.code = code
}
