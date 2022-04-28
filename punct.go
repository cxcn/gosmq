package smq

import (
	"strconv"
)

var puncts = genPuncts()

func genPuncts() *order {
	ret := NewOrder()

	// 一般的
	en := "`-=[];',./"
	cn := []rune(`·-=【】；‘，。、`)
	for i, v := range en {
		ret.Insert(string(v), string(v), 1)
		ret.Insert(string(cn[i]), string(v), 1)
	}
	ret.Insert("’", "'", 1)

	// shift =
	shiftEN := `~_+{}:"<>?)!@#$%^&*(`
	shiftCN := []rune(`~_+{}：“《》？）！@#￥%^&*（`)
	for i, v := range shiftEN {
		key := ""
		if i >= len(en) {
			key = strconv.Itoa(i - len(en))
		} else {
			key = string(en[i])
		}
		ret.Insert(string(v), "="+key, 1)
		ret.Insert(string(shiftCN[i]), "="+key, 1)
	}

	// 其他的
	ret.Insert("”", "='", 1)
	ret.Insert("——", "=-", 1)
	ret.Insert("……", "=6", 1)

	// 大小写字母
	for i := 0; i < 26; i++ {
		ret.Insert(string(byte(i+97)), string(byte(i+97)), 1)
		ret.Insert(string(byte(i+65)), "="+string(byte(i+65)), 1)
	}
	return ret
}
