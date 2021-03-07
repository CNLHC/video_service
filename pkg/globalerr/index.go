package globalerr

import "github.com/rs/zerolog/log"

var errchan chan error

func init() {
	errchan = make(chan error, 2)
}

func GetGlobalErrorChan() chan error {
	return errchan
}

func Listen() {
	var err error
	for {
		err = <-errchan
		log.Error().Msgf("Global Error %+v", err)
	}
}
