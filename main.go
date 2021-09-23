package smq

type SmqIn struct { // Smq input
	Fpm  string // 赛码表路径
	Ding int    // 普通码表起顶码长，码长大于等于此数，首选不会追加空格
	IsS  bool   // 是否只跑单字
	IsW  bool   // 是否输出赛码表

	Fpt string // 文本路径
	Csk string // 自定义选重键(2重开始)，custom select keys
	Fpo string // 输出编码路径

	As bool // 空格是否互击

	combs map[string]*comb
	keys  [128]int // 按键所用手指
}

func NewSmq(si *SmqIn) *SmqOut {

	// defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()
	si.combs = newCombs(si.As)
	keys := "1qaz2wsx3edc4rfv5tgb_6yhn7ujm8ik,9ol.0p;/'"
	fins := "111122223333444444445666666667777888899999"
	for i := range keys {
		si.keys[keys[i]] = int(fins[i] - 48)
	}

	so := newSmqOut(si)
	return so
}
