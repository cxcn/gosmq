package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	smq "github.com/cxcn/gosmq"
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
		isS  bool   // 是否只跑单字
		isW  bool   // 是否输出赛码表
		fpt  string // file path text
		aF   bool   // 是否关闭手感统计
		aS   bool   // 空格是否互击
		csk  string // custom select keys
		fpo  string // output file path
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
	flag.BoolVar(&help, "h", false, "显示帮助")
	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Println("请在命令行中运行此程序\n按Enter键退出...")
		fmt.Scanln(&fpo)
		return
	}
	if help || len(fpm) == 0 {
		fmt.Print("saimaqi version: 0.7, github: https://github.com/cxcn/saimaqi\n\n")
		fmt.Print("Usage: saimaqi.exe [-i mb] [-n int] [-d] [-w] [-t text] [-f] [-s] [-k string] [-o output]\n\n")
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

	out := ""
	out += fmt.Sprintln("----------------------")

	out += fmt.Sprintf("非汉字：%s\n", so.NotHan)
	out += fmt.Sprintf("缺字：%s\n", so.Lack)
	out += "\n"

	out += fmt.Sprintf("码长统计：%v\n", so.CodeStat)
	out += fmt.Sprintf("词长统计：%v\n", so.WordStat)
	out += fmt.Sprintf("选重统计：%v\n", so.RepeatStat)
	out += "\n"

	t1 := table.NewWriter()
	t1.AppendHeader(table.Row{"文本字数", "总键数", "码长", "十击速度", "非汉字数", "缺字数"})
	t1.AppendRow([]interface{}{
		so.TextLen, so.CodeLen,
		fmt.Sprintf("%.4f", so.CodeAvg),
		fmt.Sprintf("%.2f", 600/so.CodeAvg),
		so.NotHanCount, so.LackCount,
	})
	t1.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t1.Render())
	out += "\n"

	t2 := table.NewWriter()
	t2.AppendHeader(table.Row{"打词", "选重", "打词字数", "选重字数"})
	t2.AppendRow([]interface{}{
		so.WordCount,
		so.RepeatCount,
		so.WordLen,
		so.RepeatLen,
	})
	t2.AppendRow([]interface{}{
		fmt.Sprintf("%.3f%%", 100*so.WordRate),
		fmt.Sprintf("%.3f%%", 100*so.RepeatRate),
		fmt.Sprintf("%.3f%%", 100*so.WordLenRate),
		fmt.Sprintf("%.3f%%", 100*so.RepeatLenRate),
	})
	t2.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t2.Render())
	out += "\n"

	if aF {
		out += fmt.Sprintln("指法统计已关闭...")
		out += fmt.Sprintln("----------------------")
		fmt.Print(out)
		return
	}
	feel := smq.NewFeel(so.Code, isS)
	t3 := table.NewWriter()
	t3.AppendHeader(table.Row{"左右", "右左", "左左", "右右"})
	t3.AppendRow([]interface{}{
		feel.HandCount[0],
		feel.HandCount[1],
		feel.HandCount[2],
		feel.HandCount[3],
	})
	t3.AppendRow([]interface{}{
		fmt.Sprintf("%.3f%%", 100*feel.HandRate[0]),
		fmt.Sprintf("%.3f%%", 100*feel.HandRate[1]),
		fmt.Sprintf("%.3f%%", 100*feel.HandRate[2]),
		fmt.Sprintf("%.3f%%", 100*feel.HandRate[3]),
	})
	t3.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t3.Render())
	out += "\n"

	t4 := table.NewWriter()
	t4.AppendHeader(table.Row{"当量", "左手", "右手", "异手", "同指", "大跨排", "小跨排", "异指", "小指干扰", "错手"})
	t4.AppendRow([]interface{}{
		fmt.Sprintf("%.4f", feel.Eq),
		fmt.Sprintf("%.3f%%", 100*feel.LeftHand),
		fmt.Sprintf("%.3f%%", 100*feel.RightHand),
		fmt.Sprintf("%.3f%%", 100*feel.DiffHandRate),
		fmt.Sprintf("%.3f%%", 100*feel.SameFinRate),
		fmt.Sprintf("%.3f%%", 100*feel.Dkp),
		fmt.Sprintf("%.3f%%", 100*feel.Xkp),
		fmt.Sprintf("%.3f%%", 100*feel.DiffFinRate),
		fmt.Sprintf("%.3f%%", 100*feel.Lfd),
		fmt.Sprintf("%.3f%%", 100*feel.Cs),
	})
	t4.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t4.Render())
	out += "\n"

	t5 := table.NewWriter()
	t5.AppendHeader(table.Row{"小指", "无名指", "中指", "食指", "大拇指", "食指", "中指", "无名指", "小指", "其他"})
	t5_row_1 := []interface{}{}
	t5_row_2 := []interface{}{}
	for i := 1; i < 10; i++ {
		t5_row_1 = append(t5_row_1, feel.FinCount[i])
		t5_row_2 = append(t5_row_2, fmt.Sprintf("%.3f%%", 100*feel.FinRate[i]))
	}
	t5_row_1 = append(t5_row_1, feel.FinCount[0])
	t5_row_2 = append(t5_row_2, fmt.Sprintf("%.3f%%", 100*feel.FinRate[0]))
	t5.AppendRow(t5_row_1)
	// t.AppendSeparator()
	t5.AppendRow(t5_row_2)
	t5.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t5.Render())
	out += "\n"

	keys := "1234567890qwertyuiopasdfghjkl;zxcvbnm,./"
	t6_1_header := []interface{}{}
	t6_2_header := []interface{}{}
	t6_3_header := []interface{}{}
	t6_4_header := []interface{}{}
	t6_1_row := []interface{}{}
	t6_2_row := []interface{}{}
	t6_3_row := []interface{}{}
	t6_4_row := []interface{}{}
	for i := 0; i < 10; i++ {
		t6_1_header = append(t6_1_header, string(keys[i]))
		t6_2_header = append(t6_2_header, string(keys[i+10]))
		t6_3_header = append(t6_3_header, string(keys[i+20]))
		t6_4_header = append(t6_4_header, string(keys[i+30]))
		t6_1_row = append(t6_1_row, fmt.Sprintf("%.2f%%", 100*feel.KeyRate[keys[i]]))
		t6_2_row = append(t6_2_row, fmt.Sprintf("%.2f%%", 100*feel.KeyRate[keys[i+10]]))
		t6_3_row = append(t6_3_row, fmt.Sprintf("%.2f%%", 100*feel.KeyRate[keys[i+20]]))
		t6_4_row = append(t6_4_row, fmt.Sprintf("%.2f%%", 100*feel.KeyRate[keys[i+30]]))
	}
	t6_3_header = append(t6_3_header, "'")
	t6_3_row = append(t6_3_row, fmt.Sprintf("%.2f%%", 100*feel.KeyRate[keys[30]]))

	t6_1 := table.NewWriter()
	t6_1.AppendHeader(t6_1_header)
	t6_1.AppendRow(t6_1_row)
	t6_1.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t6_1.Render())

	t6_2 := table.NewWriter()
	t6_2.AppendHeader(t6_2_header)
	t6_2.AppendRow(t6_2_row)
	t6_2.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t6_2.Render())

	t6_3 := table.NewWriter()
	t6_3.AppendHeader(t6_3_header)
	t6_3.AppendRow(t6_3_row)
	t6_3.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t6_3.Render())

	t6_4 := table.NewWriter()
	t6_4.AppendHeader(t6_4_header)
	t6_4.AppendRow(t6_4_row)
	t6_4.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t6_4.Render())

	out += fmt.Sprintln("----------------------")
	fmt.Print(out)
	// time.Sleep(5 * time.Second)
}
