package main

import (
	"fmt"

	"github.com/ygzaydn/golang-sip/logger"
	"github.com/ygzaydn/golang-sip/models/sip"
	"github.com/ygzaydn/golang-sip/transport/client"
	"github.com/ygzaydn/golang-sip/transport/server"
)

func main() {
	logger := logger.New(1)
	server.UDPEngine("127.0.0.1", 5060, 1024, logger)
	clientA, err := client.New("127.0.0.1", 5060, logger)
	if err != nil {
		fmt.Println(err)
	}

	requestHeaders := map[string][]string{
		"Via": {
			"SIP/2.0/UDP first.example.com;branch=z9hG4bK1",
			"SIP/2.0/UDP second.example.com;branch=z9hG4bK2",
		},
		"From":           {"<sip:alice@example.com>;tag=12345"},
		"To":             {"<sip:alice@example.com>"},
		"Call-ID":        {"1234567890@example.com"},
		"CSeq":           {"1 REGISTER"},
		"Contact":        {"<sip:alice@client.example.com>"},
		"Content-Length": {"0"},
		"Max-Forwards":   {"70"},
		"User-Agent":     {"MySIPClient/1.0"},
		"Expires":        {"3600"},
		"Authorization":  {"Digest username=\"alice\", realm=\"example.com\", nonce=\"xyz\", uri=\"sip:example.com\", response=\"abc123\""},
	}

	sipRequest := sip.NewRequest("REGISTER", requestHeaders, "")
	clientA.SendMessage(sipRequest)
	defer clientA.Connection.Close()
	for {
	}
}
