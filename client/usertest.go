package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
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

func suggestTrade(credentials communication.Credentials, tradeId, cardId int) communication.Message {
	msg := communication.TradeMessage{
		Credentials: credentials,
		TradeID:     tradeId,
		CardID:      cardId,
	}
	bytes, _ := json.Marshal(msg)
	return communication.Message{
		Cmd: communication.PUBLISH,
		Tpc: "suggestTrade",
		Msg: bytes,
	}
}

func acceptTrade(credentials communication.Credentials, tradeId int) communication.Message {
	msg := communication.TradeMessage{
		Credentials: credentials,
		TradeID:     tradeId,
	}
	bytes, _ := json.Marshal(msg)
	return communication.Message{
		Cmd: communication.PUBLISH,
		Tpc: "acceptTrade",
		Msg: bytes,
	}
}
func denyTrade(credentials communication.Credentials, tradeId int) communication.Message {
	msg := communication.TradeMessage{
		Credentials: credentials,
		TradeID:     tradeId,
	}
	bytes, _ := json.Marshal(msg)
	return communication.Message{
		Cmd: communication.PUBLISH,
		Tpc: "denyTrade",
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

func match(credentials communication.Credentials) communication.Message {
	bytes, _ := json.Marshal(credentials)
	return communication.Message{
		Cmd: communication.PUBLISH,
		Tpc: "enqueue",
		Msg: bytes,
	}
}
func playCard(credentials communication.Credentials, cardId, matchId int) communication.Message {
	matchMessage := communication.MatchMessage{
		Credentials: credentials,
		MatchID:     matchId,
		CardID:      cardId,
	}
	bytes, _ := json.Marshal(matchMessage)
	return communication.Message{
		Cmd: communication.PUBLISH,
		Tpc: "playCard",
		Msg: bytes,
	}
}
func surrender(credentials communication.Credentials, matchId int) communication.Message {
	matchMessage := communication.MatchMessage{
		Credentials: credentials,
		MatchID:     matchId,
	}
	bytes, _ := json.Marshal(matchMessage)
	return communication.Message{
		Cmd: communication.PUBLISH,
		Tpc: "surrender",
		Msg: bytes,
	}
}

var IOMUTEX sync.Mutex

func main() {
	conn, err := net.Dial("tcp", HOSTNAME+communication.BROKERPORT)
	if err != nil {
		panic(err)
	}
	cred := getCredentials()
	msg := communication.Message{
		Tpc: cred.Username,
		Cmd: communication.SUBSCRIBE,
	}
	communication.SendMessage(conn, msg) // me inscrevendo no topico
	go func() {
		fmt.Println("received message")
		bytes := communication.ReceiveBytes(conn)
		IOMUTEX.Lock()
		fmt.Println("\nReceived: ", string(bytes))
		IOMUTEX.Unlock()
	}()

	for {
		IOMUTEX.Lock()
		test := utils.Input("o que testar? ")
		IOMUTEX.Unlock()
		switch test {
		case "create":
			communication.SendMessage(conn, createUser(cred))
		case "cards":
			communication.SendMessage(conn, getCards(cred))
		case "ctrade":
			c := utils.Input("id: ")
			cardId, _ := strconv.ParseInt(c, 10, 0)
			communication.SendMessage(conn, createTrade(cred, int(cardId)))
		case "buy":
			communication.SendMessage(conn, getBooster(cred))
		case "trades":
			communication.SendMessage(conn, getTrades(cred))
		case "suggest":
			c := utils.Input("id: ")
			cardId, _ := strconv.ParseInt(c, 10, 0)
			c = utils.Input("id: ")
			tradeId, _ := strconv.ParseInt(c, 10, 0)
			communication.SendMessage(conn, suggestTrade(cred, int(tradeId), int(cardId)))
		case "accept":
			c := utils.Input("id: ")
			tradeId, _ := strconv.ParseInt(c, 10, 0)
			communication.SendMessage(conn, acceptTrade(cred, int(tradeId)))
		case "deny":
			c := utils.Input("id: ")
			tradeId, _ := strconv.ParseInt(c, 10, 0)
			communication.SendMessage(conn, denyTrade(cred, int(tradeId)))
		}
	}
}
