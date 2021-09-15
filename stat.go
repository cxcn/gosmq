package smq

import (
	"strings"
)

type smqOut struct {
	TextLen     int    //文本字数
	NotHan      string //非汉字
	NotHanCount int    //非汉字数
	Lack        string //缺字
	LackCount   int    //缺字数

	CodeSep   string //空格间隔的全部编码
	Code      string //全部编码
	UnitCount int    //上屏数

	//以下可由上面计算得
	CodeLen int     //总键数
	CodeAvg float64 //码长

	WordCount   int     //打词数
	WordLen     int     //打词字数
	WordRate    float64 //打词率（上屏）
	WordLenRate float64 //打词率（字数）

	RepeatCount   int     //选重数
	RepeatLen     int     //选重字数
	RepeatRate    float64 //选重率（上屏）
	RepeatLenRate float64 //选重率（字数）

	CodeStat   map[int]int //码长统计
	WordStat   map[int]int //词长统计
	RepeatStat map[int]int //选重统计
}

func (so *smqOut) stat() {
	so.Code = strings.ReplaceAll(so.CodeSep, " ", "")
	so.CodeLen = len(so.Code)
	so.CodeAvg = div(so.CodeLen, so.TextLen)
	so.WordRate = div(so.WordCount, so.UnitCount)
	so.WordLenRate = div(so.WordLen, so.TextLen)
	so.RepeatRate = div(so.RepeatCount, so.UnitCount)
	so.RepeatLenRate = div(so.RepeatLen, so.TextLen)
}

func div(x, y int) float64 {
	return float64(x) / float64(y)
}
