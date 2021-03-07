package models

import (
	"argus/video/pkg/config"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"fmt"

	"github.com/gofrs/uuid"
)

type Task struct {
	Id        uuid.UUID `sql:"type:uuid;primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
	Type      string
	Result    string
}

func (Task) TableName() string {
	return "async_task"
}

var _db *gorm.DB
var loader sync.Once

func GetDB() *gorm.DB {
	loader.Do(func() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			config.Get("PG_HOST"),
			config.Get("PG_USER"),
			config.Get("PG_PASSWORD"),
			config.Get("PG_DB"),
			config.Get("PG_PORT"),
		)

		if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}); err != nil {
			panic(err)
		} else {
			_db = db
			_db.Logger.LogMode(logger.Info)
		}
	})
	return _db
}
