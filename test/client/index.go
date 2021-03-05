package main

import (
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	var (
		msg_file   string
		input_data []byte
	)

	flag.StringVar(&msg_file, "input", "", "")
	flag.Parse()
	input_data, err := ioutil.ReadFile(msg_file)
	if err != nil {
		log.Printf("can not read input from (%s)", msg_file)
	}

	nc, err := nats.Connect("core1.cnworkshop.xyz:24222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Send the request
	err = nc.PublishRequest("updates", "test", input_data)
	sub, err := nc.SubscribeSync("test")

	if err != nil {
		log.Fatal(err)
	}

	max := 100 * time.Second
	start := time.Now()
	for time.Now().Sub(start) < max {
		msg, err := sub.NextMsg(1 * time.Second)
		if err != nil {
			break
		}

		log.Printf("Reply: %s", string(msg.Data))
	}

	sub.Unsubscribe()
	// Use the response
	// Close the connection
	nc.Close()
}
