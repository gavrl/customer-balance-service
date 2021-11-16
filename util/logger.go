package util

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// InitLogger set JSON formatter (Log as JSON instead of the default ASCII formatter),
// sets log level
func InitLogger() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	cfg := GetConfig()
	// Set proper loglevel based on config
	level, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
}
