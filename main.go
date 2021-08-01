package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {
	// start := time.Now()
	// defer func() {
	// 	cost := time.Since(start)
	// 	fmt.Println("main cost time = ", cost)
	// }()

	var fpm string
	var fpt string
	var fpo string
	var help bool
	// var isCode *bool

	flag.StringVar(&fpm, "m", "", "码表路径")
	flag.StringVar(&fpt, "t", "", "文本路径")
	flag.StringVar(&fpo, "o", "", "输出路径")
	flag.BoolVar(&help, "h", false, "帮助")
	flag.Parse()

	if help {
		fmt.Println("saimaqi version: 0.1.1\n\nUsage: saimaqi.exe [-m mb] [-t text] [-o output]")
		flag.PrintDefaults()
		return
	}

	if fpm == "" || fpt == "" {
		fmt.Println("缺少路径")
		return
	}

	res := cacl(fpm, fpt)
	if fpo != "" {
		// fmt.Println(fpo)
		err := ioutil.WriteFile(fpo, []byte(res.codeSep), 0777)
		errHandler(err)
	}
	fmt.Println(
		"文本字数：", res.lenText,
		// "\n码：", res.code,
		"\n非汉字：", res.notHan,
		"\n非汉字数：", res.countNotHan,
		"\n缺字：", res.lack,
		"\n缺字数：", res.countLack,
		// "\n选重：", res.choose,
		"\n总键数：", res.lenCode,
		"\n码长：", res.avlCode,
		"\n空格数：", res.countSpace,
		"\n打词数：", res.countWord,
		"\n打词率（上屏）：", res.rateWord,
		"\n打词字数：", res.lenWord,
		"\n打词率（字数）：", res.rLenWord,
		"\n选重数：", res.countChoose,
		"\n选重率（上屏）：", res.rateChoose,
		"\n选重字数：", res.lenChoose,
		"\n选重率（字数）：", res.rLenChoose,
		"\n键数：", res.stat,
	)
	test()
	// time.Sleep(5 * time.Second)
}

func test() {

}
