package main

import (
	"fmt"
	"time"
)

type freq struct {
	code  string
	times int
}

type result struct {
	textLen     int    //文本字数
	notHan      string //非汉字
	notHanCount int    //非汉字数
	lack        string //缺字
	lackCount   int    //缺字数

	repeat   map[string]struct{} //选重
	freqStat map[string]freq     //字词：频率

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

	keyCount []int
	keyRate  []float64
	posCount []int     // LR RL LL RR
	posRate  []float64 // LR RL LL RR

	diffHandRate float64 // 异手
	sameFinRate  float64 // 同指
	diffFinRate  float64 // 同手异指
}

func div(x, y int) float64 {
	return float64(x) / float64(y)
}

func (res *result) stat() {

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

func (res *result) fingering(space bool) {
	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("fingering cost time = ", cost)
	}()

	pos := make(map[byte]int)
	aaa := "`1qaz2wsx3edc4rfv5tgb_6yhn7ujm8ik,9ol.0p;/'"
	bbb := "1111122223333444444445666666667777888899999"
	for i := range aaa {
		v := int(bbb[i] - 48)
		pos[aaa[i]] = v
	}
	res.keyCount = make([]int, 10)
	res.posCount = make([]int, 4)
	var countSameFin int
	a := pos[res.code[0]]
	res.keyCount[a]++

	L := func(x int) bool {
		if space {
			return x <= 5
		}
		return x < 5
	}
	R := func(x int) bool {
		if space {
			return x >= 5
		}
		return x > 5
	}
	for i := 1; i < len(res.code); i++ {
		b, ok := pos[res.code[i]]
		if !ok {
			b = 0
		}
		res.keyCount[b]++
		if a == b {
			countSameFin++
		}

		if L(a) && R(b) { // LR
			res.posCount[0]++
		} else if R(a) && L(b) { // RL
			res.posCount[1]++
		} else if L(a) && L(b) { // LL
			res.posCount[2]++
		} else if R(a) && R(b) { // RR
			res.posCount[3]++
		}
		a = b
	}

	res.posRate = make([]float64, 4)
	posSum := 0
	for _, v := range res.posCount {
		posSum += v
	}
	for i, v := range res.posCount {
		res.posRate[i] = div(v, posSum)
	}

	res.keyRate = make([]float64, 10)
	for i, v := range res.keyCount {
		res.keyRate[i] = div(v, res.codeLen)
	}
	res.diffHandRate = res.posRate[0] + res.posRate[1]
	res.sameFinRate = div(countSameFin, posSum)
	res.diffFinRate = res.posRate[2] + res.posRate[3] - res.sameFinRate
}
