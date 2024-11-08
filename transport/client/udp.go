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
	u.logger.BuildLogMessage("Client Sent\t- " + utils.FormatLogMessage(message.Method))
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
