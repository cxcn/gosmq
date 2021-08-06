package main

import (
	"bytes"
	"fmt"
	"strconv"
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
	var buf bytes.Buffer

	part := func(w, c string) {
		// 选重
		num := c[len(c)-1]
		if 48 <= num && num <= 57 {
			res.choose[w] = struct{}{}
			if conf.isConf {
				s, err := strconv.Atoi(string(num))
				if err != nil {
					errHandler(err)
				}
				key, ok := conf.ak[s]
				if ok {
					c = c[:len(c)-1] + key
				}
			}
		}
		buf.WriteString(c + " ")

		tmp := res.mapFreq[w]
		if tmp.code == "" {
			tmp.code = c
		}
		tmp.times++
		res.mapFreq[w] = tmp
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
		var a, b Trie
		b = dict
		i, j := 0, 0
		for p+j < res.lenText {
			if b.children[text[p+j]] == nil {
				break
			}
			b = *b.children[text[p+j]]
			if b.isWord {
				a = b
				i = j
			}
			j++
		}

		word := string(text[p : p+i+1])
		code := a.code
		if len(code) != 0 {
			part(word, code)
		} else {
			part(word, word)
		}
		p += i + 1
	}

	res.codeSep = buf.String()
	res.stat()
	return *res
}
