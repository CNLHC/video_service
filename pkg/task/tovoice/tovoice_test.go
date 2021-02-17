package tovoice

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestTovoice(t *testing.T) {

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	cfg := ToVoiceCfg{
		Src: filepath.Join(dir, "..", "..", "..", "data", "index.mp4"),

		Dest: filepath.Join(dir, "..", "..", "..", "data", "index.mp3"),
	}
	task := NewToVoiceTask(cfg)
	err := task.Init(cfg)
	t.Logf("err %+v", err)
	err = task.Start()
	t.Logf("err %+v", err)
	if err != nil {
		t.Errorf("%+v", err)
	}
}
