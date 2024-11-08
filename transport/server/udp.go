package server

import (
	"errors"
	"fmt"
	"net"

	"github.com/ygzaydn/golang-sip/logger"
	"github.com/ygzaydn/golang-sip/models/sip"
	"github.com/ygzaydn/golang-sip/utils"
)

func UDPEngine(ip string, port int, bufferSize int, logger *logger.Logger) {
	udpServer, err := createUDPServer(ip, port)
	if err != nil {
		panic(err)
	}
	go udpListener(udpServer, bufferSize, logger)

}

func createUDPServer(ip string, port int) (*net.UDPConn, error) {

	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}

	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		return nil, errors.New("error creating UDP server")
	}

	fmt.Printf("UDP server listening on port %d...\n", port)

	return conn, err
}

func udpListener(conn *net.UDPConn, bufferSize int, logger *logger.Logger) {

	// Not sure if I should make bufferSize as a parameter
	buffer := make([]byte, bufferSize)
	defer conn.Close()
	for {

		n, clientAddr, err := conn.ReadFromUDP(buffer)

		if err != nil {
			fmt.Println("Error reading from client:", err)
			continue
		}

		msg := string(buffer[:n])
		isValid := sip.ISSIPMessage(msg)

		if !isValid {
			logger.BuildLogMessage("Server received a message, but format is wrong, message skipped.")
			continue
		}

		message, err := sip.ToSIP(msg)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if logger != nil {
			logger.BuildLogMessage("Server Received\t- " + utils.FormatLogMessage(message.Method))
		}

		// Example response - will change it later on
		response := []byte("Message received!")

		_, err = conn.WriteToUDP(response, clientAddr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}

	}

}
