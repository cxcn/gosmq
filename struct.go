package smq

import (
	"io"
)

type Matcher interface {
	// 插入一个词条 word code order
	Insert(string, string, int)
	// 读取完码表后的操作
	Handle()
	// 匹配下一个词
	Match([]rune, int) (int, string, int)
}

type Dict struct {
	Name   string // 码表名
	Single bool   // 单字模式

	Format string /* 码表格式
	default:默认 本程序赛码表 词\t编码选重\t选重
	jisu:极速赛码表 词\t编码选重
	duoduo:多多格式码表 词\t编码
	jidian:极点格式 编码\t词1 词2 词3
	*/
	SavePath   string // 读取非默认码表格式时自动转换并保存的路径，默认保存在 dict 目录下
	SelectKeys string // 普通码表自定义选重键(默认为_;')
	PushStart  int    // 普通码表起顶码长(码长大于等于此数，首选不会追加空格)

	// 初始化 Matcher
	Algorithm string // 匹配算法 trie:前缀树 order:顺序匹配（极速跟打器） longest:最长匹配
	Matcher   Matcher

	PressSpaceBy   string // 空格按键方式 left|right|both
	ReturnSegments bool   // 是否输出赛码分词结果

	reader  io.Reader // 赛码表 io 流
	length  int       // 词条数
	illegal bool      // 非法输入
}

type Result struct {
	Name      string
	Basic     basic
	Words     words     // 打词
	Collision collision // 选重
	CodeLen   codeLen   // 码长

	Keys    keys  // 按键统计
	Combs   combs // 按键组合
	Fingers fingers
	Hands   hands

	Data export

	toTalEq10 int // 总当量*10
	mapKeys   map[byte]int
	mapNotHan map[rune]struct{}
	mapLack   map[rune]struct{}
	// codes     string
}

// count and rate
type CaR struct {
	Count int
	Rate  float64
}

// 可能要导出的数据
type export struct {
	WordSlice [][]rune // 分词
	CodeSlice []string // 编码
}

// 基础
type basic struct {
	DictLen     int    // 词条数
	TextLen     int    // 文本字数
	NotHan      string // 非汉字
	NotHans     int    // 非汉字数（去重）
	NotHanCount int    // 非汉字计数
	Lack        string // 缺字
	Lacks       int    // 缺字数（去重）
	LackCount   int    // 缺字计数
	Commits     int    // 上屏数
}

// 打词
type words struct {
	Commits CaR         // 打词数
	Chars   CaR         // 打词字数
	Dist    map[int]int // 词长分布统计
}

// 选重
type collision struct {
	Commits CaR         // 选重数
	Chars   CaR         // 选重字数
	Dist    map[int]int // 选重分布统计
}

// 码长
type codeLen struct {
	Total   int         // 全部码长
	PerChar float64     // 字均码长
	Dist    map[int]int // 码长分布统计
}

// 按键 左空格_，右空格+
type keys map[string]*CaR

// 按键组合
type combs struct {
	Count      int     // 按键组合数
	Equivalent float64 // 当量

	DoubleHit  CaR // 同键双击
	TribleHit  CaR // 同键三连击
	SingleSpan CaR // 小跨排
	MultiSpan  CaR // 大跨排

	LongFingersDisturb   CaR // 错手
	LittleFingersDisturb CaR // 小指干扰
}

type fingers struct {
	Dist [11]*CaR // 手指分布，按键盘上的列，第11个是41键以外的
	Same CaR      // 同指
	Diff CaR      // 异指（同手）
}

type hands struct {
	Left  CaR // 左手
	Right CaR // 右手
	Same  CaR // 同手
	Diff  CaR // 异手

	LL CaR `json:"LeftToLeft"`   // 左左
	LR CaR `json:"LeftToRight"`  // 左右
	RL CaR `json:"RightToLeft"`  // 右左
	RR CaR `json:"RightToRight"` // 右右
}
