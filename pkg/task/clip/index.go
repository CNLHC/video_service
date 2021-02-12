package clip

import (
	"argus/video/pkg/task"
	"argus/video/pkg/task/ffmpeg"
	"argus/video/pkg/utils"
	"errors"
	_ "errors"
	"fmt"
	_ "io"
	_ "io/ioutil"
	"net"
	"os/exec"

	_ "github.com/rs/zerolog/log"
)

var (
	ErrTaskNotStart = errors.New("Task not start")
)

type ClipTask struct {
	ffmpeg.FFMPEGTask
	Cfg           ClipTaskCfg
	progress_sock net.Listener
	cmd           *exec.Cmd
	Stats         utils.FFMpegStats
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

	}
	return nil
}

func NewClipTask(cfg ClipTaskCfg) (res *ClipTask) {
	res = &ClipTask{Cfg: cfg}
	res.Init(cfg)
	return res
}
