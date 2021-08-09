package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode"
)

func (res *result) calc(dict *Trie, fpt string, csk string) {

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("calc cost time = ", cost)
	}()

	f, err := os.Open(fpt)
	if err != nil {
		fmt.Println("文本读取错误:", err)
		return
	}
	defer f.Close()
	buff := bufio.NewReader(f)

	res.freqStat = make(map[string]*freq)
	res.repeat = make(map[string]struct{})
	var builder strings.Builder

	// 逐行读取 text
	for {
		line, _, eof := buff.ReadLine()
		if eof == io.EOF {
			break
		}
		text := []rune(strings.TrimSpace(string(line)))
		res.textLen += len(text)

		p := 0 // point
		for p < len(text) {
			// 非汉字
			if !unicode.Is(unicode.Han, text[p]) {
				res.notHanCount++
				if !strings.Contains(res.notHan, string(text[p])) {
					res.notHan += string(text[p])
				}
			} else if _, ok := dict.children[text[p]]; !ok { // 缺字
				if !strings.Contains(res.lack, string(text[p])) {
					res.lack += string(text[p])
					res.lackCount++
				}
				p++
				continue
			} else if len(dict.children[text[p]].code) == 0 { // 缺字 有词没字
				if !strings.Contains(res.lack, string(text[p])) {
					res.lack += string(text[p])
					res.lackCount++
				}
			}
			// 最长匹配
			var a *Trie
			var i int
			for b, j := dict, 0; p+j < len(text); j++ {
				if _, ok := b.children[text[p+j]]; !ok {
					break
				}
				b = b.children[text[p+j]]
				if len(b.code) != 0 {
					a, i = b, j
				}
			}

			w := string(text[p : p+i+1])
			var c string
			if a == nil {
				c = w
			} else {
				c = a.code
				// 选重，替换选重键 ascii 50: 2
				if i := c[len(c)-1]; 50 <= i && i <= 57 {
					res.repeat[w] = struct{}{}
					if len(csk) > int(i-50) {
						tmp := []byte(c)
						tmp[len(c)-1] = csk[int(i-50)]
						c = string(tmp)
					}
				}
			}
			builder.WriteString(c)
			builder.WriteString(" ")
			if res.freqStat[w] == nil {
				res.freqStat[w] = new(freq)
			}
			res.freqStat[w].code = c
			res.freqStat[w].times++
			res.unitCount++
			p += i + 1
		}
	}
	res.codeSep = builder.String()
}
