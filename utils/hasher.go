package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func MD5Hasher(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GenerateHA1(username, realm, password string) string {
	return MD5Hasher(fmt.Sprintf("%s:%s:%s", username, realm, password))
}

func GenerateHA2(method, uri string) string {
	return MD5Hasher(fmt.Sprintf("%s:sip:%s", method, uri))
}

func GenerateResponse(HA1, nonce, nonce_count, opaque, qop, HA2 string) string {
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s", HA1, nonce, nonce_count, opaque, qop, HA2)
}
