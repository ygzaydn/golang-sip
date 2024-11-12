package udp

import (
	"net"

	"github.com/ygzaydn/golang-sip/logger"
	"github.com/ygzaydn/golang-sip/models/sip"
)

type AuthenticationParameters struct {
	Authentication string // auth or auth-int or none
	Schema         string // digest or basic or None
}

type ServerParameters struct {
	Uri            string
	Realm          string
	Domain         string
	Authentication AuthenticationParameters
	ServerType     string // server or proxy
}

type ClientCredentials struct {
	Username string
	Password string
}

type ClientParameters struct {
	Uri          string
	Realm        string
	Domain       string
	Credentials  ClientCredentials
	RegistrarURI string
	Contact      string
	DisplayName  string
	UserAgent    string
}

type UDPEntity struct {
	Connection     *net.UDPConn
	logger         *logger.Logger
	entityType     string
	Address        *net.UDPAddr
	LastMessage    *sip.SIPMessage
	MessageChannel chan *sip.SIPMessage
}

type UDPServer struct {
	Entity     UDPEntity
	parameters ServerParameters
}

type UDPClient struct {
	Entity     UDPEntity
	Parameters ClientParameters
}

type UDPListener interface {
	udpListener(bufferSize int)
}
