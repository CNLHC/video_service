package task

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

type TaskMeta map[string]interface{}
type TaskStatus struct {
	RequestID string
	IsRunning bool
	Progress  float32
	Status    string
	StartAt   time.Time
	ETA       time.Duration
}

type TaskResult struct {
	Err  error
	Data interface{}
}

type TaskCallback func(c AsyncTask)

const (
	EventPrepare  = "Prepare"
	EventProgress = "Progress"
	EventDone     = "Done"
	EventFail     = "Fail"
)

const (
	StatusPreparing = "Preparing"
	StatusRunning   = "Running"
	StatusDone      = "Done"
	StatusFail      = "Fail"
)

var (
	ErrWrongCfg      = errors.New("WrongCfg")
	ErrNotAvailable  = errors.New("Not available")
	ErrUpstreamError = errors.New("Upstream Error")
	ErrTimeout       = errors.New("Timeout")
)

type BaseAsyncTaskResp struct {
	RequestID string      `json:"RequestID"`
	State     string      `json:"State"`
	ErrorMsg  string      `json:"ErrorMsg"`
	Data      interface{} `json:"Data"`
}
type AsyncTask interface {
	GetId() uuid.UUID
	GetMeta(key string) (interface{}, bool)

	SetMeta(key string, item interface{})
	SetCallback(event string, fn TaskCallback)

	Init(cfg interface{}) error
	GetTaskType() string

	Start() error
	Terminate() error
	GetStatus() TaskStatus
	GetResult() TaskResult
}
