package smq

import (
	"bufio"
	"strconv"
	"strings"
)

func (dict *Dict) fromJidian() {
	t := new(trie)
	scan := bufio.NewScanner(dict.reader)
	var wb []byte
	// 生成字典
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
			wb = append(wb, scan.Bytes()...)
			if len(wc[1]) >= dict.PushStart && order == 1 {
			} else {
				if int(order) <= len(dict.SelectKeys) {
					wb = append(wb, dict.SelectKeys[order-1])
				} else {
					wb = append(wb, []byte(strconv.Itoa(int(order)))...)
				}
			}
			wb = append(wb, '\t')
			wb = append(wb, []byte(strconv.Itoa(int(order)))...)
			wb = append(wb, '\n')

			t.Insert(wc[0], wc[1], order)
			dict.length++
		}
	}
	// 添加符号
	for _, v := range puncts.o {
		t.Insert(v.word, v.code, v.order)
	}
	dict.Matcher = t
}
