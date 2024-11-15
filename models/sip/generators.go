package sip

import (
	"errors"
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

func (s *SIPMessage) handle401RegisterMessage(info ServerParameters) (*SIPMessage, error) {
	//var err error
	requestHeaders := s.Headers
	var err error = nil
	if requestHeaders["Authorization"] == nil {
		return nil, errors.New("authorization header must be present")
	}

	authorizationHeader, err := utils.ParseWWWAuthenticateandAuthorizationHeader(requestHeaders["Authorization"][0])

	if err != nil {
		return nil, err
	}

	fromHeader, err := utils.ParseFromandToHeader(requestHeaders["From"][0])

	if err != nil {
		return nil, err
	}

	HA1 := utils.GenerateHA1(fromHeader["User"].(string), authorizationHeader["Realm"].(string), fromHeader["User"].(string))

	HA2 := utils.GenerateHA2("REGISTER", fromHeader["Host"].(string))

	userState := info.State[fromHeader["User"].(string)]

	response := utils.GenerateResponse(HA1, userState.Nonce, authorizationHeader["Nc"].(string), userState.Opaque, authorizationHeader["Qop"].(string), HA2)

	cSeqHeader, err := utils.ParseCSeqHeader(requestHeaders["CSeq"][0])

	if err != nil {
		return nil, err
	}

	if response == authorizationHeader["Response"].(string) {
		newCSeq := cSeqHeader["CSeq"].(int) + 1
		requestHeaders["CSeq"] = []string{fmt.Sprintf("%d %s", newCSeq, cSeqHeader["Method"])}
		return NewResponse(200, "OK", requestHeaders, ""), nil
	}

	return s.generate401UnauthorizedMessage(info), err
}

func (s *SIPMessage) GenerateRegisterAfter401() (*SIPMessage, error) {
	requestHeaders := s.Headers
	var err error = nil
	if requestHeaders["WWW-Authenticate"] == nil {
		return nil, errors.New("www-authenticate header must be present")
	}

	parsedWWWHeader, err := utils.ParseWWWAuthenticateandAuthorizationHeader(requestHeaders["WWW-Authenticate"][0])
	if err != nil {
		return nil, err
	}

	parsedFromHeader, err := utils.ParseFromandToHeader(requestHeaders["From"][0])
	if err != nil {
		return nil, err
	}

	HA1 := utils.GenerateHA1(parsedFromHeader["User"].(string), parsedWWWHeader["Realm"].(string), parsedFromHeader["User"].(string))

	HA2 := utils.GenerateHA2("REGISTER", parsedFromHeader["Host"].(string))

	nonce_count := "00000001"

	response := utils.GenerateResponse(HA1, parsedWWWHeader["Nonce"].(string), nonce_count, parsedWWWHeader["Opaque"].(string), parsedWWWHeader["Qop"].(string), HA2)

	authorizationHeader := utils.GenerateAuthorizationHeader(parsedWWWHeader["Schema"].(string), parsedFromHeader["User"].(string), parsedWWWHeader["Realm"].(string), parsedWWWHeader["Nonce"].(string), parsedFromHeader["Host"].(string), response, parsedWWWHeader["Opaque"].(string), parsedWWWHeader["Qop"].(string), parsedWWWHeader["Nonce"].(string), nonce_count, parsedWWWHeader["Algorithm"].(string))

	delete(requestHeaders, "WWW-Authenticate")
	requestHeaders["Authorization"] = []string{authorizationHeader}

	return NewRequest("REGISTER", requestHeaders, ""), err
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

// func (s *SIPMessage) GenerateRegisterAfter401() map[string]string {
// TODO
// }
