package util

import (
	"os"
	"path/filepath"
)

func GetAbsRunDirPath() (dir string) {
	dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	return
}

func GetHLCRoot() (dir string) {
	dir = os.Getenv("HLC_ROOT")
	if dir == "" {
		dir = GetAbsRunDirPath()
	}
	return
}
