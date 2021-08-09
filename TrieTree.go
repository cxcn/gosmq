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

func (t *Trie) Insert(word, code string) {
	for _, v := range word {
		if _, ok := t.children[v]; !ok {
			//子节点
			node := NewTrie()
			t.children[v] = node
			t = node
		} else {
			t = t.children[v]
		}
	}
	if len(t.code) == 0 {
		t.code = code
	}
}
