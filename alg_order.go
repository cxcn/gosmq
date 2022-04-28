package smq

type entry struct {
	word  string
	code  string
	order int
}

// 顺序匹配
type order struct {
	o []entry
}

func NewOrder() *order {
	return new(order)
}

func (o *order) Insert(word, code string, order int) {
	o.o = append(o.o, entry{word, code, order})
}

func (o order) Handle() {}

// 顺序匹配
func (o order) Match(text []rune, p int) (int, string, int) {
	for _, v := range o.o {
		if v.word == string(text[p:p+len(v.word)]) {
			return len(v.word), v.code, v.order
		}
	}
	return 0, "", 1
}
