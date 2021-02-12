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

	nc, err := nats.Connect("localhost:24222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Send the request
	msg, err := nc.Request("updates", input_data, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	// Use the response
	log.Printf("Reply: %s", string(msg.Data))

	// Close the connection
	nc.Close()

}
