package logger

import (
	"github.com/rs/zerolog"
	"os"
	"sync"
)

var once sync.Once

var log zerolog.Logger

func Get() zerolog.Logger {
	once.Do(func() {
		log = zerolog.New(os.Stdout).With().Timestamp().Str("service", "go-buxclient").Logger()

	})
	return log
}
