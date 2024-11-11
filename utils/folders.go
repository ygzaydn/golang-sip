package utils

import (
	"fmt"
	"os"
)

func WriteToLogFile(logfile string, itemToWrite string) {
	file, err := os.OpenFile(logfile+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	text := itemToWrite

	_, err = file.WriteString(text)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
