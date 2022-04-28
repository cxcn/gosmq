package smq

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"sync"
)

type Smq struct {
	Name   string    // 文本名
	Text   io.Reader // 文本
	Inputs []*Dict   // 码表选项
}

// 初始化一个赛码器
func New(name string, rd io.Reader) Smq {
	log.Println("从字节流初始化赛码器，文本名：", name)
	return Smq{name, rd, []*Dict{}}
}

func NewFromString(name, text string) Smq {
	rd := readFromString(text)
	log.Println("从字符串初始化赛码器，文本名：", name)
	return Smq{name, rd, []*Dict{}}
}

func NewFromPath(name, path string) Smq {
	rd, err := readFromPath(path)
	if err != nil {
		log.Panicln("Error! 从文件初始化赛码器，文本名：", name)
	}
	if name == "" {
		name = GetFileName(path)
	}
	log.Println("从文件初始化赛码器，文本名：", name)
	return Smq{name, rd, []*Dict{}}
}

// 添加一个码表
func (smq *Smq) Add(dict *Dict) {
	dict.init()
	// 合法输入
	if !dict.illegal {
		smq.Inputs = append(smq.Inputs, dict)
		log.Println("添加了一个码表：", dict.Name)
	}
	log.Println("怎么回事")
}

func newResult() *Result {
	res := new(Result)
	res.mapKeys = make(map[byte]int)
	res.mapLack = make(map[rune]struct{})
	res.mapNotHan = make(map[rune]struct{})
	res.Words.Dist = make(map[int]int)
	res.Collision.Dist = make(map[int]int)
	res.CodeLen.Dist = make(map[int]int)
	res.Keys = make(keys)
	return res
}

// 开始计算
func (smq *Smq) Run() []*Result {
	smqLen := len(smq.Inputs)
	ret := make([]*Result, 0, smqLen)
	for i := 0; i < len(smq.Inputs); i++ {
		ret = append(ret, newResult())
		log.Println(smq.Inputs[i])
	}
	brd := bufio.NewReader(smq.Text)
	// 逐行读取文本文件
	for {
		line, err := brd.ReadString('\n')
		var wg sync.WaitGroup
		for i, v := range smq.Inputs {
			wg.Add(1)
			tmp := ret[i]
			go func(arg *Dict) {
				codes := tmp.match([]rune(line), arg.Matcher)
				tmp.feel(codes, arg)
				wg.Done()
			}(v)
		}
		wg.Wait()
		if err != nil {
			break
		}
	}
	for i, v := range ret {
		v.stat(smq.Inputs[i])
	}
	return ret
}

// 输出 json
func ResToJson(res []*Result) ([]byte, error) {
	return json.Marshal(res)
}

// 输出 json
func (smq *Smq) ToJson() ([]byte, error) {
	res := smq.Run()
	return ResToJson(res)
}
