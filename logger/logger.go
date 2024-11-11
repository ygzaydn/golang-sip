package logger

import (
	"fmt"
	"time"

	"github.com/ygzaydn/golang-sip/utils"
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
		Time:    time.Now().UTC(),
	}
	logLine := fmt.Sprintf("%s\t- %s\n", l.Message.Time, l.Message.Message)
	fmt.Print(logLine)
	utils.WriteToLogFile("log", logLine)
}
