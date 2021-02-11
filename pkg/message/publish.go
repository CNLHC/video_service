package message

import (
	"argus/video/pkg/task"
	"encoding/json"
	"fmt"
	"time"
)

func Publish(topic string, message []byte) error {
	return nil
}

type Publisher struct {
	first_msg_published bool
	last_published      time.Time
	Limit               time.Duration
}

func (p *Publisher) ShouldPublish() bool {
	if p.first_msg_published == false {
		return true
	}
	return time.Now().Sub(p.last_published) > p.Limit

}

func (p *Publisher) GetCallback() task.TaskCallback {
	return func(c task.AsyncTask) {
		if p.ShouldPublish() {
			status := c.GetStatus()
			topic := fmt.Sprintf("%s.status", c.GetId())
			msg, _ := json.Marshal(status)
			if err := Publish(topic, msg); err != nil {
				p.first_msg_published = true
			}
		}
	}
}
