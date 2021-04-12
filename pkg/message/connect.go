package message

import (
	"argus/video/pkg/config"
	"sync"

	nats "github.com/nats-io/nats.go"
	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog/log"
)

var _consumer *nsq.Consumer
var load_nsq_producer_once sync.Once
var _nc *nats.Conn
var load_nc_once sync.Once

func GetNSQConsumer(subject string) *nsq.Consumer {
	var err error
	nsq_cfg := nsq.NewConfig()
	_consumer, err = nsq.NewConsumer(
		subject,
		config.Get("NSQ_CHANNEL"),
		nsq_cfg,
	)
	if err != nil {
		panic(err)
	}

	return _consumer
}

func errorHandler(nc *nats.Conn, s *nats.Subscription, err error) {
	log.Error().Msgf("nast error %+v", err)
}

func GetNATSConn() *nats.Conn {
	load_nc_once.Do(func() {
		var err error
		_nc, err = nats.Connect(
			config.Get("NATS_URL"),
			nats.ErrorHandler(errorHandler))
		if err != nil {
			panic(err)
		}
	})
	return _nc
}
