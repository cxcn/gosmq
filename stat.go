package main

import (
	"fmt"
	"strings"
	"time"
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

	codeSep   string              //空格间隔的全部编码
	choose    map[string]struct{} //选重
	countCode int                 //上屏数
	mapFreq   map[string]freq     //词：频率

	//以下可由上面计算得
	code    string  //全部编码
	lenCode int     //总键数
	avlCode float64 //码长

	countWord   int     //打词数
	lenWord     int     //打词字数
	rateWord    float64 //打词率（上屏）
	rateLenWord float64 //打词率（字数）

	countChoose   int     //选重数
	lenChoose     int     //选重字数
	rateChoose    float64 //选重率（上屏）
	rateLenChoose float64 //选重率（字数）

	statCode map[int]int //码长统计
	statWord map[int]int //词长统计

	countKey []int
	rateKey  []float64
	countPos []int     // LR RL LL RR
	ratePos  []float64 // LR RL LL RR

	rateDiffHand float64 // 异手
	rateSameFin  float64 // 同指
	rateDiffFin  float64 // 同手异指
}

func (res *result) stat() {
	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("stat cost time = ", cost)
	}()

	res.code = strings.ReplaceAll(res.codeSep, " ", "")
	res.lenCode = len(res.code)

	rateCode := func(x int) float64 {
		return float64(x) / float64(res.countCode)
	}
	rateText := func(x int) float64 {
		return float64(x) / float64(res.lenText)
	}

	res.avlCode = rateText(res.lenCode)

	res.statCode = make(map[int]int)
	res.statWord = make(map[int]int)
	for k, v := range res.mapFreq {
		l := len([]rune(k))
		if l > 1 {
			res.countWord += v.times
			res.lenWord += l * v.times
		}
		res.statCode[len(v.code)] += v.times
		res.statWord[l] += v.times
	}
	res.rateWord = rateCode(res.countWord)
	res.rateLenWord = rateText(res.lenWord)
	for k := range res.choose {
		l := len([]rune(k))
		res.countChoose += res.mapFreq[k].times
		res.lenChoose += l * res.mapFreq[k].times
	}
	res.rateChoose = rateCode(res.countChoose)
	res.rateLenChoose = rateText(res.lenChoose)
	res.fingering()
}

func (res *result) fingering() {

	pos := make(map[byte]int)
	aaa := "`1qaz2wsx3edc4rfv5tgb_6yhn7ujm8ik,9ol.0p;/'"
	bbb := "0000011112222333333334666666667777888899999"
	for i := range aaa {
		v := int(bbb[i] - 48)
		pos[aaa[i]] = v
	}
	fmt.Println(pos)
	res.countKey = make([]int, 10)
	res.countPos = make([]int, 4)
	var countSameFin int
	a := pos[res.code[0]]
	res.countKey[a]++
	for i := 1; i < len(res.code); i++ {
		b, ok := pos[res.code[i]]
		if !ok {
			b = 5
		}
		res.countKey[b]++
		if a == b {
			countSameFin++
		}
		if a < 4 && b > 5 { // LR
			res.countPos[0]++
		} else if a > 5 && b < 4 { // RL
			res.countPos[1]++
		} else if a < 4 && b < 4 { // LL
			res.countPos[2]++
		} else if a > 5 && b > 5 { // RR
			res.countPos[3]++
		}
		a = b
	}

	rate := func(x int) float64 {
		return float64(x) / float64(res.lenCode)
	}

	res.rateKey = make([]float64, 10)
	for i, v := range res.countKey {
		res.rateKey[i] = rate(v)
	}
	res.ratePos = make([]float64, 4)
	for i, v := range res.countPos {
		res.ratePos[i] = rate(v)
	}
	res.rateDiffHand = res.ratePos[0] + res.ratePos[1]
	res.rateSameFin = rate(countSameFin)
	res.rateDiffFin = res.ratePos[2] + res.ratePos[3] - res.rateSameFin
}
