package capture

import (
	"argus/video/pkg/task"
	"argus/video/pkg/task/ffmpeg"
	"argus/video/pkg/utils/video"
	"errors"
	_ "errors"
	"fmt"
	_ "io"
	_ "io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/rs/zerolog/log"
)

var (
	ErrTaskNotStart = errors.New("Task not start")
)

type CaptureTask struct {
	ffmpeg.FFMPEGTask
}

type CaptureTaskCfg struct {
	Src             string
	Dest            string
	ThumbnailCounts int
}

func (c *CaptureTask) GetTaskType() string {
	return "Capture"
}

func (c *CaptureTask) countsToFPS(cfg CaptureTaskCfg) (fps float64, err error) {
	var (
		prober   = video.Prober{}
		format   video.ProberResp
		duration float64
	)
	format, err = prober.Probe(cfg.Src)
	duration, err = strconv.ParseFloat(format.Format.Duration, 64)
	if err != nil {
		return fps, err
	}
	fps = float64(cfg.ThumbnailCounts) / duration
	return fps, err

}

func (c *CaptureTask) Init(cfg interface{}) (err error) {
	var fps float64

	switch cfg.(type) {
	case CaptureTaskCfg:
		t := cfg.(CaptureTaskCfg)
		fps, err = c.countsToFPS(t)
		if err != nil {
			return err
		}
		if t.Dest == "" {
			log.Warn().Msgf("set default dest")
			t.Dest = "/home/cn/Project/video_service/data/%05d.jpg"
		}
		err = os.MkdirAll(path.Dir(t.Dest), 0777)
		c.FFMPEGTask.Source = t.Src
		c.FFMPEGTask.Flags = []string{
			"-i", fmt.Sprintf("%s", t.Src),
			"-r", fmt.Sprintf("%f", fps),
			"-f", "image2",
			"-y", t.Dest,
		}
		c.BaseTask = task.NewBaseTask()
	}
	return err
}

func NewCaptureTask(cfg CaptureTaskCfg) (res *CaptureTask) {
	res = &CaptureTask{}
	res.Init(cfg)
	return res
}
