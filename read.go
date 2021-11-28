package smq

import (
	"io"
	"os"

	"github.com/gogs/chardet"
	"golang.org/x/net/html/charset"
)

func ReadFile(name string) (*os.File, io.Reader, error) {

	f, err := os.Open(name)
	if err != nil {
		return f, nil, err
	}
	buf := make([]byte, 1024)
	f.Read(buf)
	detector := chardet.NewTextDetector()
	cs, err := detector.DetectBest(buf) // 检测编码格式
	if err != nil {
		return f, nil, err
	}
	if cs.Confidence != 100 && cs.Charset != "UTF-8" {
		cs.Charset = "GB18030"
	}

	f.Seek(0, 0)
	rd, err := charset.NewReaderLabel(cs.Charset, f) // 转换字节流
	return f, rd, err
}
