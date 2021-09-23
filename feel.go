package smq

func (so *SmqOut) feel(code string, combs map[string]*comb) {

	if len(code) == 0 {
		return
	}
	for i := 0; i <= len(code)-2; i++ {

		// 处理单键
		if code[i] < 128 {
			so.keyCount[code[i]]++
		} else {
			so.finCount[0]++
			continue
		}

		// 处理按键组合
		comb, ok := combs[code[i:i+2]]
		if !ok {
			continue
		}
		so.eqSum += comb.eq
		so.combLen++
		// 同手
		switch comb.sh {
		case 0:
		case 2:
			so.xkpCount++
		case 3:
			so.dkpCount++
		case 4:
			so.csCount++
		case 1:
			so.skCount++
		}
		if comb.lfd { // 小指干扰
			so.lfdCount++
		}
		// 异手
		switch comb.dist {
		case 1:
			so.handCount[0]++
		case 2:
			so.handCount[1]++
		case 3:
			so.handCount[2]++
		case 4:
			so.handCount[3]++
		}
	}
}
