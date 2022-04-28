package smq

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
)

type jidian struct{}

func (j *jidian) Read(dict *Dict) []byte {
	var buf bytes.Buffer

	scan := bufio.NewScanner(dict.reader)
	for scan.Scan() {
		wc := strings.Split(scan.Text(), " ")
		if len(wc) < 2 {
			continue
		}
		// 单字模式修正
		revise := 0
		for i := 1; i < len(wc); i++ {
			if dict.Single && len([]rune(wc[i])) != 1 {
				revise++
				continue
			}
			order := i - revise
			// 生成赛码表
			buf.WriteString(wc[i])
			buf.WriteByte('\t')
			buf.WriteString(wc[0])
			if len(wc[1]) >= dict.PushStart && order == 1 {
			} else {
				if int(order) <= len(dict.SelectKeys) {
					buf.WriteByte(dict.SelectKeys[order-1])
				} else {
					buf.WriteString(strconv.Itoa(int(order)))
				}
			}
			buf.WriteByte('\t')
			buf.WriteString(strconv.Itoa(int(order)))
			buf.WriteByte('\n')
		}
	}
	return buf.Bytes()
}
