package smq

import (
	_ "embed"
	"strconv"
	"strings"
)

type comb struct {
	eq  int  // 当量*10 equivalent
	sh  int  // 同手 same hand: 同键 小跨排 大跨排 错手
	lfd bool // 小指干扰 little finger disturb
}

//go:embed assets/equivalent.txt
var equivalent string

//go:embed assets/fingering.txt
var fingering string

var combData = newCombData()

func newCombData() map[string]*comb {

	var ret = make(map[string]*comb, 1800)

	// 当量
	for _, v := range strings.Split(equivalent, "\n") {
		tmp := strings.Split(v, "\t")
		if len(tmp) != 2 {
			continue
		}
		c := new(comb)
		c.eq, _ = strconv.Atoi(tmp[1])
		ret[tmp[0]] = c
	}

	// 指法
	fg := strings.Split(fingering, "\n")
	// 小跨排
	xkp := strings.Split(fg[0], " ")
	for _, v := range xkp {
		ret[v].sh = 2
	}
	// 大跨排
	dkp := strings.Split(fg[1], " ")
	for _, v := range dkp {
		ret[v].sh = 3
	}
	// 错手
	cs := strings.Split(fg[2], " ")
	for _, v := range cs {
		ret[v].sh = 4
	}
	// 小指干扰
	lf := strings.Split(fg[3], " ")
	for _, v := range lf {
		ret[v].lfd = true
	}
	return ret
}
