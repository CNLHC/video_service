package xunfeistt

import (
	"argus/video/pkg/config"
	"argus/video/pkg/poller"
	"argus/video/pkg/task"
	"argus/video/pkg/task/xunfeistt/sdk"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

type XunFeiSTTTask struct {
	task.BaseTask
	cfg    XunFeiSTTCfg
	xfsdk  *sdk.XunfeiSDK
	result string
	err    error
	status task.TaskStatus
}

type XunFeiSTTCfg struct {
	Source   string
	Language string
	TaskId   string
}

func (c *XunFeiSTTTask) Init(cfg interface{}) (err error) {
	switch cfg.(type) {
	case XunFeiSTTCfg:
		c.cfg = cfg.(XunFeiSTTCfg)
		AppId := config.Get("Xunfei_AID")
		AppSk := config.Get("Xunfei_SK")
		log.Debug().Msgf("Xunfei AID:%s, Xunfei AK: %s", AppId, AppSk)
		c.xfsdk, err = sdk.GetXunfeiSDK(AppId, AppSk, c.cfg.Source)
		c.BaseTask = task.NewBaseTask()

		return
	default:
		return task.ErrWrongCfg
	}
}

func (c *XunFeiSTTTask) Terminate() (err error) {
	return task.ErrNotAvailable
}

func (c *XunFeiSTTTask) GetTaskType() string {
	return "STT"
}

func (c *XunFeiSTTTask) Start() (err error) {

	c.status.Status = task.StatusPreparing
	c.RunCallback(task.EventPrepare, c.status, c)
	var p poller.Poller

	//for test
	if c.cfg.TaskId != "" {
		c.xfsdk.SetTaskID(c.cfg.TaskId)
		goto getResult
	}

	log.Info().Msgf("%s prepare start", c.GetId())
	_, err = c.xfsdk.Prepare(sdk.PrepareReq{Language: c.cfg.Language})
	if err != nil {
		return err
	}
	log.Info().Msgf("%s upload start", c.GetId())
	_, err = c.xfsdk.Upload()
	if err != nil {
		return err
	}
	log.Info().Msgf("%s merge start", c.GetId())
	_, err = c.xfsdk.Merge()
	if err != nil {
		return err
	}
	p = poller.GetBasicPoller(
		5*time.Second,
		time.Now().Add(200*time.Second),
		func() (resp interface{}, err error) {
			resp, err = c.xfsdk.GetProgress()
			return

		},
		func(resp interface{}) (err error) {
			log.Debug().Msgf("task %s check xunfei resp %+v", c.GetId(), resp)
			t := resp.(sdk.Status)
			if t.Status == 9 {
				c.status.Status = task.StatusDone
				c.RunCallback(task.EventDone, c.status, c)
				c.status.Progress = 100
				return nil
			}
			c.status.IsRunning = true
			c.status.Status = task.StatusRunning
			c.status.Progress = float32(t.Status * 10)
			err = errors.New("unfinished")
			return
		})

	err = p.Start()

	if err != nil {
		c.status.Status = task.StatusFail
		c.status.IsRunning = false
		c.RunCallback(task.EventFail, c.status, c)
		return err
	}

getResult:
	resp, err := c.xfsdk.GetResult()
	if err != nil {
		return err
	}
	log.Debug().Msgf("get result message task(%s) %+v", c.GetId(), resp)
	c.result = resp.Data
	return nil
}
func (c *XunFeiSTTTask) GetStatus() task.TaskStatus {
	return c.status
}
func (c *XunFeiSTTTask) GetResult() (resp task.TaskResult) {
	resp.Err = nil
	resp.Data = c.result
	return
}
