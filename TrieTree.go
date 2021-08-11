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
		if t.children[v] == nil {
			//子节点
			t.children[v] = NewTrie()
		}
		t = t.children[v]
	}
	if t.code == "" {
		t.code = code
	}
}
