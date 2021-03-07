package main

import (
	"argus/video/pkg/globalerr"
	"argus/video/pkg/message"
	"argus/video/pkg/models"
	"argus/video/pkg/monitor"
	"sync"

	"github.com/joho/godotenv"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	godotenv.Load("./.env")
	var mon monitor.MonitorContext

	_db := models.GetDB()
	_db.AutoMigrate(&models.Task{})
	suber := message.Subscriber{}
	if sub, err := suber.Subscribe(); err != nil {
		panic(err.Error())
	} else {
		mon.NcSub = sub
		go mon.Report()
		go globalerr.Listen()
	}
	wg.Wait()
}
