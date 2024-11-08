package client

import (
	"net"

	"github.com/ygzaydn/golang-sip/logger"
	"github.com/ygzaydn/golang-sip/models/sip"
	"github.com/ygzaydn/golang-sip/utils"
)

type UDPClient struct {
	Connection *net.UDPConn
	logger     *logger.Logger
}

func New(ip string, port int, logger *logger.Logger) (*UDPClient, error) {
	conn, err := connectUDPServer(ip, port)
	if err != nil {
		return nil, err
	}

	return &UDPClient{
		Connection: conn,
		logger:     logger,
	}, err
}

func (u *UDPClient) SendMessage(message *sip.SIPMessage) {
	u.logger.BuildLogMessage(utils.FormatLogMessage(message.ToString()))
	u.Connection.Write([]byte(message.ToString()))

}

func connectUDPServer(ip string, port int) (*net.UDPConn, error) {
	serverAddr := net.UDPAddr{
		IP:   net.ParseIP(ip), // Server IP address
		Port: port,            // SIP server port
	}
	conn, err := net.DialUDP("udp", nil, &serverAddr)
	if err != nil {
		return nil, err
	}

	return conn, err
}

/*
`REGISTER sip:example.com SIP/2.0
Via: SIP/2.0/UDP first.example.com;branch=z9hG4bK1
Via: SIP/2.0/UDP second.example.com;branch=z9hG4bK2
From: <sip:alice@example.com>;tag=12345
To: <sip:alice@example.com>
Call-ID: 1234567890@example.com
CSeq: 1 REGISTER
Contact: <sip:alice@client.example.com>
Max-Forwards: 70
Expires: 3600
Authorization: Digest username="alice", realm="example.com", nonce="xyz", uri="sip:example.com", response="abc123"
Content-Length: 0
User-Agent: MySIPClient/1.0
`
*/
