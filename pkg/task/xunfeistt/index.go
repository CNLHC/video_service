package xunfeistt

import (
	"argus/video/pkg/task"
	"argus/video/pkg/task/xunfeistt/sdk"
)

type XunFeiSTTTask struct {
	cfg   XunFeiSTTCfg
	xfsdk *sdk.XunfeiSDK
}

type XunFeiSTTCfg struct {
	Source   string
	Language string
}

func (c *XunFeiSTTTask) Init(cfg interface{}) (err error) {
	switch cfg.(type) {
	case XunFeiSTTCfg:
		c.cfg = cfg.(XunFeiSTTCfg)
		AppId := ""
		AppSk := ""
		c.xfsdk, err = sdk.GetXunfeiSDK(AppId, AppSk, c.cfg.Source)

		return
	default:
		return task.ErrWrongCfg
	}
}

func (c *XunFeiSTTTask) Terminate() (err error) {
	return task.ErrNotAvailable
}

func (c *XunFeiSTTTask) Start() (err error) {
	_, err = c.xfsdk.Prepare(sdk.PrepareReq{Language: c.cfg.Language})
	if err != nil {
		return err
	}
	_, err = c.xfsdk.Upload()
	if err != nil {
		return err
	}
	_, err = c.xfsdk.Merge()
	if err != nil {
		return err
	}
	_, err = c.xfsdk.GetProgress()
	if err != nil {
		return err
	}
	return nil
}
