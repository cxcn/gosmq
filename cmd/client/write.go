package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func writeDict(dn string, db []byte) {
	_ = os.Mkdir("dict", 0666)
	err := ioutil.WriteFile(".\\dict\\"+dn+"_赛码表.txt", db, 0666)
	if err != nil {
		fmt.Println("Error! 赛码表写入错误:", err)
	} else {
		fmt.Println("Success! 成功写入赛码表:", ".\\dict\\"+dn+"_赛码表.txt")
	}
}

func writeResult(tn, dn string, ws [][]rune, cs []string) {
	var buf bytes.Buffer
	for i, v := range ws {
		buf.WriteString(string(v))
		buf.WriteByte('\t')
		buf.WriteString(cs[i])
		buf.WriteByte('\n')
	}
	_ = os.Mkdir("result", 0666)
	err := ioutil.WriteFile(".\\result\\"+tn+"_"+dn+".txt", buf.Bytes(), 0666)
	if err != nil {
		fmt.Println("Error! 输出结果错误:", err)
	} else {
		fmt.Println("Suceess! 成功输出结果:", ".\\result\\"+tn+"_"+dn+".txt")
	}
}
