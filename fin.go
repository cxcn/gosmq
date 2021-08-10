package main

import (
	"fmt"
	"time"
)

type Fin struct {
	finCount  []int
	finRate   []float64
	leftHand  float64 // 左手
	rightHand float64 // 右手

	handCount    []int     // LR RL LL RR
	handRate     []float64 // LR RL LL RR
	diffHandRate float64   // 异手
	sameFinRate  float64   // 同指
	diffFinRate  float64   // 同手异指

	dl   float64 // 当量
	dkp  float64 // 大跨排
	xkp  float64 // 小跨排
	xzgr float64 // 小指干扰
	cs   float64 // 错手
}

func NewFin(code string, isS bool) *Fin {

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("NewFin cost time = ", cost)
	}()

	zhifa := newZhifa(isS)
	fin := new(Fin)
	fin.finCount = make([]int, 10)
	fin.handCount = make([]int, 4)

	finger := make(map[byte]int)
	aaa := "1qaz2wsx3edc4rfv5tgb_6yhn7ujm8ik,9ol.0p;/'"
	bbb := "111122223333444444445666666667777888899999"
	for i := range aaa {
		v := int(bbb[i] - 48)
		finger[aaa[i]] = v
	}
	var (
		sameFinCount int
		dlSum        float64
		dkpCount     int
		xkpCount     int
		xzgrCount    int
		csCount      int
		combLen      int
	)

	a := 0
	for i := 1; i < len(code); i++ {

		// 处理单键
		b, ok := finger[code[i]]
		if !ok {
			b = 0
			fin.finCount[b]++
			a = b
			continue
		}
		fin.finCount[b]++
		if a == 0 {
			a = b
			continue
		}

		// 同指
		if a == b {
			sameFinCount++
		}
		a = b

		// 处理按键组合
		zf := zhifa[code[i-1]][code[i]]
		dlSum += zf.dl
		combLen++
		// 大小跨排等
		switch zf.zf {
		case 0:
		case 2:
			xkpCount++
		case 3:
			xzgrCount++
		case 1:
			dkpCount++
		case 4:
			csCount++
		}
		// 互击
		switch zf.hj {
		case 1:
			fin.handCount[0]++
		case 2:
			fin.handCount[1]++
		case 3:
			fin.handCount[2]++
		case 4:
			fin.handCount[3]++
		}
	}

	fin.finRate = make([]float64, 10)
	for i, v := range fin.finCount {
		fin.finRate[i] = div(v, len(code))
	}
	fin.leftHand = fin.finRate[1] + fin.finRate[2] + fin.finRate[3] + fin.finRate[4]
	fin.rightHand = fin.finRate[6] + fin.finRate[7] + fin.finRate[8] + fin.finRate[9]
	fin.leftHand = fin.leftHand / (fin.leftHand + fin.rightHand) // 归一
	fin.rightHand = 1 - fin.leftHand

	fin.handRate = make([]float64, 4)
	handSum := 0
	for _, v := range fin.handCount {
		handSum += v
	}
	for i, v := range fin.handCount {
		fin.handRate[i] = div(v, handSum)
	}
	fin.diffHandRate = fin.handRate[0] + fin.handRate[1]
	fin.sameFinRate = div(sameFinCount, handSum)
	fin.diffFinRate = fin.handRate[2] + fin.handRate[3] - fin.sameFinRate

	fin.dl = dlSum / float64(combLen)
	fin.dkp = div(dkpCount, combLen)
	fin.xkp = div(xkpCount, combLen)
	fin.xzgr = div(xzgrCount, combLen)
	fin.cs = div(csCount, combLen)

	return fin
}
