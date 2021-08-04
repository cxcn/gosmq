package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func readText(fp string) []rune {

	var text []rune
	// 读文件
	fb, err := ioutil.ReadFile(fp)
	if err != nil {
		fmt.Println(err)
		return text
	}
	// 去除空白字符
	fs := string(fb)
	// for _, v := range fs {
	// 	if !unicode.IsSpace(v) {
	// 		mb = append(mb, v)
	// 	}
	// }
	str := "\r\n\t "
	for _, v := range str {
		fs = strings.Replace(fs, string(v), "", -1)
	}
	text = []rune(fs)
	return text
}
