package smq

type entry struct {
	word  string
	code  string
	order int
}

// 顺序匹配
type order []entry

func NewOrder() *order {
	o := make(order, 0, 9999)
	return &o
}

func (o *order) Insert(word, code string, order int) {
	*o = append(*o, entry{word, code, order})
}

func (o order) Handle() {
}

// 顺序匹配
func (o order) Match(text []rune, p int) (int, string, int) {
	for _, v := range o {
		word := []rune(v.word)
		if p+len(word) > len(text) {
			continue
		}
		if v.word == string(text[p:p+len(word)]) {
			return len(word), v.code, v.order
		}
	}
	return 0, "", 1
}
