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

func newDict(si *SmqIn) *trie {

	// start := time.Now()
	// defer func() {
	// 	cost := time.Since(start)
	// 	fmt.Println("NewDict cost time = ", cost)
	// }()

	_, filename := filepath.Split(si.Fpm)
	// 读取码表
	dict := newTrie()
	f, err := os.Open(si.Fpm)
	if err != nil {
		fmt.Println("码表读取错误:", err)
		return dict
	}
	scan := bufio.NewScanner(f)
	if si.IsS {
		fmt.Println("只跑单字...")
	}

	if si.Ding < 1 {
		fmt.Println("检测到赛码表:", filename)
		for scan.Scan() {
			wc := strings.Split(scan.Text(), "\t")
			if len(wc) != 2 {
				continue
			}
			if si.IsS && len([]rune(wc[0])) != 1 {
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
		if si.IsS && len([]rune(wc[0])) != 1 {
			continue
		}
		c := wc[1]
		freq[c]++
		rp := freq[c]
		suf := "_"
		if rp != 1 {
			suf = strconv.Itoa(rp)
			c += suf
		} else if si.Ding > len(c) {
			c += suf
		}
		if si.IsW {
			wb = append(wb, scan.Bytes()...)
			if rp != 1 || si.Ding > len(c) {
				wb = append(wb, suf...)
			}
			wb = append(wb, '\n')
		}
		dict.insert(wc[0], c)
	}
	dict.addPunct()
	f.Close()

	// 写入赛码表
	if si.IsW {
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
