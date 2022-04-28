package smq

func (res *Result) feel(codes string, dict *Dict) {
	if len(codes) == 0 {
		return
	}
	if res.Keys[codes[0]] != nil {
		res.Keys[codes[0]].Count++
	} else {
		res.Keys[codes[0]] = new(CaR)
	}
	last := keyData[codes[0]]
	if last == nil {
		last = new(key)
		last.key = codes[0]
	}
	for i := 1; i < len(codes); i++ {
		current := codes[i]
		// for key
		if current == '_' {
			switch dict.PressSpaceBy {
			case "left":
			case "right":
				current = '+'
			default: // "both"
				// 如果上一个键是左手
				if !last.lor {
					current = '+'
				}
			}
		}
		if res.Keys[current] != nil {
			res.Keys[current].Count++
		} else {
			res.Keys[current] = new(CaR)
		}

		if keyData[current] == nil {
			last.key = current
			continue
		}
		currentData := keyData[current]

		// for comb
		comb := combData[string([]byte{last.key, current})]
		if comb == nil {
			// log.Printf(`comb nil"%v"%v"%v`, last.key, current, comb)

			last.key = currentData.key
			last.fin = currentData.fin
			last.lor = currentData.lor
			continue
		}
		res.Combs.Count++
		// for finger
		if currentData.fin == last.fin {
			res.Fingers.Same.Count++
		}
		// for hands
		if currentData.lor {
			if last.lor { // RR
				res.Hands.RR.Count++
			} else { // RL
				res.Hands.RL.Count++
			}
		} else {
			if last.lor { // LR
				res.Hands.LR.Count++
			} else { // LL
				res.Hands.LL.Count++
			}
		}
		// for comb
		if current == last.key {
			res.Combs.DoubleHit.Count++
			if i < len(codes)-1 {
				if current == codes[i+1] {
					res.Combs.TribleHit.Count++
				}
			}
		}
		res.toTalEq10 += comb.eq
		switch comb.sh {
		case 2: // 小跨排
			res.Combs.SingleSpan.Count++
		case 3: // 大跨排
			res.Combs.MultiSpan.Count++
		case 4: // 错手
			res.Combs.LongFingersDisturb.Count++
		}
		if comb.lfd { // 小拇指干扰
			res.Combs.LittleFingersDisturb.Count++
		}
		last.key = currentData.key
		last.fin = currentData.fin
		last.lor = currentData.lor
	}

}
