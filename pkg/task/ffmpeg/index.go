package ffmpeg

import (
	"argus/video/pkg/task"
	"argus/video/pkg/utils"
	"errors"
	_ "errors"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	ErrTaskNotStart = errors.New("Task not start")
)

type FFMPEGTask struct {
	task.BaseTask
	progress_sock net.Listener
	cmd           *exec.Cmd
	Stats         utils.FFMpegStats
	Flags         []string
}

func (c *FFMPEGTask) Terminate() error {
	if c.cmd == nil {
		return ErrTaskNotStart
	}
	return c.cmd.Process.Kill()
}

func (c *FFMPEGTask) Start() error {
	c.StartAt = time.Now()
	log.Info().Msgf("%s start at %s", c.GetId(), c.StartAt)
	var (
		ln  net.Listener
		err error
	)

	ln, err = net.Listen("tcp", "127.0.0.1:0")
	c.progress_sock = ln
	if err != nil {
		return err
	}

	url := fmt.Sprintf("tcp://%s", ln.Addr().String())

	c.Flags = append(
		c.Flags,
		"-progress", url,
		"-hide_banner",
	)

	cmd := exec.Command("ffmpeg", c.Flags...)
	c.cmd = cmd

	log.Info().Msgf("FFMPeg Task %s:  Listen at %s, cmd is %s", c.GetId().String(), url, cmd.String())

	log.Printf("cmd:%s", cmd.String())

	var outBuf []byte
	var errBuf []byte
	reader, err := cmd.StderrPipe()
	if err != nil {
		log.Error().Msgf("%s", err.Error())
	}
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Error().Msgf("%s", err.Error())
	}
	err = cmd.Start()
	if err != nil {
		log.Error().Msgf("%s", err.Error())
	}
	go c.wait()

	if err != nil {
		log.Error().Msgf("%s", err.Error())
	}
	outBuf, _ = ioutil.ReadAll(outReader)
	errBuf, _ = ioutil.ReadAll(reader)
	_ = outBuf
	_ = errBuf
	log.Printf("output %s", string(outBuf))
	log.Printf("error %s", string(errBuf))
	cmd.Wait()

	if strings.Contains(strings.ToLower(string(errBuf)), "error") {
		err = errors.New("ffmpeg output contains error")
		return err
	}
	exitcode := cmd.ProcessState.ExitCode()
	log.Printf("code %d", exitcode)
	if exitcode != 0 {
		err = errors.New("ffmpeg has non-zero exit code")
		return err
	}

	return nil
}

func (c *FFMPEGTask) Init(cfg interface{}) error {
	switch cfg.(type) {
	case []string:
		c.Flags = cfg.([]string)
	}
	return nil
}

func (c *FFMPEGTask) wait() {
	var (
		buf = make([]byte, 1024)
	)
	conn, err := c.progress_sock.Accept()
	defer conn.Close()
	log.Info().Msgf("Accept Connection %s", conn.RemoteAddr().String())
	if err != nil {
		// handle error
		log.Error().Msg(err.Error())
	}
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			c.Stats = c.Stats.Parse(string(buf[:n]))
			if fns, ok := c.Callback[task.EventProgress]; ok {
				for _, fn := range fns {
					if fn != nil {
						fn(c)
					}
				}
			}
			log.Printf("%+v", c.Stats)
		} else {
			continue
		}
		if err != nil {
			log.Printf("err: %s", err.Error())
			return
		}
	}
}

func (c *FFMPEGTask) isRunned() bool {
	if c.cmd == nil {
		return false
	}
	if c.cmd.ProcessState == nil {
		return false
	}
	return true
}

func (c *FFMPEGTask) isRunning() bool {
	if !c.isRunned() {
		return false
	}
	return !c.cmd.ProcessState.Exited()
}

func (c *FFMPEGTask) getProgress() int {
	return 0
}

func (c *FFMPEGTask) getETA() time.Duration {
	return 0
}

func (c *FFMPEGTask) GetStatus() task.TaskStatus {
	return task.TaskStatus{
		IsRunning: c.isRunning(),
		Progress:  c.getProgress(),
		StartAt:   c.StartAt,
		Status:    "",
		ETA:       c.getETA(),
	}
}

func (c *FFMPEGTask) GetResult() (resp task.TaskResult) {
	resp.Err = nil
	return
}
