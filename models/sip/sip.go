package sip

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ygzaydn/golang-sip/utils"
)

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

func (s *SIPMessage) ClientHandler(channel chan *SIPMessage) {
	// Will work as SIP Parser
	//output := make([]*SIPMessage, 0)
	switch s.Method {
	case "REGISTER":
		//output = append(output, s.generateTryingMessage())
		channel <- s.generateTryingMessage()
		time.Sleep(2 * time.Second)
		if len(s.Headers["Authorization"]) < 1 {
			//output = append(output, s.generate401UnauthorizedMessage())
		} else {
			//output = append(output, s.generateOKMessage())
			channel <- s.generateOKMessage()
		}

	}
	switch s.StatusCode {
	case 100:
	case 200:
	case 401:

	}
}

func (s *SIPMessage) ServerHandler(channel chan *SIPMessage, info ServerParameters) {
	// Will work as SIP Parser
	//output := make([]*SIPMessage, 0)
	switch s.Method {
	case "REGISTER":
		//output = append(output, s.generateTryingMessage())
		channel <- s.generateTryingMessage()
		time.Sleep(2 * time.Second)
		if len(s.Headers["Authorization"]) < 1 {
			//output = append(output, s.generate401UnauthorizedMessage())
			channel <- s.generate401UnauthorizedMessage(info)
		} else {
			//output = append(output, s.generateOKMessage())
			channel <- s.generateOKMessage()
		}

	}
	switch s.StatusCode {
	case 100:
	case 200:
	case 401:

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

func (s *SIPMessage) generate401UnauthorizedMessage(info ServerParameters) *SIPMessage {
	//var err error
	responseHeaders := s.Headers
	toTag := utils.CheckTag(responseHeaders["To"][0])

	if toTag == "" {
		responseHeaders["To"] = []string{responseHeaders["To"][0] + ";tag=" + utils.GenerateTag()}
	}

	fromTag := utils.CheckTag(responseHeaders["From"][0])
	if fromTag == "" {
		responseHeaders["From"] = []string{responseHeaders["From"][0] + ";tag=" + utils.GenerateTag()}
	}

	delete(responseHeaders, "Contact")
	delete(responseHeaders, "Expires")
	delete(responseHeaders, "Max-Forwards")
	delete(responseHeaders, "User-Agent")

	responseHeaders["WWW-Authenticate"] = []string{fmt.Sprintf("%s realm=\"%s\", nonce=\"%s\", opaque=\"%s\", algorithm=%s, qop=\"%s\"", info.Authentication.Schema, info.Realm, utils.GenerateNonce(), utils.GenerateOpaque(), info.Authentication.Algorithm, info.Authentication.Authentication)}

	return NewResponse(401, "Unauthorized", responseHeaders, "")
}

func (s *SIPMessage) ShouldCloseResponseChannel() bool {
	if s.StatusCode == 200 {
		return true
	}
	if s.StatusCode == 401 {
		return true
	}

	return false
}

func (s *SIPMessage) ShouldCloseRequestChannel() bool {
	if s.Method == "REGISTER" {
		return true
	}
	return false
}

func GenerateInitialRegisterHeaders(port int, parameters ClientParameters) map[string][]string {
	portString := fmt.Sprintf("%d", port)
	return map[string][]string{
		"Via": {
			"SIP/2.0/UDP " + parameters.Realm + ":" + portString + ";branch=" + utils.GenerateBranch(),
		},
		"From":           {"<" + parameters.Uri + ">;tag=" + utils.GenerateTag()},
		"To":             {"<" + parameters.Uri + ">"},
		"Call-ID":        {utils.GenerateCallID() + "@" + parameters.Domain},
		"CSeq":           {utils.GenerateCSeq() + " REGISTER"},
		"Contact":        {parameters.Contact},
		"Content-Length": {"0"}, // No body in this request
		"Max-Forwards":   {"70"},
		"User-Agent":     {parameters.UserAgent},
		"Expires":        {"3600"},
	}
}
