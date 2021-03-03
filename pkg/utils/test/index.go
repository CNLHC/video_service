package testutil

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetGoModuleRoot() string {
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(b)
	)
	for {
		if _, err := os.Stat(filepath.Join(basepath, "go.mod")); os.IsNotExist(err) {
			basepath, _ = filepath.Abs(filepath.Join(basepath, ".."))
		} else {
			return basepath
		}
	}

}
