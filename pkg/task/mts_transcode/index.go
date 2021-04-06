package mts_transcode

import (
	"argus/video/pkg/config"
	"argus/video/pkg/poller"
	"argus/video/pkg/task"
	"argus/video/pkg/task/alivod"
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/mts"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
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

func (c *MTSTranscode) Init(cfg interface{}) (err error) {
	switch t := cfg.(type) {
	case MTSTranscodeCfg:
		c.cfg = t
		c.location = "oss-cn-beijing"
		c.AliVodTask.Init(alivod.AliVodTaskCfg{})
		return
	default:
		return task.ErrWrongCfg
	}
}

func (c *MTSTranscode) GetId() uuid.UUID {
	return c.BaseTask.TaskId
}

func (c *MTSTranscode) Start() (err error) {
	var (
		_           *mts.AddMediaResponse
		submit_resp *mts.SubmitJobsResponse
		jobid       string
		job         mts.JobResult
		p           poller.Poller
	)
	_, err = c.AddMedia()
	if err != nil {
		goto errHandle
	}

	// if strings.HasSuffix(strings.ToLower(c.cfg.Src), ".mp4") {
	// 	c.result = c.cfg.Src
	// 	return
	// }

	submit_resp, err = c.SubmitJob()

	if err != nil {
		goto errHandle
	}
	job = submit_resp.JobResultList.JobResult[0]
	jobid = job.Job.JobId
	if !job.Success {
		err = errors.Wrap(task.ErrUpstreamError, job.Message)
		goto errHandle

	}

	log.Info().Str("MTS_JOB_ID", jobid).Msgf("Submit Job To AliYun MTS %+v", submit_resp)

	p = poller.GetBasicPoller(
		time.Second*10,
		time.Now().Add(600*time.Second),
		func() (resp interface{}, err error) {

			log.Debug().Str("MTS_JOB_ID", jobid).Msg("Send Query")
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
				log.Error().Msg("Error Response")
				return errors.Wrap(poller.ErrFatal, "Error Response")
			}
			job := t.JobList.Job[0]

			state := job.State
			prog := job.Percent

			log.Debug().Str("MTS_JOB_ID", jobid).
				Str("MTS_JOB_STATE", state).
				Str("MTS_ERROR_CODE", job.Code).
				Msg("Query Response")

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

				if job.Code == "ConditionTranscoding.AudioBitrateNotSatisfied" {
					c.status.Status = task.StatusDone
					return nil
				} else {
					return errors.Wrap(poller.ErrFatal, job.Message)
				}

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

	furl := fmt.Sprintf("http://%s.%s.aliyuncs.com/%s", c.cfg.Bucket, c.location, c.cfg.Src)
	log.Printf("media url %s", furl)
	req.FileURL = (furl)
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
	type MTSOutput struct {
		OutputObject string
		TemplateId   string
	}

	type MTSOutputs []MTSOutput

	outputs := MTSOutputs{
		MTSOutput{
			OutputObject: fmt.Sprintf("mts/mp4sd/%s", path.Base(c.cfg.Src)),
			TemplateId:   config.Get("Ali_MTS_TEMPLATE_ID"),
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

	log.Debug().Msgf("job input %s", string(input_bytes))
	log.Debug().Msgf("job output %s", string(outputs_bytes))

	req.Scheme = "https"
	req.PipelineId = config.Get("Ali_MTS_PIPELINE_ID")
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
	resp.Data = c.result
	return
}

func (c *MTSTranscode) GetStatus() task.TaskStatus {
	return c.status
}
