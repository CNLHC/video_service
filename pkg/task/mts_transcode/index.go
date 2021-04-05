package mts_transcode

import (
	"argus/video/pkg/config"
	"argus/video/pkg/poller"
	"argus/video/pkg/task"
	"argus/video/pkg/task/alivod"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/mts"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

// 作业状态：
// Submitted表示作业已提交，
// Transcoding表示转码中，
// TranscodeSuccess表示转码成功，
// TranscodeFail表示转码失败，
// TranscodeCancelled表示转码取消。
type MTSTranscode struct {
	task.BaseTask
	alivod.AliVodTask
	status     task.TaskStatus
	templateID string
	cfg        MTSTranscodeCfg
	result     string
	location   string
}

type MTSTranscodeCfg struct {
	Bucket string
	Src    string
}

func (c *MTSTranscode) Init() {

}

func (c *MTSTranscode) GetId() uuid.UUID {
	return c.BaseTask.TaskId
}

func (c *MTSTranscode) Start() (err error) {
	var (
		_           *mts.AddMediaResponse
		submit_resp *mts.SubmitJobsResponse
		jobid       string
		p           poller.Poller
	)
	_, err = c.AddMedia()
	if err != nil {
		goto errHandle
	}
	if strings.HasSuffix(strings.ToLower(c.cfg.Src), ".mp4") {
		c.result = c.cfg.Src
		return
	}

	submit_resp, err = c.SubmitJob()
	if err != nil {
		goto errHandle
	}

	jobid = submit_resp.JobResultList.JobResult[0].Job.JobId
	p = poller.GetBasicPoller(
		time.Second*10,
		time.Now().Add(600*time.Second),
		func() (resp interface{}, err error) {
			return c.QueryJob(jobid)
		},
		// 		作业状态：
		// Submitted表示作业已提交，
		// Transcoding表示转码中，
		// TranscodeSuccess表示转码成功，
		// TranscodeFail表示转码失败，
		// TranscodeCancelled表示转码取消。
		func(resp interface{}) (e error) {
			t, ok := resp.(*mts.QueryJobListResponse)
			if !ok {
				return errors.New("Error Response")
			}

			state := t.JobList.Job[0].State
			prog := t.JobList.Job[0].Percent
			if state == "TranscodeSuccess" {
				c.status.Status = task.StatusDone
				c.status.IsRunning = false
				c.status.Progress = 100
				return nil
			}

			if state == "TranscodeFail" {
				c.status.Status = task.StatusFail
				c.status.IsRunning = false
				c.status.Progress = 100
			}

			c.status.IsRunning = true
			c.status.Status = task.StatusRunning
			c.status.Progress = float32(prog)
			return errors.New("unfinished")
		},
	)

	err = p.Start()
	return

errHandle:

	log.Error().Msgf("TransCode Error %+v", err)
	return

}

func (c *MTSTranscode) QueryJob(ids string) (resp *mts.QueryJobListResponse, err error) {
	cli := c.GetVodCli()
	req := mts.CreateQueryJobListRequest()
	req.Scheme = "https"
	req.JobIds = ids
	return cli.QueryJobList(req)

}

func (c *MTSTranscode) AddMedia() (resp *mts.AddMediaResponse, err error) {
	cli := c.GetVodCli()
	req := mts.CreateAddMediaRequest()
	req.Scheme = "https"
	req.FileURL = fmt.Sprintf("http://%s.%s.aliyuncs.com/%s", c.cfg.Bucket, c.location, c.cfg.Src)
	return cli.AddMedia(req)
}

func (c *MTSTranscode) SubmitJob() (resp *mts.SubmitJobsResponse, err error) {
	cli := c.GetVodCli()
	req := mts.CreateSubmitJobsRequest()

	input := mts.Input{
		Bucket:   c.cfg.Bucket,
		Location: c.location,
		Object:   c.cfg.Src,
	}
	type MTSOutputs []mts.Output
	outputs := MTSOutputs{
		mts.Output{
			OutputFile: mts.OutputFile{
				Bucket:   c.cfg.Bucket,
				Object:   fmt.Sprintf("mts/mp4sd/%s", path.Base(c.cfg.Src)),
				Location: c.location,
			},
			TemplateId: config.Get("Ali_MTS_TEMPLATE_ID"),
		},
	}
	input_bytes, err := json.Marshal(input)
	if err != nil {
		return
	}
	outputs_bytes, err := json.Marshal(outputs)
	if err != nil {
		return
	}

	req.Scheme = "https"
	req.PipelineId = config.Get("123ad79742da4dbd8aa99038c067d4c1")
	req.OutputBucket = c.cfg.Bucket
	req.OutputLocation = c.location
	req.Input = string(input_bytes)
	req.Outputs = string(outputs_bytes)
	return cli.SubmitJobs(req)

}

func (c *MTSTranscode) Terminate() error {
	return task.ErrNotAvailable
}

func (c *MTSTranscode) GetResult() (resp task.TaskResult) {
	resp.Err = nil
	return
}

func (c *MTSTranscode) GetStatus() task.TaskStatus {
	return c.status
}
