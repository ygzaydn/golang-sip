package logger

import (
	"fmt"
	"time"
)

type LogMessage struct {
	Message string
	Time    time.Time
}

type Logger struct {
	Level   int
	Message LogMessage
	Error   error
}

func New(level int) *Logger {
	return &Logger{
		Level: level,
	}
}

func (l *Logger) BuildLogMessage(message string) {
	l.Message = LogMessage{
		Message: message,
		Time:    time.Now(),
	}

	fmt.Printf("%s - %s", l.Message.Time, l.Message.Message)
}
