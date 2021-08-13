package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func NewDict(fpm string, ding int, isW bool, isD bool) *Trie {

	// start := time.Now()
	// defer func() {
	// 	cost := time.Since(start)
	// 	fmt.Println("NewDict cost time = ", cost)
	// }()

	_, filename := filepath.Split(fpm)
	// 读取码表
	dict := NewTrie()
	f, err := os.Open(fpm)
	if err != nil {
		fmt.Println("码表读取错误:", err)
		return dict
	}
	scan := bufio.NewScanner(f)

	if ding < 1 {
		fmt.Println("检测到赛码表:", filename)
		for scan.Scan() {
			wc := strings.Split(scan.Text(), "\t")
			if len(wc) != 2 {
				continue
			}
			if isD && len([]rune(wc[0])) != 1 {
				continue
			}
			dict.Insert(wc[0], wc[1])
		}
		dict.addPunct()
		return dict
	}

	fmt.Println("检测到普通码表:", filename)
	var wb []byte
	freq := make(map[string]int)
	// 生成字典
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) != 2 {
			continue
		}
		if isD && len([]rune(wc[0])) != 1 {
			continue
		}
		c := wc[1]
		freq[c]++
		suf := ""
		if freq[c] != 1 {
			suf = strconv.Itoa(freq[c])
			c = c + suf
		} else if len(c) < ding {
			suf = "_"
			c = c + suf
		}
		if isW {
			wb = append(wb, scan.Bytes()...)
			wb = append(wb, suf...)
			wb = append(wb, '\n')
		}
		dict.Insert(wc[0], c)
	}
	f.Close()
	dict.addPunct()

	// 写入赛码表
	if isW {
		err := ioutil.WriteFile(".\\smb\\"+filename, wb, 0777)
		if err != nil {
			fmt.Println("赛码表写入错误:", err)
			return dict
		}
		fmt.Println("赛码表写入成功:", ".\\smb\\"+filename)
	}
	return dict
}
