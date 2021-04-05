package tovoice

import (
	testutil "argus/video/pkg/utils/test"
	"path"
	"testing"
)

func TestTovoice(t *testing.T) {
	base := testutil.GetGoModuleRoot()

	dest := path.Join(base, "/data/index.mp3")
	src := path.Join(base, "/data/index.mp4")

	cfg := ToVoiceCfg{
		Src:  src,
		Dest: dest,
	}
	t.Logf("%+v", cfg)
	task := &ToVoiceTask{}
	err := task.Init(cfg)
	if err != nil {
		t.Errorf("init err %+v", err)
		return
	}

	err = task.Start()
	if err != nil {
		t.Errorf("run err %+v", err)
		return
	}
}
