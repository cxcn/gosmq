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

	var fpm string
	var ding int
	var fpt string
	var csk string
	var fpo string
	var help bool

	flag.BoolVar(&help, "h", false, "显示帮助")
	flag.StringVar(&fpm, "i", "", "码表路径，可以是rime格式码表 或 极速跟打器赛码表")
	flag.IntVar(&ding, "d", 0, "普通码表起顶码长，码长大于等于此数，首选不会追加空格")
	flag.StringVar(&fpt, "t", "", "文本路径，utf8编码格式文本，会自动去除空白符")
	flag.StringVar(&csk, "s", ";'", "custom_select_key: 自定义选重键(2重开始)")
	flag.StringVar(&fpo, "o", "", "输出路径")
	flag.Parse()

	if help {
		fmt.Print("saimaqi version: 0.4\n\nUsage: saimaqi.exe [-i mb] [-d int] [-t text] [-s string] [-o output]\n\n")
		flag.PrintDefaults()
		return
	}

	if fpm == "" || fpt == "" {
		fmt.Println("缺少路径")
		return
	}

	dict := read(fpm, ding)
	text := readText(fpt)
	res := calc(dict, text, csk)
	res.fingering()
	if fpo != "" {
		// fmt.Println(fpo)
		err := ioutil.WriteFile(fpo, []byte(res.codeSep), 0777)
		errHandler(err)
	}
	out := ""
	out += fmt.Sprintln("----------------------")
	out += fmt.Sprintf("文本字数：%d\n", res.textLen)
	out += fmt.Sprintf("非汉字：%s\n", res.notHan)
	out += fmt.Sprintf("非汉字数：%d\n", res.notHanCount)
	out += fmt.Sprintf("缺字：%s\n", res.lack)
	out += fmt.Sprintf("缺字数：%d\n", res.lackCount)
	out += fmt.Sprintln("------------")
	out += fmt.Sprintf("总键数：%d\n", res.codeLen)
	out += fmt.Sprintf("码长：%.4f\n", res.codeAvg)
	out += fmt.Sprintf("打词：    %d\t%.3f%%\n", res.wordCount, 100*res.wordRate)
	out += fmt.Sprintf("打词字数：%d\t%.3f%%\n", res.wordLen, 100*res.wordLenRate)
	out += fmt.Sprintf("选重：    %d\t%.3f%%\n", res.repeatCount, 100*res.repeatRate)
	out += fmt.Sprintf("选重字数：%d\t%.3f%%\n", res.repeatLen, 100*res.repeatLenRate)
	out += fmt.Sprintf("码长统计：%v\n", res.codeStat)
	out += fmt.Sprintf("词长统计：%v\n", res.wordStat)
	out += fmt.Sprintln("------------")
	out += fmt.Sprintf("左右：%d\t%.3f%%\n", res.posCount[0], 100*res.posRate[0])
	out += fmt.Sprintf("右左：%d\t%.3f%%\n", res.posCount[1], 100*res.posRate[1])
	out += fmt.Sprintf("左左：%d\t%.3f%%\n", res.posCount[2], 100*res.posRate[2])
	out += fmt.Sprintf("右右：%d\t%.3f%%\n", res.posCount[3], 100*res.posRate[3])
	out += fmt.Sprintf("异手：%.3f%%\n", 100*res.diffHandRate)
	out += fmt.Sprintf("同指：%.3f%%\n", 100*res.sameFinRate)
	out += fmt.Sprintf("异指：%.3f%%\n", 100*res.diffFinRate)
	for i, v := range res.keyCount {
		out += fmt.Sprintf("第%d列：%d\t%.3f%%\n", i, v, 100*res.keyRate[i])
	}
	out += fmt.Sprintln("----------------------")
	fmt.Print(out)
	// time.Sleep(5 * time.Second)
}

func errHandler(err error) {
	if err != nil {
		fmt.Println("error: ", err)
	}
}
