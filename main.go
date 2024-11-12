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
		Uri:            "sip:127.0.0.1",
		Realm:          "127.0.0.1",
		Domain:         "127.0.0.1",
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
		Uri:          "sip:alice@127.0.0.1",
		Realm:        "127.0.0.1",
		Domain:       "127.0.0.1",
		Credentials:  clientCredentials,
		RegistrarURI: serverParameters.Uri,
		Contact:      "<sip:alice@127.0.0.1:5065>",
		DisplayName:  "Alice",
		UserAgent:    "MySIPClient/1.0",
	}

	clientA, err := udp.Client("127.0.0.1", 5065, 1024, logger, clientParameters)
	if err != nil {
		fmt.Println(err)
	}

	requestHeaders := map[string][]string{
		"Via": {
			"SIP/2.0/UDP 127.0.0.1:5065;branch=z9hG4bK776asdhds",
		},
		"From":           {"<" + clientA.Parameters.Uri + ">;tag=12345"},
		"To":             {"<" + clientA.Parameters.Uri + ">"},
		"Call-ID":        {"a84b4c76e66710@127.0.0.1"},
		"CSeq":           {"371920 REGISTER"},
		"Contact":        {clientA.Parameters.Contact},
		"Content-Length": {"0"}, // No body in this request
		"Max-Forwards":   {"70"},
		"User-Agent":     {clientA.Parameters.UserAgent},
		"Expires":        {"3600"},
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
				"SIP/2.0/UDP alice.127.0.0.1;branch=z9hG4bK1",
				"SIP/2.0/UDP second.127.0.0.1;branch=z9hG4bK2",
			},
			"From":           {"Alice <sip:alice@127.0.0.1>;tag=12345"},
			"To":             {"Alice <sip:alice@127.0.0.1>"},
			"Call-ID":        {"1234567890@127.0.0.1"},
			"CSeq":           {"1 REGISTER"},
			"Contact":        {"<sip:alice@client.127.0.0.1>"},
			"Content-Length": {"0"},
			"Max-Forwards":   {"70"},
			"User-Agent":     {"MySIPClient/1.0"},
			"Expires":        {"3600"},
			"Authorization":  {"Digest username=\"alice\", realm=\"127.0.0.1\", nonce=\"xyz\", uri=\"sip:127.0.0.1\", response=\"abc123\""},
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
