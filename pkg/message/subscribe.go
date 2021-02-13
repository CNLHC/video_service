package message

import (
	"argus/video/pkg/task"
	"argus/video/pkg/task/capture"
	"argus/video/pkg/task/clip"
	"encoding/json"
	"time"

	"github.com/mitchellh/mapstructure"
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
		} else {
			log.Error().Msgf("error on msg err:%v msg:%s", err, string(msg.Data))
		}
	}

	return err
}

func LaunchTaskAndWait(doorbell *TaskDoorbell, publisher Publisher) {
	var err error

	switch doorbell.Type {
	case TypeClip:
		log.Info().Msgf("handle clip task %+v", doorbell)
		err = runTask(
			&clip.ClipTask{},
			clip.ClipTaskCfg{},
			doorbell,
			publisher)

	case TypeCapture:
		log.Info().Msgf("handle clip task %+v", doorbell)
		err = runTask(
			&capture.CaptureTask{},
			capture.CaptureTaskCfg{},
			doorbell,
			publisher)

	default:
		goto errHandle
	}
	if err != nil {
		goto errHandle
	}
	return
errHandle:
	buf, _ := json.Marshal(&MsgResp{Code: -1, Msg: err.Error()})
	publisher.Publish(buf)
	return
}

func runTask(t task.AsyncTask, cfg interface{}, doorbell *TaskDoorbell, publisher Publisher) (err error) {
	err = mapstructure.Decode(doorbell.Cfg, &cfg)
	if err != nil {
		return err
	}
	err = t.Init(cfg)
	if err != nil {
		return err
	}
	cb := publisher.GetCallback()
	t.SetCallback(task.EventDone, cb)
	t.SetCallback(task.EventProgress, cb)
	err = t.Start()
	return err
}
