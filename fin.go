package main

import (
	"fmt"
	"time"
)

type Fin struct {
	keyCount []int
	keyRate  []float64
	posCount []int     // LR RL LL RR
	posRate  []float64 // LR RL LL RR

	leftHand     float64
	rightHand    float64
	diffHandRate float64 // 异手
	sameFinRate  float64 // 同指
	diffFinRate  float64 // 同手异指
}

func NewFin(code string, isS bool) *Fin {
	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("NewFin cost time = ", cost)
	}()

	fin := new(Fin)
	pos := make(map[byte]int)
	aaa := "`1qaz2wsx3edc4rfv5tgb_6yhn7ujm8ik,9ol.0p;/'"
	bbb := "1111122223333444444445666666667777888899999"
	for i := range aaa {
		v := int(bbb[i] - 48)
		pos[aaa[i]] = v
	}
	fin.keyCount = make([]int, 10)
	fin.posCount = make([]int, 4)
	var countSameFin int
	a := pos[code[0]]
	fin.keyCount[a]++

	L := func(x int) bool {
		if isS {
			return x <= 5
		}
		return x < 5
	}
	R := func(x int) bool {
		if isS {
			return x >= 5
		}
		return x > 5
	}
	for i := 1; i < len(code); i++ {
		b, ok := pos[code[i]]
		if !ok {
			b = 0
		}
		fin.keyCount[b]++
		if a == b {
			countSameFin++
		}

		if L(a) && R(b) { // LR
			fin.posCount[0]++
		} else if R(a) && L(b) { // RL
			fin.posCount[1]++
		} else if L(a) && L(b) { // LL
			fin.posCount[2]++
		} else if R(a) && R(b) { // RR
			fin.posCount[3]++
		}
		a = b
	}

	fin.posRate = make([]float64, 4)
	posSum := 0
	for _, v := range fin.posCount {
		posSum += v
	}
	for i, v := range fin.posCount {
		fin.posRate[i] = div(v, posSum)
	}

	fin.keyRate = make([]float64, 10)
	for i, v := range fin.keyCount {
		fin.keyRate[i] = div(v, len(code))
	}
	fin.leftHand = fin.keyRate[1] + fin.keyRate[2] + fin.keyRate[3] + fin.keyRate[4]
	fin.rightHand = fin.keyRate[6] + fin.keyRate[7] + fin.keyRate[8] + fin.keyRate[9]
	fin.leftHand = fin.leftHand / (fin.leftHand + fin.rightHand) // 归一
	fin.rightHand = 1 - fin.leftHand
	fin.diffHandRate = fin.posRate[0] + fin.posRate[1]
	fin.sameFinRate = div(countSameFin, posSum)
	fin.diffFinRate = fin.posRate[2] + fin.posRate[3] - fin.sameFinRate
	return fin
}
