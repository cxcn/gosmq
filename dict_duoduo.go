package smq

import (
	"bufio"
	"strconv"
	"strings"
)

func (dict *Dict) fromDuoduo() {
	t := new(trie)
	scan := bufio.NewScanner(dict.Reader)
	mapOrder := make(map[string]int)
	var wb []byte
	// 生成字典
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) != 2 {
			continue
		}
		if dict.Single && len([]rune(wc[0])) != 1 {
			continue
		}

		mapOrder[wc[1]]++
		order := mapOrder[wc[1]]
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
	// 添加符号
	for _, v := range puncts.o {
		t.Insert(v.word, v.code, v.order)
	}
	dict.Matcher = t
}
