package udp

import (
	"errors"
	"fmt"
	"net"

	"github.com/ygzaydn/golang-sip/logger"
	"github.com/ygzaydn/golang-sip/models/sip"
	"github.com/ygzaydn/golang-sip/utils"
)

type UDPEntity struct {
	Connection     *net.UDPConn
	logger         *logger.Logger
	entityType     string
	Address        *net.UDPAddr
	LastMessage    *sip.SIPMessage
	MessageChannel chan *sip.SIPMessage
	Wg             *utils.CustomWaitGroup
}

func New(entityType string, ip string, port int, bufferSize int, logger *logger.Logger) (*UDPEntity, error) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}

	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		return nil, errors.New("error creating UDP server")
	}

	fmt.Printf("UDP server listening on port %d...\n", port)

	entity := &UDPEntity{
		Connection:     conn,
		logger:         logger,
		entityType:     entityType,
		Address:        &addr,
		LastMessage:    nil,
		MessageChannel: make(chan *sip.SIPMessage, 50),
		Wg:             &utils.CustomWaitGroup{},
	}

	go entity.udpListener(bufferSize, entityType)

	return entity, err
}

func (u *UDPEntity) SendMessage(address *net.UDPAddr, message *sip.SIPMessage) error {
	var err error
	if message != nil {
		if u.logger != nil {
			if message.StatusCode != 0 {
				u.logger.BuildLogMessage(u.entityType + " sent \t- " + fmt.Sprint(message.StatusCode) + " " + utils.FormatLogMessage(message.Reason))
			} else {
				u.logger.BuildLogMessage(u.entityType + " sent\t- " + utils.FormatLogMessage(message.Method))
			}

		}
		_, err = u.Connection.WriteToUDP([]byte(message.ToString()), address)
	}

	return err
}

func (u *UDPEntity) udpListener(bufferSize int, entityType string) {

	// Not sure if I should make bufferSize as a parameter
	buffer := make([]byte, bufferSize)

	defer u.Connection.Close()

	for {

		n, clientAddr, err := u.Connection.ReadFromUDP(buffer)

		if err != nil {
			fmt.Println("Error reading: ", err)
			continue
		}

		msg := string(buffer[:n])
		isValid := sip.ISSIPMessage(msg)

		if !isValid {
			u.logger.BuildLogMessage("Server received a message, but format is wrong, message skipped.")
			continue
		}

		message, err := sip.ToSIP(msg)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if u.logger != nil {
			if message.StatusCode != 0 {
				u.logger.BuildLogMessage(entityType + " received \t- " + fmt.Sprint(message.StatusCode) + " " + utils.FormatLogMessage(message.Reason))
			} else {
				u.logger.BuildLogMessage(entityType + " received \t- " + utils.FormatLogMessage(message.Method))
			}
		}

		responses := message.HandleRequest()

		if len(responses) < 1 {
			continue
		}

		for _, response := range responses {
			err = u.SendMessage(clientAddr, response)
			if err != nil {
				fmt.Println("Error sending response:", err)
				continue
			}
			if response != nil {
				u.MessageChannel <- response
				if response.StatusCode == 200 {
					close(u.MessageChannel)
				}
				if response.StatusCode == 401 {
					close(u.MessageChannel)
				}
			}
		}

	}
}

func (u *UDPEntity) ReadLastMessage() {
	for value := range u.MessageChannel {

		if value != nil {
			u.LastMessage = value
		}
	}

	u.MessageChannel = make(chan *sip.SIPMessage, 50)
	fmt.Println("Last message: ", u.LastMessage)
}
