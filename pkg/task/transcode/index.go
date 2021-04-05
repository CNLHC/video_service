package transcode

import (
	"argus/video/pkg/task"
	"argus/video/pkg/task/alivod"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vod"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type VodTranscode struct {
	task.BaseTask
	alivod.AliVodTask
	status task.TaskStatus
}

func (c *VodTranscode) GetId() uuid.UUID {
	return c.BaseTask.TaskId
}

func (c *VodTranscode) Start() error {

	cli := c.GetVodCli()
	request := vod.CreateSubmitTranscodeJobsRequest()

	request.Scheme = "https"
	request.VideoId = c.VideoId
	request.TemplateGroupId = "gid"
	response, err := cli.SubmitTranscodeJobs(request)

	if err != nil {
		log.Error().Msgf("upstream error %s", err.Error)
		return errors.Wrap(task.ErrUpstreamError, err.Error())
	}
	c.VodJobId = response.TranscodeTaskId
	c.poll()
	return nil
}

func (c *VodTranscode) poll() (err error) {
	var (
		counter  = 0
		max_iter = 500
		interval = 2 * time.Second
	)
	cli := c.GetVodCli()
	request := vod.CreateGetTranscodeTaskRequest()
	request.Scheme = "https"
	request.TranscodeTaskId = c.VodJobId

	ticker := time.NewTicker(interval)
	for ; true; <-ticker.C {
		response, err := cli.GetTranscodeTask(request)
		if counter > max_iter {
			goto finish
		}
		if err != nil {
			log.Error().Err(err)
			goto finish
		}
		c.parseVodStatus(response)
		if !c.status.IsRunning {
			goto finish
		}
		counter++
	}
finish:
	return err
}

func (c *VodTranscode) parseVodStatus(resp *vod.GetTranscodeTaskResponse) (isTerminated bool) {
	var totalProcess int64
	for _, process := range resp.TranscodeTask.TranscodeJobInfoList {
		if process.TranscodeJobStatus == "Transcoding" {
			totalProcess += process.TranscodeProgress
		} else if process.TranscodeJobStatus == "TranscodeSuccess" {
			totalProcess += 100
		}
	}

	c.status.Progress = (int(totalProcess) / len(resp.TranscodeTask.TranscodeJobInfoList))
	c.status.Status = resp.TranscodeTask.TaskStatus
	ts := c.status.Status
	if ts == "Processing" || ts == "Partial" {
		c.status.IsRunning = true
	} else {
		c.status.IsRunning = false
	}
	return isTerminated
}

func (c *VodTranscode) Terminate() error {
	return task.ErrNotAvailable
}

func (c *VodTranscode) GetResult() (resp task.TaskResult) {
	resp.Err = nil
	return
}

func (c *VodTranscode) GetStatus() task.TaskStatus {
	return c.status
}
