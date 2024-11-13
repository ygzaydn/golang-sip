package utils

import (
	"fmt"
	"math/rand"
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

func CheckTag(field string) string {
	split := strings.SplitN(field, ">;tag=", 2)
	if len(split) > 1 {
		return split[1]
	}
	return ""
}
