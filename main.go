package main

import (
	"fmt"

	"github.com/ygzaydn/golang-sip/logger"
	"github.com/ygzaydn/golang-sip/models/sip"
	"github.com/ygzaydn/golang-sip/transport/udp"
)

func main() {
	logger := logger.New(1)

	authenticationParameters := sip.AuthenticationParameters{
		Authentication: "auth",
		Schema:         "Digest",
		Algorithm:      "MD5",
	}

	serverParameters := sip.ServerParameters{
		Uri:            "sip:127.0.0.1",
		Realm:          "127.0.0.1",
		Domain:         "127.0.0.1",
		Authentication: authenticationParameters,
		ServerType:     "server",
		State:          make(map[string]sip.ClientInfo),
	}

	server, err := udp.Server("127.0.0.1", 5060, 1024, logger, serverParameters)
	if err != nil {
		fmt.Println(err)
	}

	clientCredentials := sip.ClientCredentials{
		Username: "alice",
		Password: "alice",
	}

	clientParameters := sip.ClientParameters{
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

	sipRequest := sip.NewRequest("REGISTER", sip.GenerateInitialRegisterHeaders(clientA.Entity.Address.Port, clientA.Parameters), "")

	err = clientA.SendMessage(server.Entity.Address, sipRequest)

	if err != nil {
		fmt.Println("Error sending SIP Message")
	}

	clientMsg := clientA.ReadLastMessage()
	_ = server.ReadLastMessage()

	//fmt.Println(clientMsg, serverMsg)

	if clientMsg.StatusCode == 401 {
		requestHeaders := clientMsg.Headers
		requestHeaders["Authorization"] = []string{"Digest username=\"alice\", realm=\"127.0.0.1\", nonce=\"xyz\", uri=\"sip:127.0.0.1\", response=\"abc123\""}
		delete(requestHeaders, "WWW-Authenticate")

		sipRequest = sip.NewRequest("REGISTER", requestHeaders, "")
		err = clientA.SendMessage(server.Entity.Address, sipRequest)

		if err != nil {
			fmt.Println("Error sending SIP Message")
		}

		_ = clientA.ReadLastMessage()
		_ = server.ReadLastMessage()

		// fmt.Println(clientMsg, serverMsg)
	}

	//fmt.Println(server.Parameters.State["sip:alice@127.0.0.1"].IsRegistered)

	defer clientA.Entity.Connection.Close()
	defer server.Entity.Connection.Close()

}
