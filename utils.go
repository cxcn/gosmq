package smq

import (
	"path/filepath"
	"strings"
)

func GetFileName(fp string) string {
	name := filepath.Base(fp)
	ext := filepath.Ext(fp)
	return strings.TrimSuffix(name, ext)
}
