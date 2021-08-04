package main

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
)

type freq struct {
	code  string
	times int
}

type result struct {
	lenText     int    //文本字数
	notHan      string //非汉字
	countNotHan int    //非汉字数
	lack        string //缺字
	countLack   int    //缺字数

	codeSep string //空格间隔的全部编码
	mapFreq map[string]freq
	choose  map[string]struct{} //选重

	//以下可由上面计算得
	code       string   //全部编码
	codeSlice  []string //编码切片
	lenCode    int      //总键数
	avlCode    float64  //码长
	countSpace int      //空格数

	countWord   int     //打词数
	lenWord     int     //打词字数
	rateWord    float64 //打词率（上屏）
	rateLenWord float64 //打词率（字数）

	countChoose   int     //选重数
	lenChoose     int     //选重字数
	rateChoose    float64 //选重率（上屏）
	rateLenChoose float64 //选重率（字数）

	statCode map[int]int
	statWord map[int]int
}

var res result
var buf bytes.Buffer

func part(w, c string) {

	num := c[len(c)-1]
	if 48 <= num && num <= 57 {
		res.choose[w] = struct{}{}
		if conf.isConf {
			s, err := strconv.Atoi(string(num))
			if err != nil {
				errHandler(err)
			}
			key, ok := conf.ak[s]
			if ok {
				c = c[:len(c)-1] + key
			}
		}
	}
	buf.WriteString(c + " ")

	tmp := res.mapFreq[w]
	if tmp.code == "" {
		tmp.code = c
	}
	tmp.times++
	res.mapFreq[w] = tmp

}

func cacl(fpm, fpt string) {
	text := readText(fpt)
	dict := read(fpm)

	// start := time.Now()
	// defer func() {
	// 	cost := time.Since(start)
	// 	fmt.Println("cacl cost time = ", cost)
	// }()

	res.lenText = len(text)
	res.mapFreq = make(map[string]freq)
	res.choose = make(map[string]struct{})
	max := 0 // 最大词长
	for k := range dict {
		if k > max {
			max = k
		}
	}

	p := 0 // point
	for p < res.lenText {
		if max > res.lenText-p {
			max = res.lenText - p
		}
		for i := max; i > 0; i-- {
			word := string(text[p : p+i])
			code, ok := dict[i][word]
			if i == 1 {
				// 非汉字
				if !unicode.Is(unicode.Han, text[p]) {
					res.countNotHan++
					if !strings.Contains(res.notHan, word) {
						res.notHan += word
					}
					if ok { // 码表中的非汉字
						part(word, code)
					} else {
						part(word, word)
					}
				} else if ok { // 码表中的汉字
					part(word, code)
				} else if !strings.Contains(res.lack, word) { // 缺字
					res.lack += word
					res.countLack++
				}
				p++
			} else if ok { // 码表中的词
				part(word, code)
				p += i
				break
			}
		}
	}
	res.codeSep = buf.String()
	write()
}
