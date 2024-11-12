package main

import (
	"fmt"

	"github.com/ygzaydn/golang-sip/logger"
	"github.com/ygzaydn/golang-sip/models/sip"
	"github.com/ygzaydn/golang-sip/transport/udp"
)

func main() {
	logger := logger.New(1)

	authenticationParameters := udp.AuthenticationParameters{
		Authentication: "auth",
		Schema:         "digest",
	}

	serverParameters := udp.ServerParameters{
		Uri:            "sip:example.com",
		Realm:          "example.com",
		Domain:         "example.com",
		Authentication: authenticationParameters,
		ServerType:     "server",
	}

	server, err := udp.Server("127.0.0.1", 5060, 1024, logger, serverParameters)
	if err != nil {
		fmt.Println(err)
	}

	clientCredentials := udp.ClientCredentials{
		Username: "alice",
		Password: "alice",
	}

	clientParameters := udp.ClientParameters{
		Uri:          "sip:alice@example.com",
		Realm:        "example.com",
		Domain:       "example.com",
		Credentials:  clientCredentials,
		RegistrarURI: serverParameters.Uri,
		Contact:      "<sip:alice@client.example.com:5065>",
		DisplayName:  "Alice",
	}

	clientA, err := udp.Client("127.0.0.1", 5065, 1024, logger, clientParameters)
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

	err = clientA.SendMessage(server.Entity.Address, sipRequest)

	if err != nil {
		fmt.Println("Error sending SIP Message")
	}

	clientMsg := clientA.ReadLastMessage()
	serverMsg := server.ReadLastMessage()

	fmt.Println(clientMsg, serverMsg)

	if clientMsg.StatusCode == 401 {
		requestHeaders = map[string][]string{
			"Via": {
				"SIP/2.0/UDP alice.example.com;branch=z9hG4bK1",
				"SIP/2.0/UDP second.example.com;branch=z9hG4bK2",
			},
			"From":           {"Alice <sip:alice@example.com>;tag=12345"},
			"To":             {"Alice <sip:alice@example.com>"},
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
		err = clientA.SendMessage(server.Entity.Address, sipRequest)

		if err != nil {
			fmt.Println("Error sending SIP Message")
		}

		clientMsg = clientA.ReadLastMessage()
		serverMsg = server.ReadLastMessage()

		fmt.Println(clientMsg, serverMsg)
	}

	defer clientA.Entity.Connection.Close()
	defer server.Entity.Connection.Close()

}
