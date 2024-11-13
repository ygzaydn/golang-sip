package sip

import (
	"fmt"

	"github.com/ygzaydn/golang-sip/utils"
)

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
