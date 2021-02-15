package baiduvca

import (
	"argus/video/pkg/task"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/http"
	"github.com/baidubce/bce-sdk-go/services/vca"
	"github.com/baidubce/bce-sdk-go/services/vca/api"
)

type VCATask struct {
	task.BaseTask
	status task.TaskStatus
	cfg    VCATaskCfg
	cli    *vca.Client
}

type VCATaskCfg struct {
	Source       string
	Preset       string
	Notification string
}

type VCATaskResp map[string]interface{}

func (c *VCATask) Init(cfg interface{}) (err error) {
	switch cfg.(type) {
	case VCATaskCfg:
		c.cfg = cfg.(VCATaskCfg)
		AK, SK := "<your-access-key-id>", "<your-secret-access-key>"
		ENDPOINT := "<domain-name>"
		c.cli, err = vca.NewClient(AK, SK, ENDPOINT)
		return err
	default:
		return task.ErrWrongCfg
	}
}

func (c *VCATask) Start() (err error) {
	var req api.PutMediaArgs

	req.Source = c.cfg.Source
	req.Preset = c.cfg.Preset
	req.Notification = c.cfg.Notification

	c.cli.PutMedia(&req)

	return err
}

type MiddleTypeResp map[string]interface{}

func GetMediaTempResult(cli bce.Client, source string, middle_type string) (mresp MiddleTypeResp, err error) {
	req := &bce.BceRequest{}
	req.SetUri("v2/media/" + middle_type)
	req.SetMethod(http.GET)
	req.SetParam("source", source)

	// Send request and get response
	resp := &bce.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	if err := resp.ParseJsonBody(mresp); err != nil {
		return nil, err
	}
	return mresp, err
}

func (c *VCATask) poll() (resp MiddleTypeResp, err error) {
	var (
		interval = 3 * time.Second
		max_iter = 1500
		counter  = 0
	)

	ticker := time.NewTicker(interval)

	for ; true; <-ticker.C {
		counter += 1
		resp, err = GetMediaTempResult(c.cli, c.cfg.Source, "character")
		if err != nil {
			return resp, err
		}
		if resp["status"] == "FINISHED" {
			return resp, nil
		}
		if resp["status"] == "ERROR" {
			return resp, task.ErrUpstreamError
		}
		if counter > max_iter {
			return resp, task.ErrTimeout
		}
	}
	return resp, err

}

func (c *VCATask) Terminate() (err error) {
	return err
}

func (c *VCATask) GetStatus() task.TaskStatus {
	return c.status
}
