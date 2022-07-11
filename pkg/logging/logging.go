package logging

import (
	"log"
	"os"
)

type Level int

const (
	Error Level = iota
	Info
	Debug
)

func (l Level) String() string {
	switch l {
	case Error:
		return "error"
	case Info:
		return "info"
	case Debug:
		return "debug"
	}

	return ""
}

type Logger interface {
	Error(logMessage any)
	Info(logMessage any)
	Debug(logMessage any)
}

type logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	debugLog *log.Logger
	level    Level
}

func New(l Level) Logger {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lmicroseconds)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)
	debugLog := log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	return &logger{
		infoLog, errorLog, debugLog, l,
	}
}

func (l logger) Error(logMessage any) {
	l.errorLog.Println(logMessage)
}

func (l logger) Info(logMessage any) {
	if l.level > 0 {
		l.infoLog.Println(logMessage)
	}
}

func (l logger) Debug(logMessage any) {
	if l.level > 1 {
		l.debugLog.Println(logMessage)
	}
}
