package smq

import (
	"bufio"
	"log"
	"strconv"
	"strings"
)

func (dict *Dict) init() {
	// 读取码表
	if dict.Reader == nil {
		if dict.String == "" {
			if dict.Path == "" {
				log.Println("没有输入码表")
				dict.illegal = true
				return
			} else {
				rd, err := readFromPath(dict.Path)
				if err != nil {
					log.Println("Warning! 从文件读取码表失败，路径：", dict.Path)
					dict.illegal = true
					return
				}
				dict.Reader = rd
			}
		} else {
			dict.Reader = readFromString(dict.String)
		}
	}
	if dict.SelectKeys == "" {
		dict.SelectKeys = "_;'"
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
	scan := bufio.NewScanner(dict.Reader)
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
