package smq

type code_order struct {
	code  string
	order int
}

// 最长匹配
type longest struct {
	max int
	m   map[int]map[string]code_order
}

func NewLongest() *longest {
	ret := new(longest)
	ret.m = make(map[int]map[string]code_order)
	return ret
}

func (l *longest) Insert(word, code string, order int) {
	length := len([]rune(word))
	if l.m[length] == nil {
		l.m[length] = make(map[string]code_order)
	}
	// 不替换原有的
	if l.m[length][word].code == "" {
		l.m[length][word] = code_order{code, order}
	}
}

func (l *longest) Handle() {
	for k := range l.m {
		if k > l.max {
			l.max = k
		}
	}
	for i := 1; i < l.max; i++ {
		if l.m[i] == nil {
			l.m[i] = make(map[string]code_order)
		}
	}
}

// 最长匹配
func (l *longest) Match(text []rune, p int) (int, string, int) {
	max := l.max
	if len(text)-p < max {
		max = len(text) - p
	}
	for i := max; i > 0; i-- {
		v, ok := l.m[i][string(text[p:p+i])]
		if ok {
			return i, v.code, v.order
		}
	}
	return 0, "", 1
}
