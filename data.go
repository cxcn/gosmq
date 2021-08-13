package main

import (
	_ "embed"
	"fmt"
	"io"
	"strings"
)

//go:embed data\\dangliang
var dangliang string

type Zhifa struct {
	dl float64 // 当量
	zf int     // 指法：大小跨排等
	hj int     // 互击 LR RL LL RR
	tz bool    // 同指
}

func newZhifa(isS bool) map[string]*Zhifa {

	var zhifa = make(map[string]*Zhifa, 999)
	r := strings.NewReader(dangliang)
	var key string
	var dl float64
	for {
		_, err := fmt.Fscanln(r, &key, &dl)
		if err == io.EOF {
			break
		}
		z := new(Zhifa)
		z.dl = dl
		zhifa[key] = z
	}

	lr := make(map[byte]bool)
	lKey := "12345qwertasdfgzxcvb"
	for i := 0; i < len(lKey); i++ {
		lr[lKey[i]] = true
	}
	rKey := "67890yuiophjkl;'nm,./"
	for i := 0; i < len(rKey); i++ {
		lr[rKey[i]] = false
	}

	for k, v := range zhifa {

		if k[0] == k[1] {
			v.tz = true
		}
		if k[0] == '_' {
			if isS {
				v.hj = 1
			}
			continue
		}
		if k[1] == '_' {
			if isS {
				v.hj = 2
			}
			continue
		}

		if lr[k[0]] && !lr[k[1]] {
			v.hj = 1 // LR
		} else if !lr[k[0]] && lr[k[1]] {
			v.hj = 2 // RL
		} else if lr[k[0]] && lr[k[1]] {
			v.hj = 3 // LL
		} else if !lr[k[0]] && !lr[k[1]] {
			v.hj = 4 // RR
		}

	}

	d := []string{"br", "bt", "ce", "ec", "mu", "my", "nu", "ny",
		"p/", "qz", "rb", "rv", "tb", "tv", "um", "un", "vr", "vt",
		"wx", "xw", "ym", "yn", "zq", ",i", "/p"}
	x := []string{"qa", "za", "fb", "gb", "vb", "dc", "cd", "ed",
		"de", "bf", "gf", "rf", "tf", "vf", "bg", "fg", "rg", "tg",
		"vg", "jh", "mh", "nh", "uh", "yh", "ki", "hj", "mj", "nj",
		"uj", "yj", "ik", "ol", "hm", "jm", "nm", "hn", "jn", "mn",
		"lo", ";p", "aq", "fr", "gr", "tr", "ws", "xs", "ft", "gt",
		"rt", "hu", "ju", "yu", "bv", "fv", "gv", "sw", "sx", "hy",
		"jy", "uy", "az", "k,", ";/", "p;", "/;"}
	gr := []string{"aa", "ac", "ad", "ae", "aq", "as", "aw", "ax",
		"az", "ca", "cq", "cz", "da", "dq", "dz", "ea", "eq", "ez",
		"ip", "i/", "i;", "kp", "k/", "k;", "lp", "l/", "l;", "op",
		"o/", "o;", "pi", "pk", "pl", "po", "pp", "p;", "qa", "qc",
		"qd", "qe", "qq", "qs", "qw", "qx", "sa", "sq", "sz", "wa",
		"wq", "wz", "xa", "xq", "xz", "za", "zc", "zd", "ze", "zs",
		"zw", "zx", "zz", ",p", ",/", ",;", "/i", "/k", "/l", "/o",
		"//", "/;", ";i", ";k", ";l", ";o", ";p", ";/", ";;"}
	c := []string{"ct", ",y", "tc", "y,", "cr", ",u", "rc", "u,",
		"cw", ",o", "wc", "o,", "qc", ",p", "cq", "p,", "qx", "p.",
		"xq", ".p", "xe", ".i", "ex", "i.", "xr", ".u", "rx", "u.",
		"xt", ".y", "tx", "y."}

	for _, v := range d {
		zhifa[v].zf = 1
	}
	for _, v := range x {
		zhifa[v].zf = 2
	}
	for _, v := range gr {
		zhifa[v].zf = 3
	}
	for _, v := range c {
		zhifa[v].zf = 4
	}
	return zhifa
}
