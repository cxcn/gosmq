package smq

func (so *SmqOut) feel(codeSlice []string, combMap map[string]*comb) {

	var keyComb string
	for i, keys := range codeSlice {
		for j := 0; j < len(keys); j++ {
			// 处理单键
			if keys[j] < 128 {
				so.keyCount[keys[j]]++
			} else {
				so.finCount[0]++
				continue
			}

			if j < len(keys)-1 {
				keyComb = keys[j : j+2]
			} else if i < len(codeSlice)-1 && len(codeSlice[i+1]) > 0 {
				keyComb = string([]byte{keys[j], codeSlice[i+1][0]})
			} else {
				continue
			}

			// 处理按键组合
			comb, ok := combMap[keyComb]
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
}
