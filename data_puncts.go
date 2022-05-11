package smq

import (
	_ "embed"
	"strings"
)

//go:embed assets/puncts.txt
var punctsData string

var puncts = genPuncts()

func genPuncts() map[string]string {
	ret := make(map[string]string)
	punctsData = strings.ReplaceAll(punctsData, "\r\n", "\n")

	lines := strings.Split(punctsData, "\n")
	for _, line := range lines {
		wc := strings.Split(line, "\t")
		if len(wc) != 2 {
			continue
		}
		ret[wc[0]] = wc[1]
	}
	return ret
}
