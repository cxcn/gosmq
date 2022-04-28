package smq

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func GetFileName(fp string) string {
	name := filepath.Base(fp)
	ext := filepath.Ext(fp)
	return strings.TrimSuffix(name, ext)
}

func div(x, y int) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", float64(x)/float64(y)), 64)
	return value
}
