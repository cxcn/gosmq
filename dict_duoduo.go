package smq

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func (dict *Dict) fromDuoduo() {
	t := new(trie)
	scan := bufio.NewScanner(dict.reader)
	mapOrder := make(map[string]int)
	var buf bytes.Buffer
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
		buf.WriteString(scan.Text())
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

		t.Insert(wc[0], wc[1], order)
		dict.length++
	}
	// 添加符号
	for _, v := range puncts.o {
		t.Insert(v.word, v.code, v.order)
	}
	// 输出赛码表
	_ = os.Mkdir("dict", 0666)
	err := ioutil.WriteFile("dict/"+dict.Name+".txt", buf.Bytes(), 0666)
	if err != nil {
		log.Println(err)
	}
	dict.Matcher = t
}
