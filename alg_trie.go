package smq

// trie 树
type trie struct {
	children map[rune]*trie
	code     string
	order    int
}

func NewTrie() *trie {
	return new(trie)
}

func (t *trie) Insert(word, code string, order int) {
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
		t.order = order
	}
}

func (t *trie) Handle() {}

// 前缀树最长匹配
func (t *trie) Match(text []rune, p int) (int, string, int) {
	j := 0 // 已匹配的字数
	i := 0 // 有编码的匹配
	dict := t
	code := "" // 编码
	order := 0
	for p+j < len(text) {
		dict = dict.children[text[p+j]]
		j++
		if dict == nil {
			break
		}
		if dict.code != "" {
			i = j
			code = dict.code
			order = dict.order
		}
	}
	return i, code, order
}
