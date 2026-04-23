package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelNames = []string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
}

type Logger struct {
	mu       sync.Mutex
	writers  []io.Writer
	level    Level
	prefix   string
	timeFmt  string
}

var defaultLogger *Logger
var once sync.Once

func Default() *Logger {
	once.Do(func() {
		defaultLogger = New(os.Stderr, InfoLevel)
	})
	return defaultLogger
}

func New(w io.Writer, level Level) *Logger {
	return &Logger{
		writers: []io.Writer{w},
		level:   level,
		timeFmt: "2006-01-02 15:04:05",
	}
}

func (l *Logger) AddWriter(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writers = append(l.writers, w)
}

func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format(l.timeFmt)
	levelName := levelNames[level]
	
	if l.prefix != "" {
		msg = fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, levelName, l.prefix, msg)
	} else {
		msg = fmt.Sprintf("[%s] [%s] %s", timestamp, levelName, msg)
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	
	for _, w := range l.writers {
		fmt.Fprintln(w, msg)
	}
	
	if level == FatalLevel {
		os.Exit(1)
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DebugLevel, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WarnLevel, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FatalLevel, format, args...)
}

// Convenience functions for default logger
func Debug(format string, args ...interface{}) {
	Default().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	Default().Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	Default().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	Default().Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	Default().Fatal(format, args...)
}

func SetLevel(level Level) {
	Default().SetLevel(level)
}

func SetPrefix(prefix string) {
	Default().SetPrefix(prefix)
}