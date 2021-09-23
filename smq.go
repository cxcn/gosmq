package smq

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func newSmqOut(si *SmqIn) *SmqOut {

	so := new(SmqOut)
	// 读取文本
	f, err := os.Open(si.Fpt)
	if err != nil {
		fmt.Println("文本读取错误:", err)
		return so
	}
	_, filename := filepath.Split(si.Fpt)
	fmt.Println("文本读取成功:", filename)
	defer f.Close()

	// 读取码表
	dict := newDict(si)
	if dict.children == nil {
		return so
	}

	so.RepeatStat = make(map[int]int)
	so.CodeStat = make(map[int]int)
	so.WordStat = make(map[int]int)

	buf := bufio.NewReader(f)
	var codeSep bytes.Buffer
	var notHan []rune
	var lack []rune
	for {
		line, err := buf.ReadString('\n')
		text := []rune(line)
		so.TextLen += len(text)
		var code strings.Builder

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
			j := 0  // 已匹配的字数
			i := 0  // 有编码的匹配
			c := "" // 编码
			for p+j < len(text) {
				t = t.children[text[p+j]]
				j++
				if t != nil {
					if t.code != "" {
						i = j
						c = t.code
					}
				} else {
					break
				}
			}
			if i == 0 { // 缺字
				if isHan {
					lack = append(lack, text[p])
				}
				p++
				continue
			}

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
				so.RepeatLen += i
				// 替换选重键
				if rp-2 <= len(si.Csk)-1 {
					tmp := []byte(c)
					tmp[len(c)-1] = si.Csk[rp-2]
					c = string(tmp)
				}
			}
			so.CodeStat[len(c)]++
			so.WordStat[i]++
			if i > 1 {
				so.WordCount++
				so.WordLen += i
			}
			so.UnitCount++
			so.CodeLen += len(c)
			code.WriteString(c)
			if si.Fpo != "" {
				codeSep.WriteString(c)
				codeSep.WriteByte(' ')
			}
			p += i
		}
		if si.Fpo != "" {
			codeSep.WriteByte('\n')
		}
		so.feel(code.String(), si.combs)
		if err != nil {
			break
		}
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
	so.stat(si)

	// 输出编码
	if si.Fpo != "" {
		_ = ioutil.WriteFile(si.Fpo, codeSep.Bytes(), 0666)
	}
	return so
}

// ascii 转数字
func btoi(d byte) int {
	// ascii 48: 0
	if 48 <= d && d <= 57 {
		return int(d - 48)
	}
	return -1
}
