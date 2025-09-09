package logger

import (
	"log"
	"os"

	"github.com/nhutphat1203/hestia-backend/internal/config"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	level Level
	std   *log.Logger
	cfg   *config.Config
}

func New(cfg *config.Config) *Logger {
	if cfg == nil {
		log.Fatal("Config cannot be nil")
	}
	var level Level = DEBUG
	switch cfg.LogLevel {
	case "info":
		level = INFO
	case "warn":
		level = WARN
	case "error":
		level = ERROR
	default:
		level = DEBUG
	}

	return &Logger{
		level: level,
		std:   log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		cfg:   cfg,
	}
}

func (l *Logger) log(level Level, prefix string, msg string) {
	if level >= l.level {
		l.std.Printf("[%s] %s", prefix, msg)
	}
}

func (l *Logger) Debug(msg string) { l.log(DEBUG, "DEBUG", msg) }
func (l *Logger) Info(msg string)  { l.log(INFO, "INFO", msg) }
func (l *Logger) Warn(msg string)  { l.log(WARN, "WARN", msg) }
func (l *Logger) Error(msg string) { l.log(ERROR, "ERROR", msg) }
