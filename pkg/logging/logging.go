// Package logging provides a succinct logger that supports log levels. The logger is preconfigured to output log entry's
// with time, date and log level prefixes.
package logging

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
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

func (l *Level) UnmarshalJSON(b []byte) error {
	level, err := stringToLevel(string(b))
	if err != nil {
		return fmt.Errorf("Level.UnmarshalJSON: failed to unmarshal: %w", err)
	}

	*l = level

	return nil
}

func (l *Level) UnmarshalYAML(n *yaml.Node) error {
	level, err := stringToLevel(n.Value)
	if err != nil {
		return fmt.Errorf("Level.UnmarshalYAML: failed to unmarshal: %w", err)
	}

	*l = level

	return nil
}

func stringToLevel(s string) (Level, error) {
	switch s {
	case "error":
		return Error, nil
	case "info":
		return Info, nil
	case "debug":
		return Debug, nil
	}

	return 0, fmt.Errorf("%s is not a valid log level - valid options are error, info or debug", s)
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

// New returns an instance of Logger configured to output log entry of the passed Level or higher.
func New(l Level) Logger {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lmicroseconds)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)
	debugLog := log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	lo := &logger{
		infoLog, errorLog, debugLog, l,
	}

	lo.Info(fmt.Sprintf("Creating Logger with log level: %s", l))

	return lo
}

// Error prints message to Stderr
func (l logger) Error(logMessage any) {
	l.errorLog.Println(logMessage)
}

// Info prints message to Stdout
func (l logger) Info(logMessage any) {
	if l.level > 0 {
		l.infoLog.Println(logMessage)
	}
}

// Debug prints message to Stdout
func (l logger) Debug(logMessage any) {
	if l.level > 1 {
		l.debugLog.Println(logMessage)
	}
}
