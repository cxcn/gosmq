package main

import (
	"strings"
)

type freq struct {
	code  string
	times int
}

type Smq struct {
	textLen     int    //文本字数
	notHan      string //非汉字
	notHanCount int    //非汉字数
	lack        string //缺字
	lackCount   int    //缺字数

	repeat   map[string]struct{} //选重
	freqStat map[string]*freq    //字词：频率

	codeSep   string //空格间隔的全部编码
	code      string //全部编码
	unitCount int    //上屏数

	//以下可由上面计算得
	codeLen int     //总键数
	codeAvg float64 //码长

	wordCount   int     //打词数
	wordLen     int     //打词字数
	wordRate    float64 //打词率（上屏）
	wordLenRate float64 //打词率（字数）

	repeatCount   int     //选重数
	repeatLen     int     //选重字数
	repeatRate    float64 //选重率（上屏）
	repeatLenRate float64 //选重率（字数）

	codeStat map[int]int //码长统计
	wordStat map[int]int //词长统计
}

func (res *Smq) stat() {

	res.code = strings.ReplaceAll(res.codeSep, " ", "")
	res.codeLen = len(res.code)
	res.codeAvg = div(res.codeLen, res.textLen)

	res.codeStat = make(map[int]int)
	res.wordStat = make(map[int]int)
	for k, v := range res.freqStat {
		l := len([]rune(k))
		if l > 1 {
			res.wordCount += v.times
			res.wordLen += l * v.times
		}
		res.codeStat[len(v.code)] += v.times
		res.wordStat[l] += v.times
	}
	res.wordRate = div(res.wordCount, res.unitCount)
	res.wordLenRate = div(res.wordLen, res.textLen)
	for k := range res.repeat {
		l := len([]rune(k))
		res.repeatCount += res.freqStat[k].times
		res.repeatLen += l * res.freqStat[k].times
	}
	res.repeatRate = div(res.repeatCount, res.unitCount)
	res.repeatLenRate = div(res.repeatLen, res.textLen)
}

func div(x, y int) float64 {
	return float64(x) / float64(y)
}
