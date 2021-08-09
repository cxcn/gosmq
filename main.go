package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
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
		ding int    // 是否输出赛码表
		isW  bool   // file path write smb
		fpt  string // file path text
		isS  bool   // 空格是否互击
		csk  string // custom select keys
		fpo  string // output file path
		help bool
	)

	flag.StringVar(&fpm, "i", "", "码表路径，可以是rime格式码表 或 极速跟打器赛码表")
	flag.IntVar(&ding, "d", 0, "普通码表起顶码长，码长大于等于此数，首选不会追加空格")
	flag.BoolVar(&isW, "w", false, "是否输出赛码表(保存在.\\smb\\文件夹下)")
	flag.StringVar(&fpt, "t", "", "文本路径，utf8编码格式文本，会自动去除空白符")
	flag.BoolVar(&isS, "s", false, "空格是否互击")
	flag.StringVar(&csk, "k", ";'", "自定义选重键(2重开始)")
	flag.StringVar(&fpo, "o", "", "输出路径")
	flag.BoolVar(&help, "h", false, "显示帮助")
	flag.Parse()

	if help {
		fmt.Print("saimaqi version: 0.4\n\n")
		fmt.Print("Usage: saimaqi.exe [-i mb] [-d int] [-w] [-t text] [-s] [-k string] [-o output]\n\n")
		flag.PrintDefaults()
		return
	}

	dict := NewDict(fpm, ding, isW)
	if len(dict.children) == 0 {
		return
	}
	smq := NewSmq(dict, fpt, csk)
	if smq.textLen == 0 {
		return
	}
	fin := NewFin(smq.code, isS)

	if fpo != "" {
		// fmt.Println(fpo)
		_ = ioutil.WriteFile(fpo, []byte(smq.codeSep), 0777)
	}
	out := ""
	out += fmt.Sprintln("----------------------")

	t1 := table.NewWriter()
	t1.AppendHeader(table.Row{"文本字数", "总键数", "码长", "非汉字数", "缺字数"})
	t1.AppendRow([]interface{}{
		smq.textLen, smq.codeLen,
		fmt.Sprintf("%.4f", smq.codeAvg),
		smq.notHanCount, smq.lackCount,
	})
	t1.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t1.Render())
	out += "\n"

	out += fmt.Sprintf("非汉字：%s\n", smq.notHan)
	out += fmt.Sprintf("缺字：%s\n", smq.lack)
	out += "\n"

	t2 := table.NewWriter()
	t2.AppendHeader(table.Row{"打词", "选重", "打词字数", "选重字数"})
	t2.AppendRow([]interface{}{
		smq.wordCount,
		smq.repeatCount,
		smq.wordLen,
		smq.repeatLen,
	})
	t2.AppendRow([]interface{}{
		fmt.Sprintf("%.3f%%", 100*smq.wordRate),
		fmt.Sprintf("%.3f%%", 100*smq.repeatRate),
		fmt.Sprintf("%.3f%%", 100*smq.wordLenRate),
		fmt.Sprintf("%.3f%%", 100*smq.repeatLenRate),
	})
	t2.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t2.Render())
	out += "\n"

	out += fmt.Sprintf("码长统计：%v\n", smq.codeStat)
	out += fmt.Sprintf("词长统计：%v\n", smq.wordStat)
	out += "\n"

	t3 := table.NewWriter()
	t3.AppendHeader(table.Row{"左右", "右左", "左左", "右右"})
	t3.AppendRow([]interface{}{
		fin.posCount[0],
		fin.posCount[1],
		fin.posCount[2],
		fin.posCount[3],
	})
	t3.AppendRow([]interface{}{
		fmt.Sprintf("%.3f%%", 100*fin.posRate[0]),
		fmt.Sprintf("%.3f%%", 100*fin.posRate[1]),
		fmt.Sprintf("%.3f%%", 100*fin.posRate[2]),
		fmt.Sprintf("%.3f%%", 100*fin.posRate[3]),
	})
	t3.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t3.Render())
	out += "\n"

	t4 := table.NewWriter()
	t4.AppendHeader(table.Row{"异手", "同指", "异指"})
	t4.AppendRow([]interface{}{
		fmt.Sprintf("%.3f%%", 100*fin.diffHandRate),
		fmt.Sprintf("%.3f%%", 100*fin.sameFinRate),
		fmt.Sprintf("%.3f%%", 100*fin.diffFinRate),
	})
	t4.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t4.Render())
	out += "\n"

	t5 := table.NewWriter()
	t5.AppendHeader(table.Row{"小指", "无名指", "中指", "食指", "大拇指", "食指", "中指", "无名指", "小指", "其他"})
	t5.AppendRow([]interface{}{
		fin.keyCount[1],
		fin.keyCount[2],
		fin.keyCount[3],
		fin.keyCount[4],
		fin.keyCount[5],
		fin.keyCount[6],
		fin.keyCount[7],
		fin.keyCount[8],
		fin.keyCount[9],
		fin.keyCount[0],
	})
	// t.AppendSeparator()
	t5.AppendRow([]interface{}{
		fmt.Sprintf("%.3f%%", 100*fin.keyRate[1]),
		fmt.Sprintf("%.3f%%", 100*fin.keyRate[2]),
		fmt.Sprintf("%.3f%%", 100*fin.keyRate[3]),
		fmt.Sprintf("%.3f%%", 100*fin.keyRate[4]),
		fmt.Sprintf("%.3f%%", 100*fin.keyRate[5]),
		fmt.Sprintf("%.3f%%", 100*fin.keyRate[6]),
		fmt.Sprintf("%.3f%%", 100*fin.keyRate[7]),
		fmt.Sprintf("%.3f%%", 100*fin.keyRate[8]),
		fmt.Sprintf("%.3f%%", 100*fin.keyRate[9]),
		fmt.Sprintf("%.3f%%", 100*fin.keyRate[0]),
	})
	t5.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t5.Render())

	out += fmt.Sprintln("----------------------")
	fmt.Print(out)
	// time.Sleep(5 * time.Second)
}
