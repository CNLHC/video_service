package task

import (
	"time"

	"github.com/gofrs/uuid"
)

type TaskMeta map[string]interface{}
type TaskStatus struct {
	Progress int
	Status   string
	StartAt  time.Time
	ETA      time.Duration
}
type TaskCallback func(c AsyncTask)

type AsyncTask interface {
	GetId() uuid.UUID
	GetMeta(key string) (interface{}, bool)

	SetMeta(key string, item interface{})
	SetCallback(fn TaskCallback)

	Start() error
	Pause() error
	Terminate() error
	GetStatus() TaskStatus
}
