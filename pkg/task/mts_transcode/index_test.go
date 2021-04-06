package mts_transcode

import (
	testutil "argus/video/pkg/utils/test"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func TestTranscoding(t *testing.T) {
	var (
		err  error
		root = testutil.GetGoModuleRoot()
	)

	err = godotenv.Load(filepath.Join(root, ".env"))

	if err != nil {
		t.Error(err)
	}
	cfg := MTSTranscodeCfg{
		Bucket: "argustest",
		Src:    "index.mp4",
	}

	_task := MTSTranscode{}
	_task.Init(cfg)
	err = _task.Start()
	if err != nil {
		t.Errorf("start error %+v", err)
	}
	result := _task.GetResult()
	t.Logf("Result %+v", result)
}
