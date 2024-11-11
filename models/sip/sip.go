package sip

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ygzaydn/golang-sip/utils"
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
		if s.Method == "REGISTER" {
			builder.WriteString(fmt.Sprintf("%s %s SIP/2.0\r\n", s.Method, utils.ExtractInfoFromSigns(s.Headers["To"][0])))
		} else {
			builder.WriteString(fmt.Sprintf("%s SIP/2.0\r\n", s.Method))
		}
		// Request format

	} else {
		// Response format
		builder.WriteString(fmt.Sprintf("SIP/2.0 %d %s\r\n", s.StatusCode, s.Reason))
	}

	for key, values := range s.Headers {
		for _, value := range values {
			builder.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
		}
	}
	builder.WriteString("\r\n")
	if s.Body != "" {
		builder.WriteString(s.Body)
	}

	return builder.String()
}

func ToSIP(rawMessage string) (*SIPMessage, error) {
	// Will handle incoming requests
	isRequest := ISSIPRequest(rawMessage)
	message := &SIPMessage{
		Headers: make(map[string][]string),
	}
	var err error

	lines := strings.Split(rawMessage, "\r\n")
	firstLine := lines[0]

	if isRequest {
		lineParts := strings.SplitN(firstLine, " ", 3)
		if len(lineParts) < 2 {
			return nil, errors.New("invalid SIP request start line")
		}
		message.StatusCode = 0
		message.Method = lineParts[0]

	} else {
		lineParts := strings.SplitN(firstLine, " ", 3)
		if len(lineParts) < 2 {
			return nil, errors.New("invalid SIP response start line")
		}
		message.StatusCode, _ = strconv.Atoi(lineParts[1])
		message.Reason = lineParts[2]
	}

	i := 1
	for ; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			break // End of headers section, body section must be covered aswell
		}
		headerParts := strings.SplitN(line, ": ", 2)
		if len(headerParts) != 2 {
			return nil, fmt.Errorf("invalid header format: %s", line)
		}
		headerName := headerParts[0]
		headerValue := headerParts[1]
		message.Headers[headerName] = append(message.Headers[headerName], headerValue)
	}

	return message, err
}

func ISSIPRequest(message string) bool {
	lines := strings.Split(message, "\r\n")
	if len(lines) > 0 {
		firstLine := lines[0]
		return !strings.HasPrefix(firstLine, "SIP/")
	}
	return false
}

func ISSIPMessage(message string) bool {

	lines := strings.Split(message, "\r\n")
	startLine := lines[0]

	if strings.HasPrefix(startLine, "SIP/2.0") {
		parts := strings.SplitN(startLine, " ", 3)

		if len(parts) >= 3 && strings.HasPrefix(parts[0], "SIP/2.0") {
			return true
		}
	} else {
		parts := strings.SplitN(startLine, " ", 3)

		if len(parts) == 3 && strings.HasSuffix(parts[2], "SIP/2.0") {
			return true
		}
	}
	return false
}

func (s *SIPMessage) HandleRequest(responseChannel chan *SIPMessage) {
	// Will work as SIP Parser

	switch s.Method {
	case "REGISTER":
		responseChannel <- s.generateTryingMessage()
		//time.Sleep(2 * time.Second)
		if len(s.Headers["Authorization"]) < 1 {
			responseChannel <- s.generate401UnauthorizedMessage()
		} else {
			responseChannel <- s.generateOKMessage()
		}

	}
	switch s.StatusCode {
	case 100:
	case 200:
	case 401:
		responseChannel <- s.handle401UnauthorizedMessage()
	}

}

func (s *SIPMessage) generateTryingMessage() *SIPMessage {
	responseHeaders := map[string][]string{
		"Via":     s.Headers["Via"],
		"From":    s.Headers["From"],
		"To":      s.Headers["To"],
		"Call-ID": s.Headers["Call-ID"],
		"CSeq":    s.Headers["CSeq"],
	}

	return NewResponse(100, "Trying", responseHeaders, "")
}

func (s *SIPMessage) generateOKMessage() *SIPMessage {
	responseHeaders := map[string][]string{
		"Via":     s.Headers["Via"],
		"From":    s.Headers["From"],
		"To":      s.Headers["To"],
		"Call-ID": s.Headers["Call-ID"],
		"CSeq":    s.Headers["CSeq"],
	}

	return NewResponse(200, "OK", responseHeaders, "")
}

func (s *SIPMessage) generate401UnauthorizedMessage() *SIPMessage {
	responseHeaders := map[string][]string{
		"Via":     s.Headers["Via"],
		"From":    s.Headers["From"],
		"To":      s.Headers["To"],
		"Call-ID": s.Headers["Call-ID"],
		"CSeq":    s.Headers["CSeq"],
	}

	return NewResponse(401, "Unauthorized", responseHeaders, "")
}

func (s *SIPMessage) handle401UnauthorizedMessage() *SIPMessage {
	s.Headers["Authorization"] = []string{"Digest username=\"alice\", realm=\"example.com\", nonce=\"xyz\", uri=\"sip:example.com\", response=\"abc123\""}

	parsedCSeq := strings.SplitN(s.Headers["CSeq"][0], " ", 2)

	CSeqNum, err := strconv.Atoi(parsedCSeq[0])

	if err != nil {
		fmt.Println("Wrong CSeq value")
	}

	updatedCSeq := fmt.Sprintf("%d %s", CSeqNum+1, parsedCSeq[1])

	s.Headers["CSeq"] = []string{updatedCSeq}

	return NewRequest("REGISTER", s.Headers, s.Body)

}
