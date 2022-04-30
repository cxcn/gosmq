package smq

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
)

type duoduo struct {
	reverse bool
}

func (j *duoduo) Read(dict *Dict) []byte {
	var buf bytes.Buffer
	buf.Grow(1e6)
	mapOrder := make(map[string]int)

	scan := bufio.NewScanner(dict.reader)
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) < 2 {
			continue
		}
		var w, c string
		if j.reverse {
			w, c = wc[1], wc[0]
		} else {
			w, c = wc[0], wc[1]
		}

		mapOrder[c]++
		order := mapOrder[c]
		// 生成赛码表
		buf.WriteString(w)
		buf.WriteByte('\t')
		buf.WriteString(c)
		if len(c) >= dict.PushStart && order == 1 {
		} else {
			if order <= len(dict.SelectKeys) {
				buf.WriteByte(dict.SelectKeys[order-1])
			} else {
				buf.WriteString(strconv.Itoa(order))
			}
		}
		buf.WriteByte('\t')
		buf.WriteString(strconv.Itoa(order))
		buf.WriteByte('\n')
	}
	// fmt.Println(string(buf.Bytes()))
	return buf.Bytes()
}
