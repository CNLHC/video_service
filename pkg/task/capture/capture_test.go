package capture

import (
	"os"
	"testing"
)

func TestCaptureBasic(t *testing.T) {

	src := "/home/cn/Project/video_service/data/index.mp4"

	dest_base := "/home/cn/Project/video_service/data/thumbnail/"
	dest := "%05d.png"
	os.MkdirAll(dest_base, os.FileMode(0777))
	cfg := CaptureTaskCfg{
		Src:             src,
		Dest:            dest_base + dest,
		ThumbnailCounts: 15,
	}

	task := &CaptureTask{}
	err := task.Init(cfg)
	if err != nil {
		t.Errorf("Init error %+v", err.Error())
		return
	}

	err = task.FFMPEGTask.Start()
	if err != nil {
		t.Errorf("start error %+v", err.Error())
	}
	_, err = os.Stat(dest_base + "00017.png")

	if os.IsNotExist(err) {
		t.Errorf("can not get capture")
	}
	os.RemoveAll(dest_base)

}
