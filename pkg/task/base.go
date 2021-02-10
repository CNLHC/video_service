package task

import "github.com/gofrs/uuid"

type BaseTask struct {
	TaskId   uuid.UUID
	Meta     map[string]interface{}
	Callback []TaskCallback
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

func (c BaseTask) SetCallback(fn TaskCallback) {
	c.Callback = append(c.Callback, fn)
}
