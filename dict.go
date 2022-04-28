package smq

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"
)

// 从 io 流加载码表
func (dict *Dict) Load(rd io.Reader) {
	dict.reader = Tranformer(rd)
}

// 从字符串流加载码表
func (dict *Dict) LoadFromString(s string) {
	dict.reader = readFromString(s)
}

// 从文件加载码表
func (dict *Dict) LoadFromPath(path string) {
	rd, err := readFromPath(path)
	if err != nil {
		log.Println("Warning! 从文件读取码表失败，路径：", path)
		dict.illegal = true
		return
	}
	if dict.Name == "" {
		dict.Name = GetFileName(path)
	}
	dict.reader = rd
}

func (dict *Dict) init() {
	// 读取码表
	if dict.SelectKeys == "" {
		dict.SelectKeys = "_;'"
	}
	if dict.PushStart == 0 {
		dict.PushStart = 4
	}
	// 转换、输出赛码表
	// 非本程序格式只支持前缀树算法
	switch dict.Format {
	case "jisu":
		dict.fromJisu()
		return
	case "duoduo":
		dict.fromDuoduo()
		return
	case "jidian":
		dict.fromJidian()
		return
	}
	// 本程序格式支持所有算法
	// 外部算法
	if dict.Matcher != nil {
		dict.read()
		return
	}
	switch dict.Algorithm {
	case "order":
		dict.Matcher = NewOrder()
	case "longest":
		dict.Matcher = NewLongest()
	default: // "trie"
		dict.Matcher = NewTrie()
	}
	dict.read()
}

func (dict *Dict) read() {
	m := dict.Matcher
	scan := bufio.NewScanner(dict.reader)
	// 生成字典
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) != 3 {
			continue
		}
		if dict.Single && len([]rune(wc[0])) != 1 {
			continue
		}
		order, err := strconv.Atoi(wc[2])
		if err != nil {
			order = 1
		}
		m.Insert(wc[0], wc[1], order)
		dict.length++
	}
	// 添加符号
	for _, v := range puncts.o {
		m.Insert(v.word, v.code, v.order)
	}
	m.Handle()
}
