package capture

import (
	"os"
	"testing"
)

func TestCaptureBasic(t *testing.T) {

	src := "/root/Project/argus_video_management/data/index.mp4"

	dest_base := "/root/Project/argus_video_management/data/thumbnail/"
	dest := "%05d.png"
	os.MkdirAll(dest_base, os.FileMode(0777))

	task := NewCaptureTask(CaptureTaskCfg{
		Src:             src,
		Dest:            dest_base + dest,
		ThumbnailCounts: 15,
	})

	task.Start()
	_, err := os.Stat(dest_base + "00017.png")

	if os.IsNotExist(err) {
		t.Errorf("can not get capture")
	}
	os.RemoveAll(dest_base)

}
