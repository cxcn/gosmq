package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	smq "github.com/cxcn/gosmq"

	"github.com/jessevdk/go-flags"
)

func cli() {

	var opts struct {
		Text         string   `short:"t" long:"text" description:"string\t文本路径"`
		Dict         []string `short:"i" long:"input" description:"[]string 码表路径"`
		Single       bool     `short:"s" long:"single" description:"bool\t单字模式"`
		Format       string   `short:"f" long:"format" description:"string\t码表格式 default|jisu,js|duoduo,dd|jidian,jd|bingling,bl"`
		SelectKeys   string   `long:"select" description:"string\t自定义选重键"`
		PushStart    int      `short:"p" long:"push" description:"int\t普通码表起顶码长，码长大于等于此数，首选不会追加空格"`
		Algorithm    string   `short:"a" long:"alg" description:"string\t匹配算法 trie,t|order,o|longest,l"`
		PressSpaceBy string   `short:"k" long:"space" description:"string\t空格按键方式 left|right|both"`
		Details      bool     `short:"d" long:"details" description:"bool\t详细数据"`

		Ver bool `short:"v" long:"version" description:"bool\t查看版本信息"`
	}

	flags.Parse(&opts)
	if opts.Ver {
		printInfo()
		return
	}
	if len(opts.Dict) == 0 {
		return
	}

	var s smq.Smq
	start := time.Now()
	// 不输入文本，直接转换码表
	isEmpty := false
	if opts.Text == "" {
		s = smq.NewFromString("没有输入文本，仅转换码表", "none")
		isEmpty = true
	} else {
		s = smq.NewFromPath("", opts.Text)
	}
	for _, v := range opts.Dict {
		dict := &smq.Dict{
			Single:       opts.Single,
			Format:       opts.Format,
			SelectKeys:   opts.SelectKeys,
			PushStart:    opts.PushStart,
			Algorithm:    opts.Algorithm,
			PressSpaceBy: opts.PressSpaceBy,
			Details:      opts.Details,
		}
		dict.LoadFromPath(v)
		s.Add(dict)
	}

	fmt.Printf("耗时：%v\n", time.Since(start))
	fmt.Printf("比赛开始，一共 %d 个码表\n", len(s.Inputs))
	mid := time.Now()
	res := s.Run()
	fmt.Printf("比赛结束，耗时：%v\n", time.Since(mid))
	fmt.Printf("累计耗时：%v\n", time.Since(start))
	if len(res) == 0 || isEmpty {
		return
	}
	output(res)

	for _, v := range res {
		// 创建文件夹
		os.Mkdir("result", 0666)
		// 输出分词结果
		var buf strings.Builder
		for i := 0; i < len(v.Data.CodeSlice); i++ {
			buf.WriteString(fmt.Sprintf("%s\t%s\n", v.Data.WordSlice[i], string(v.Data.CodeSlice[i])))
		}
		ioutil.WriteFile(fmt.Sprintf("result/%s_%s_分词结果.txt", s.Name, v.Name), []byte(buf.String()), 0666)
		// 输出词条数据
		buf.Reset()
		buf.WriteString("词条\t编码\t顺序\t次数\n")
		type details struct {
			smq.CoC
			word string
		}
		tmp := make([]details, 0, len(v.Data.Details))
		for k, v := range v.Data.Details {
			tmp = append(tmp, details{
				*v,
				k,
			})
		}
		sort.Slice(tmp, func(i, j int) bool {
			return tmp[i].Count > tmp[j].Count
		})
		for _, v := range tmp {
			buf.WriteString(v.word)
			buf.WriteByte('\t')
			buf.WriteString(v.Code)
			buf.WriteByte('\t')
			buf.WriteString(strconv.Itoa(v.Order))
			buf.WriteByte('\t')
			buf.WriteString(strconv.Itoa(v.Count))
			buf.WriteByte('\n')
		}
		ioutil.WriteFile(fmt.Sprintf("result/%s_%s_词条数据.txt", s.Name, v.Name), []byte(buf.String()), 0666)
		// 输出 json 数据
		v.Data.CodeSlice = []string{}
		v.Data.WordSlice = []string{}
		v.Data.Details = make(map[string]*smq.CoC)
		tmp2, _ := json.MarshalIndent(v, "", "  ")
		ioutil.WriteFile(fmt.Sprintf("result/%s_%s.json", s.Name, v.Name), tmp2, 0666)
	}
}
