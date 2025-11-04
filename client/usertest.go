package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"wizard-duel-distributed/communication"
	"wizard-duel-distributed/utils"
)

var HOSTNAME = utils.GetSelfAddres()

func getCredentials() communication.Credentials {

	username := utils.Input("Digite seu username: ")
	password := utils.Input("Digite seu password: ")
	return communication.Credentials{
		Username: username,
		Password: password,
	}
}

func createUser(credentials communication.Credentials) communication.Message {
	bytes, _ := json.Marshal(credentials)
	return communication.Message{
		Cmd: communication.PUBLISH,
		Tpc: "signup",
		Msg: bytes,
	}
}

func getCards(credentials communication.Credentials) communication.Message {
	bytes, _ := json.Marshal(credentials)
	return communication.Message{
		Cmd: communication.PUBLISH,
		Tpc: "getCards",
		Msg: bytes,
	}
}

func createTrade(credentials communication.Credentials, cardId int) communication.Message {
	msg := communication.TradeMessage{
		Credentials: credentials,
		CardID:      cardId,
	}
	bytes, _ := json.Marshal(msg)
	return communication.Message{
		Cmd: communication.PUBLISH,
		Tpc: "createTrade",
		Msg: bytes,
	}
}

func getTrades(credentials communication.Credentials) communication.Message {
	bytes, _ := json.Marshal(credentials)
	return communication.Message{
		Cmd: communication.PUBLISH,
		Tpc: "tradableCards",
		Msg: bytes,
	}
}

func getBooster(credentials communication.Credentials) communication.Message {
	bytes, _ := json.Marshal(credentials)
	return communication.Message{
		Cmd: communication.PUBLISH,
		Tpc: "buy",
		Msg: bytes,
	}
}

func main() {
	conn, err := net.Dial("tcp", HOSTNAME+communication.BROKERPORT)
	if err != nil {
		panic(err)
	}
	cred := getCredentials()
	msg := communication.Message {
		Tpc: cred.Username,
		Cmd: communication.SUBSCRIBE,
	}
	communication.SendMessage(conn, msg) // me inscrevendo no topico
	go func() {
		bytes := communication.ReceiveBytes(conn)
		fmt.Println("Received: ", string(bytes))
	}()

	for {
		test := utils.Input("o que testar? ")
		switch test {
		case "create":
			communication.SendMessage(conn, createUser(cred))
		case "cards":
			communication.SendMessage(conn, getCards(cred))
		case "ctrade":
			c := utils.Input("id: ")
			cardId, _ := strconv.ParseInt(c, 10, 0)
			communication.SendMessage(conn, createTrade(cred, int(cardId)))
		case "trades":
			communication.SendMessage(conn, getTrades(cred))
		case "buy":
			communication.SendMessage(conn, getBooster(cred))
		}
	}
}
