package sip

import "time"

type AuthenticationParameters struct {
	Authentication string // auth or auth-int or none
	Schema         string // Digest or Basic or None
	Algorithm      string // MD5
}

type ServerParameters struct {
	Uri            string
	Realm          string
	Domain         string
	ServerType     string // server or proxy
	Authentication AuthenticationParameters
	State          map[string]ClientInfo
}

type ClientCredentials struct {
	Username string
	Password string
}

type ClientParameters struct {
	Uri          string
	Realm        string
	Domain       string
	RegistrarURI string
	Contact      string
	DisplayName  string
	UserAgent    string
	Credentials  ClientCredentials
}

type SIPMessage struct {
	Method     string
	StatusCode int
	Reason     string
	Headers    map[string][]string
	Body       string // Optional, I could use map[string][]string in case of SDP
}

type ClientInfo struct {
	CSeq               int
	IsRegistered       bool
	Contact            string
	AuthToken          string
	TransportType      string    // e.g., "UDP", "TCP", "TLS"
	RegistrationExpiry time.Time // When the registration expires
}
