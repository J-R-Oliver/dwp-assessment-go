package logging

import (
	"log"
	"os"
)

type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	debugLog *log.Logger
}

func New() Logger {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	debugLog := log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	return Logger{
		infoLog, errorLog, debugLog,
	}
}

func (l Logger) Info(logMessage any) {
	l.infoLog.Println(logMessage)
}

func (l Logger) Error(logMessage any) {
	l.errorLog.Println(logMessage)
}

func (l Logger) Fatal(logMessage any) {
	l.errorLog.Fatal(logMessage)
}

func (l Logger) Debug(logMessage any) {
	l.debugLog.Println(logMessage)
}
