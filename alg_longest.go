package smq

type code_order struct {
	code  string
	order int
}

// 最长匹配
type longest struct {
	l []map[string]code_order
}

func NewLongest() *longest {
	l := new(longest)
	l.l = make([]map[string]code_order, 0, 30)
	return l
}

func (l *longest) Insert(word, code string, order int) {
	i := len([]rune(word)) // 词长
	for i+1 > len(l.l) {   // 扩容
		l.l = append(l.l, make(map[string]code_order))
	}
	// 不替换原有的
	if l.l[i][word].code == "" {
		l.l[i][word] = code_order{code, order}
	}
}

func (l *longest) Handle() {
}

// 最长匹配
func (l *longest) Match(text []rune, p int) (int, string, int) {
	max := len(l.l) - 1
	if len(text)-p < max {
		max = len(text) - p
	}
	for i := max; i > 0; i-- {
		v, ok := l.l[i][string(text[p:p+i])]
		if ok {
			return i, v.code, v.order
		}
	}
	return 0, "", 1
}
