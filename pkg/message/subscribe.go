package message

import (
	"argus/video/pkg/task"
	"argus/video/pkg/task/clip"
	"encoding/json"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

const (
	TypeClip    = "Clip"
	TypeCapture = "Capture"
)

type Subscriber struct {
	sub *nats.Subscription
}

type MsgResp struct {
	Code int
	Msg  string
}
type TaskDoorbell struct {
	Type string
	Cfg  map[string]interface{}
}

func errorHandler(nc *nats.Conn, s *nats.Subscription, err error) {
	log.Error().Msgf("msg error %+v", err)

}

func (c *Subscriber) Subscribe() (err error) {
	nc, err := nats.Connect(
		"localhost:24222",
		nats.ErrorHandler(errorHandler))

	c.sub, err = nc.SubscribeSync("updates")

	for {
		msg, err := c.sub.NextMsg(10 * time.Second)
		if err == nil {
			log.Info().Msgf("receive msg %s", string(msg.Data))
			var doorbell TaskDoorbell
			p := Publisher{Msg: msg}
			err = json.Unmarshal(msg.Data, &doorbell)
			if err != nil {
				log.Error().Msgf("error msg %s", err.Error())
			}
			LaunchTaskAndWait(&doorbell, p)
		}
	}

	return err
}

func LaunchTaskAndWait(doorbell *TaskDoorbell, publisher Publisher) {
	var err error
	switch doorbell.Type {
	case TypeClip:
		log.Info().Msgf("handle clip task %+v", doorbell)
		var (
			cfg clip.ClipTaskCfg
		)
		temp, _ := json.Marshal(doorbell.Cfg)
		err = json.Unmarshal(temp, &cfg)
		if err != nil {
			goto errHandle
		}
		t := &clip.ClipTask{}
		t.Init(cfg)
		cb := publisher.GetCallback()
		t.SetCallback(task.EventDone, cb)
		t.SetCallback(task.EventProgress, cb)
		t.Start()
	case TypeCapture:
		goto errHandle

	default:
		goto errHandle
	}
	return
errHandle:
	buf, _ := json.Marshal(&MsgResp{Code: -1, Msg: err.Error()})
	publisher.Publish(buf)
	return
}
