package baiduvca

import (
	"argus/video/pkg/config"
	"argus/video/pkg/poller"
	"argus/video/pkg/task"
	"argus/video/pkg/utils"
	"encoding/json"
	"io/ioutil"
	_ "os"
	"time"

	"github.com/pkg/errors"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/http"
	"github.com/baidubce/bce-sdk-go/services/vca"
	"github.com/baidubce/bce-sdk-go/services/vca/api"
	"github.com/rs/zerolog/log"
)

// the \lambda of 100(1-exp(x*\lambda)) makes f(60)=90;
// which means the average time to finish 90% of task is 60s
const MagicNumber = -0.038376418

type VCATask struct {
	task.BaseTask
	status task.TaskStatus
	cfg    VCATaskCfg
	cli    *vca.Client
	result string
	err    error
}

type VCATaskCfg struct {
	Source       string
	Preset       string
	Notification string
}

type VCATaskResp map[string]interface{}

func (c *VCATask) GetTaskType() string {
	return "VCA"
}

func (c *VCATask) GetResult() task.TaskResult {
	return task.TaskResult{
		Err:  c.err,
		Data: c.result}
}
func (c *VCATask) Init(cfg interface{}) (err error) {
	switch cfg.(type) {
	case VCATaskCfg:
		c.cfg = cfg.(VCATaskCfg)
		if c.cfg.Source == "" || c.cfg.Preset == "" {
			return task.ErrWrongCfg
		}
		AK, SK := config.Get("Baidu_AK"), config.Get("Baidu_SK")
		ENDPOINT := config.Get("Baidu_endpoint")
		c.cli, err = vca.NewClient(AK, SK, ENDPOINT)
		c.BaseTask = task.NewBaseTask()

		return err
	default:
		return task.ErrWrongCfg
	}
}

func (c *VCATask) Start() (err error) {
	var req api.PutMediaArgs
	c.RunCallback(task.EventPrepare, c.status, c)

	c.status.StartAt = time.Now()
	c.status.IsRunning = true
	c.status.Status = task.StatusPreparing
	req.Source = c.cfg.Source
	req.Preset = c.cfg.Preset
	req.Notification = c.cfg.Notification

	estimator := utils.GetDefaultEstimator()

	log.Info().Msgf("baidu_vca start putting media")
	c.cli.PutMedia(&req)

	log.Info().Msgf("baidu_vca start polling")

	p := poller.GetBasicPoller(
		5*time.Second,
		time.Now().Add(200*time.Second),
		func() (resp interface{}, err error) {
			resp, err = GetMediaTempResult(c.cli, c.cfg.Source, "character")
			return
		},
		func(resp interface{}) (err error) {
			t := resp.(MiddleTypeResp)
			if s, ok := t["status"]; ok {
				if s == "FINISHED" {
					c.status.Progress = 100
					c.status.IsRunning = false
					c.status.Status = task.StatusDone
					var buf []byte
					buf, err = json.Marshal(t)
					c.err = err
					c.result = string(buf)
					c.RunCallback(task.EventDone, c.status, c)
					return nil
				}
			} else {
				c.status.IsRunning = true
				c.status.Status = task.StatusRunning
				c.status.Progress = estimator.EstimatePercentage()
				log.Debug().Msgf("task %s check  baidu resp %+v", c.GetId(), resp)

				c.RunCallback(task.EventProgress, c.status, c)
				err = errors.New("unfinished")
			}
			return
		})

	err = p.Start()

	return err
}

type MiddleTypeResp map[string]interface{}

func GetMediaTempResult(cli bce.Client, source string, middle_type string) (mresp MiddleTypeResp, err error) {
	req := &bce.BceRequest{}
	req.SetUri("/v1/media/" + middle_type)
	req.SetMethod(http.GET)
	req.SetParam("source", source)
	var (
		msg []byte
	)

	// Send request and get response
	resp := &bce.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		t, _ := json.Marshal(resp.ServiceError())
		json.Unmarshal(t, &mresp)

		if resp.ServiceError().Message == "invalid media: media is PROCESSING" ||
			resp.ServiceError().Message == "invalid media: media is PROVISIONING" {
			return mresp, nil

		}
		return nil, err
	}

	msg, err = ioutil.ReadAll(resp.Body())
	if err != nil {
		return
	}
	defer resp.Body().Close()
	if err := json.Unmarshal(msg, &mresp); err != nil {
		return nil, err
	}

	return mresp, err
}
func (c *VCATask) Terminate() (err error) {
	return err
}

func (c *VCATask) GetStatus() task.TaskStatus {
	return c.status
}
