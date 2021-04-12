package message

import (
	"argus/video/pkg/config"
	_ "argus/video/pkg/controller"
	"argus/video/pkg/task"
	"argus/video/pkg/task/baiduvca"
	"argus/video/pkg/task/capture"
	"argus/video/pkg/task/clip"
	"argus/video/pkg/task/mts_transcode"
	"argus/video/pkg/task/tovoice"
	"argus/video/pkg/task/xunfeistt"
	"encoding/json"
	"errors"

	"github.com/mitchellh/mapstructure"
	nats "github.com/nats-io/nats.go"
	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog/log"
)

const (
	TypeClip      = "Clip"
	TypeCapture   = "Capture"
	TypeVCA       = "VCA"
	TypeSTT       = "STT"
	TypeTranscode = "Transcode"
	TypeToVoice   = "ToVoice"
)

var (
	ErrUnknownTaskType = errors.New("Unknown TaskType")
)

type Subscriber struct {
	sub           *nats.Subscription
	TaskType      string
	MaxConcurrent int
	Task          task.AsyncTask
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
}

func (c *NSQMsgHandler) HandleMessage(msg *nsq.Message) (err error) {
	var door_bell TaskDoorbell

	err = json.Unmarshal(msg.Body, &door_bell)
	log.Info().
		Str("Action", "Handle NSQ Message").
		Msgf("Request: %s", door_bell)

	if err != nil {
		log.Error().
			Str("Action", "Handle NSQ Message").
			Err(err)
	}

	p := Publisher{Reply: door_bell.Reply}

	log.Error().
		Str("Action", "Handle NSQ Message").
		Str("Stage", "LaunchTaskAndWait").Msg("")

	err = LaunchTaskAndWait(&door_bell, p)

	if err != nil {
		log.Error().
			Str("Action", "Handle NSQ Message").
			Err(err)
		p.Publish(task.BaseAsyncTaskResp{
			State:    task.StatusFail,
			ErrorMsg: err.Error(),
		})
	}

	log.Error().
		Str("Action", "Handle NSQ Message").
		Str("Stage", "Finish").Msg("")
	msg.Finish()
	return nil
}

func (c *Subscriber) Subscribe() (err error) {
	log.Info().Str("TaskType", c.TaskType).Msgf("Register Service")
	_nsq := GetNSQConsumer(c.TaskType)
	handler := NSQMsgHandler{}
	_nsq.AddConcurrentHandlers(&handler, c.MaxConcurrent)

	err = _consumer.ConnectToNSQD(config.Get("NSQ_URL"))
	if err != nil {
		panic(err)
	}
	return
}

func LaunchTaskAndWait(doorbell *TaskDoorbell, publisher Publisher) (err error) {
	log.Info().
		Str("Action", "LaunchTaskAndWait").
		Str("TaskType", TypeClip).Msg("Run")
	switch doorbell.Type {
	case TypeClip:
		err = runTask(&clip.ClipTask{}, clip.ClipTaskCfg{}, doorbell, publisher)
	case TypeCapture:
		err = runTask(&capture.CaptureTask{}, capture.CaptureTaskCfg{}, doorbell, publisher)
	case TypeVCA:
		err = runTask(&baiduvca.VCATask{}, baiduvca.VCATaskCfg{}, doorbell, publisher)
	case TypeToVoice:
		err = runTask(&tovoice.ToVoiceTask{}, tovoice.ToVoiceCfg{}, doorbell, publisher)
	case TypeTranscode:
		err = runTask(&mts_transcode.MTSTranscode{}, &mts_transcode.MTSTranscodeCfg{}, doorbell, publisher)
	case TypeSTT:
		err = runTask(&xunfeistt.XunFeiSTTTask{}, &xunfeistt.XunFeiSTTCfg{}, doorbell, publisher)
	default:
		err = ErrUnknownTaskType
		goto errHandle
	}
	if err != nil {
		goto errHandle
	}
	return
errHandle:
	log.Error().
		Str("Action", "LaunchTaskAndWait").
		Err(err)
	return
}

func runTask(t task.AsyncTask, cfg interface{}, doorbell *TaskDoorbell, publisher Publisher) (err error) {

	log.Info().
		Str("RequestID", doorbell.RequestID).
		Str("Action", "Decode Struct")

	err = mapstructure.Decode(doorbell.Cfg, &cfg)
	if err != nil {
		return err
	}

	log.Info().
		Str("RequestID", doorbell.RequestID).
		Str("Action", "Init Task")
	err = t.Init(cfg)

	if err != nil {
		log.Error().
			Str("RequestID", doorbell.RequestID).
			Msgf("Task Inited Failed", t.GetId())
		return err
	}
	log.Info().
		Str("RequestID", doorbell.RequestID).
		Msgf("Task Inited: ID(%v)", t.GetId())

	cb := publisher.GetCallback()
	//t.SetCallback(task.EventPrepare, controller.CreateInstanceInDB)
	t.SetCallback(task.EventProgress, cb)
	t.SetCallback(task.EventDone, cb)
	//t.SetCallback(task.EventDone, controller.PersistResult)
	err = t.Start()
	if err != nil {
		log.Error().
			Str("RequestID", doorbell.RequestID).
			Str("Action", "Run Task").
			Err(err)
		return err
	}
	return err
}
