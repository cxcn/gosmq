package main

import (
	"fmt"

	smq "github.com/cxcn/gosmq"
	"github.com/jedib0t/go-pretty/v6/table"
)

func output(so *smq.SmqOut) {

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

	t3 := table.NewWriter()
	t3.AppendHeader(table.Row{"左手", "右手", "左右", "右左", "左左", "右右"})
	t3.AppendRow([]interface{}{
		fmt.Sprintf("%.3f%%", 100*so.LeftHand),
		fmt.Sprintf("%.3f%%", 100*so.RightHand),
		fmt.Sprintf("%.3f%%", 100*so.HandRate[0]),
		fmt.Sprintf("%.3f%%", 100*so.HandRate[1]),
		fmt.Sprintf("%.3f%%", 100*so.HandRate[2]),
		fmt.Sprintf("%.3f%%", 100*so.HandRate[3]),
	})
	t3.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t3.Render())
	out += "\n"

	t4 := table.NewWriter()
	t4.AppendHeader(table.Row{"当量", "异手", "同指", "同键", "小跨排", "大跨排", "异指", "小指干扰", "错手"})
	t4.AppendRow([]interface{}{
		fmt.Sprintf("%.4f", so.Eq),
		fmt.Sprintf("%.3f%%", 100*so.DiffHandRate),
		fmt.Sprintf("%.3f%%", 100*so.SameFinRate),
		fmt.Sprintf("%.3f%%", 100*so.Sk),
		fmt.Sprintf("%.3f%%", 100*so.Xkp),
		fmt.Sprintf("%.3f%%", 100*so.Dkp),
		fmt.Sprintf("%.3f%%", 100*so.DiffFinRate),
		fmt.Sprintf("%.3f%%", 100*so.Lfd),
		fmt.Sprintf("%.3f%%", 100*so.Cs),
	})
	t4.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t4.Render())
	out += "\n"

	t5 := table.NewWriter()
	t5.AppendHeader(table.Row{"小指", "无名指", "中指", "食指", "大拇指", "食指", "中指", "无名指", "小指", "其他"})
	t5_row_2 := []interface{}{}
	for i := 1; i < 10; i++ {
		t5_row_2 = append(t5_row_2, fmt.Sprintf("%.3f%%", 100*so.FinRate[i]))
	}
	t5_row_2 = append(t5_row_2, fmt.Sprintf("%.3f%%", 100*so.FinRate[0]))
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
		t6_1_row = append(t6_1_row, fmt.Sprintf("%.2f%%", 100*so.KeyRate[keys[i]]))
		t6_2_row = append(t6_2_row, fmt.Sprintf("%.2f%%", 100*so.KeyRate[keys[i+10]]))
		t6_3_row = append(t6_3_row, fmt.Sprintf("%.2f%%", 100*so.KeyRate[keys[i+20]]))
		t6_4_row = append(t6_4_row, fmt.Sprintf("%.2f%%", 100*so.KeyRate[keys[i+30]]))
	}
	t6_3_header = append(t6_3_header, "'")
	t6_3_row = append(t6_3_row, fmt.Sprintf("%.2f%%", 100*so.KeyRate[keys[30]]))

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
}
