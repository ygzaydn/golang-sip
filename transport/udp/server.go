package udp

import (
	"errors"
	"fmt"
	"net"

	"github.com/ygzaydn/golang-sip/logger"
	"github.com/ygzaydn/golang-sip/models/sip"
	"github.com/ygzaydn/golang-sip/utils"
)

func Server(ip string, port int, bufferSize int, logger *logger.Logger, serverParameters sip.ServerParameters) (*UDPServer, error) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}

	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		return nil, errors.New("error creating UDP server")
	}

	fmt.Printf("UDP server listening on port %d...\n", port)

	server := &UDPServer{
		Entity: UDPEntity{
			Connection:     conn,
			logger:         logger,
			entityType:     "server",
			Address:        &addr,
			LastMessage:    nil,
			MessageChannel: make(chan *sip.SIPMessage, 50),
		},
		Parameters: serverParameters,
	}

	go server.udpListener(bufferSize)

	return server, err
}

func (u *UDPServer) udpListener(bufferSize int) {

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
			message.ServerHandler(messageChannel, u.Parameters)
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
			err = u.checkState(response)
			if err != nil {
				fmt.Println("Error on state:", err)
				continue
			}
			err = u.updateState(response)
			if err != nil {
				fmt.Println("Error updating state:", err)
				continue
			}
			err = u.SendMessage(clientAddr, response)
			if err != nil {
				fmt.Println("Error sending response:", err)
				continue
			}
		}

	}
}

func (u *UDPServer) SendMessage(address *net.UDPAddr, message *sip.SIPMessage) error {
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

func (u *UDPServer) ReadLastMessage() *sip.SIPMessage {
	for value := range u.Entity.MessageChannel {

		if value != nil {
			u.Entity.LastMessage = value
		}
	}

	u.Entity.MessageChannel = make(chan *sip.SIPMessage, 50)
	return u.Entity.LastMessage
}

func (u *UDPServer) updateState(parsedMessage *sip.SIPMessage) error {
	//TODO

	switch parsedMessage.StatusCode {

	case 401:
		contact, err := utils.ParseFromandToHeader(parsedMessage.Headers["From"][0])
		if err != nil {
			return err
		}

		CSeq, err := utils.ParseCSeqHeader(parsedMessage.Headers["CSeq"][0])
		if err != nil {
			return err
		}
		var newState sip.ClientInfo
		username := fmt.Sprintf("sip:%s@%s", contact["User"].(string), contact["Host"].(string))

		authenticate, err := utils.ParseWWWAuthenticateandAuthorizationHeader(parsedMessage.Headers["WWW-Authenticate"][0])
		if err != nil {
			return err
		}

		isPresent := u.Parameters.State[username].IsPresent

		if isPresent {
			newState = u.Parameters.State[username]
			u.Parameters.State[username] = newState
		} else {
			newState.IsPresent = true
			newState.IsRegistered = false
			newState.Contact = contact["User"].(string)
			newState.AuthToken = ""
			newState.TransportType = "UDP"
			newState.CSeq = CSeq["CSeq"].(int)
			newState.Nonce = authenticate["Nonce"].(string)
			newState.Opaque = authenticate["Opaque"].(string)
			u.Parameters.State[username] = newState
		}
	case 200:
		CSeq, err := utils.ParseCSeqHeader(parsedMessage.Headers["CSeq"][0])
		if err != nil {
			return err
		}

		if CSeq["Method"] == "REGISTER" {

			var newState sip.ClientInfo
			contact, err := utils.ParseFromandToHeader(parsedMessage.Headers["From"][0])
			if err != nil {
				return err
			}

			username := fmt.Sprintf("sip:%s@%s", contact["User"].(string), contact["Host"].(string))

			isPresent := u.Parameters.State[username].IsPresent

			if isPresent {

				newState = u.Parameters.State[username]
				newState.IsRegistered = true
				u.Parameters.State[username] = newState

			} else {
				newState.IsPresent = true
				newState.IsRegistered = true
				newState.Contact = contact["User"].(string)
				newState.AuthToken = ""
				newState.TransportType = "UDP"
				newState.CSeq = CSeq["CSeq"].(int)
				u.Parameters.State[username] = newState
			}
		}

	}

	return nil
}

func (u *UDPServer) checkState(parsedMessage *sip.SIPMessage) error {
	//TODO
	return nil
}
