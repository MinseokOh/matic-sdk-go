package types

import (
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	PanicLevel log.Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type Logger struct {
	enable bool
	log    *log.Entry
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: false,
		ForceColors:   true,
	})
}

func NewLogger(prefix string, config DebugConfig) *Logger {
	log.SetLevel(config.Level)
	return &Logger{
		enable: config.Enable,
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
