package main

import (
	"fmt"

	"github.com/ygzaydn/golang-sip/transport"
)

func main() {
	transport.UDPEngine("127.0.0.1", 8081, 1024)
	fmt.Printf("1234")
	for {

	}
}
