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
	buff := bufio.NewReader(f)

	smq.freqStat = make(map[string]*freq)
	smq.repeat = make(map[string]struct{})
	notHan := make(map[rune]struct{})
	lack := make(map[rune]struct{})
	var builder strings.Builder

	// 读取 text
	for {
		line, err := buff.ReadString('\n')
		text := []rune(line)
		smq.textLen += len(text)
		p := 0 // point
		for p < len(text) {
			// 删掉空白字符
			switch text[p] {
			case 65533, '\n', '\r', '\t', ' ':
				p++
				smq.textLen--
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
				notHan[text[p]] = struct{}{}
				if i == -1 { // 缺非汉字￥
					builder.WriteString(string(text[p]))
					builder.WriteString(" ")
					p++
					continue
				}
			} else if i == -1 { // 缺字
				lack[text[p]] = struct{}{}
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
		if err != nil {
			break
		}
	}

	smq.codeSep = builder.String()
	for k := range notHan {
		smq.notHan += string(k)
	}
	for k := range lack {
		smq.lack += string(k)
		smq.lackCount++
	}
	smq.stat()
	return smq
}
