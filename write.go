package main

import (
	"strings"
	"unicode/utf8"
)

func write() {
	// start := time.Now()
	// defer func() {
	// 	cost := time.Since(start)
	// 	fmt.Println("write cost time = ", cost)
	// }()
	Res.code = strings.Replace(Res.codeSep, " ", "", -1)
	Res.codeSlice = strings.Split(Res.codeSep, " ")
	Res.lenCode = len(Res.code)
	Res.avlCode = float64(Res.lenCode) / float64(Res.lenText)
	Res.countSpace = strings.Count(Res.code, "_")

	Res.stat = make(map[int]int)
	for k, v := range Res.freq {
		l := utf8.RuneCountInString(k)
		// fmt.Println(l)
		if l > 1 {
			Res.countWord += v.times
			Res.lenWord += l * v.times
			// fmt.Println(v)
		}
		Res.stat[len(v.code)] += v.times
		// if len(v.code) < 4 {
		// 	fmt.Println(k, v)
		// }
	}
	// fmt.Println(len(Res.code))
	Res.rateWord = float64(Res.countWord) / float64(len(Res.codeSlice))
	Res.rLenWord = float64(Res.lenWord) / float64(Res.lenText)
	for k := range Res.choose {
		l := utf8.RuneCountInString(k)
		Res.countChoose += Res.freq[k].times
		Res.lenChoose += l * Res.freq[k].times
	}
	Res.rateChoose = float64(Res.countChoose) / float64(len(Res.codeSlice))
	Res.rLenChoose = float64(Res.lenChoose) / float64(Res.lenText)

}
