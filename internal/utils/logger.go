package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

type Logger struct {
	level  LogLevel
	format string
	output io.Writer
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

func NewLogger(level, format string) *Logger {
	logLevel := INFO
	switch level {
	case "debug":
		logLevel = DEBUG
	case "info":
		logLevel = INFO
	case "warn":
		logLevel = WARN
	case "error":
		logLevel = ERROR
	}

	return &Logger{
		level:  logLevel,
		format: format,
		output: os.Stdout,
	}
}

func (l *Logger) log(level LogLevel, message string, fields map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     string(level),
		Message:   message,
		Fields:    fields,
	}

	if l.format == "json" {
		data, _ := json.Marshal(entry)
		fmt.Fprintln(l.output, string(data))
	} else {
		fieldsStr := ""
		if len(fields) > 0 {
			fieldsData, _ := json.Marshal(fields)
			fieldsStr = " " + string(fieldsData)
		}
		fmt.Fprintf(l.output, "[%s] %s: %s%s\n", entry.Timestamp, entry.Level, message, fieldsStr)
	}
}

func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		DEBUG: 0,
		INFO:  1,
		WARN:  2,
		ERROR: 3,
	}
	return levels[level] >= levels[l.level]
}

func (l *Logger) Debug(message string) {
	l.log(DEBUG, message, nil)
}

func (l *Logger) DebugWithFields(message string, fields map[string]interface{}) {
	l.log(DEBUG, message, fields)
}

func (l *Logger) Info(message string) {
	l.log(INFO, message, nil)
}

func (l *Logger) InfoWithFields(message string, fields map[string]interface{}) {
	l.log(INFO, message, fields)
}

func (l *Logger) Warn(message string) {
	l.log(WARN, message, nil)
}

func (l *Logger) WarnWithFields(message string, fields map[string]interface{}) {
	l.log(WARN, message, fields)
}

func (l *Logger) Error(message string) {
	l.log(ERROR, message, nil)
}

func (l *Logger) ErrorWithFields(message string, fields map[string]interface{}) {
	l.log(ERROR, message, fields)
}

// Fatal logs an error message and exits the program
func (l *Logger) Fatal(message string) {
	l.Error(message)
	log.Fatal(message)
}
