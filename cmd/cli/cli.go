package main

import (
	smq "github.com/cxcn/gosmq"

	"github.com/jessevdk/go-flags"
)

func cli() {

	var opts struct {
		Text         string `short:"t" long:"text" description:"string\t文本路径"`
		Dict         string `short:"i" long:"input" description:"string\t码表路径"`
		Single       bool   `short:"s" long:"single" description:"bool\t单字模式"`
		Format       string `short:"f" long:"format" description:"string\t码表格式 default|jisu|duoduo|jidian"`
		SelectKeys   string `long:"select" description:"string\t自定义选重键"`
		PushStart    int    `short:"p" long:"push" description:"int\t普通码表起顶码长，码长大于等于此数，首选不会追加空格"`
		Algorithm    string `short:"a" long:"alg" description:"string\t匹配算法 trie|order|longest"`
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
	s.Add(dict)
	res := s.Run()
	for _, v := range res {
		output(v)
	}
}
