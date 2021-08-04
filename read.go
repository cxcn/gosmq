package main

import (
	"bufio"
	"fmt"
	"io"
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
		el := strings.Split(string(b), "\t")
		if len(el) != 2 {
			continue
		}

		word, code := el[0], el[1]
		if conf.isConf {
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
		} else {
			addWord(word, code)
		}
	}
	return smb
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
