package smq

import "sort"

type code_order struct {
	code  string
	order int
}

type length_m struct {
	length int
	m      map[string]code_order
}

// 最长匹配
type longest struct {
	m     map[int]map[string]code_order
	slice []length_m
}

func NewLongest() *longest {
	ret := new(longest)
	ret.m = make(map[int]map[string]code_order)
	ret.slice = make([]length_m, 0, 10)
	return ret
}

func (l *longest) Insert(word, code string, order int) {
	length := len([]rune(word))
	l.m[length][word] = code_order{code, order}
}

func (l *longest) Handle() {
	for k, v := range l.m {
		l.slice = append(l.slice, length_m{k, v})
	}
	sort.Slice(l.slice, func(i, j int) bool {
		return l.slice[i].length > l.slice[j].length
	})
	l.m = nil
}

// 最长匹配
func (l *longest) Match(text []rune, p int) (int, string, int) {
	for _, lm := range l.slice {
		for k, v := range lm.m {
			if k == string(text[p:p+lm.length]) {
				return lm.length, v.code, v.order
			}
		}
	}
	return 0, "", 1
}
