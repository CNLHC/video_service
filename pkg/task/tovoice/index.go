package tovoice

import (
	"argus/video/pkg/task"
	"argus/video/pkg/task/ffmpeg"
	"errors"
	_ "errors"
	_ "fmt"
	_ "io"
	_ "io/ioutil"
	"reflect"

	"github.com/rs/zerolog/log"
)

var (
	ErrTaskNotStart = errors.New("Task not start")
)

type ToVoiceTask struct {
	ffmpeg.FFMPEGTask
	Cfg    ToVoiceCfg
	result string
	err    error
}

type ToVoiceCfg struct {
	Src  string
	Dest string
}

func (c *ToVoiceTask) Init(cfg interface{}) error {
	switch cfg.(type) {
	case ToVoiceCfg:
		c.Cfg = cfg.(ToVoiceCfg)
		c.FFMPEGTask.Flags = []string{
			"-i", c.Cfg.Src,
			"-vn",
			"-b:a", "128k",
			"-ar", "16000",
			"-ac", "1",
			"-y", c.Cfg.Dest,
		}
		c.BaseTask = task.NewBaseTask()
		return nil
	default:
		log.Error().Msgf("wrong type for clip %+v", reflect.TypeOf(cfg))

	}
	return task.ErrWrongCfg
}

func (c *ToVoiceTask) GetResult() (resp task.TaskResult) {
	resp.Err = nil
	return
}

func NewToVoiceTask(cfg ToVoiceCfg) (res *ToVoiceTask) {
	res = &ToVoiceTask{Cfg: cfg}
	res.Init(cfg)
	return res
}
