package task

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

type TaskMeta map[string]interface{}
type TaskStatus struct {
	IsRunning bool
	Progress  int
	Status    string
	StartAt   time.Time
	ETA       time.Duration
}

type TaskCallback func(c AsyncTask)

const (
	EventProgress = "Progress"
	EventDone     = "Done"
)

var (
	ErrWrongCfg      = errors.New("WrongCfg")
	ErrNotAvailable  = errors.New("Not available")
	ErrUpstreamError = errors.New("Upstream Error")
)

type AsyncTask interface {
	GetId() uuid.UUID
	GetMeta(key string) (interface{}, bool)

	SetMeta(key string, item interface{})
	SetCallback(event string, fn TaskCallback)

	Init(cfg interface{}) error

	Start() error
	Terminate() error
	GetStatus() TaskStatus
}
