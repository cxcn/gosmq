package smq

func (so *SmqOut) stat() {
	so.LackCount = len([]rune(so.Lack))
	so.CodeAvg = div(so.CodeLen, so.TextLen)
	so.WordRate = div(so.WordCount, so.UnitCount)
	so.WordLenRate = div(so.WordLen, so.TextLen)
	so.RepeatRate = div(so.RepeatCount, so.UnitCount)
	so.RepeatLenRate = div(so.RepeatLen, so.TextLen)

	keys := "1234567890qwertyuiopasdfghjkl;zxcvbnm,./'_"
	shiftKeys := "!@#$%^&*()QWERTYUIOPASDFGHJKL:ZXCVBNM<>?\""
	for i, v := range shiftKeys {
		so.keyCount[keys[i]] += so.keyCount[v]
	}

	a := "1qaz2wsx3edc4rfv5tgb_6yhn7ujm8ik,9ol.0p;/'"
	b := "111122223333444444445666666667777888899999"
	var fins [128]int
	for i := range a {
		fins[a[i]] = int(b[i] - 48)
	}

	for i, v := range keys {
		so.KeyRate[i] = div(so.keyCount[v], so.CodeLen)
		so.finCount[fins[v]] += so.keyCount[v]
	}

	for i, v := range so.finCount {
		so.FinRate[i] = div(v, so.CodeLen)
	}

	so.LeftHand = div(so.finCount[1]+so.finCount[2]+so.finCount[3]+so.finCount[4], so.CodeLen-so.finCount[5]-so.finCount[0])
	so.RightHand = 1 - so.LeftHand

	noSpace := so.handCount[0] + so.handCount[1] + so.handCount[2] + so.handCount[3]
	for i, v := range so.handCount {
		so.HandRate[i] = div(v, noSpace)
	}
	so.DiffHandRate = so.HandRate[0] + so.HandRate[1]
	so.SameFinRate = div(so.skCount+so.xkpCount+so.dkpCount, so.combLen)
	so.DiffFinRate = 1 - so.DiffHandRate - so.SameFinRate

	so.Eq = div(so.eqSum, 10) / float64(so.combLen)
	so.Sk = div(so.skCount, so.combLen)
	so.Xkp = div(so.xkpCount, so.combLen)
	so.Dkp = div(so.dkpCount, so.combLen)
	so.Lfd = div(so.lfdCount, so.combLen)
	so.Cs = div(so.csCount, so.combLen)
}

func div(x, y int) float64 {
	return float64(x) / float64(y)
}
