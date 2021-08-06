package main

import (
	"strings"
)

func (res *result) write() {
	res.code = strings.Replace(res.codeSep, " ", "", -1)
	res.codeSlice = strings.Split(res.codeSep, " ")
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
	res.rateWord = float64(res.countWord) / float64(len(res.codeSlice))
	res.rateLenWord = float64(res.lenWord) / float64(res.lenText)
	for k := range res.choose {
		l := len([]rune(k))
		res.countChoose += res.mapFreq[k].times
		res.lenChoose += l * res.mapFreq[k].times
	}
	res.rateChoose = float64(res.countChoose) / float64(len(res.codeSlice))
	res.rateLenChoose = float64(res.lenChoose) / float64(res.lenText)
}
