package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateBranch() string {
	uniqueID := uuid.New().String()

	//Typically, the prefix is z9hG4bK (which is a well-known constant used to identify the branch in a SIP transaction).

	branch := "z9hG4bK" + uniqueID

	return branch
}

func GenerateTag() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	return fmt.Sprintf("%d", rng.Int())
}

func GenerateCSeq() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	return fmt.Sprintf("%d", rng.Int())
}

func GenerateCallID() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	return fmt.Sprintf("%x", rng.Int63())
}

func GenerateNonce() string {
	nonce := uuid.New()
	return nonce.String()[:7]
}

func GenerateOpaque() string {
	nonce := uuid.New()
	return nonce.String()[:12]
}

func GenerateAuthorizationHeader(schema, username, realm, nonce, uri, response, opaque, qop, cnonce, nonce_count, algorithm string) string {
	return fmt.Sprintf("%s username=\"%s\", realm=\"%s\", nonce=\"%s\", uri=\"sip:%s\", response=\"%s\", opaque=\"%s\", qop=\"%s\", cnonce=\"%s\", nc=%s, algorithm=%s", schema, username, realm, nonce, uri, response, opaque, qop, cnonce, nonce_count, algorithm)
}

func CheckTag(field string) string {
	split := strings.SplitN(field, ">;tag=", 2)
	if len(split) > 1 {
		return split[1]
	}
	return ""
}

func ParseCSeqHeader(header string) (map[string]any, error) {
	output := make(map[string]any)
	parsedHeader := strings.Split(header, " ")
	if len(parsedHeader) != 2 {
		return nil, errors.New("invalid CSeq header")
	}

	output["Method"] = parsedHeader[1]

	CSeq, err := strconv.Atoi(parsedHeader[0])

	if err != nil {
		return nil, errors.New("fail to parse CSeq number")
	}
	output["CSeq"] = CSeq

	return output, nil
}

func ParseWWWAuthenticateandAuthorizationHeader(header string) (map[string]any, error) {
	output := make(map[string]any)
	parsedHeader := strings.SplitAfterN(header, " ", 2)

	if len(parsedHeader) < 1 {
		return nil, errors.New("fail to parse WWW-Authenticate header")
	}

	output["Schema"] = strings.TrimSpace(parsedHeader[0])

	feedInfo := strings.Split(parsedHeader[1], ", ")

	if len(feedInfo) < 1 {
		return nil, errors.New("fail to parse WWW-Authenticate header")
	}

	restInfo, err := extractHeaderInfo(feedInfo)
	return appendMaps(output, restInfo), err
}

func ParseFromandToHeader(header string) (map[string]any, error) {
	output := make(map[string]any)
	parsedHeader := strings.Split(header, ";")
	if len(parsedHeader) < 1 {
		return nil, errors.New("fail to parse From header")
	}

	if len(parsedHeader) > 1 {
		parseTag := strings.Split(parsedHeader[1], "=")
		output[capitalizeFirstLetter(parseTag[0])] = parseTag[1]
	}

	infoBetweenSign := ExtractInfoFromSigns(parsedHeader[0])

	if len(infoBetweenSign) < 1 {
		return nil, errors.New("fail to parse From header")
	}

	parsedInfoBetweenSign := strings.Split(infoBetweenSign, "@")

	if len(parsedInfoBetweenSign) < 1 {
		return nil, errors.New("fail to parse From header")
	}

	output["User"] = parsedInfoBetweenSign[0][4:]
	output["Host"] = parsedInfoBetweenSign[1]

	return output, nil
}
