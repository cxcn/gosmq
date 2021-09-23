package smq

import (
	_ "embed"
	"fmt"
	"io"
	"strings"
)

//go:embed assets\\punct
var punct string

//go:embed assets\\equivalent
var equivalent string

//go:embed assets\\fingering
var fingering string

type comb struct {
	eq   float64 // 当量 equivalent
	dist int     // 分布 distribution: LR RL LL RR
	sh   int     // 同手 same hand: 同键 小跨排 大跨排 错手
	lfd  bool    // 小指干扰 little finger disturb
}

func (t *trie) addPunct() {
	r := strings.NewReader(punct)
	var p string // punct
	var k string // key
	for {
		_, err := fmt.Fscanln(r, &p, &k)
		if err == io.EOF {
			break
		}
		t.insert(p, k)
	}
	// ascii
	var i byte = 33
	for i < 127 {
		t.insert(string(i), string(i))
		i++
	}
}

// allow space -> map
func newCombs(aS bool) map[string]*comb {

	var c = make(map[string]*comb, 999)

	// 当量
	r := strings.NewReader(equivalent)
	var key string
	var eq float64 // equivalent
	for {
		_, err := fmt.Fscanln(r, &key, &eq)
		if err == io.EOF {
			break
		}
		tmp := new(comb)
		tmp.eq = eq
		c[key] = tmp
	}

	// 互击
	isL := make(map[byte]bool)   // is left keys
	lk := "12345qwertasdfgzxcvb" // left keys
	for i := 0; i < len(lk); i++ {
		isL[lk[i]] = true
	}
	rk := "67890yuiophjkl;'nm,./" // right keys
	for i := 0; i < len(rk); i++ {
		isL[rk[i]] = false
	}
	for k, v := range c {
		// 同键
		if k[0] == k[1] {
			v.sh = 1
		}
		// 处理空格
		if k[0] == '_' {
			if aS {
				v.dist = 1
			}
			continue
		}
		if k[1] == '_' {
			if aS {
				v.dist = 2
			}
			continue
		}

		if isL[k[0]] && !isL[k[1]] {
			v.dist = 1 // LR
		} else if !isL[k[0]] && isL[k[1]] {
			v.dist = 2 // RL
		} else if isL[k[0]] && isL[k[1]] {
			v.dist = 3 // LL
		} else if !isL[k[0]] && !isL[k[1]] {
			v.dist = 4 // RR
		}
	}

	// 指法
	fg := strings.Split(fingering, "\r\n")
	// 小跨排
	xkp := strings.Split(fg[0], " ")
	for _, v := range xkp {
		c[v].sh = 2
	}
	// 大跨排
	dkp := strings.Split(fg[1], " ")
	for _, v := range dkp {
		c[v].sh = 3
	}
	// 错手
	cs := strings.Split(fg[2], " ")
	for _, v := range cs {
		c[v].sh = 4
	}
	// 小指干扰
	lf := strings.Split(fg[3], " ")
	for _, v := range lf {
		c[v].lfd = true
	}
	return c
}
