package xunfeistt

import "argus/video/pkg/task"

type XunFeiSTTTask struct {
	cfg XunFeiSTTCfg
}

type XunFeiSTTCfg struct {
	Source string
}

func (c *XunFeiSTTTask) Init(cfg interface{}) (err error) {
	switch cfg.(type) {
	case XunFeiSTTCfg:
		c.cfg = cfg.(XunFeiSTTCfg)
		return err
	default:
		return task.ErrWrongCfg
	}
}

func (c *XunFeiSTTTask) Terminate() (err error) {
	return task.ErrNotAvailable
}

func (c *XunFeiSTTTask) Start() (err error) {

	return nil
}
