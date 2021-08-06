package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"
)

func main() {
	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("main cost time = ", cost)
	}()

	var fpm string
	var fpt string
	var fpo string
	var help bool

	flag.StringVar(&fpm, "m", "", "码表路径")
	flag.StringVar(&fpt, "t", "", "文本路径")
	flag.StringVar(&fpo, "o", "", "输出路径")
	flag.BoolVar(&help, "h", false, "帮助")
	flag.Parse()

	if help {
		fmt.Println("saimaqi version: 0.3\n\nUsage: saimaqi.exe [-m mb] [-t text] [-o output]")
		flag.PrintDefaults()
		return
	}

	if fpm == "" || fpt == "" {
		fmt.Println("缺少路径")
		return
	}

	res := calc(fpm, fpt)
	if fpo != "" {
		// fmt.Println(fpo)
		err := ioutil.WriteFile(fpo, []byte(res.codeSep), 0777)
		errHandler(err)
	}
	out := ""
	out += fmt.Sprintln("----------------------")
	out += fmt.Sprintf("文本字数：%d\n", res.lenText)
	out += fmt.Sprintf("非汉字：%s\n", res.notHan)
	out += fmt.Sprintf("非汉字数：%d\n", res.countNotHan)
	out += fmt.Sprintf("缺字：%s\n", res.lack)
	out += fmt.Sprintf("缺字数：%d\n", res.countLack)
	out += fmt.Sprintln("------------")
	out += fmt.Sprintf("总键数：%d\n", res.lenCode)
	out += fmt.Sprintf("码长：%.4f\n", res.avlCode)
	out += fmt.Sprintf("空格数：%d\n", res.countSpace)
	out += fmt.Sprintf("打词：    %d\t%.3f%%\n", res.countWord, 100*res.rateWord)
	out += fmt.Sprintf("打词字数：%d\t%.3f%%\n", res.lenWord, 100*res.rateLenWord)
	out += fmt.Sprintf("选重：    %d\t%.3f%%\n", res.countChoose, 100*res.rateChoose)
	out += fmt.Sprintf("选重字数：%d\t%.3f%%\n", res.lenChoose, 100*res.rateLenChoose)
	out += fmt.Sprintf("码长统计：%v\n", res.statCode)
	out += fmt.Sprintf("词长统计：%v\n", res.statWord)
	out += fmt.Sprintln("------------")
	out += fmt.Sprintf("异手：%.3f%%\n", 100*res.rateDiffHand)
	out += fmt.Sprintf("同指：%.3f%%\n", 100*res.rateSameFin)
	out += fmt.Sprintf("异指：%.3f%%\n", 100*res.rateDiffFin)
	for i, v := range res.countKey {
		out += fmt.Sprintf("第%d列：%d\t%.3f%%\n", i+1, v, 100*res.rateKey[i])
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
