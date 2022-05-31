package smq

import (
	"strings"
	"unicode"
)

func (res *Result) match(text []rune, dict *Dict) string {
	var sb strings.Builder
	sb.Grow(len(text))
	res.Basic.TextLen += len(text)
	for p := 0; p < len(text); {
		// 删掉空白字符
		switch text[p] {
		case 65533, '\n', '\r', '\t', ' ', '　':
			p++
			res.Basic.TextLen--
			continue
		}
		// 非汉字
		isHan := unicode.Is(unicode.Han, text[p])
		if !isHan {
			res.Basic.NotHanCount++
			res.mapNotHan[text[p]] = struct{}{}
		}

		i, code, order := dict.Matcher.Match(text, p)
		// 缺字
		if i == 0 {
			if isHan {
				res.Basic.LackCount++
				res.mapLack[text[p]] = struct{}{}
			}
			sb.WriteByte(' ')
			p++
			continue
		}

		sb.WriteString(code)
		AddTo(&res.wordsDist, i) // 词长分布
		if order != 1 {
			res.Collision.Chars.Count += i // 选重字数
		} else if i != 1 {
			res.Words.FirstCount++ // 首选词
		}
		AddTo(&res.collDist, order)     // 选重分布
		AddTo(&res.codeDist, len(code)) // 码长分布

		if dict.Details {
			word := string(text[p : p+i])
			res.Data.WordSlice = append(res.Data.WordSlice, word)
			res.Data.CodeSlice = append(res.Data.CodeSlice, code)
			if _, ok := res.Data.Details[word]; !ok {
				res.Data.Details[word] = &CoC{Code: code, Order: i}
			}
			res.Data.Details[word].Count++
		}
		p += i
	}
	return sb.String()
}
