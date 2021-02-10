package clip

import (
	"testing"
	"time"
)

func TestClipBasic(t *testing.T) {

	task := NewClipTask(ClipTaskCfg{
		Src:       "/root/Project/argus_video_management/data/index.mp4",
		Dest:      "/root/Project/argus_video_management/data/out.mp4",
		ClipStart: time.Duration(0 * time.Second),
		ClipEnd:   time.Duration(30 * time.Second),
	})

	task.Start()
}
