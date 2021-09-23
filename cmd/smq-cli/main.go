package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	smq "github.com/cxcn/gosmq"
	"github.com/jessevdk/go-flags"
)

func main() {

	// defer profile.Start().Stop()
	// defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()

	if len(os.Args) == 1 {
		fmt.Println("请在命令行中运行此程序\n按Enter键退出...")
		fmt.Scanln()
		return
	}

	type option struct {
		Fpm  []string `short:"i" long:"input" description:"[]string\t码表路径，可设置多个"`
		Ding int      `short:"d" long:"ding" description:"int\t普通码表起顶码长，码长大于等于此数，首选不会追加空格"`
		IsS  bool     `short:"s" long:"single" description:"bool\t是否只跑单字"`
		IsW  bool     `short:"w" description:"bool\t是否输出赛码表(保存在.\\smb\\文件夹下)"`
		Fpt  string   `short:"t" long:"text" description:"string\tutf8编码格式文本"`
		AS   bool     `short:"k" description:"bool\t空格是否互击"`
		Csk  string   `short:"c" default:";'" description:"string\t自定义选重键(2重开始)"`
		Fpo  string   `short:"o" long:"output" description:"string\t输出编码路径"`
		Ver  bool     `short:"v" long:"version" description:"bool\t查看版本信息"`
	}

	var opt option
	flags.Parse(&opt)
	if opt.Ver {
		fmt.Printf("smq-cli version 0.10 %s/%s\n\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println("repo address: https://github.com/cxcn/gosmq/")
		return
	}

	if len(opt.Fpm) == 0 {
		return
	}

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("main cost time = ", cost)
	}()

	_, textFileName := filepath.Split(opt.Fpt)
	a := NewHTMLOutputInfo(textFileName)
	for _, v := range opt.Fpm {
		si := smq.SmqIn{
			Fpm:  v,
			Ding: opt.Ding,
			IsS:  opt.IsS,
			IsW:  opt.IsW,
			Fpt:  opt.Fpt,
			Csk:  opt.Csk,
			Fpo:  opt.Fpo,
			As:   opt.AS,
		}
		so := smq.NewSmq(&si)
		if so.CodeLen == 0 {
			continue
		}
		_, mbFileName := filepath.Split(v)
		a.AddSchema(mbFileName, so)
		output(so)
	}
	a.OutputHTMLFile("result.html")

	// time.Sleep(5 * time.Second)
}
