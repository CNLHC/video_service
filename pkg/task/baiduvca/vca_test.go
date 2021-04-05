package baiduvca

import (
	"argus/video/pkg/task"
	testutil "argus/video/pkg/utils/test"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func TestBaiduVCABasic(t *testing.T) {
	var (
		err   error
		_task = VCATask{}
		root  = testutil.GetGoModuleRoot()
		res   = task.TaskResult{}
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
	err = _task.Init(cfg)
	if err != nil {
		goto error_handle
	}
	err = _task.Start()
	if err != nil {
		goto error_handle
	}
	res = _task.GetResult()
	t.Logf("result %+v", res)

error_handle:
	if err != nil {
		t.Errorf("%v", err)
	}
}
