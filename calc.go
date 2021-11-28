package smq

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

func (so *SmqOut) calc(rd io.Reader, dict *trie, csk string, as, isO bool) {

	buf := bufio.NewReader(rd)
	var notHan []rune
	var lack []rune
	combMap := newCombMap(as)

	for {
		// 逐行读取文本文件
		line, err := buf.ReadString('\n')
		text := []rune(line)
		so.TextLen += len(text)
		var codeSlice []string

		for p := 0; p < len(text); {
			// 删掉空白字符
			switch text[p] {
			case 65533, '\n', '\r', '\t', ' ', '　':
				p++
				so.TextLen--
				continue
			}
			// 非汉字
			isHan := unicode.Is(unicode.Han, text[p])
			if !isHan {
				so.NotHanCount++
				notHan = append(notHan, text[p])
			}
			// 最长匹配
			t := dict
			j := 0 // 已匹配的字数
			i := 0 // 有编码的匹配

			c := string(text[p]) // 编码
			for p+j < len(text) {
				t = t.children[text[p+j]]
				j++
				if t == nil {
					break
				}
				if t.code != "" {
					i = j
					c = t.code
				}
			}
			if i == 0 { // 缺字
				if isHan {
					lack = append(lack, text[p])
				}
				if isO {
					so.WordSlice = append(so.WordSlice, text[p:p+1])
					so.CodeSlice = append(so.CodeSlice, c)
				} else {
					codeSlice = append(codeSlice, c)
				}
				p++
				continue
			}

			so.UnitCount++   // 上屏数
			so.WordStat[i]++ // 词长
			if i > 1 {
				so.WordCount++
				so.WordLen += i
			}
			c, rp := repeat(c, csk) // 选重
			so.RepeatStat[rp]++
			if rp > 1 {
				so.RepeatCount++
				so.RepeatLen += i
			}
			so.CodeStat[len(c)]++ // 码长
			so.CodeLen += len(c)
			if isO {
				so.WordSlice = append(so.WordSlice, text[p:p+i])
				so.CodeSlice = append(so.CodeSlice, c)
			} else {
				codeSlice = append(codeSlice, c)
			}

			p += i
		}

		if !isO {
			so.feel(codeSlice, combMap)
		}
		if err != nil {
			break
		}
	}

	if isO {
		so.feel(so.CodeSlice, combMap)
	}
	for _, v := range notHan {
		if !strings.ContainsRune(so.NotHan, v) {
			so.NotHan += string(v)
		}
	}
	for _, v := range lack {
		if !strings.ContainsRune(so.Lack, v) {
			so.Lack += string(v)
		}
	}
	so.stat()
}

func repeat(c, csk string) (string, int) {

	// ascii 转数字
	btoi := func(d byte) int {
		// ascii 48: 0
		if 48 <= d && d <= 57 {
			return int(d - 48)
		}
		return -1
	}
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
		// 替换选重键
		if rp-2 <= len(csk)-1 {
			tmp := []byte(c)
			tmp[len(c)-1] = csk[rp-2]
			c = string(tmp)
		}
	} else {
		rp = 1
	}
	return c, rp
}
