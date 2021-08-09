package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func read(fp string, ding int) *Trie {

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("read cost time = ", cost)
	}()

	var dict = NewTrie()
	dict.addPunct()

	f, err := os.Open(fp)
	errHandler(err)
	defer f.Close()
	buff := bufio.NewReader(f)

	if ding > 0 {
		fmt.Println("检测到普通码表")
	} else {
		fmt.Println("检测到赛码表")
	}

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
		if ding > 0 {
			freq[c] += 1
			if freq[c] != 1 {
				c = c + strconv.Itoa(freq[c])
			} else if len(c) < ding {
				c = c + "_"
			}
		}
		dict.Insert(wc[0], c)
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
