package smq

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func newSmqOut(dict *trie, fpt, csk, fpo string) *SmqOut {
	so := new(SmqOut)
	f, err := os.Open(fpt)
	if err != nil {
		fmt.Println("文本读取错误:", err)
		return so
	}
	_, filename := filepath.Split(fpt)
	fmt.Println("文本读取成功:", filename)
	defer f.Close()
	buff := bufio.NewReader(f)

	so.CodeStat = make(map[int]int)
	so.WordStat = make(map[int]int)
	so.RepeatStat = make(map[int]int)

	notHan := make(map[rune]struct{})
	lack := make(map[rune]struct{})
	var builder strings.Builder

	// ascii 转数字
	btoi := func(d byte) int {
		// ascii 48: 0
		if 48 <= d && d <= 57 {
			return int(d - 48)
		}
		return -1
	}
	// 读取 text
	for {
		line, err := buff.ReadString('\n')
		text := []rune(line)
		so.TextLen += len(text)
		p := 0 // point
		for p < len(text) {
			// 删掉空白字符
			switch text[p] {
			case 65533, '\n', '\r', '\t', ' ', '　':
				p++
				so.TextLen--
				continue
			}
			// 最长匹配
			a := new(trie)
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
				so.NotHanCount++
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

			c := a.code
			// 选重
			rp := 0
			pow := 1
			for n := len(c) - 1; n >= 0; n-- {
				if d := btoi(c[n]); d >= 0 && len(c) > 1 {
					rp += d * pow
					pow *= 10
				} else {
					break
				}
			}

			if rp > 1 {
				so.RepeatStat[rp]++
				so.RepeatCount++
				so.RepeatLen += i + 1
				// 替换选重键
				if rp-2 <= len(csk)-1 {
					tmp := []byte(c)
					tmp[len(c)-1] = csk[rp-2]
					c = string(tmp)
				}
			}
			so.CodeStat[len(c)]++
			so.WordStat[i+1]++
			if i > 0 {
				so.WordCount++
				so.WordLen += i + 1
			}
			so.UnitCount++
			builder.WriteString(c)
			builder.WriteString(" ")
			p += i + 1
		}
		if err != nil {
			break
		}
	}

	so.CodeSep = builder.String()
	// 输出编码
	if fpo != "" {
		_ = ioutil.WriteFile(fpo, []byte(so.CodeSep), 0666)
	}
	for k := range notHan {
		so.NotHan += string(k)
	}
	so.LackCount = len(lack)
	for k := range lack {
		so.Lack += string(k)
	}
	so.stat()
	return so
}
