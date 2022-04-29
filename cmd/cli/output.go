package main

import (
	"fmt"

	smq "github.com/cxcn/gosmq"
	"github.com/jedib0t/go-pretty/table"
)

func output(res *smq.Result) {
	// r, _ := json.MarshalIndent(res, "", "  ")
	// fmt.Printf(string(r))

	out := ""
	out += fmt.Sprintln("----------------------")

	out += fmt.Sprintf("文本字数：%d\n", res.Basic.TextLen)
	out += fmt.Sprintf("词条数：%d\n", res.Basic.DictLen)
	out += fmt.Sprintf("非汉字：%s\n", res.Basic.NotHan)
	if res.Basic.Lack != "" {
		out += fmt.Sprintf("缺字：%s\n", res.Basic.Lack)
	}
	out += "\n"

	out += fmt.Sprintf("码长分布：%v\n", res.CodeLen.Dist)
	out += fmt.Sprintf("词长分布：%v\n", res.Words.Dist)
	out += fmt.Sprintf("选重分布：%v\n", res.Collision.Dist)
	out += "\n"

	t1 := table.NewWriter()
	t1.AppendHeader(table.Row{"总键数", "码长", "十击速度", "非汉字数", "非汉字计数", "缺字数", "缺字计数"})
	t1.AppendRow([]interface{}{
		res.CodeLen.Total,
		fmt.Sprintf("%.4f", res.CodeLen.PerChar),
		fmt.Sprintf("%.2f", 600/res.CodeLen.PerChar),
		res.Basic.NotHans, res.Basic.NotHanCount,
		res.Basic.Lacks, res.Basic.LackCount,
	})
	t1.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t1.Render())
	out += "\n"

	t2 := table.NewWriter()
	t2.AppendHeader(table.Row{"打词", "打词字数", "选重", "选重字数"})
	t2.AppendRow([]interface{}{
		res.Words.Commits.Count,
		res.Words.Chars.Count,
		res.Collision.Commits.Count,
		res.Collision.Chars.Count,
	})
	t2.AppendRow([]interface{}{
		fmt.Sprintf("%.3f%%", 100*res.Words.Commits.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Words.Chars.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Collision.Commits.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Collision.Chars.Rate),
	})
	t2.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t2.Render())
	out += "\n"

	t3 := table.NewWriter()
	t3.AppendHeader(table.Row{"左手", "右手", "左右", "右左", "左左", "右右"})
	t3.AppendRow([]interface{}{
		fmt.Sprintf("%.3f%%", 100*res.Hands.Left.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Hands.Right.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Hands.LL.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Hands.LR.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Hands.RL.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Hands.RR.Rate),
	})
	t3.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t3.Render())
	out += "\n"

	t4 := table.NewWriter()
	t4.AppendHeader(table.Row{"当量", "异手", "同指", "同键", "小跨排", "大跨排", "异指", "小指干扰", "错手"})
	t4.AppendRow([]interface{}{
		fmt.Sprintf("%.4f", res.Combs.Equivalent),
		fmt.Sprintf("%.3f%%", 100*res.Hands.Diff.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Fingers.Same.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Combs.DoubleHit.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Combs.SingleSpan.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Combs.MultiSpan.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Fingers.Diff.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Combs.LittleFingersDisturb.Rate),
		fmt.Sprintf("%.3f%%", 100*res.Combs.LongFingersDisturb.Rate),
	})
	t4.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t4.Render())
	out += "\n"

	t5 := table.NewWriter()
	t5.AppendHeader(table.Row{"小指", "无名指", "中指", "食指", "左大拇指", "右大拇指", "食指", "中指", "无名指", "小指"})
	t5_row_2 := []interface{}{}
	for i := 1; i < 10; i++ {
		t5_row_2 = append(t5_row_2, fmt.Sprintf("%.3f%%", 100*res.Fingers.Dist[i].Rate))
	}
	t5_row_2 = append(t5_row_2, fmt.Sprintf("%.3f%%", 100*res.Fingers.Dist[0].Rate))
	t5.AppendRow(t5_row_2)
	t5.SetStyle(table.StyleColoredBright)
	out += fmt.Sprintln(t5.Render())
	out += "\n"

	t6_1_header := []interface{}{}
	t6_2_header := []interface{}{}
	t6_3_header := []interface{}{}
	t6_4_header := []interface{}{}
	t6_1_row := []interface{}{}
	t6_2_row := []interface{}{}
	t6_3_row := []interface{}{}
	t6_4_row := []interface{}{}
	// keys := "1234567890qwertyuiopasdfghjkl;zxcvbnm,./"
	keys1 := "1234567890"
	keys2 := "qwertyuiop"
	keys3 := "asdfghjkl;" // 单引号单独加
	keys4 := "zxcvbnm,./"
	for i := 0; i < 10; i++ {
		a, b, c, d := string(keys1[i]), string(keys2[i]), string(keys3[i]), string(keys4[i])
		t6_1_header = append(t6_1_header, a)
		t6_2_header = append(t6_2_header, b)
		t6_3_header = append(t6_3_header, c)
		t6_4_header = append(t6_4_header, d)
		t6_1_row = append(t6_1_row, fmt.Sprintf("%.2f%%", 100*res.Keys[a].Rate))
		t6_2_row = append(t6_2_row, fmt.Sprintf("%.2f%%", 100*res.Keys[b].Rate))
		t6_3_row = append(t6_3_row, fmt.Sprintf("%.2f%%", 100*res.Keys[c].Rate))
		t6_4_row = append(t6_4_row, fmt.Sprintf("%.2f%%", 100*res.Keys[d].Rate))
	}
	t6_3_header = append(t6_3_header, "'")
	t6_3_row = append(t6_3_row, fmt.Sprintf("%.2f%%", 100*res.Keys["'"].Rate))

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
