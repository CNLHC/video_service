package message

import (
	_ "argus/video/pkg/controller"
	"argus/video/pkg/globalerr"
	"argus/video/pkg/task"
	"argus/video/pkg/task/baiduvca"
	"argus/video/pkg/task/capture"
	"argus/video/pkg/task/clip"
	"encoding/json"
	"errors"

	"github.com/mitchellh/mapstructure"
	nats "github.com/nats-io/nats.go"
	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog/log"
)

const (
	TypeClip    = "Clip"
	TypeCapture = "Capture"
	TypeVCA     = "VCA"
	TypeSTT     = "STT"
)

var (
	ErrUnknownTaskType = errors.New("Unknown TaskType")
)

type Subscriber struct {
	sub *nats.Subscription
}

type MsgResp struct {
	RequestID string `json:"RequestID"`
	Msg       string `json:"Msg"`
	State     string `json:"State"`
	ErrorMsg  string `json:"ErrorMsg"`
}
type TaskDoorbell struct {
	RequestID string `json:"RequestID"`
	Type      string `json:"Type"`
	Reply     string `json:"Reply"`
	Cfg       map[string]interface{}
}

type NSQMsgHandler struct {
	ErrChan chan error
}

func (c *NSQMsgHandler) HandleMessage(msg *nsq.Message) (err error) {
	var door_bell TaskDoorbell

	err = json.Unmarshal(msg.Body, &door_bell)
	if err != nil {
		c.ErrChan <- err
		msg.Finish()
	}
	p := Publisher{Reply: door_bell.Reply}
	err = LaunchTaskAndWait(&door_bell, p)

	if err != nil {
		buf, _ := json.Marshal(&MsgResp{ErrorMsg: err.Error()})
		p.Publish(buf)
		msg.Requeue(-1)
	}

	msg.Finish()
	return
}

func (c *Subscriber) Subscribe() (err error) {
	nsq := GetNSQConsumer()
	handler := NSQMsgHandler{}
	nsq.AddConcurrentHandlers(&handler, 20)

	if err != nil {
		return
	}
	return
}

func LaunchTaskAndWait(doorbell *TaskDoorbell, publisher Publisher) (err error) {
	switch doorbell.Type {
	case TypeClip:
		log.Info().Msgf("handle clip task %+v", doorbell)
		err = runTask(
			&clip.ClipTask{},
			clip.ClipTaskCfg{},
			doorbell,
			publisher)
	case TypeCapture:
		log.Info().Msgf("handle capture task %+v", doorbell)
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
		err = ErrUnknownTaskType
		goto errHandle
	}
	if err != nil {
		goto errHandle
	}
	return
errHandle:
	log.Error().Msgf("run task error %+v", err)
	globalerr.GetGlobalErrorChan() <- err
	return
}

func runTask(t task.AsyncTask, cfg interface{}, doorbell *TaskDoorbell, publisher Publisher) (err error) {

	log.Info().Msgf("Start Decode")
	err = mapstructure.Decode(doorbell.Cfg, &cfg)
	if err != nil {

		return err
	}
	log.Info().Msgf("Start Init Task", t.GetId())
	err = t.Init(cfg)
	log.Info().Msgf("Task Inited: ID(%v)", t.GetId())

	if err != nil {
		return err
	}
	cb := publisher.GetCallback()
	//t.SetCallback(task.EventPrepare, controller.CreateInstanceInDB)
	t.SetCallback(task.EventProgress, cb)
	t.SetCallback(task.EventDone, cb)
	//t.SetCallback(task.EventDone, controller.PersistResult)
	err = t.Start()
	return err
}
