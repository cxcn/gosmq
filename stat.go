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

	codeSep string              //空格间隔的全部编码
	mapFreq map[string]freq     //词：频率
	choose  map[string]struct{} //选重

	//以下可由上面计算得
	code       string  //全部编码
	lenCode    int     //总键数
	avlCode    float64 //码长
	countSpace int     //空格数

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
	countCode := len(strings.Split(res.codeSep, " ")) // 上屏数
	res.lenCode = len(res.code)
	res.avlCode = float64(res.lenCode) / float64(res.lenText)
	res.countSpace = strings.Count(res.code, "_")

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
	res.rateWord = float64(res.countWord) / float64(countCode)
	res.rateLenWord = float64(res.lenWord) / float64(res.lenText)
	for k := range res.choose {
		l := len([]rune(k))
		res.countChoose += res.mapFreq[k].times
		res.lenChoose += l * res.mapFreq[k].times
	}
	res.rateChoose = float64(res.countChoose) / float64(countCode)
	res.rateLenChoose = float64(res.lenChoose) / float64(res.lenText)
	res.fingering()
}

func (res *result) fingering() {
	loc := make(map[rune]bool)
	for _, v := range "`12345qwertasdfgzxcvb" {
		loc[v] = true // 左手 true
	}
	// for _, v := range "67890yuiophjkl;'nm,./" {
	// 	loc[v] = false // 右手 false
	// }
	rcode := []rune(res.code)
	res.countKey = make([]int, 10)
	var countDiffHand, countSameFin, countDiffFin int

	for i := range rcode {
		switch rcode[i] {
		case 'q', 'a', 'z', '1':
			res.countKey[0]++
		case 'w', 's', 'x', '2':
			res.countKey[1]++
		case 'e', 'd', 'c', '3':
			res.countKey[2]++
		case 'v', 'f', 'r', '4':
			res.countKey[3]++
		case 't', 'g', 'b', '5':
			res.countKey[4]++
		case 'n', 'h', 'y', '6':
			res.countKey[5]++
		case 'u', 'j', 'm', '7':
			res.countKey[6]++
		case 'i', 'k', ',', '8':
			res.countKey[7]++
		case 'o', 'l', '.', '9':
			res.countKey[8]++
		case 'p', ';', '/', '0', '\'':
			res.countKey[9]++
		}
		if i+1 == len(rcode) {
			break
		}
		if rcode[i] == rcode[i+1] {
			countSameFin++
		} else if loc[rcode[i]] == loc[rcode[i+1]] {
			countDiffFin++
		} else {
			countDiffHand++
		}
	}
	res.rateDiffHand = float64(countDiffHand) / float64(res.lenCode)
	res.rateSameFin = float64(countSameFin) / float64(res.lenCode)
	res.rateDiffFin = float64(countDiffFin) / float64(res.lenCode)
	res.rateKey = make([]float64, 10)
	for i, v := range res.countKey {
		res.rateKey[i] = float64(v) / float64(res.lenCode)
	}
}
