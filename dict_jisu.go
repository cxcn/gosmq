package smq

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func (dict *Dict) fromJisu() {
	scan := bufio.NewScanner(dict.reader)
	var buf bytes.Buffer
	// 生成字典
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) < 2 {
			continue
		}
		if dict.Single && len([]rune(wc[0])) != 1 {
			continue
		}
		c := wc[1]
		code := ""
		order := 0
		// a_ aa_
		if len(c)-1 > 0 && c[len(c)-1] == '_' {
			code = c[:len(c)-1]
			order = 1
		} else {
			re := regexp.MustCompile(`\d+$`)
			match := re.FindString(c)
			// a1 aa3
			if match != "" {
				code = c[:len(c)-len(match)]
				order, _ = strconv.Atoi(match)
			} else { // akdb ksdw
				code = c
				order = 0 // 和前面区分开
			}
		}
		// 生成赛码表
		buf.WriteString(wc[0])
		buf.WriteByte('\t')
		buf.WriteString(code)
		// 自定义选重键
		if order != 0 {
			if order <= len(dict.SelectKeys) {
				buf.WriteByte(dict.SelectKeys[order-1])
			} else {
				buf.WriteString(strconv.Itoa(order))
			}
		}
		buf.WriteByte('\t')
		if order == 0 {
			order = 1
		}
		buf.WriteString(strconv.Itoa(order))
		buf.WriteByte('\n')
	}
	// 输出赛码表
	err := ioutil.WriteFile(dict.SavePath, buf.Bytes(), 0666)
	if err != nil {
		// SavePath 不对则保存在 dict 目录下
		os.Mkdir("dict", 0666)
		err = ioutil.WriteFile("./dict/"+dict.Name+".txt", buf.Bytes(), 0666)
		if err != nil {
			log.Println(err)
		}
	}
	dict.reader = bytes.NewReader(buf.Bytes())
}
