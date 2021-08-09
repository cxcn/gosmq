package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func read(fpm string, ding int, isW bool) *Trie {

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("read cost time = ", cost)
	}()

	// 读取码表
	var dict = NewTrie()
	f, err := os.Open(fpm)
	if err != nil {
		fmt.Println("码表读取错误:", err)
		return dict
	}

	buff := bufio.NewReader(f)
	if ding < 1 {
		fmt.Println("检测到赛码表...")
		for {
			b, _, eof := buff.ReadLine()
			if eof == io.EOF {
				break
			}
			wc := strings.Split(string(b), "\t")
			if len(wc) == 2 {
				dict.Insert(wc[0], wc[1])
			}
		}
		dict.addPunct()
		return dict
	}

	fmt.Println("检测到普通码表...")
	var wb []byte
	freq := make(map[string]int)
	// 生成字典
	for {
		b, _, eof := buff.ReadLine()
		if eof == io.EOF {
			break
		}
		wc := strings.Split(string(b), "\t")
		if len(wc) != 2 {
			continue
		}
		c := wc[1]
		freq[c] += 1
		if freq[c] != 1 {
			c = c + strconv.Itoa(freq[c])
		} else if len(c) < ding {
			c = c + "_"
		}
		if isW {
			wb = append(wb, []byte(wc[0]+"\t"+c+"\n")...)
		}
		dict.Insert(wc[0], c)
	}
	f.Close()
	dict.addPunct()

	// 写入赛码表
	if isW {
		_, filename := filepath.Split(fpm)
		err := ioutil.WriteFile(".\\smb\\smb_"+filename, wb, 0777)
		if err != nil {
			fmt.Println("赛码表写入错误:", err)
			return dict
		}
		fmt.Println("赛码表写入成功:", ".\\smb\\smb_"+filename)
	}
	return dict
}

func (t *Trie) addPunct() {
	// 符号
	punct := map[string]string{
		"·": "`",
		"【": "[",
		"】": "]",
		"；": ";",
		"‘": "\"",
		"’": "\"",
		"，": ",",
		"。": ".",
		"、": "/",
		// 用 ~ 表示换挡
		"——": "~-",
		"：":  "~;",
		"“":  "~\"",
		"”":  "~\"",
		"《":  "~,",
		"》":  "~.",
		"？":  "~/",
	}
	for k, v := range punct {
		t.Insert(k, v)
	}
}
