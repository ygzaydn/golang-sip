package utils

import (
	"strings"
	"unicode"
)

func FormatLogMessage(message string) string {
	return strings.ReplaceAll(message, " ", "")
}

func ExtractInfoFromSigns(message string) string {
	return strings.SplitN(strings.SplitN(message, "<", 2)[1], ">", 2)[0]
}

func extractHeaderInfo(headerArray []string) (map[string]any, error) {
	output := make(map[string]any)

	for _, values := range headerArray {
		parsedValues := strings.Split(values, "=")
		key := parsedValues[0]
		value := extractValuesBetweenQuotes(parsedValues[1])
		output[key] = value
	}
	return output, nil
}

func appendMaps(map1, map2 map[string]any) map[string]any {

	result := make(map[string]any)

	for key, value := range map1 {
		result[capitalizeFirstLetter(key)] = value
	}

	for key, value := range map2 {
		result[capitalizeFirstLetter(key)] = value
	}

	return result
}

func extractValuesBetweenQuotes(input string) string {
	var results []string

	if !(len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"') {
		return input
	}

	for {
		start := strings.Index(input, `"`)
		if start == -1 {
			break
		}
		end := strings.Index(input[start+1:], `"`)
		if end == -1 {
			break
		}

		results = append(results, input[start+1:start+end+1])
		input = input[start+end+2:]
	}

	return strings.Join(results, "")
}

func capitalizeFirstLetter(s string) string {
	// If the string is empty, return it as is
	if len(s) == 0 {
		return s
	}

	// Convert the first letter to uppercase and concatenate with the rest of the string
	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}
