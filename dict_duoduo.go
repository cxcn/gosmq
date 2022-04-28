package smq

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
)

type duoduo struct{}

func (j *duoduo) Read(dict *Dict) []byte {
	var buf bytes.Buffer
	mapOrder := make(map[string]int)

	scan := bufio.NewScanner(dict.reader)
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) < 2 {
			continue
		}
		mapOrder[wc[1]]++
		order := mapOrder[wc[1]]
		// 生成赛码表
		buf.WriteString(scan.Text())
		if len(wc[1]) >= dict.PushStart && order == 1 {
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
	return buf.Bytes()
}
