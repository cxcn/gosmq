package smq

import (
	"fmt"
	"strings"
	"testing"
)

func TestJisu(t *testing.T) {
	src := `人心不足蛇吞象	ahnh
人民邮电出版社	amjd4
咯	gjf_
嗝	gjfe2
骼	gjg_
`
	dict := &Dict{
		reader:     strings.NewReader(src),
		Transfer:   new(jisu),
		SelectKeys: "_;'",
	}
	got := dict.Transfer.Read(dict)
	fmt.Println("极速格式转换\n", string(got))
}

func TestDuoduo(t *testing.T) {
	src := `人	a
然	a
如果	a
瑞	aa
仍然	aa
睿	aaa
仍然是	aad
瑞士	aadi
锐升	aado
`
	dict := &Dict{
		reader:     strings.NewReader(src),
		Transfer:   new(duoduo),
		PushStart:  4,
		SelectKeys: "_;'",
	}
	got := dict.Transfer.Read(dict)
	fmt.Println("多多格式转换\n", string(got))
}

func TestJidian(t *testing.T) {
	src := `a 人 然 如果
aa 瑞 仍然
aaa 睿
aad 仍然是
aadi 瑞士
`
	dict := &Dict{
		reader:     strings.NewReader(src),
		Transfer:   new(jidian),
		PushStart:  4,
		SelectKeys: "_;'",
	}
	got := dict.Transfer.Read(dict)
	fmt.Println("极点格式转换\n", string(got))
}
