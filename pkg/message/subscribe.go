package message

import (
	"argus/video/pkg/controller"
	"argus/video/pkg/task"
	"argus/video/pkg/task/baiduvca"
	"argus/video/pkg/task/capture"
	"argus/video/pkg/task/clip"
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	nats "github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

const (
	TypeClip    = "Clip"
	TypeCapture = "Capture"
	TypeVCA     = "VCA"
	TypeSTT     = "STT"
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
func msgHandler(msg *nats.Msg) {
	var (
		err error
	)

	log.Info().Msgf("receive msg %s", string(msg.Data))
	var doorbell TaskDoorbell
	p := Publisher{Msg: msg}
	err = json.Unmarshal(msg.Data, &doorbell)
	if err != nil {
		log.Error().Msgf("error msg %s", err.Error())
	}
	LaunchTaskAndWait(&doorbell, p)
}

func (c *Subscriber) Subscribe() (err error) {
	nc, err := nats.Connect(
		"core1.cnworkshop.xyz:24222",
		nats.ErrorHandler(errorHandler))

	nc.QueueSubscribe("updates", "default", msgHandler)

	c.sub, err = nc.SubscribeSync("updates")
	if err != nil {
		return
	}
	return
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
	case TypeVCA:
		log.Info().Msgf("handle vca task %+v", doorbell)
		err = runTask(
			&baiduvca.VCATask{},
			baiduvca.VCATaskCfg{},
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
	t.SetCallback(task.EventPrepare, controller.CreateInstanceInDB)
	t.SetCallback(task.EventProgress, cb)
	t.SetCallback(task.EventDone, cb)
	t.SetCallback(task.EventDone, controller.PersistResult)
	err = t.Start()
	return err
}
