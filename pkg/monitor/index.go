package monitor

import (
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type MonitorContext struct {
	NcSub *nats.Subscription
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
func (c MonitorContext) Report() {
	ticker := time.NewTicker(time.Second * 2)
	for ; true; <-ticker.C {
		a, b, m := c.NcSub.Pending()
		log.Info().Msgf("NC Pending %+v %+v %+v", a, b, m)
	}

}
