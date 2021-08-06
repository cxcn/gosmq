package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

func calc(fpm, fpt string) result {
	text := readText(fpt)
	dict := read(fpm)
	conf := readConf(fpm)

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("calc cost time = ", cost)
	}()

	res := new(result)
	res.lenText = len(text)
	res.mapFreq = make(map[string]freq)
	res.choose = make(map[string]struct{})
	var builder strings.Builder
	builder.Grow(res.lenText / 4)

	// 选重，替换选重键
	replace := func(w, c string) string {
		if i := c[len(c)-1]; 48 <= i && i <= 57 {
			res.choose[w] = struct{}{}
			if conf.isConf {
				key, ok := conf.ak[int(i-48)]
				if ok {
					c = c[:len(c)-1] + key
				}
			}
		}
		return c
	}

	p := 0 // point
	for p < res.lenText {
		// 非汉字
		if !unicode.Is(unicode.Han, text[p]) {
			res.countNotHan++
			if !strings.Contains(res.notHan, string(text[p])) {
				res.notHan += string(text[p])
			}
		} else if dict.children[text[p]] == nil { // 缺字
			if !strings.Contains(res.lack, string(text[p])) {
				res.lack += string(text[p])
				res.countLack++
			}
			p++
			continue
		} else if !dict.children[text[p]].isWord { // 缺字 有词没字
			if !strings.Contains(res.lack, string(text[p])) {
				res.lack += string(text[p])
				res.countLack++
			}
		}
		// 最长匹配
		var a Trie
		var i int
		for b, j := dict, 0; p+j < res.lenText; j++ {
			if b.children[text[p+j]] == nil {
				break
			}
			b = *b.children[text[p+j]]
			if b.isWord {
				a, i = b, j
			}
		}

		word := string(text[p : p+i+1])
		var code string
		if len(a.code) == 0 {
			code = word
		} else {
			code = replace(word, a.code)
		}
		builder.WriteString(code)
		builder.WriteString(" ")
		tmp := res.mapFreq[word]
		tmp.code = code
		tmp.times++
		res.mapFreq[word] = tmp

		p += i + 1
	}

	res.codeSep = builder.String()
	res.stat()
	return *res
}
