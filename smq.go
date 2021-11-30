package smq

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (si *SmqIn) Smq() (*SmqOut, error) {

	so := new(SmqOut)
	// 读取码表
	dict := new(trie)
	rd, err := ReadFile(si.DictReader)
	if err != nil {
		fmt.Println("Error! 码表读取错误:", err)
		return so, err
	}
	if si.BeginPush < 1 {
		fmt.Println("Success! 检测到赛码表:")
		so.DictLen = dict.read(rd, si.IsSingleOnly)
	} else {
		fmt.Println("Success! 检测到普通码表:")
		var wb []byte
		so.DictLen, wb = dict.readC(rd, si.IsSingleOnly, si.IsOutputDict, si.BeginPush)
		so.DictBytes = wb
	}

	if si.IsSingleOnly {
		fmt.Println("Option: 只跑单字。。。")
	}

	so.RepeatStat = make(map[int]int)
	so.CodeStat = make(map[int]int)
	so.WordStat = make(map[int]int)
	//读取文本
	rd, err = ReadFile(si.TextReader)
	if err != nil {
		fmt.Println("Error! 文本读取错误:", err)
		return so, err
	}
	fmt.Println("Success! 成功读取文本:")
	so.calc(rd, dict, si.SelectKeys, si.IsSpaceDiffHand, si.IsOutputResult)

	return so, nil
}

func GetFileName(fp string) string {
	name := filepath.Base(fp)
	ext := filepath.Ext(fp)
	return strings.TrimSuffix(name, ext)
}

func GetDictName(s string) string {
	s = strings.TrimSuffix(s, "赛码表")
	return strings.TrimSuffix(s, "_")
}
