package smq

type Feel struct {
	KeyRate   map[byte]float64
	FinCount  []int
	FinRate   []float64
	LeftHand  float64 // 左手
	RightHand float64 // 右手

	HandCount    []int     // LR RL LL RR
	HandRate     []float64 // LR RL LL RR
	DiffHandRate float64   // 异手
	SameFinRate  float64   // 同指
	DiffFinRate  float64   // 同手异指

	Eq  float64 // 当量 equivalent
	Sk  float64 // 同键 same key
	Xkp float64 // 小跨排
	Dkp float64 // 大跨排
	Cs  float64 // 错手
	Lfd float64 // 小指干扰 little finger disturb
}

func NewFeel(code string, aS bool) *Feel {

	// start := time.Now()
	// defer func() {
	// 	cost := time.Since(start)
	// 	fmt.Println("NewFeel cost time = ", cost)
	// }()

	data := newData(aS)
	feel := new(Feel)
	if len(code) == 0 {
		return feel
	}
	feel.FinCount = make([]int, 10)
	feel.HandCount = make([]int, 4)

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
		eqSum    float64
		skCount  int
		xkpCount int
		dkpCount int
		csCount  int
		lfdCount int
		keyLen   int
		combLen  int
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
			feel.FinCount[0]++
			continue
		}

		// 处理按键组合
		comb := data[code[i-1:i+1]]
		eqSum += comb.eq
		combLen++
		// 同手
		switch comb.sh {
		case 0:
		case 2:
			xkpCount++
		case 3:
			dkpCount++
		case 4:
			csCount++
		case 1:
			skCount++
		}
		if comb.lfd { // 小指干扰
			lfdCount++
		}
		// 异手
		switch comb.dist {
		case 1:
			feel.HandCount[0]++
		case 2:
			feel.HandCount[1]++
		case 3:
			feel.HandCount[2]++
		case 4:
			feel.HandCount[3]++
		}
	}

	feel.KeyRate = make(map[byte]float64)
	for k, v := range keyCount {
		feel.FinCount[finger[k]] += v
		feel.KeyRate[k] = div(v, keyLen)
	}

	feel.FinRate = make([]float64, 10)
	for i, v := range feel.FinCount {
		feel.FinRate[i] = div(v, len(code))
	}
	feel.LeftHand = feel.FinRate[1] + feel.FinRate[2] + feel.FinRate[3] + feel.FinRate[4]
	feel.RightHand = feel.FinRate[6] + feel.FinRate[7] + feel.FinRate[8] + feel.FinRate[9]
	feel.LeftHand = feel.LeftHand / (feel.LeftHand + feel.RightHand) // 归一
	feel.RightHand = 1 - feel.LeftHand

	feel.HandRate = make([]float64, 4)
	handSum := 0
	for _, v := range feel.HandCount {
		handSum += v
	}
	for i, v := range feel.HandCount {
		feel.HandRate[i] = div(v, handSum)
	}
	feel.DiffHandRate = feel.HandRate[0] + feel.HandRate[1]
	feel.SameFinRate = div(skCount+xkpCount+dkpCount, handSum)
	feel.DiffFinRate = feel.HandRate[2] + feel.HandRate[3] - feel.SameFinRate

	feel.Eq = eqSum / float64(combLen)
	feel.Sk = div(skCount, combLen)
	feel.Xkp = div(xkpCount, combLen)
	feel.Dkp = div(dkpCount, combLen)
	feel.Lfd = div(lfdCount, combLen)
	feel.Cs = div(csCount, combLen)

	return feel
}
