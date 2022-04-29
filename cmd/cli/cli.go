package main

import (
	"fmt"
	"time"

	smq "github.com/cxcn/gosmq"

	"github.com/jessevdk/go-flags"
)

func cli() {

	var opts struct {
		Text         string `short:"t" long:"text" description:"string\t文本路径"`
		Dict         string `short:"i" long:"input" description:"string\t码表路径"`
		Single       bool   `short:"s" long:"single" description:"bool\t单字模式"`
		Format       string `short:"f" long:"format" description:"string\t码表格式 default|jisu,js|duoduo,dd|jidian,jd|bingling,bl"`
		SelectKeys   string `long:"select" description:"string\t自定义选重键"`
		PushStart    int    `short:"p" long:"push" description:"int\t普通码表起顶码长，码长大于等于此数，首选不会追加空格"`
		Algorithm    string `short:"a" long:"alg" description:"string\t匹配算法 trie,t|order,o|longest,l"`
		PressSpaceBy string `short:"k" long:"space" description:"string\t空格按键方式 left|right|both"`

		Ver bool `short:"v" long:"version" description:"bool\t查看版本信息"`
	}

	flags.Parse(&opts)
	if opts.Ver {
		printInfo()
		return
	}
	if opts.Dict == "" {
		return
	}

	dict := &smq.Dict{
		Single:       opts.Single,
		Format:       opts.Format,
		SelectKeys:   opts.SelectKeys,
		PushStart:    opts.PushStart,
		Algorithm:    opts.Algorithm,
		PressSpaceBy: opts.PressSpaceBy,
	}
	dict.LoadFromPath(opts.Dict)
	s := smq.NewFromPath("", opts.Text)

	start := time.Now()
	s.Add(dict)
	fmt.Printf("耗时：%v\n", time.Since(start))
	fmt.Printf("比赛开始，一共 %d 个码表\n", len(s.Inputs))
	mid := time.Now()
	res := s.Run()
	fmt.Printf("比赛结束，耗时：%v\n", time.Since(mid))
	fmt.Printf("累计耗时：%v\n", time.Since(start))
	for _, v := range res {
		output(v)
		// var buf strings.Builder
		// for i := 0; i < len(v.Data.CodeSlice); i++ {
		// 	buf.WriteString(fmt.Sprintf("%s\t%s\n", v.Data.CodeSlice[i], string(v.Data.WordSlice[i])))
		// }
		// ioutil.WriteFile("res2.txt", []byte(buf.String()), 0666)
	}
}
