package controller

import (
	"argus/video/pkg/globalerr"
	"argus/video/pkg/models"
	"argus/video/pkg/task"
	"encoding/json"

	"github.com/rs/zerolog/log"
)

func CreateInstanceInDB(s task.AsyncTask) {
	_db := models.GetDB()
	log.Info().Msgf("Invoke CreateInstance")
	_db = _db.Create(&models.Task{
		Type: s.GetTaskType(),
		Id:   s.GetId(),
	})

	if _db.Error != nil {
		globalerr.GetGlobalErrorChan() <- _db.Error
	}
}

func PersistResult(s task.AsyncTask) {
	var (
		err error
		buf []byte
	)
	_db := models.GetDB()
	res := s.GetResult()
	buf, err = json.Marshal(res.Data)

	if err != nil {
		goto errhandle
	}

	_db = _db.Model(models.Task{}).
		Where("id = ? ", s.GetId()).
		Update("result", string(buf))

	if _db.Error != nil {
		goto errhandle
	}
errhandle:
	globalerr.GetGlobalErrorChan() <- _db.Error
}
