package smq

type SmqOut struct {
	TextLen     int     //文本字数
	MbLen       int     //词条数
	NotHan      string  //非汉字
	NotHanCount int     //非汉字数
	Lack        string  //缺字
	LackCount   int     //缺字数
	UnitCount   int     //上屏数
	CodeLen     int     //总键数
	CodeAvg     float64 //码长

	CodeStat   map[int]int //码长统计
	WordStat   map[int]int //词长统计
	RepeatStat map[int]int //选重统计

	WordCount   int     //打词数
	WordLen     int     //打词字数
	WordRate    float64 //打词率（上屏）
	WordLenRate float64 //打词率（字数）

	RepeatCount   int     //选重数
	RepeatLen     int     //选重字数
	RepeatRate    float64 //选重率（上屏）
	RepeatLenRate float64 //选重率（字数）

	// 下面是手感部分

	eqSum     int // 总当量*10
	skCount   int // 同键
	xkpCount  int // 小跨排
	dkpCount  int // 大跨排
	csCount   int // 错手
	lfdCount  int // 小指干扰
	combLen   int // 按键组合数
	keyCount  [128]int
	finCount  [10]int
	handCount [4]int // LR RL LL RR

	KeyRate   [42]float64
	FinRate   [10]float64
	LeftHand  float64 // 左手
	RightHand float64 // 右手

	HandRate     [4]float64 // LR RL LL RR
	DiffHandRate float64    // 异手
	SameFinRate  float64    // 同指
	DiffFinRate  float64    // 同手异指

	Eq  float64 // 当量 equivalent
	Sk  float64 // 同键 same key
	Xkp float64 // 小跨排
	Dkp float64 // 大跨排
	Cs  float64 // 错手
	Lfd float64 // 小指干扰 little finger disturb
}

func (so *SmqOut) stat(si *SmqIn) {
	so.LackCount = len([]rune(so.Lack))
	so.CodeAvg = div(so.CodeLen, so.TextLen)
	so.WordRate = div(so.WordCount, so.UnitCount)
	so.WordLenRate = div(so.WordLen, so.TextLen)
	so.RepeatRate = div(so.RepeatCount, so.UnitCount)
	so.RepeatLenRate = div(so.RepeatLen, so.TextLen)

	keyLen := 0
	for i, v := range so.keyCount {
		if key := si.keys[i]; key != 0 {
			so.finCount[si.keys[i]] += v
			keyLen += v
		} else {
			so.finCount[0] += v
		}
	}
	for i, v := range "1234567890qwertyuiopasdfghjkl;zxcvbnm,./'_" {
		so.KeyRate[i] = div(so.keyCount[v], keyLen)
	}
	for i, v := range so.finCount {
		so.FinRate[i] = div(v, so.CodeLen)
	}

	so.LeftHand = div(so.finCount[1]+so.finCount[2]+so.finCount[3]+so.finCount[4], keyLen-so.finCount[5]-so.finCount[0])
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
