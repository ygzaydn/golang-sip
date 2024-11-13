package udp

import (
	"net"

	"github.com/ygzaydn/golang-sip/logger"
	"github.com/ygzaydn/golang-sip/models/sip"
)

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
	Parameters sip.ServerParameters
}

type UDPClient struct {
	Entity     UDPEntity
	Parameters sip.ClientParameters
}

type UDPListener interface {
	udpListener(bufferSize int)
}
