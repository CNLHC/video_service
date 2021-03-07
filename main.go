package main

import (
	"argus/video/pkg/message"
	"argus/video/pkg/models"
	"sync"

	"github.com/joho/godotenv"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	godotenv.Load("./.env")

	_db := models.GetDB()
	_db.AutoMigrate(&models.Task{})
	sub := message.Subscriber{}
	if err := sub.Subscribe(); err != nil {
		panic(err.Error())
	}
	wg.Wait()
}
