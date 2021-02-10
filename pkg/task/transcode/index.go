package transcode

import (
	"argus/video/pkg/task"

	"github.com/gofrs/uuid"
)

type Transcode struct {
	task.BaseTask
}

func (c *Transcode) GetId() uuid.UUID {
	return c.BaseTask.TaskId
}
