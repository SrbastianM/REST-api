package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// Return a new Logger instance wich writes log entries at or above a minimun severity
// level to a specific output destination.
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// This helpers writes log entries at difrent levels.
func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

// This internal method writes the log entry
func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	// If the severity level of the log entry is below the minimun severity for the logger,
	// then return with no further action.
	if level < l.minLevel {
		return 0, nil
	}
	// Declare an anonymous struct holding the data for the entry.
	aux := struct {
		Level      string
		Time       string
		Message    string
		Properties map[string]string
		Trace      string
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}
	// Include the stack trace for entries at the ERROR and FATAL levels.
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}
	// Declare a line variable for holding the actual log entry text.
	var line []byte
	// Marshal the annonymous struct to JSON and store it in the line variable. If there was a problem
	// creating the JSON, set the contents of the log entry to be that plain-text error message instead.
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}
	// Lock the mutext so that no two writes to the output destination cannot happen concurrently.
	l.mu.Lock()
	defer l.mu.Unlock()
	// Write the log entry followed by a newline.
	return l.out.Write(append(line, '\n'))
}

// This helper implements a Write() method to the Logger type so it satisfies the
// io.Writer interface. Write a log entry at the ERROR level with no additional
// properties
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
