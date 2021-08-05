package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type Dict map[string]string

var dict = Constructor()

func read(fp string) {

	addPunct()
	readConf(fp)

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("read cost time = ", cost)
	}()

	if conf.isConf {
		fmt.Println("检测到普通码表", conf)
	} else {
		fmt.Println("检测到赛码表")
		readSMB(fp)
		return
	}

	f, err := os.Open(fp)
	errHandler(err)
	defer f.Close()

	freq := make(map[string]int)
	buff := bufio.NewReader(f)
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
		key := ""
		if freq[code] == 1 {
			if len(code) < conf.as {
				key = "_"
			}
		} else {
			key = strconv.Itoa(freq[code])
		}
		dict.Insert(word, code+key)
	}
}

func readSMB(fp string) {
	f, err := ioutil.ReadFile(fp)
	errHandler(err)
	fs := string(f)
	fs = strings.Replace(fs, "\r", "", -1)
	words := strings.Split(fs, "\n")
	for _, v := range words {
		wc := strings.Split(v, "\t")
		if len(wc) == 2 {
			dict.Insert(wc[0], wc[1])
		}
	}
}

func addPunct() {
	// 符号
	punct := Dict{
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
		dict.Insert(k, v)
	}
}
