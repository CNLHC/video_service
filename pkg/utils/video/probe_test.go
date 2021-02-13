package video

import "testing"

func TestProbeVideo(t *testing.T) {
	prober := Prober{}
	format, err := prober.Probe("/root/Project/argus_video_management/data/index.mp4")
	if err != nil || format.Format.Duration != "702.766000" {
		t.Errorf("Unexpected output")
	}
	t.Logf("%+v (err: %+v)", format, err)
}
