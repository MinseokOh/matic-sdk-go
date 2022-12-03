package types

import (
	log "github.com/sirupsen/logrus"
)

type Logger struct {
	enable bool
	log    *log.Entry
}

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: false,
		ForceColors:   true,
	})
}

func NewLogger(prefix string, enable bool) *Logger {
	return &Logger{
		enable: enable,
		log: log.WithFields(log.Fields{
			"@client": prefix,
		}),
	}
}

func (logger *Logger) Info(msg string, field map[string]interface{}) {
	if !logger.enable {
		return
	}
	logger.log.WithFields(field).Info(msg)
}

func (logger *Logger) Warn(msg string, field map[string]interface{}) {
	if !logger.enable {
		return
	}
	logger.log.WithFields(field).Warn(msg)
}

func (logger *Logger) Debug(msg string, field map[string]interface{}) {
	if !logger.enable {
		return
	}
	logger.log.WithFields(field).Debug(msg)
}

func (logger *Logger) Error(msg string, field map[string]interface{}) {
	if !logger.enable {
		return
	}
	logger.log.WithFields(field).Error(msg)
}
