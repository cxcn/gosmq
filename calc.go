package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

func (res *result) calc(dict *Trie, text []rune, csk string) {

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("calc cost time = ", cost)
	}()

	res.textLen = len(text)
	res.freqStat = make(map[string]freq)
	res.repeat = make(map[string]struct{})
	var builder, builderCode strings.Builder
	builder.Grow(res.textLen / 4)
	builderCode.Grow(res.textLen / 4)

	// 选重，替换选重键
	replace := func(w, c string) string {
		// 50: 2
		if i := c[len(c)-1]; 50 <= i && i <= 57 {
			res.repeat[w] = struct{}{}
			if len(csk) > int(i-50) {
				key := csk[int(i-50)]
				tmp := []byte(c)
				tmp[len(c)-1] = key
				c = string(tmp)
			}
		}
		return c
	}

	p := 0 // point
	for p < res.textLen {
		// 非汉字
		if !unicode.Is(unicode.Han, text[p]) {
			res.notHanCount++
			if !strings.Contains(res.notHan, string(text[p])) {
				res.notHan += string(text[p])
			}
		} else if dict.children[text[p]] == nil { // 缺字
			if !strings.Contains(res.lack, string(text[p])) {
				res.lack += string(text[p])
				res.lackCount++
			}
			p++
			continue
		} else if !dict.children[text[p]].isWord { // 缺字 有词没字
			if !strings.Contains(res.lack, string(text[p])) {
				res.lack += string(text[p])
				res.lackCount++
			}
		}
		// 最长匹配
		var a *Trie
		var i int
		for b, j := dict, 0; p+j < res.textLen; j++ {
			if b.children[text[p+j]] == nil {
				break
			}
			b = b.children[text[p+j]]
			if b.isWord {
				a, i = b, j
			}
		}

		word := string(text[p : p+i+1])
		var code string
		if a == nil {
			code = word
		} else {
			code = replace(word, a.code)
			builderCode.WriteString(code)
		}
		builder.WriteString(code)
		builder.WriteString(" ")
		tmp := res.freqStat[word]
		tmp.code = code
		tmp.times++
		res.freqStat[word] = tmp
		res.unitCount++
		p += i + 1
	}
	res.codeSep = builder.String()
	res.code = builderCode.String()
}
