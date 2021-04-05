package xunfeistt

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/joho/godotenv"
)

func TestXunfeiBasic(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	root := filepath.Join(dir, "..", "..", "..")
	err := godotenv.Load(filepath.Join(root, ".env"))
	cfg := XunFeiSTTCfg{
		Source:   filepath.Join(root, "data", "en.mp3"),
		Language: "en",
	}
	task := XunFeiSTTTask{}
	t.Logf("cfg %+v", cfg)
	task.Init(cfg)
	err = task.Start()
	t.Logf("result :%+v", task.GetResult())
	t.Logf("error :%+v", err)
}
