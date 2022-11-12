package pkg

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func setupLogs() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
