package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// SetLogLevel configures the global log level based on the LOG_LEVEL environment variable, defaulting to Info if undefined or invalid.
func SetLogLevel() {
	lv, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil || lv == zerolog.NoLevel {
		lv = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lv)
}
