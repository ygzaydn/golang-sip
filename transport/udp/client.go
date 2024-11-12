package udp

import (
	"errors"
	"fmt"
	"net"

	"github.com/ygzaydn/golang-sip/logger"
	"github.com/ygzaydn/golang-sip/models/sip"
	"github.com/ygzaydn/golang-sip/utils"
)

func Client(ip string, port int, bufferSize int, logger *logger.Logger, clientParameters ClientParameters) (*UDPClient, error) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}

	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		return nil, errors.New("error creating UDP client")
	}

	fmt.Printf("UDP client listening on port %d...\n", port)

	client := &UDPClient{
		Entity: UDPEntity{
			Connection:     conn,
			logger:         logger,
			entityType:     "client",
			Address:        &addr,
			LastMessage:    nil,
			MessageChannel: make(chan *sip.SIPMessage, 50),
		},
		Parameters: clientParameters,
	}

	go client.udpListener(bufferSize)
	return client, err
}

func (u *UDPClient) udpListener(bufferSize int) {

	// Not sure if I should make bufferSize as a parameter
	buffer := make([]byte, bufferSize)

	defer u.Entity.Connection.Close()

	for {

		n, clientAddr, err := u.Entity.Connection.ReadFromUDP(buffer)

		if err != nil {
			fmt.Println("Error reading: ", err)
			continue
		}

		msg := string(buffer[:n])
		isValid := sip.ISSIPMessage(msg)

		if !isValid {
			u.Entity.logger.BuildLogMessage("Server received a message, but format is wrong, message skipped.")
			continue
		}

		message, err := sip.ToSIP(msg)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if u.Entity.logger != nil {
			if message.StatusCode != 0 {
				u.Entity.logger.BuildLogMessage(u.Entity.entityType + " received \t- " + fmt.Sprint(message.StatusCode) + " " + utils.FormatLogMessage(message.Reason))
			} else {
				u.Entity.logger.BuildLogMessage(u.Entity.entityType + " received \t- " + utils.FormatLogMessage(message.Method))
			}
		}

		//responses := message.HandleRequest()

		messageChannel := make(chan *sip.SIPMessage, 50)
		go func() {
			message.HandleRequest(messageChannel)
			close(messageChannel)
		}()

		if message != nil {

			u.Entity.MessageChannel <- message

			if u.Entity.entityType == "client" {
				// fmt.Println(clientAddr.String() + " send msg -> " + message.Method + " " + message.Reason)
				if message.ShouldCloseResponseChannel() {
					// fmt.Println(clientAddr.String() + " closed")
					close(u.Entity.MessageChannel)
				}
			} else if u.Entity.entityType == "server" {
				if message.ShouldCloseRequestChannel() {
					// fmt.Println(clientAddr.String() + " closed")
					close(u.Entity.MessageChannel)
				}
			}

		}

		// if len(responses) < 1 {
		// 	continue
		// }

		for response := range messageChannel {
			err = u.SendMessage(clientAddr, response)
			if err != nil {
				fmt.Println("Error sending response:", err)
				continue
			}

		}

	}
}

func (u *UDPClient) SendMessage(address *net.UDPAddr, message *sip.SIPMessage) error {
	var err error
	if message != nil {
		if u.Entity.logger != nil {
			if message.StatusCode != 0 {
				u.Entity.logger.BuildLogMessage(u.Entity.entityType + " sent \t- " + fmt.Sprint(message.StatusCode) + " " + utils.FormatLogMessage(message.Reason))
			} else {
				u.Entity.logger.BuildLogMessage(u.Entity.entityType + " sent\t- " + utils.FormatLogMessage(message.Method))
			}

		}
		_, err = u.Entity.Connection.WriteToUDP([]byte(message.ToString()), address)
	}

	return err
}

func (u *UDPClient) ReadLastMessage() *sip.SIPMessage {

	for value := range u.Entity.MessageChannel {

		if value != nil {
			u.Entity.LastMessage = value
		}
	}

	u.Entity.MessageChannel = make(chan *sip.SIPMessage, 50)
	return u.Entity.LastMessage
}
