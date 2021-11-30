package smq

import (
	"bufio"
	"io"

	"github.com/gogs/chardet"
	"golang.org/x/net/html/charset"
)

func ReadFile(f io.Reader) (io.Reader, error) {

	brd := bufio.NewReader(f)
	buf, _ := brd.Peek(1024)
	detector := chardet.NewTextDetector()
	cs, err := detector.DetectBest(buf) // 检测编码格式
	if err != nil {
		return brd, err
	}
	if cs.Confidence != 100 && cs.Charset != "UTF-8" {
		cs.Charset = "GB18030"
	}
	rd, err := charset.NewReaderLabel(cs.Charset, brd) // 转换字节流
	return rd, err
}
