package main

import (
	"bytes"
	"strings"
	"unicode"
)

type void struct{}

type element struct {
	code  string
	times int
}

type Result struct {
	lenText     int    //文本字数
	notHan      string //非汉字
	countNotHan int    //非汉字数
	lack        string //缺字
	countLack   int    //缺字数

	codeSep string //空格间隔的全部编码
	freq    map[string]element
	choose  map[string]void //选重
	// dict    map[string]string //新的词典
	// freq    map[string]int    //字词频率统计

	//以下可由上面计算得
	code       string   //全部编码
	codeSlice  []string //编码切片
	lenCode    int      //总键数
	avlCode    float64  //码长
	countSpace int      //空格数

	countWord int     //打词数
	lenWord   int     //打词字数
	rateWord  float64 //打词率（上屏）
	rLenWord  float64 //打词率（字数）

	countChoose int     //选重数
	lenChoose   int     //选重字数
	rateChoose  float64 //选重率（上屏）
	rLenChoose  float64 //选重率（字数）

	stat map[int]int
}

var Res Result

var buf bytes.Buffer

func part(w, c string) {
	// Res.code = append(Res.code, c)
	buf.WriteString(c + " ")
	tmp := Res.freq[w]
	if tmp.code == "" {
		tmp.code = c
	}
	tmp.times++
	Res.freq[w] = tmp

	// fmt.Println(tmp)
	// Res.dict[w] = c
	// Res.freq[w]++
	if 48 <= c[len(c)-1] && c[len(c)-1] <= 57 {
		Res.choose[w] = void{}
		// Res.choose = append(Res.choose, w)
	}
}

func cacl(fpm, fpt string) Result {
	text := readText(fpt)
	dict := readMB(fpm)
	Res.lenText = len(text)
	Res.freq = make(map[string]element)
	Res.choose = make(map[string]void)
	// Res.dict = make(map[string]string)
	// Res.freq = make(map[string]int)
	max := 0 // 最大词长
	for k := range dict {
		if k > max {
			max = k
		}
	}

	// start := time.Now()
	// defer func() {
	// 	cost := time.Since(start)
	// 	fmt.Println("calc cost time = ", cost)
	// }()

	p := 0 // point
	for p < Res.lenText {
		if max > Res.lenText-p {
			max = Res.lenText - p
		}
		for i := max; i > 0; i-- {
			word := string(text[p : p+i])
			code, ok := dict[i][word]
			if i == 1 {
				// 非汉字
				if !unicode.Is(unicode.Han, text[p]) {
					Res.countNotHan++
					if !strings.Contains(Res.notHan, word) {
						Res.notHan += word
					}
					if ok { // 码表中的非汉字
						part(word, code)
					} else {
						part(word, word)
					}
				} else if ok { // 码表中的汉字
					part(word, code)
				} else if !strings.Contains(Res.lack, word) { // 缺字
					Res.lack += word
					Res.countLack++
				}
				p++
			} else if ok { // 码表中的词
				part(word, code)
				p += i
				break
			}
		}
	}
	Res.codeSep = buf.String()
	write()
	return Res
}
