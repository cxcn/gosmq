package smq

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gogs/chardet"
	"golang.org/x/net/html/charset"
)

// 将 io流 转换为 utf-8
func Tranformer(f io.Reader) io.Reader {

	brd := bufio.NewReader(f)
	buf, _ := brd.Peek(1024)
	detector := chardet.NewTextDetector()
	cs, err := detector.DetectBest(buf) // 检测编码格式
	if err != nil {
		return brd
	}
	if cs.Confidence != 100 && cs.Charset != "UTF-8" {
		cs.Charset = "GB18030"
	}
	// 删除 UTF-8 BOM 文件头
	if cs.Charset == "UTF-8" {
		bom, _ := brd.Peek(3)
		if string(bom) == string(rune(65279)) {
			brd.ReadRune()
		}
	}
	rd, _ := charset.NewReaderLabel(cs.Charset, brd) // 转换字节流
	return rd
}

func readFromString(s string) io.Reader {
	rd := strings.NewReader(s)
	return rd
}

func readFromPath(path string) (io.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return f, err
	}
	return Tranformer(f), nil
}

func GetFileName(fp string) string {
	name := filepath.Base(fp)
	ext := filepath.Ext(fp)
	return strings.TrimSuffix(name, ext)
}

func div(x, y int) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", float64(x)/float64(y)), 64)
	return value
}

type autoSlice struct {
	a []int
}

func newAutoSlice() *autoSlice {
	sli := make([]int, 0, 15)
	return &autoSlice{sli}
}

func (a *autoSlice) AddTo(i int) {
	for i > len(a.a)-1 {
		a.a = append(a.a, 0)
	}
	a.a[i]++
}
