package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"
)

func main() {

	// defer profile.Start().Stop()
	// defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("main cost time = ", cost)
	}()

	var (
		fpm  string // file path mb
		ding int    // 是否输出赛码表
		isW  bool   // file path write smb
		fpt  string // file path text
		isS  bool   // 空格是否互击
		csk  string // custom select keys
		fpo  string // output file path
		help bool
	)

	flag.StringVar(&fpm, "i", "", "码表路径，可以是rime格式码表 或 极速跟打器赛码表")
	flag.IntVar(&ding, "d", 0, "普通码表起顶码长，码长大于等于此数，首选不会追加空格")
	flag.BoolVar(&isW, "w", false, "是否输出赛码表(保存在.\\smb\\文件夹下)")
	flag.StringVar(&fpt, "t", "", "文本路径，utf8编码格式文本，会自动去除空白符")
	flag.BoolVar(&isS, "s", false, "空格是否互击")
	flag.StringVar(&csk, "k", ";'", "自定义选重键(2重开始)")
	flag.StringVar(&fpo, "o", "", "输出路径")
	flag.BoolVar(&help, "h", false, "显示帮助")
	flag.Parse()

	if help {
		fmt.Print("saimaqi version: 0.4\n\n")
		fmt.Print("Usage: saimaqi.exe [-i mb] [-d int] [-w] [-t text] [-s] [-k string] [-o output]\n\n")
		flag.PrintDefaults()
		return
	}

	dict := NewDict(fpm, ding, isW)
	if len(dict.children) == 0 {
		return
	}
	smq := NewSmq(dict, fpt, csk)
	if smq.textLen == 0 {
		return
	}
	fin := NewFin(smq.code, isS)

	if fpo != "" {
		// fmt.Println(fpo)
		_ = ioutil.WriteFile(fpo, []byte(smq.codeSep), 0777)
	}
	out := ""
	out += fmt.Sprintln("----------------------")
	out += fmt.Sprintf("文本字数：%d\n", smq.textLen)
	out += fmt.Sprintf("非汉字：%s\n", smq.notHan)
	out += fmt.Sprintf("非汉字数：%d\n", smq.notHanCount)
	out += fmt.Sprintf("缺字：%s\n", smq.lack)
	out += fmt.Sprintf("缺字数：%d\n", smq.lackCount)
	out += fmt.Sprintln("------------")
	out += fmt.Sprintf("总键数：%d\n", smq.codeLen)
	out += fmt.Sprintf("码长：%.4f\n", smq.codeAvg)
	out += fmt.Sprintf("打词：    %d\t%.3f%%\n", smq.wordCount, 100*smq.wordRate)
	out += fmt.Sprintf("打词字数：%d\t%.3f%%\n", smq.wordLen, 100*smq.wordLenRate)
	out += fmt.Sprintf("选重：    %d\t%.3f%%\n", smq.repeatCount, 100*smq.repeatRate)
	out += fmt.Sprintf("选重字数：%d\t%.3f%%\n", smq.repeatLen, 100*smq.repeatLenRate)
	out += fmt.Sprintf("码长统计：%v\n", smq.codeStat)
	out += fmt.Sprintf("词长统计：%v\n", smq.wordStat)
	out += fmt.Sprintln("------------")

	out += fmt.Sprintf("左右：%d\t%.3f%%\n", fin.posCount[0], 100*fin.posRate[0])
	out += fmt.Sprintf("右左：%d\t%.3f%%\n", fin.posCount[1], 100*fin.posRate[1])
	out += fmt.Sprintf("左左：%d\t%.3f%%\n", fin.posCount[2], 100*fin.posRate[2])
	out += fmt.Sprintf("右右：%d\t%.3f%%\n", fin.posCount[3], 100*fin.posRate[3])
	out += fmt.Sprintf("异手：%.3f%%\n", 100*fin.diffHandRate)
	out += fmt.Sprintf("同指：%.3f%%\n", 100*fin.sameFinRate)
	out += fmt.Sprintf("异指：%.3f%%\n", 100*fin.diffFinRate)
	for i, v := range fin.keyCount {
		out += fmt.Sprintf("第%d列：%d\t%.3f%%\n", i, v, 100*fin.keyRate[i])
	}
	out += fmt.Sprintln("----------------------")
	fmt.Print(out)
	// time.Sleep(5 * time.Second)
}
