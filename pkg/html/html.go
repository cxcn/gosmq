package html

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	_ "embed"

	smq "github.com/cxcn/gosmq"
)

//go:embed tmpl.html
var tmpl string

// 赛码结果
type Result struct {
	*smq.SmqOut
	KeyHeatMap [][]template.HTML
	FinHeatMap [10]template.HTML
}

// 供模版使用的数据
type TmplData struct {
	TextName    string
	TextLen     int
	NotHanCount int
	Results     []*Result
}

func NewHTML() *TmplData {
	return new(TmplData)
}

// 添加一个结果
func (d *TmplData) AddResult(so *smq.SmqOut) {

	if !strings.ContainsRune(so.TextName, '《') {
		d.TextName = "《" + so.TextName + "》"
	} else {
		d.TextName = so.TextName
	}

	d.TextLen = so.TextLen
	d.NotHanCount = so.NotHanCount
	s := strings.TrimSuffix(so.DictName, "赛码表")
	so.DictName = strings.TrimSuffix(s, "_")

	tmp := Result{
		SmqOut: so,
	}
	tmp.genKeyHeatMap()
	tmp.genFinHeatMap()
	d.Results = append(d.Results, &tmp)
}

// 输出 html 文件
func (d *TmplData) OutputHTMLFile(fileName string) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	d.OutputHTML(file)
}

func (d *TmplData) OutputHTML(w io.Writer) {
	funcMap := template.FuncMap{"toPer": toPercentage}
	t := template.New("tmpl.html").Funcs(funcMap)
	_, err := t.Parse(tmpl)
	if err != nil {
		panic(err)
	}
	t.Execute(w, d)
}

// float64转换成百分数
func toPercentage(src float64) string {
	return fmt.Sprintf("%.2f%%", src*100.0)
}

// 生成按键热力图
func (res *Result) genKeyHeatMap() {
	src := res.KeyRate
	max := 0.07
	res.KeyHeatMap = make([][]template.HTML, 4)
	line := "1234567890"
	for i, v := range line {
		res.KeyHeatMap[0] = append(res.KeyHeatMap[0], genKeyHeatCode(src[i], max, v))
	}
	line = "QWERTYUIOP"
	for i, v := range line {
		res.KeyHeatMap[1] = append(res.KeyHeatMap[1], genKeyHeatCode(src[i+10], max, v))
	}
	line = "ASDFGHJKL;"
	for i, v := range line {
		res.KeyHeatMap[2] = append(res.KeyHeatMap[2], genKeyHeatCode(src[i+20], max, v))
	}
	res.KeyHeatMap[2] = append(res.KeyHeatMap[2], genKeyHeatCode(src[40], max, rune('\'')))
	line = "ZXCVBNM,./"
	for i, v := range line {
		res.KeyHeatMap[3] = append(res.KeyHeatMap[3], genKeyHeatCode(src[i+30], max, v))
	}
}

// 按键颜色代码片段
func genKeyHeatCode(freq, max float64, key rune) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<td class="key" style="background-color: rgba(255,0,0,%.4f);">%s <div class="heatMapRate">%.2f</div></td>`,
		freq/max*0.6, string(key), freq*100))
}

// 生成手指热力图
func (res *Result) genFinHeatMap() {
	src := res.FinRate
	max := 0.25
	fins := []string{"左小", "左无", "左中", "左食", "大拇指", "右食", "右中", "右无", "右小"}
	for i := 0; i < 9; i++ {
		res.FinHeatMap[i] = genFinHeatCode(src[i+1], max, i, fins[i])
	}
}

// 手指颜色代码片段
func genFinHeatCode(freq, max float64, id int, fin string) template.HTML {
	if id == 4 {
		return template.HTML(fmt.Sprintf(
			`<td class="key fin" colspan="2" style="background-color: rgba(0,0,255,%.4f);">%s <div class="heatMapRate">%.2f</div></td>`,
			freq/max*0.6, fin, freq*100))
	}
	return template.HTML(fmt.Sprintf(
		`<td class="key fin" style="background-color: rgba(0,0,255,%.4f);">%s <div class="heatMapRate">%.2f</div></td>`,
		freq/max*0.6, fin, freq*100))
}
