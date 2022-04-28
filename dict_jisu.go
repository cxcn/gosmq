package smq

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
)

func (dict *Dict) fromJisu() {
	t := new(trie)
	scan := bufio.NewScanner(dict.Reader)
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
		c := wc[1]
		order := 0
		for i := 0; i < len(dict.SelectKeys); i++ {
			if c[len(c)-1] == dict.SelectKeys[i] {
				order = i + 1
				break
			}
		}
		if order == 0 {
			re := regexp.MustCompile(`\d+$`)
			match := re.FindString(c)
			if match != "" {
				order, _ = strconv.Atoi(match)
			} else {
				order = 1
			}
		} else {
			order = 1
		}
		// 生成赛码表
		wb = append(wb, scan.Bytes()...)
		wb = append(wb, '\t')
		wb = append(wb, []byte(strconv.Itoa(order))...)
		wb = append(wb, '\n')

		t.Insert(wc[0], c, order)
		dict.length++
	}
	// 添加符号
	for _, v := range puncts.o {
		t.Insert(v.word, v.code, v.order)
	}
	dict.Matcher = t
}
