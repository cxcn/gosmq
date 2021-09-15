package smq

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func newDict(fpm string, ding int, isW bool, isS bool) *trie {

	// start := time.Now()
	// defer func() {
	// 	cost := time.Since(start)
	// 	fmt.Println("NewDict cost time = ", cost)
	// }()

	_, filename := filepath.Split(fpm)
	// 读取码表
	dict := newTrie()
	f, err := os.Open(fpm)
	if err != nil {
		fmt.Println("码表读取错误:", err)
		return dict
	}
	scan := bufio.NewScanner(f)

	if isS {
		fmt.Println("只跑单字...")
	}
	if ding < 1 {
		fmt.Println("检测到赛码表:", filename)
		for scan.Scan() {
			wc := strings.Split(scan.Text(), "\t")
			if len(wc) != 2 {
				continue
			}
			if isS && len([]rune(wc[0])) != 1 {
				continue
			}
			dict.insert(wc[0], wc[1])
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
		if isS && len([]rune(wc[0])) != 1 {
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
		dict.insert(wc[0], c)
	}
	dict.addPunct()
	f.Close()

	// 写入赛码表
	if isW {
		_ = os.Mkdir("smb", 0666)
		err := ioutil.WriteFile(".\\smb\\"+filename, wb, 0666)
		if err != nil {
			fmt.Println("赛码表写入错误:", err)
		} else {
			fmt.Println("赛码表写入成功:", ".\\smb\\"+filename)
		}
	}
	return dict
}
