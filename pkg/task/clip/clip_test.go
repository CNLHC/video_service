package clip

import (
	testutil "argus/video/pkg/utils/test"
	"argus/video/pkg/utils/video"
	"os"
	"path"
	"testing"
)

func TestClipBasic(t *testing.T) {
	base := testutil.GetGoModuleRoot()
	dest := path.Join(base, "/data/out.mp4")
	src := path.Join(base, "/data/index.mp4")

	task := NewClipTask(ClipTaskCfg{
		Src:       src,
		Dest:      dest,
		ClipStart: "0",
		ClipEnd:   "30",
	})

	task.Start()

	prober := video.Prober{}
	format, err := prober.Probe(dest)

	if err != nil || format.Format.Duration != "30.001000" {
		t.Errorf("Unexpected output %+v", format.Format)
	}
	os.Remove(dest)
	t.Logf("%+v (err: %+v)", format, err)
}
