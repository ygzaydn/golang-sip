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

func (s *SIPMessage) GenerateRegisterAfter401() (*SIPMessage, error) {
	requestHeaders := s.Headers
	var err error = nil
	if requestHeaders["WWW-Authenticate"] == nil {
		return nil, errors.New("www-authenticate header must be present")
	}

	parsedWWWHeader, err := utils.ParseWWWAuthenticateHeader(requestHeaders["WWW-Authenticate"][0])
	if err != nil {
		return nil, err
	}

	parsedFromHeader, err := utils.ParseFromandToHeader(requestHeaders["From"][0])
	if err != nil {
		return nil, err
	}

	fmt.Println(parsedWWWHeader)
	fmt.Println(parsedFromHeader)

	HA1 := utils.MD5Hasher(fmt.Sprintf("%s:%s:%s", parsedFromHeader["User"], parsedWWWHeader["Realm"], parsedFromHeader["User"]))

	HA2 := utils.MD5Hasher(fmt.Sprintf("REGISTER:sip:%s", parsedFromHeader["Host"]))

	nonce_count := "00000001"

	response := utils.MD5Hasher(fmt.Sprintf("%s:%s:%s:%s:%s:%s", HA1, parsedWWWHeader["Nonce"], nonce_count, parsedWWWHeader["Opaque"], parsedWWWHeader["Qop"], HA2))

	authorizationHeader := fmt.Sprintf("%s username=\"%s\", realm=\"%s\", nonce=\"%s\", uri=\"sip:%s\", response=\"%s\", opaque=\"%s\", qop=\"%s\", cnonce=\"%s\", nc=%s, algorithm=%s", parsedWWWHeader["Schema"], parsedFromHeader["User"], parsedWWWHeader["Realm"], parsedWWWHeader["Nonce"], parsedFromHeader["Host"], response, parsedWWWHeader["Opaque"], parsedWWWHeader["Qop"], parsedWWWHeader["Nonce"], nonce_count, parsedWWWHeader["Algorithm"])

	delete(requestHeaders, "WWW-Authenticate")
	requestHeaders["Authorization"] = []string{authorizationHeader}

	fmt.Println(authorizationHeader)

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
