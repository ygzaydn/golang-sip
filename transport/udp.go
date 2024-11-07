package transport

import (
	"errors"
	"fmt"
	"net"
)

func UDPEngine(ip string, port int, bufferSize int) {
	udpServer, err := createUDPServer(ip, port)
	if err != nil {
		panic(err)
	}
	go udpListener(udpServer, bufferSize)

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

func udpListener(conn *net.UDPConn, bufferSize int) {
	buffer := make([]byte, bufferSize)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)

		if err != nil {
			fmt.Println("Error reading from client:", err)
			break
		}
		fmt.Printf("Received %s from %s\n", string(buffer[:n]), clientAddr)

		response := []byte("Message received!")

		_, err = conn.WriteToUDP(response, clientAddr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}
	}
	defer conn.Close()
}
