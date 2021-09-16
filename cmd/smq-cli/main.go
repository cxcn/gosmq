package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	smq "github.com/cxcn/gosmq"
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
		ding int    // 普通码表起顶码长
		isS  bool   // 是否只跑单字
		isW  bool   // 是否输出赛码表
		fpt  string // file path text
		aF   bool   // 是否关闭手感统计
		aS   bool   // 空格是否互击
		csk  string // custom select keys
		fpo  string // output file path
		fpc  string // compare
		help bool
	)

	flag.StringVar(&fpm, "i", "", "码表路径，可以是rime格式码表 或 极速跟打器赛码表")
	flag.IntVar(&ding, "n", 0, "普通码表起顶码长，码长大于等于此数，首选不会追加空格")
	flag.BoolVar(&isS, "d", false, "是否只跑单字")
	flag.BoolVar(&isW, "w", false, "是否输出赛码表(保存在.\\smb\\文件夹下)")
	flag.StringVar(&fpt, "t", "", "文本路径，utf8编码格式文本，会自动去除空白符")
	flag.BoolVar(&aF, "f", false, "是否关闭手感统计")
	flag.BoolVar(&aS, "s", false, "空格是否互击")
	flag.StringVar(&csk, "k", ";'", "自定义选重键(2重开始)")
	flag.StringVar(&fpo, "o", "", "输出路径")
	flag.StringVar(&fpc, "c", "", "对比码表路径，只能是赛码表")
	flag.BoolVar(&help, "h", false, "显示帮助")
	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Println("请在命令行中运行此程序\n按Enter键退出...")
		fmt.Scanln(&fpo)
		return
	}
	if help || len(fpm) == 0 {
		fmt.Print("smq-cli version: 0.7, github: https://github.com/cxcn/gosmq\n\n")
		fmt.Print("Usage: smq-cli.exe [-i mb] [-n int] [-d] [-w] [-t text] [-f] [-s] [-k string] [-o output] [-c compare]\n\n")
		flag.PrintDefaults()
		return
	}
	fmt.Println()

	si := smq.SmqIn{
		Fpm:  fpm,
		Ding: ding,
		IsS:  isS,
		IsW:  isW,
		Fpt:  fpt,
		Csk:  csk,
		Fpo:  fpo,
	}
	so := smq.NewSmq(si)

	if so.CodeLen == 0 {
		return
	}

	soc := new(smq.SmqOut)
	if len(fpc) != 0 {
		sic := smq.SmqIn{
			Fpm:  fpc,
			Ding: 0,
			IsS:  isS,
			IsW:  false,
			Fpt:  fpt,
			Csk:  csk,
			Fpo:  "",
		}
		soc = smq.NewSmq(sic)
	}

	output(so, soc, aF, aS)
	// time.Sleep(5 * time.Second)
}
