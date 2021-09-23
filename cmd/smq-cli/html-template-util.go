/*
本文件定义了预处理类。赛码器输出到单个静态网页时，需要清洗数据，包括：
收集文件名，
多个方案里，找出最优的数据，
数据展示时，要修辞
*/

package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	smq "github.com/cxcn/gosmq"
)

//go:embed smq-out-template.html
var smqOutTemplate string

// 该对象会自动初始化
var templateObj *template.Template

func init() {
	templateObj = template.Must(
		template.New("smq-out").Funcs(
			template.FuncMap{
				"calcFingerHeatColor":   calcFingerHeatColor,
				"calcKeyHeatColor":      calcKeyHeatColor,
				"calcHandHeatColor":     calcHandHeatColor,
				"removeTxt":             removeTxt,
				"getNowDate":            getNowDate,
				"convertToPercentage":   convertToPercentage,
				"convertToNoPercentage": convertToNoPercentage,
			}).Parse(smqOutTemplate))
}

// 每个方案的结果信息
type SchemaInfo struct {
	SchemaName string
	SmqOut     *smq.SmqOut // 该方案的效率信息
}

type HTMLOutputInfo struct {
	TextFileName string       // 赛文的文件名
	Schemas      []SchemaInfo // 多个方案的结果
}

func NewHTMLOutputInfo(textFileName string) HTMLOutputInfo {
	return HTMLOutputInfo{TextFileName: textFileName}
}

func (hoi *HTMLOutputInfo) AddSchema(schemaName string, smqOut *smq.SmqOut) *HTMLOutputInfo {
	hoi.Schemas = append(hoi.Schemas, SchemaInfo{schemaName, smqOut})
	return hoi
}

func (hoi *HTMLOutputInfo) OutputHTMLFile(fileName string) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = templateObj.Execute(file, hoi)
	if err != nil {
		panic(err)
	}
}

//模板自定义通道方法：传入某键位的使用率，返回对应的按键的热度颜色
func calcKeyHeatColor(keyRate float64) template.CSS {
	const heatestRate float64 = .13
	alpha := 1.0
	if keyRate < heatestRate {
		alpha = keyRate / heatestRate
	}
	return template.CSS(fmt.Sprintf("rgba(255,0,0,%.4f)", alpha))
}

//模板自定义通道方法：传入某键位的使用率，返回的用指频率的热度颜色
func calcFingerHeatColor(finRate float64) template.CSS {
	return template.CSS(fmt.Sprintf("rgba(0,0,255,%.4f)", finRate))
}

//模板自定义通道方法：传入某键位的使用率，返回的用手的频率的热度颜色
func calcHandHeatColor(finRate float64) template.CSS {
	if finRate > 0.5 {
		return template.CSS(fmt.Sprintf("rgba(255,155,0,%.4f)", (finRate-0.5)*2))
	} else {
		return template.CSS(fmt.Sprintf("rgba(0,200,255,%.4f)", (0.5-finRate)*2))
	}
}

//模板自定义通道方法：去除结尾的“.txt”
func removeTxt(rawFileName string) string {
	return strings.TrimSuffix(rawFileName, ".txt")
}

//模板自定义通道方法：返回当前的日期、时间
func getNowDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//模板自定义通道方法：float64转换成百分数
func convertToPercentage(src float64) string {
	return fmt.Sprintf("%.2f%%", src*100.0)
}

//模板自定义通道方法：float64转换成百分数不带百分号
func convertToNoPercentage(src float64) string {
	return fmt.Sprintf("%.2f", src*100.0)
}
