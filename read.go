package smq

import (
	"bufio"
	"io"
	"os"
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
	rd, _ := charset.NewReaderLabel(cs.Charset, brd) // 转换字节流
	return rd
}

func readFromString(s string) io.Reader {
	rd := strings.NewReader(s)
	return rd
}

func readFromPath(s string) (io.Reader, error) {
	f, err := os.Open(s)
	if err != nil {
		return f, err
	}
	return Tranformer(f), nil
}
