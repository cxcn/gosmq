package smq

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func (dict *Dict) fromDuoduo() {
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
	// 输出赛码表
	err := ioutil.WriteFile(dict.SavePath, buf.Bytes(), 0666)
	if err != nil {
		log.Println(err)
	}
	dict.reader = bytes.NewReader(buf.Bytes())
}
