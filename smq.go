package smq

import (
	"bytes"
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
		fmt.Println("Error! 码表读取错误:", err)
		return so
	}
	if si.Ding < 1 {
		fmt.Println("Success! 检测到赛码表:", so.DictName)
		so.DictLen = dict.read(rd, si.IsS)
	} else {
		fmt.Println("Success! 检测到普通码表:", so.DictName)
		var wb []byte
		so.DictLen, wb = dict.readC(rd, si.IsS, si.Ding)
		// 写入赛码表
		if si.IsW {
			_ = os.Mkdir("dict", 0666)
			err := ioutil.WriteFile(".\\dict\\"+so.DictName+"_赛码表.txt", wb, 0666)
			if err != nil {
				fmt.Println("Error! 赛码表写入错误:", err)
			} else {
				fmt.Println("Success! 成功写入赛码表:", ".\\dict\\"+so.DictName+"_赛码表.txt")
			}
		}
	}
	f.Close()

	if si.IsS {
		fmt.Println("Option: 只跑单字。。。")
	}

	so.TextName = GetFileName(si.Fpt)
	so.RepeatStat = make(map[int]int)
	so.CodeStat = make(map[int]int)
	so.WordStat = make(map[int]int)
	//读取文本
	f, rd, err = ReadFile(si.Fpt)
	if err != nil {
		fmt.Println("Error! 文本读取错误:", err)
		return so
	}
	fmt.Println("Success! 成功读取文本:", so.TextName)
	so.calc(rd, dict, si.Csk, si.As, si.IsO)
	f.Close()

	if si.IsO {
		var wb bytes.Buffer
		for i, v := range so.WordSlice {
			wb.WriteString(string(v))
			wb.WriteByte('\t')
			wb.WriteString(so.CodeSlice[i])
			wb.WriteByte('\n')
		}
		_ = os.Mkdir("result", 0666)
		err := ioutil.WriteFile(".\\result\\"+so.TextName+"_"+so.DictName+".txt", wb.Bytes(), 0666)
		if err != nil {
			fmt.Println("Error! 输出结果错误:", err)
		} else {
			fmt.Println("Suceess! 成功输出结果:", ".\\result\\"+so.TextName+"_"+so.DictName+".txt")
		}
	}
	return so
}

func GetFileName(fp string) string {
	name := filepath.Base(fp)
	ext := filepath.Ext(fp)
	return strings.TrimSuffix(name, ext)
}
