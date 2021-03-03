package baiduvca

import (
	"argus/video/pkg/config"
	"argus/video/pkg/poller"
	"argus/video/pkg/task"
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
		AK, SK := config.Get("Baidu_AK"), config.Get("Baidu_SK")
		ENDPOINT := config.Get("Baidu_endpoint")
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
			log.Debug().Msgf("task %s check  baidu resp %+v", c.GetId(), resp)
			t := resp.(MiddleTypeResp)
			if s, ok := t["source"]; ok {
				if s == "FINISHED" {
					return nil
				}
			} else {
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
