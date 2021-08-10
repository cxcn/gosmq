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
		ding int    // 普通码表起顶码长
		isD  bool   // 是否只跑单字
		isW  bool   // 是否输出赛码表
		fpt  string // file path text
		isnF bool   // 是否关闭手感统计
		isS  bool   // 空格是否互击
		csk  string // custom select keys
		fpo  string // output file path
		help bool
	)

	flag.StringVar(&fpm, "i", "", "码表路径，可以是rime格式码表 或 极速跟打器赛码表")
	flag.IntVar(&ding, "n", 0, "普通码表起顶码长，码长大于等于此数，首选不会追加空格")
	flag.BoolVar(&isD, "d", false, "是否只跑单字")
	flag.BoolVar(&isW, "w", false, "是否输出赛码表(保存在.\\smb\\文件夹下)")
	flag.StringVar(&fpt, "t", "", "文本路径，utf8编码格式文本，会自动去除空白符")
	flag.BoolVar(&isnF, "f", false, "是否关闭手感统计")
	flag.BoolVar(&isS, "s", false, "空格是否互击")
	flag.StringVar(&csk, "k", ";'", "自定义选重键(2重开始)")
	flag.StringVar(&fpo, "o", "", "输出路径")
	flag.BoolVar(&help, "h", false, "显示帮助")
	flag.Parse()

	if help {
		fmt.Print("saimaqi version: 0.5\n\n")
		fmt.Print("Usage: saimaqi.exe [-i mb] [-n int] [-d] [-w] [-t text] [-f] [-s] [-k string] [-o output]\n\n")
		flag.PrintDefaults()
		return
	}
	fmt.Println()
	if isD {
		fmt.Println("只跑单字...")
	}

	dict := NewDict(fpm, ding, isW, isD)
	if len(dict.children) == 0 || len(fpt) == 0 {
		return
	}
	smq := NewSmq(dict, fpt, csk)
	if smq.textLen == 0 {
		fmt.Println("文本为空...")
		return
	}

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

	if isnF {
		out += fmt.Sprintln("指法统计已关闭...")
		out += fmt.Sprintln("----------------------")
		fmt.Print(out)
		return
	}
	feel := NewFeel(smq.code, isS)
	t3 := table.NewWriter()
	t3.AppendHeader(table.Row{"左右", "右左", "左左", "右右"})
	t3.AppendRow([]interface{}{
		feel.handCount[0],
		feel.handCount[1],
		feel.handCount[2],
		feel.handCount[3],
	})
	t3.AppendRow([]interface{}{
		fmt.Sprintf("%.3f%%", 100*feel.handRate[0]),
		fmt.Sprintf("%.3f%%", 100*feel.handRate[1]),
		fmt.Sprintf("%.3f%%", 100*feel.handRate[2]),
		fmt.Sprintf("%.3f%%", 100*feel.handRate[3]),
	})
	t3.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t3.Render())
	out += "\n"

	t4 := table.NewWriter()
	t4.AppendHeader(table.Row{"当量", "左手", "右手", "异手", "同指", "大跨排", "小跨排", "异指", "小指干扰", "错手"})
	t4.AppendRow([]interface{}{
		fmt.Sprintf("%.4f", feel.dl),
		fmt.Sprintf("%.3f%%", 100*feel.leftHand),
		fmt.Sprintf("%.3f%%", 100*feel.rightHand),
		fmt.Sprintf("%.3f%%", 100*feel.diffHandRate),
		fmt.Sprintf("%.3f%%", 100*feel.sameFinRate),
		fmt.Sprintf("%.3f%%", 100*feel.dkp),
		fmt.Sprintf("%.3f%%", 100*feel.xkp),
		fmt.Sprintf("%.3f%%", 100*feel.diffFinRate),
		fmt.Sprintf("%.3f%%", 100*feel.xzgr),
		fmt.Sprintf("%.3f%%", 100*feel.cs),
	})
	t4.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t4.Render())
	out += "\n"

	t5 := table.NewWriter()
	t5.AppendHeader(table.Row{"小指", "无名指", "中指", "食指", "大拇指", "食指", "中指", "无名指", "小指", "其他"})
	t5_row_1 := []interface{}{}
	t5_row_2 := []interface{}{}
	for i := 1; i < 10; i++ {
		t5_row_1 = append(t5_row_1, feel.finCount[i])
		t5_row_2 = append(t5_row_2, fmt.Sprintf("%.3f%%", 100*feel.finRate[i]))
	}
	t5_row_1 = append(t5_row_1, feel.finCount[0])
	t5_row_2 = append(t5_row_2, fmt.Sprintf("%.3f%%", 100*feel.finRate[0]))
	t5.AppendRow(t5_row_1)
	// t.AppendSeparator()
	t5.AppendRow(t5_row_2)
	t5.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t5.Render())

	out += fmt.Sprintln("----------------------")
	fmt.Print(out)
	// time.Sleep(5 * time.Second)
}
