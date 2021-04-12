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
	STT_Sub := message.Subscriber{
		TaskType:      "STT",
		MaxConcurrent: 20,
	}
	VCA_Sub := message.Subscriber{
		TaskType:      "VCA",
		MaxConcurrent: 20,
	}
	Transcode_Sub := message.Subscriber{
		TaskType:      "Transcode",
		MaxConcurrent: 20,
	}
	ToVoice_Sub := message.Subscriber{
		TaskType:      "ToVoice",
		MaxConcurrent: 20,
	}

	if err := STT_Sub.Subscribe(); err != nil {
		panic(err.Error())
	}
	if err := VCA_Sub.Subscribe(); err != nil {
		panic(err.Error())
	}
	if err := Transcode_Sub.Subscribe(); err != nil {
		panic(err.Error())
	}
	if err := ToVoice_Sub.Subscribe(); err != nil {
		panic(err.Error())
	}
	wg.Wait()
}
