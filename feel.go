package main

type Feel struct {
	keyRate   map[byte]float64
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

func NewFeel(code string, isS bool) *Feel {

	// start := time.Now()
	// defer func() {
	// 	cost := time.Since(start)
	// 	fmt.Println("NewFeel cost time = ", cost)
	// }()

	zhifa := newZhifa(isS)
	feel := new(Feel)
	feel.finCount = make([]int, 10)
	feel.handCount = make([]int, 4)

	keyCount := make(map[byte]int)
	finger := make(map[byte]int)
	keys := "1qaz2wsx3edc4rfv5tgb_6yhn7ujm8ik,9ol.0p;/'"
	fins := "111122223333444444445666666667777888899999"
	for i := range keys {
		keyCount[keys[i]] = 0
		v := int(fins[i] - 48)
		finger[keys[i]] = v
	}
	var (
		sameFinCount int
		dlSum        float64
		dkpCount     int
		xkpCount     int
		xzgrCount    int
		csCount      int
		keyLen       int
		combLen      int
	)

	_, ok := keyCount[code[0]]
	for i := 1; i < len(code); i++ {

		// 处理单键
		if _, okk := keyCount[code[i]]; okk {
			keyLen++
			keyCount[code[i]]++
			if !ok {
				ok = okk
				continue
			}
		} else {
			ok = okk
			feel.finCount[0]++
			continue
		}

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
			feel.handCount[0]++
		case 2:
			feel.handCount[1]++
		case 3:
			feel.handCount[2]++
		case 4:
			feel.handCount[3]++
		}
		// 同指
		if zf.tz {
			sameFinCount++
		}
	}

	feel.keyRate = make(map[byte]float64)
	for k, v := range keyCount {
		feel.finCount[finger[k]] += v
		feel.keyRate[k] = div(v, keyLen)
	}

	feel.finRate = make([]float64, 10)
	for i, v := range feel.finCount {
		feel.finRate[i] = div(v, len(code))
	}
	feel.leftHand = feel.finRate[1] + feel.finRate[2] + feel.finRate[3] + feel.finRate[4]
	feel.rightHand = feel.finRate[6] + feel.finRate[7] + feel.finRate[8] + feel.finRate[9]
	feel.leftHand = feel.leftHand / (feel.leftHand + feel.rightHand) // 归一
	feel.rightHand = 1 - feel.leftHand

	feel.handRate = make([]float64, 4)
	handSum := 0
	for _, v := range feel.handCount {
		handSum += v
	}
	for i, v := range feel.handCount {
		feel.handRate[i] = div(v, handSum)
	}
	feel.diffHandRate = feel.handRate[0] + feel.handRate[1]
	feel.sameFinRate = div(sameFinCount, handSum)
	feel.diffFinRate = feel.handRate[2] + feel.handRate[3] - feel.sameFinRate

	feel.dl = dlSum / float64(combLen)
	feel.dkp = div(dkpCount, combLen)
	feel.xkp = div(xkpCount, combLen)
	feel.xzgr = div(xzgrCount, combLen)
	feel.cs = div(csCount, combLen)

	return feel
}
