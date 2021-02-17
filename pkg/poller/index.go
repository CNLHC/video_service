package poller

import (
	"time"
)

type Poller interface {
	Start() error
}

type ReqFn func() (res interface{}, err error)
type CheckRespFn func(resp interface{}) error

type basicPoller struct {
	Interval  time.Duration
	Deadline  time.Time
	DoReq     ReqFn
	CheckResp CheckRespFn
}

func (c *basicPoller) Start() (err error) {
	ticker := time.NewTicker(c.Interval)
	var resp interface{}
	for ; true; <-ticker.C {
		resp, err = c.DoReq()
		if err != nil {
			return
		}
		if err2 := c.CheckResp(resp); err2 == nil {
			return
		}
	}
	return
}

func GetBasicPoller(
	interval time.Duration,
	dead time.Time,
	reqfn ReqFn,
	respfn CheckRespFn,
) Poller {
	return &basicPoller{
		Interval:  interval,
		Deadline:  dead,
		DoReq:     reqfn,
		CheckResp: respfn,
	}
}
