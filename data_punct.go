package smq

import (
	"strconv"
)

var puncts = genPuncts()

func genPuncts() map[string]string {
	ret := make(map[string]string)

	// 一般的
	en := "`-=[];',./"
	cn := []rune(`·-=【】；‘，。、`)
	for i, v := range en {
		ret[string(v)] = string(v)
		ret[string(cn[i])] = string(v)
	}
	ret["’"] = "'"

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
		ret[string(v)] = "=" + key
		ret[string(shiftCN[i])] = "=" + key
	}

	// 其他的
	ret["”"] = "='"
	ret["——"] = "=-"
	ret["……"] = "=6"

	// 大小写字母
	for i := 0; i < 26; i++ {
		ret[string(byte(i+97))] = string(byte(i + 97))
		ret[string(byte(i+65))] = "=" + string(byte(i+65))
	}
	// 数字
	for i := 0; i < 10; i++ {
		ret[strconv.Itoa(i)] = strconv.Itoa(i)
	}

	return ret
}
