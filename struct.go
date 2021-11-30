package smq

import (
	"io"
)

type SmqIn struct {
	TextReader     io.Reader // 文本
	DictReader     io.Reader // 赛码表
	IsOutputDict   bool      // 是否输出赛码表
	IsOutputResult bool      // 是否输出赛码分词结果

	BeginPush       int    // 普通码表起顶码长(码长大于等于此数，首选不会追加空格)
	SelectKeys      string // 自定义选重键(2重开始，默认为;')
	IsSingleOnly    bool   // 是否只跑单字
	IsSpaceDiffHand bool   // 空格是否互击
}

type SmqOut struct {
	DictBytes []byte   //赛码表
	WordSlice [][]rune //分词
	CodeSlice []string //编码

	TextLen int //文本字数
	DictLen int //词条数

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
