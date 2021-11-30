package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	smq "github.com/cxcn/gosmq"
	"github.com/cxcn/gosmq/pkg/html"

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
		fmt.Printf("smq-cli version 0.13 %s/%s\n\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println("repo address: https://github.com/cxcn/gosmq/")
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
	h := html.NewHTML(tn)
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
			IsOutputDict:   true,
			IsOutputResult: opt.IsO,

			BeginPush:       opt.Ding,
			SelectKeys:      opt.Csk,
			IsSingleOnly:    opt.IsS,
			IsSpaceDiffHand: opt.AS,
		}

		so, _ := si.Smq()
		if so.CodeLen == 0 {
			continue
		}

		dn := smq.GetFileName(v)
		// 写入赛码表
		if si.IsOutputResult {
			_ = os.Mkdir("dict", 0666)
			err := ioutil.WriteFile(".\\dict\\"+dn+"_赛码表.txt", so.DictBytes, 0666)
			if err != nil {
				fmt.Println("Error! 赛码表写入错误:", err)
			} else {
				fmt.Println("Success! 成功写入赛码表:", ".\\dict\\"+dn+"_赛码表.txt")
			}
		}

		if si.IsOutputResult {
			var buf bytes.Buffer
			for i, v := range so.WordSlice {
				buf.WriteString(string(v))
				buf.WriteByte('\t')
				buf.WriteString(so.CodeSlice[i])
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
		h.AddResult(so, dn)
		output(so)
	}
	h.OutputHTMLFile("result.html")

}
