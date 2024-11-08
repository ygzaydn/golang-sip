package utils

import (
	"strings"
)

func FormatLogMessage(message string) string {
	return strings.ReplaceAll(message, "", "")
}
