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

	// smq.freqStat = make(map[string]*freq)
	// smq.repeat = make(map[string]struct{})
	smq.codeStat = make(map[int]int)
	smq.wordStat = make(map[int]int)
	smq.repeatStat = make(map[int]int)

	notHan := make(map[rune]struct{})
	lack := make(map[rune]struct{})
	var builder strings.Builder

	btoi := func(d byte) (int, bool) {
		// ascii 48: 0
		if 48 <= d && d <= 57 {
			return int(d - 48), true
		}
		return 0, false
	}

	// 读取 text
	for {
		line, err := buff.ReadString('\n')
		text := []rune(line)
		smq.textLen += len(text)
		p := 0 // point
		for p < len(text) {
			// 删掉空白字符
			switch text[p] {
			case 65533, '\n', '\r', '\t', ' ', '　':
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
					p++
					continue
				}
			} else if i == -1 { // 缺字
				lack[text[p]] = struct{}{}
				p++
				continue
			}

			// w := string(text[p : p+i+1])
			c := a.code
			// 选重
			rp := 0
			pow := 1
			for n := len(c) - 1; n >= 0; n-- {
				if d, ok := btoi(c[n]); ok && len(c) > 1 {
					rp += d * pow
					pow *= 10
				} else {
					break
				}
			}
			if rp != 0 {
				smq.repeatStat[rp]++
				smq.repeatCount++
				smq.repeatLen += i + 1
				if rp < 10 { // 10重以内，替换选重键
					if len(csk) > rp {
						tmp := []byte(c)
						tmp[len(c)-1] = csk[rp]
						c = string(tmp)
					}
				}
			}
			smq.codeStat[len(c)]++
			smq.wordStat[i+1]++
			if i > 0 {
				smq.wordCount++
				smq.wordLen += i + 1
			}
			// if smq.freqStat[w] == nil {
			// 	smq.freqStat[w] = new(freq)
			// }
			// smq.freqStat[w].code = c
			// smq.freqStat[w].times++
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
	smq.lackCount = len(lack)
	for k := range lack {
		smq.lack += string(k)
	}
	smq.stat()
	return smq
}
