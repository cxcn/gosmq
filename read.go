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

func read(fp string) Trie {

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("read cost time = ", cost)
	}()

	var dict = Constructor()
	dict.addPunct()
	conf := readConf(fp)

	f, err := os.Open(fp)
	errHandler(err)
	defer f.Close()
	buff := bufio.NewReader(f)

	if conf.isConf {
		fmt.Println("检测到普通码表", conf)
	} else {
		fmt.Println("检测到赛码表")
		// dict.readSMB(fp)
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
		return dict
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
		word, code := wc[0], wc[1]
		freq[code] += 1
		var key string
		if freq[code] != 1 {
			key = code + strconv.Itoa(freq[code])
		} else if len(code) < conf.as {
			key = code + "_"
		} else {
			key = code
		}
		dict.Insert(word, key)
	}
	return dict
}

func (t Trie) addPunct() {
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
