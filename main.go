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
	fmt.Println("----------------------")
	fmt.Printf("文本字数：%d\n", res.lenText)
	fmt.Printf("非汉字：%s\n", res.notHan)
	fmt.Printf("非汉字数：%d\n", res.countNotHan)
	fmt.Printf("缺字：%s\n", res.lack)
	fmt.Printf("缺字数：%d\n", res.countLack)
	fmt.Printf("总键数：%d\n", res.lenCode)
	fmt.Printf("码长：%.4f\n", res.avlCode)
	fmt.Printf("空格数：%d\n", res.countSpace)
	fmt.Printf("打词数：%d\n", res.countWord)
	fmt.Printf("打词字数：%d\n", res.lenWord)
	fmt.Printf("选重数：%d\n", res.countChoose)
	fmt.Printf("选重字数：%d\n", res.lenChoose)
	fmt.Printf("打词率（上屏）：%.4f%%\n", 100*res.rateWord)
	fmt.Printf("打词率（字数）：%.4f%%\n", 100*res.rateLenWord)
	fmt.Printf("选重率（上屏）：%.4f%%\n", 100*res.rateChoose)
	fmt.Printf("选重率（字数）：%.4f%%\n", 100*res.rateLenChoose)
	fmt.Printf("码长统计：%v\n", res.statCode)
	fmt.Printf("词长统计：%v\n", res.statWord)
	fmt.Println("----------------------")
	// time.Sleep(5 * time.Second)
}

func errHandler(err error) {
	if err != nil {
		fmt.Println("error: ", err)
	}
}
