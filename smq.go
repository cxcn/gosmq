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

func NewSmq(dict *Trie, fpt string, csk string) *Smq {

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("calc cost time = ", cost)
	}()

	smq := new(Smq)
	f, err := os.Open(fpt)
	if err != nil {
		fmt.Println("文本读取错误:", err)
		return smq
	}
	defer f.Close()
	buff := bufio.NewReader(f)

	smq.freqStat = make(map[string]*freq)
	smq.repeat = make(map[string]struct{})
	var builder strings.Builder

	// 逐行读取 text
	for {
		line, _, eof := buff.ReadLine()
		if eof == io.EOF {
			break
		}
		text := []rune(strings.TrimSpace(string(line)))
		smq.textLen += len(text)

		p := 0 // point
		for p < len(text) {
			// 非汉字
			if !unicode.Is(unicode.Han, text[p]) {
				smq.notHanCount++
				if !strings.Contains(smq.notHan, string(text[p])) {
					smq.notHan += string(text[p])
				}
			} else if _, ok := dict.children[text[p]]; !ok { // 缺字
				if !strings.Contains(smq.lack, string(text[p])) {
					smq.lack += string(text[p])
					smq.lackCount++
				}
				p++
				continue
			} else if len(dict.children[text[p]].code) == 0 { // 缺字 有词没字
				if !strings.Contains(smq.lack, string(text[p])) {
					smq.lack += string(text[p])
					smq.lackCount++
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
					smq.repeat[w] = struct{}{}
					if len(csk) > int(i-50) {
						tmp := []byte(c)
						tmp[len(c)-1] = csk[int(i-50)]
						c = string(tmp)
					}
				}
			}
			builder.WriteString(c)
			builder.WriteString(" ")
			if smq.freqStat[w] == nil {
				smq.freqStat[w] = new(freq)
			}
			smq.freqStat[w].code = c
			smq.freqStat[w].times++
			smq.unitCount++
			p += i + 1
		}
	}
	smq.codeSep = builder.String()
	smq.stat()
	return smq
}
