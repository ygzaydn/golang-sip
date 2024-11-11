package utils

import (
	"strings"
)

func FormatLogMessage(message string) string {
	return strings.ReplaceAll(message, " ", "")
}

func ExtractInfoFromSigns(message string) string {
	return strings.SplitN(strings.SplitN(message, "<", 2)[1], ">", 2)[0]
}
