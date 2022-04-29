package smq

import "sort"

// trie 树
type oTrie struct {
	children map[rune]*oTrie
	code     string
	order    int
	line     int
}

func NewOTrie() *oTrie {
	return new(oTrie)
}

var orderLine = 0

func (t *oTrie) Insert(word, code string, order int) {
	for _, v := range word {
		if t.children == nil {
			t.children = make(map[rune]*oTrie)
			t.children[v] = new(oTrie)
		} else if t.children[v] == nil {
			t.children[v] = new(oTrie)
		}
		t = t.children[v]
	}
	if t.code == "" {
		t.code = code
		t.order = order
		orderLine++
		t.line = orderLine
	}
}

func (t *oTrie) Handle() {}

// 前缀树按码表序匹配
func (t *oTrie) Match(text []rune, p int) (int, string, int) {
	j := 0 // 已匹配的字数
	i := 0 // 有编码的匹配
	dict := t
	type res_tmp struct {
		i     int
		code  string
		order int
		line  int
	}
	res := make([]res_tmp, 0, 10)
	for p+j < len(text) {
		dict = dict.children[text[p+j]]
		j++
		if dict == nil {
			break
		}
		if dict.code != "" {
			i = j
			res = append(res, res_tmp{i, dict.code, dict.order, dict.line})
		}
	}
	if len(res) == 0 {
		return 0, "", 1
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].line < res[j].line
	})
	return res[0].i, res[0].code, res[0].order
}
