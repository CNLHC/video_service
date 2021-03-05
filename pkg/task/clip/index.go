package clip

import (
	"argus/video/pkg/task"
	"argus/video/pkg/task/ffmpeg"
	"errors"
	_ "errors"
	"fmt"
	_ "io"
	_ "io/ioutil"
	"reflect"

	"github.com/rs/zerolog/log"
)

var (
	ErrTaskNotStart = errors.New("Task not start")
)

type ClipTask struct {
	ffmpeg.FFMPEGTask
	Cfg ClipTaskCfg
}

type ClipTaskCfg struct {
	Src       string
	Dest      string
	ClipStart string
	ClipEnd   string
}

func (c *ClipTask) Init(cfg interface{}) error {
	switch cfg.(type) {
	case ClipTaskCfg:
		c.Cfg = cfg.(ClipTaskCfg)
		c.FFMPEGTask.Source = c.Cfg.Src
		c.FFMPEGTask.Flags = []string{
			"-ss",
			c.Cfg.ClipStart,
			"-t",
			c.Cfg.ClipEnd,
			"-i",
			fmt.Sprintf("%s", c.Cfg.Src),
			"-codec",
			"copy",
			"-y",
			c.Cfg.Dest,
		}
		c.BaseTask = task.NewBaseTask()
		return nil
	default:
		log.Error().Msgf("wrong type for clip %+v", reflect.TypeOf(cfg))

	}
	return task.ErrWrongCfg
}

func NewClipTask(cfg ClipTaskCfg) (res *ClipTask) {
	res = &ClipTask{Cfg: cfg}
	res.Init(cfg)
	return res
}
