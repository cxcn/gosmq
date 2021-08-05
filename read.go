package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Dict map[string]string

var smb = make(map[int]Dict)

func read(fp string) map[int]Dict {

	addPunct()
	readConf(fp)

	if conf.isConf {
		fmt.Println("检测到普通码表", conf)
	} else {
		fmt.Println("检测到赛码表")
		readSMB(fp)
		return smb
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
		addWord(word, code+key)
	}
	return smb
}

func readSMB(fp string) {
	f, err := ioutil.ReadFile(fp)
	errHandler(err)
	fs := string(f)
	fs = strings.ReplaceAll(fs, "\r", "")
	words := strings.Split(fs, "\n")
	for _, v := range words {
		wc := strings.Split(v, "\t")
		if len(wc) == 2 {
			addWord(wc[0], wc[1])
		}
	}
}

func addWord(w, c string) {
	l := len([]rune(w))
	if smb[l] == nil {
		smb[l] = make(Dict)
	}
	smb[l][w] = c
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
		addWord(k, v)
	}
}
