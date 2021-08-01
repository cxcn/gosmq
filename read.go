package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func readText(fp string) []rune {
	var res []rune
	// 读文件
	fb, err := ioutil.ReadFile(fp)
	if err != nil {
		fmt.Println(err)
		return res
	}
	// 去除空白字符
	fs := string(fb)
	// for _, v := range fs {
	// 	if !unicode.IsSpace(v) {
	// 		res = append(res, v)
	// 	}
	// }
	str := "\r\n\t "
	for _, v := range str {
		fs = strings.Replace(fs, string(v), "", -1)
	}
	res = []rune(fs)
	return res
}

type Dict map[string]string

func readMB(fp string) map[int]Dict {

	res := make(map[int]Dict)

	// 读取符号配置文件
	// b, err := ioutil.ReadFile("conf.yaml")
	// errHandler(err)
	// conf := make(map[string]Dict)
	// err = yaml.Unmarshal(b, &conf)
	// errHandler(err)

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
		// 用~表示换挡
		"——": "~-",
		"：":  "~;",
		"“":  "~\"",
		"”":  "~\"",
		"《":  "~,",
		"》":  "~.",
		"？":  "~/",
	}

	// 符号
	for k, v := range punct {
		l := len([]rune(k))
		if res[l] == nil {
			res[l] = make(Dict)
		}
		res[l][k] = v
	}

	// start := time.Now()
	// defer func() {
	// 	cost := time.Since(start)
	// 	fmt.Println("readMB cost time = ", cost)
	// }()

	f, err := os.Open(fp)
	errHandler(err)
	defer f.Close()

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
		l := len([]rune(el[0]))
		if res[l] == nil {
			res[l] = make(Dict)
		}
		res[l][el[0]] = el[1]
	}
	return res
}

func errHandler(err error) {
	if err != nil {
		fmt.Println("error: ", err)
	}
}
