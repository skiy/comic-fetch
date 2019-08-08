package filepath

import (
	"path/filepath"
	"regexp"
)

// Ext returns the real file name extension used by path.
func Ext(path string) string {
	preg := `^.[a-zA-Z0-9]*`
	ext := filepath.Ext(path)
	reg, err := regexp.Compile(preg)
	if err != nil {
		return ""
	}

	return reg.FindString(ext)
}
