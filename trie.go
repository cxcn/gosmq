package smq

type trie struct {
	children map[rune]*trie
	code     string
}

func (t *trie) insert(word, code string) {
	for _, v := range word {
		if t.children == nil {
			t.children = make(map[rune]*trie)
			t.children[v] = new(trie)
		} else if t.children[v] == nil {
			t.children[v] = new(trie)
		}
		t = t.children[v]
	}
	if t.code == "" {
		t.code = code
	} else if len(code) <= len(t.code) {
		t.code = code
	}
}
