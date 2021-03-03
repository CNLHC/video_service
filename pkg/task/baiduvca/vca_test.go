package baiduvca

import (
	testutil "argus/video/pkg/utils/test"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func TestBaiduVCABasic(t *testing.T) {
	var (
		err  error
		task = VCATask{}
		root = testutil.GetGoModuleRoot()
	)
	err = godotenv.Load(filepath.Join(root, ".env"))
	t.Logf("root %s", root)
	cfg := VCATaskCfg{
		Source:       "https://publicstatic.cnworkshop.xyz/index.mp4",
		Preset:       "demo",
		Notification: "",
	}
	if err != nil {
		goto error_handle
	}
	err = task.Init(cfg)
	if err != nil {
		goto error_handle
	}
	err = task.Start()
	if err != nil {
		goto error_handle
	}
error_handle:
	if err != nil {
		t.Errorf("%v", err)
	}
}
