package smq

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func (dict *trie) read(rd io.Reader, isS bool) int {
	dictLen := 0
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		// fmt.Println(scan.Text())
		if len(wc) != 2 {
			continue
		}
		if isS && len([]rune(wc[0])) != 1 {
			continue
		}
		dict.insert(wc[0], wc[1])
		dictLen++
	}
	dict.addPunct()
	return dictLen
}

func (dict *trie) readC(rd io.Reader, isS bool, ding int) (int, []byte) {
	dictLen := 0
	scan := bufio.NewScanner(rd)
	freq := make(map[string]int)
	var wb []byte
	// 生成字典
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) != 2 {
			continue
		}
		if isS && len([]rune(wc[0])) != 1 {
			continue
		}
		c := wc[1]
		freq[c]++
		rp := freq[c]

		suf := "_"
		if rp != 1 {
			suf = strconv.Itoa(rp)
			c += suf
		} else if ding > len(c) {
			c += suf
		}

		// 生成赛码表
		wb = append(wb, scan.Bytes()...)
		if rp != 1 || ding > len(c) {
			wb = append(wb, suf...)
		}
		wb = append(wb, '\n')

		dict.insert(wc[0], c)
		dictLen++
	}
	dict.addPunct()

	return dictLen, wb
}
