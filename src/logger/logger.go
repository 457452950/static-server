package logger

import (
	"github.com/codeskyblue/go-accesslog"
	"log"
)

var (
	logger = Logger{}
)

type Logger struct {
}

func GetLogger() Logger {
	return logger
}

func (l Logger) Log(record accesslog.LogRecord) {
	log.Printf("%s - %s %d %s", record.Ip, record.Method, record.Status, record.Uri)
}
