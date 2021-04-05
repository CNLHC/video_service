package message

import (
	"argus/video/pkg/task"
	"encoding/json"
	_ "fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type Publisher struct {
	first_msg_published bool
	last_published      time.Time
	Limit               time.Duration
	Reply               string
	// Msg                 *nats.Msg
}

func (p *Publisher) Publish(message []byte) (err error) {
	nc := GetNATSConn()
	err = nc.Publish(p.Reply, message)
	if err != nil {
		log.Error().Msgf("publish error %+v", err)
	}
	return err
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
			log.Info().Msgf("task(%d) callback is invoked", c.GetId())
			status := c.GetStatus()
			//topic := fmt.Sprintf("%s.status", c.GetId())
			msg, _ := json.Marshal(status)
			if err := p.Publish(msg); err != nil {
				p.first_msg_published = true
			}
		}
	}
}
