package task

import (
	"time"

	"github.com/gofrs/uuid"
)

type BaseTask struct {
	TaskId   uuid.UUID
	Meta     map[string]interface{}
	Callback map[string][]TaskCallback
	StartAt  time.Time
}

func (c BaseTask) GetId() uuid.UUID {
	return c.TaskId
}

func (c BaseTask) GetMeta(key string) (interface{}, bool) {
	a, b := c.Meta[key]
	return a, b
}

func (c BaseTask) SetMeta(key string, item interface{}) {
	c.Meta[key] = item
}

func (c BaseTask) SetCallback(event string, fn TaskCallback) {
	c.Callback[event] = append(c.Callback[event], fn)
}

func (c BaseTask) RunCallback(e string, status TaskStatus, task AsyncTask) {
	if fns, ok := c.Callback[e]; ok {
		for _, fn := range fns {
			if fn != nil {
				fn(task)
			}
		}
	}
}

func NewBaseTask() BaseTask {
	uuid, _ := uuid.NewV4()

	return BaseTask{
		TaskId:   uuid,
		Meta:     make(map[string]interface{}),
		Callback: make(map[string][]TaskCallback),
	}
}
