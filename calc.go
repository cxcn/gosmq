package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type freq struct {
	code  string
	times int
}

type result struct {
	lenText     int    //文本字数
	notHan      string //非汉字
	countNotHan int    //非汉字数
	lack        string //缺字
	countLack   int    //缺字数

	codeSep string //空格间隔的全部编码
	mapFreq map[string]freq
	choose  map[string]struct{} //选重

	//以下可由上面计算得
	code       string   //全部编码
	codeSlice  []string //编码切片
	lenCode    int      //总键数
	avlCode    float64  //码长
	countSpace int      //空格数

	countWord   int     //打词数
	lenWord     int     //打词字数
	rateWord    float64 //打词率（上屏）
	rateLenWord float64 //打词率（字数）

	countChoose   int     //选重数
	lenChoose     int     //选重字数
	rateChoose    float64 //选重率（上屏）
	rateLenChoose float64 //选重率（字数）

	statCode map[int]int
	statWord map[int]int
}

func calc(fpm, fpt string) result {
	text := readText(fpt)
	dict := read(fpm)
	conf := readConf(fpm)

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("cacl cost time = ", cost)
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
	res.write()
	return *res
}
