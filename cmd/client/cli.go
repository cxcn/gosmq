package main

import (
	"fmt"
	"os"
	"time"

	smq "github.com/cxcn/gosmq"

	"github.com/jessevdk/go-flags"
)

func cli() {

	type option struct {
		Fpd  []string `short:"i" long:"input" description:"[]string\t码表路径，可设置多个"`
		Ding int      `short:"d" long:"ding" description:"int\t普通码表起顶码长，码长大于等于此数，首选不会追加空格"`
		IsS  bool     `short:"s" long:"single" description:"bool\t是否只跑单字"`

		Fpt string `short:"t" long:"text" description:"string\t文本"`
		Csk string `short:"c" default:";'" description:"string\t自定义选重键(2重开始)"`
		AS  bool   `short:"k" description:"bool\t空格是否互击"`

		IsO bool `short:"o" long:"output" description:"bool\t是否输出结果"`
		Ver bool `short:"v" long:"version" description:"bool\t查看版本信息"`
	}

	var opt option
	flags.Parse(&opt)
	if opt.Ver {
		printInfo()
		return
	}

	if len(opt.Fpd) == 0 {
		return
	}

	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("main cost time = ", cost)
	}()

	tn := smq.GetFileName(opt.Fpt)
	for _, v := range opt.Fpd {

		text, err := os.Open(opt.Fpt)
		if err != nil {
			panic(err)
		}
		dict, err := os.Open(v)
		if err != nil {
			panic(err)
		}

		si := smq.SmqIn{
			TextReader:     text,
			DictReader:     dict,
			IsOutputResult: opt.IsO,

			BeginPush:       opt.Ding,
			SelectKeys:      opt.Csk,
			IsSingleOnly:    opt.IsS,
			IsSpaceDiffHand: opt.AS,
		}
		if si.BeginPush > 0 {
			si.IsOutputDict = true
		}

		so, _ := si.Smq()
		if so.CodeLen == 0 {
			continue
		}

		dn := smq.GetFileName(v)
		// 写入赛码表
		if si.IsOutputResult {
			writeDict(dn, so.DictBytes)
		}

		if si.IsOutputResult {
			writeResult(tn, dn, so.WordSlice, so.CodeSlice)
		}
		output(so)
	}
}
