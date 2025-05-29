package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

type ConsoleLogger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewLogger() *ConsoleLogger {
	return &ConsoleLogger{
		infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *ConsoleLogger) Info(format string, v ...interface{}) {
	l.infoLog.Printf(format, v...)
}

func (l *ConsoleLogger) Error(format string, v ...interface{}) {
	l.errorLog.Printf(format, v...)
}

func (l *ConsoleLogger) Fatal(format string, v ...interface{}) {
	l.errorLog.Fatalf(format, v...)
}
