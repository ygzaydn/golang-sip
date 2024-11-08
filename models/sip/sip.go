package sip

import (
	"fmt"
	"strings"
)

type SIPMessage struct {
	Method     string
	StatusCode int
	Reason     string
	Headers    map[string][]string
	Body       string // Optional, I could use map[string][]string in case of SDP
}

func NewRequest(method string, headers map[string][]string, body string) *SIPMessage {
	return &SIPMessage{
		Method:  method,
		Headers: headers,
		Body:    body,
	}
}

func NewResponse(statusCode int, reason string, headers map[string][]string, body string) *SIPMessage {
	return &SIPMessage{
		StatusCode: statusCode,
		Reason:     reason,
		Headers:    headers,
		Body:       body,
	}
}

func (s *SIPMessage) ToString() string {
	var builder strings.Builder

	if s.Method != "" {
		// Request format
		builder.WriteString(fmt.Sprintf("%s SIP/2.0\r\n", s.Method))
	} else {
		// Response format
		builder.WriteString(fmt.Sprintf("SIP/2.0 %d %s\r\n", s.StatusCode, s.Reason))
	}

	for key, values := range s.Headers {
		for _, value := range values {
			builder.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
		}
	}
	if s.Body != "" {
		builder.WriteString("\r\n")
		builder.WriteString(s.Body)

	}

	return builder.String()
}

func ToSIP(SIPString string) *SIPMessage {
	// Will work as SIP Parser
	return &SIPMessage{}
}

func HandleSIPRequest() {
	// Will handle incoming requests
}

func ISSIPRequest(message string) bool {
	lines := strings.Split(message, "\r\n")
	if len(lines) > 0 {
		firstLine := lines[0]
		return !strings.HasPrefix(firstLine, "SIP/")
	}
	return false
}
