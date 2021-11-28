package smq

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (si *SmqIn) Smq() *SmqOut {
	so := new(SmqOut)
	so.DictName = GetFileName(si.Fpd)

	// 读取码表
	dict := new(trie)
	f, rd, err := ReadFile(si.Fpd)
	if err != nil {
		fmt.Println("码表读取错误:", err)
		return so
	}
	if si.Ding < 1 {
		fmt.Println("检测到赛码表:", so.DictName)
		so.DictLen = dict.read(rd, si.IsS)
	} else {
		fmt.Println("检测到普通码表:", so.DictName)
		var wb []byte
		so.DictLen, wb = dict.readC(rd, si.IsS, si.Ding)
		// 写入赛码表
		_ = os.Mkdir("dict", 0666)
		err := ioutil.WriteFile(".\\dict\\"+so.DictName+"_赛码表.txt", wb, 0666)
		if err != nil {
			fmt.Println("赛码表写入错误:", err)
		} else {
			fmt.Println("赛码表写入成功:", ".\\dict\\"+so.DictName+"_赛码表.txt")
		}
	}
	if si.IsS {
		fmt.Println("只跑单字……")
	}
	f.Close()

	so.TextName = GetFileName(si.Fpt)
	so.RepeatStat = make(map[int]int)
	so.CodeStat = make(map[int]int)
	so.WordStat = make(map[int]int)
	//读取文本
	f, rd, err = ReadFile(si.Fpt)
	if err != nil {
		fmt.Println("文本读取错误:", err)
		return so
	}
	fmt.Println("文本读取成功:", so.TextName)
	so.calc(rd, dict, si.Csk, si.As)
	f.Close()

	return so
}

func GetFileName(fp string) string {
	name := filepath.Base(fp)
	ext := filepath.Ext(fp)
	return strings.TrimSuffix(name, ext)
}
