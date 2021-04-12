package ffmpeg

import (
	"argus/video/pkg/task"
	"argus/video/pkg/utils"
	"argus/video/pkg/utils/video"
	"errors"
	_ "errors"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const ()

var (
	ErrTaskNotStart = errors.New("Task not start")
)

type FFMPEGTask struct {
	task.BaseTask

	Source string
	Stats  utils.FFMpegStats
	Flags  []string

	progress_sock net.Listener
	status        task.TaskStatus
	cmd           *exec.Cmd
	meta          video.ProberResp
}

func (c *FFMPEGTask) GetTaskType() string {
	return "FFMPEG"
}

func (c *FFMPEGTask) Terminate() error {
	if c.cmd == nil {
		return ErrTaskNotStart
	}
	return c.cmd.Process.Kill()
}

func (c *FFMPEGTask) prepareFFMPEG() (err error) {
	var (
		ln net.Listener
	)

	c.status.Status = task.StatusPreparing
	c.StartAt = time.Now()
	log.Info().Msgf("%s start at %s: source(%s) %+v", c.GetId(), c.StartAt, c.Source, c)
	c.RunCallback(task.EventPrepare, c.status, c)
	prober := video.Prober{}
	if c.meta, err = prober.Probe(c.Source); err != nil {
		return err
	}
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
	log.Info().Msgf("FFMPeg Task %s:  Listen at %s, cmd is %s", c.GetId().String(), url)
	return
}

func (c *FFMPEGTask) Start() error {
	var (
		err            error
		timeout_ticker = time.NewTicker(10 * time.Minute)
		err_ch         = make(chan error, 1)
	)

	err = c.prepareFFMPEG()
	if err != nil {
		goto finish
	}

	c.cmd = exec.Command("ffmpeg", c.Flags...)
	c.status.Status = task.StatusRunning
	go c.waitProcessExit(err_ch)
	go c.wait(err_ch)

	select {
	case <-timeout_ticker.C:
		err = task.ErrTimeout
	case e := <-err_ch:
		err = e
	}

finish:
	if c.cmd != nil &&
		c.cmd.ProcessState != nil &&
		!c.cmd.ProcessState.Exited() {
		c.cmd.Process.Kill()
	}
	if err != nil {
		c.status.Status = task.StatusFail
		c.RunCallback(task.EventFail, c.status, c)
	} else {
		c.status.Status = task.StatusDone
		c.RunCallback(task.EventDone, c.status, c)
	}

	return err
}
func (c *FFMPEGTask) waitProcessExit(ch chan error) {
	cmd := c.cmd
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

	err = c.cmd.Start()
	if err != nil {
		ch <- err
	}
	outBuf, _ = ioutil.ReadAll(outReader)
	errBuf, _ = ioutil.ReadAll(reader)
	log.Printf("output %s", string(outBuf))
	log.Printf("error %s", string(errBuf))
	cmd.Wait()
	exitcode := cmd.ProcessState.ExitCode()
	log.Printf("code %d", exitcode)
	if strings.Contains(strings.ToLower(string(errBuf)), "error") {
		err = errors.New("ffmpeg output contains error")
		goto handleerror
	}
	if exitcode != 0 {
		err = errors.New("ffmpeg has non-zero exit code")
		goto handleerror
	}
	ch <- nil
	return
handleerror:
	ch <- err

}

func (c *FFMPEGTask) Init(cfg interface{}) error {
	switch cfg.(type) {
	case []string:
		c.Flags = cfg.([]string)
	}
	return nil
}

func (c *FFMPEGTask) wait(ch chan error) {
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
			c.RunCallback(task.EventProgress, c.status, c)
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

func (c *FFMPEGTask) getEndMs() int {
	tidx := -1
	for idx, element := range c.Flags {
		if element == "-t" {
			tidx = idx
		}
	}
	if tidx >= 0 {
		endInSec, err := strconv.Atoi(c.Flags[tidx+1])
		if err == nil {
			return endInSec * 1000
		}
	}

	if s, err := strconv.ParseFloat(c.meta.Format.Duration, 32); err != nil {
		return 0
	} else {
		return int(s * 1000)
	}
}

func (c *FFMPEGTask) getProgress() (progress float32) {
	cur := c.Stats.GetOutputMs()
	end := c.getEndMs()
	if end != 0 && cur >= 0 {
		progress = 100 * float32(cur) / float32(end) / 1000
		if progress >= 100 && c.status.Status != task.StatusDone {
			progress = 99.9
		}
	} else {
		progress = -1
	}
	return progress
}

func (c *FFMPEGTask) getETA() time.Duration {
	return 0
}

func (c *FFMPEGTask) GetStatus() task.TaskStatus {

	c.status.IsRunning = c.isRunning()
	c.status.Progress = c.getProgress()
	c.status.StartAt = c.StartAt
	c.status.ETA = c.getETA()
	return c.status
}

func (c *FFMPEGTask) GetResult() (resp task.TaskResult) {
	resp.Err = nil
	return
}
