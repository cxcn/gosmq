package smq

type trie struct {
	children map[rune]*trie
	code     string
}

func newTrie() *trie {
	root := new(trie)
	root.children = make(map[rune]*trie)
	return root
}

func (t *trie) insert(word, code string) {
	for _, v := range word {
		if t.children[v] == nil {
			//子节点
			t.children[v] = newTrie()
		}
		t = t.children[v]
	}
	if t.code == "" {
		t.code = code
	} else if len(code) <= len(t.code) {
		t.code = code
	}
}
