package manager

import (
	"argus/video/pkg/task"

	"github.com/gofrs/uuid"
)

type TaskManager interface {
	ListTasks() []task.AsyncTask
	GetTask(id uuid.UUID) task.AsyncTask
	AddTask(t task.AsyncTask)
	FinishTask(t task.AsyncTask)
}
