package logger

import (
	"strings"

	"github.com/fatih/color"
)

type Logger struct {
	LogLevel string
}

func NewLogger(logLevel string) *Logger {
	return &Logger{
		LogLevel: strings.ToLower(logLevel),
	}
}

// Function to check if a message should be shown based on the configured log level
func (l *Logger) isLoggable(level string) bool {
	levels := map[string]int{"debug": 1, "info": 2, "warn": 3, "error": 4}
	return levels[strings.ToLower(level)] >= levels[l.LogLevel]
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if l.isLoggable("info") {
		c := color.New(color.FgHiCyan)
		c.Printf("INFO: "+msg+"\n", args...)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.isLoggable("warn") {
		c := color.New(color.FgYellow)
		c.Printf("WARN: "+msg+"\n", args...)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if l.isLoggable("error") {
		c := color.New(color.FgHiRed)
		c.Printf("ERROR: "+msg+"\n", args...)
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.isLoggable("debug") {
		c := color.New(color.FgMagenta)
		c.Printf("DEBUG: "+msg+"\n", args...)
	}
}
