package main

import (
	"fmt"

	"github.com/ygzaydn/golang-sip/logger"
	"github.com/ygzaydn/golang-sip/models/sip"
	"github.com/ygzaydn/golang-sip/transport/udp"
)

func main() {
	logger := logger.New(1)
	server, err := udp.New("server", "127.0.0.1", 5060, 1024, logger)
	if err != nil {
		fmt.Println(err)
	}

	clientA, err := udp.New("client", "127.0.0.1", 5065, 1024, logger)
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
		//"Authorization":  {"Digest username=\"alice\", realm=\"example.com\", nonce=\"xyz\", uri=\"sip:example.com\", response=\"abc123\""},
	}

	sipRequest := sip.NewRequest("REGISTER", requestHeaders, "")

	err = clientA.SendMessage(server.Address, sipRequest)

	if err != nil {
		fmt.Println("Error sending SIP Message")
	}

	server.ReadLastMessage()

	requestHeaders = map[string][]string{
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

	sipRequest = sip.NewRequest("REGISTER", requestHeaders, "")
	err = clientA.SendMessage(server.Address, sipRequest)

	if err != nil {
		fmt.Println("Error sending SIP Message")
	}

	server.ReadLastMessage()

	defer clientA.Connection.Close()
	defer server.Connection.Close()

}
