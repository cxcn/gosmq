package smq

import (
	"strings"
	"unicode"
)

func (res *Result) match(text []rune, m Matcher) string {
	var sb strings.Builder
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

		i, code, order := m.Match(text, p)
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
		res.Words.Dist[i]++ // 词长分布
		if order != 1 {
			if i != 1 {
				res.Words.FirstCount++
			}
			res.Collision.Chars.Count += i // 选重字数
		}
		res.Collision.Dist[order]++   // 选重分布
		res.CodeLen.Dist[len(code)]++ // 码长分布

		// res.Data.CodeSlice = append(res.Data.CodeSlice, code)
		// res.Data.WordSlice = append(res.Data.WordSlice, text[p:p+i])
		p += i
	}
	return sb.String()
}
