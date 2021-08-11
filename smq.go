package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func NewSmq(dict *Trie, fpt string, csk string) *Smq {

	// start := time.Now()
	// defer func() {
	// 	cost := time.Since(start)
	// 	fmt.Println("NewSmq cost time = ", cost)
	// }()

	smq := new(Smq)
	f, err := os.Open(fpt)
	if err != nil {
		fmt.Println("文本读取错误:", err)
		return smq
	}
	_, filename := filepath.Split(fpt)
	fmt.Println("文本读取成功:", filename)
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)

	smq.freqStat = make(map[string]*freq)
	smq.repeat = make(map[string]struct{})
	var builder strings.Builder

	// 逐行读取 text
	for scan.Scan() {

		text := []rune(scan.Text())
		smq.textLen += len(text)
		p := 0 // point
		for p < len(text) {
			// 这是一个奇怪的字符
			if text[p] == 65533 {
				p++
				continue
			}
			// 最长匹配
			a := new(Trie)
			i := -1
			for b, j := dict, 0; p+j < len(text); j++ {
				if v, ok := b.children[text[p+j]]; !ok {
					break
				} else {
					b = v
				}
				if b.code != "" {
					a, i = b, j
				}
			}

			if !unicode.Is(unicode.Han, text[p]) { // 非汉字，￥
				smq.notHanCount++
				if !strings.Contains(smq.notHan, string(text[p])) {
					smq.notHan += string(text[p])
				}
				if i == -1 { // 缺非汉字￥
					builder.WriteString(string(text[p]))
					builder.WriteString(" ")
					p++
					continue
				}
			} else if i == -1 { // 缺字
				if !strings.Contains(smq.lack, string(text[p])) {
					smq.lack += string(text[p])
					smq.lackCount++
				}
				p++
				continue
			}

			w := string(text[p : p+i+1])
			c := a.code
			// 选重，替换选重键 ascii 50: 2
			if d := c[len(c)-1]; 50 <= d && d <= 57 {
				smq.repeat[w] = struct{}{}
				if len(csk) > int(d-50) {
					tmp := []byte(c)
					tmp[len(c)-1] = csk[int(d-50)]
					c = string(tmp)
				}
			}
			if smq.freqStat[w] == nil {
				smq.freqStat[w] = new(freq)
			}
			smq.freqStat[w].code = c
			smq.freqStat[w].times++
			smq.unitCount++
			builder.WriteString(c)
			builder.WriteString(" ")
			p += i + 1
		}
	}
	smq.codeSep = builder.String()
	smq.stat()
	return smq
}
