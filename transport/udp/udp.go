package udp

import (
	"errors"
	"fmt"
	"net"
	"sync"

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
	Wg             sync.WaitGroup
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
		MessageChannel: make(chan *sip.SIPMessage),
		Wg:             sync.WaitGroup{},
	}

	go entity.udpListener(bufferSize, entityType)

	return entity, err
}

func (u *UDPEntity) SendMessage(address *net.UDPAddr, message *sip.SIPMessage) error {
	u.Wg.Add(2)
	var err error
	if message != nil {
		if u.logger != nil {
			u.logger.BuildLogMessage(u.entityType + " sent\t- " + utils.FormatLogMessage(message.Method))

		}
		_, err = u.Connection.WriteToUDP([]byte(message.ToString()), address)
	}

	u.Wg.Done()
	u.Wg.Wait()
	return err
}

func (u *UDPEntity) udpListener(bufferSize int, entityType string) {

	// Not sure if I should make bufferSize as a parameter
	buffer := make([]byte, bufferSize)
	var lastMessage *sip.SIPMessage
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

		u.Wg.Add(1)
		go func() {
			go message.HandleRequest(u.MessageChannel)

			for response := range u.MessageChannel {
				_, err = u.Connection.WriteToUDP([]byte(response.ToString()), clientAddr)
				if err != nil {
					fmt.Println("Error sending response:", err)
					continue
				}

				if u.logger != nil {
					if response.StatusCode != 0 {
						u.logger.BuildLogMessage(entityType + " sent \t- " + fmt.Sprint(response.StatusCode) + " " + utils.FormatLogMessage(response.Reason))
					} else {
						u.logger.BuildLogMessage(entityType + " sent \t- " + utils.FormatLogMessage(response.Method))
					}
				}
				lastMessage = response
			}
		}()

		u.Wg.Done()

		u.LastMessage = lastMessage
	}

}
