package log

import (
	stdlog "log"
	"os"
)

var (
	logger = stdlog.New(os.Stderr, "[shared] ", stdlog.LstdFlags|stdlog.Lshortfile)
)

func Info(msg string, args ...interface{}) {
	logger.Printf("INFO: "+msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger.Printf("WARN: "+msg, args...)
}

func Error(msg string, args ...interface{}) {
	logger.Printf("ERROR: "+msg, args...)
}
