package console_tool

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConsoleInit(name string) zerolog.Logger {
	flag.Parse()

	if name != "" {
		return log.With().Str("app", name).Logger()
	} else {
		return log.Logger
	}
}
