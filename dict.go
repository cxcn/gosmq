package smq

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// 从 io 流加载码表
func (dict *Dict) Load(rd io.Reader) {
	dict.reader = Tranformer(rd)
	dict.legal = true
}

// 从字符串流加载码表
func (dict *Dict) LoadFromString(s string) {
	dict.reader = readFromString(s)
	dict.legal = true
}

// 从文件加载码表
func (dict *Dict) LoadFromPath(path string) {
	rd, err := readFromPath(path)
	if err != nil {
		log.Println("Warning! 从文件读取码表失败，路径：", path)
		return
	}
	if dict.Name == "" {
		dict.Name = GetFileName(path)
	}
	dict.reader = rd
	dict.legal = true
}

// 转换赛码表
func (dict *Dict) Convert() {
	// 转换赛码表
	if dict.Transfer == nil {
		switch dict.Format {
		case "jisu", "js":
			dict.Transfer = &jisu{}
		case "duoduo", "dd":
			dict.Transfer = &duoduo{}
		case "jidian", "jd":
			dict.Transfer = &jidian{}
		case "bingling", "bl":
			dict.Transfer = &duoduo{true}
		}
	}
	// 输出赛码表
	if dict.Transfer != nil {
		newBytes := dict.Transfer.Read(dict)
		err := ioutil.WriteFile(dict.SavePath, newBytes, 0666)
		if err != nil {
			// SavePath 不对则保存在 dict 目录下
			os.Mkdir("dict", 0666)
			err = ioutil.WriteFile("./dict/"+dict.Name+".txt", newBytes, 0666)
			if err != nil {
				log.Println(err)
			}
		}
		dict.reader = bytes.NewReader(newBytes)
	}
}

func (dict *Dict) init() {
	// 读取码表
	if dict.SelectKeys == "" {
		dict.SelectKeys = "_;'"
	}
	if dict.PushStart == 0 {
		dict.PushStart = 4
	}
	dict.Convert()
	// 匹配算法
	if dict.Matcher == nil {
		switch dict.Algorithm {
		case "order", "o":
			dict.Matcher = NewOTrie()
		case "longest", "l":
			dict.Matcher = NewLongest()
		case "old_order", "oo":
			dict.Matcher = NewOrder()
		case "trie", "t":
			dict.Matcher = NewTrie()
		default: // "trie"
			dict.Matcher = NewTrie()
		}
	}
	dict.read()
}

func (dict *Dict) read() {
	m := dict.Matcher

	scan := bufio.NewScanner(dict.reader)
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
	for k, v := range puncts {
		m.Insert(k, v, 1)
	}
	m.Handle()
}
