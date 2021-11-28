package smq

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"unicode"
)

func newSmqOut(si *SmqIn) *SmqOut {

	so := new(SmqOut)
	// 读取文本
	f, rd, err := ReadFile(si.Fpt)
	if err != nil {
		fmt.Println("文本读取错误:", err)
		return so
	}
	_, filename := filepath.Split(si.Fpt)
	fmt.Println("文本读取成功:", filename)
	defer f.Close()

	// 读取码表
	dict, count := newDict(si)
	so.MbLen = count
	if count == 0 {
		return so
	}

	so.RepeatStat = make(map[int]int)
	so.CodeStat = make(map[int]int)
	so.WordStat = make(map[int]int)

	buf := bufio.NewReader(rd)
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
				p++
				continue
			}

			so.UnitCount++   // 上屏数
			so.WordStat[i]++ // 词长
			if i > 1 {
				so.WordCount++
				so.WordLen += i
			}
			c, rp := repeat(c, si.Csk) // 选重
			if rp > 1 {
				so.RepeatStat[rp]++
				so.RepeatCount++
				so.RepeatLen += i
			}
			so.CodeStat[len(c)]++ // 码长
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
		err = ioutil.WriteFile(si.Fpo, codeSep.Bytes(), 0666)
		if err != nil {
			fmt.Printf("输出编码错误：%v", err)
		}
	}
	return so
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
	}
	return c, rp
}
