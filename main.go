package main

import (
	"argus/video/pkg/message"
	"sync"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)

	sub := message.Subscriber{}
	if err := sub.Subscribe(); err != nil {
		panic(err.Error())
	}
	wg.Wait()
}
