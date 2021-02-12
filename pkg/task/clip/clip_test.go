package clip

import (
	"argus/video/pkg/utils/video"
	"os"
	"testing"
)

func TestClipBasic(t *testing.T) {
	dest := "/root/Project/argus_video_management/data/out.mp4"

	task := NewClipTask(ClipTaskCfg{
		Src:       "/root/Project/argus_video_management/data/index.mp4",
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
