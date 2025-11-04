package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
	"wizard-duel-distributed/api"
	"wizard-duel-distributed/communication"
	"wizard-duel-distributed/models"
	"wizard-duel-distributed/utils"
)

var SERVERNAME string
var NETIP string

const SERVERPREFIX = 6                                   // quantidade de letras no prefixo do servername
var SERVERHEALTH map[string]bool = make(map[string]bool) // SERVERNAME: isAlive
var DEFAULTPORT = ":8080"
var LOGSPATH = "logs/logs.json"
var PLAYERSPATH = "database/players.json"
var MATCHESPATH = "database/matches.json"
var CARDSPATH = "database/cards.json"
var TRADESPATH = "database/trades.json"
var QUEUEPATH = "database/queue.json" // é um vetor de inteiro
var COMMANDQUEUE = make(utils.PriorityQueue, 0)
var MAPMUTEX sync.Mutex
var QUEUEMUTEX sync.Mutex
var DATAMUTEX sync.Mutex
var CHECKEDALL = false

func checkPeerHealth(peerAddr string) {
	for {
		resp, err := http.Get(peerAddr + "/api/checkhealth")

		MAPMUTEX.Lock()
		if err != nil || resp.Status != "200 OK" {
			SERVERHEALTH[peerAddr] = false
			if peerAddr == "172.16.201.7" {
				fmt.Println("[debug] - Unable to connect with peer: ", peerAddr)
				fmt.Println("[debug] - err ", err)
				if err == nil {
					fmt.Println("[debug] - resp ", resp.Status)
				}
			}
		} else {
			var msg api.Message
			err := json.NewDecoder(resp.Body).Decode(&msg)
			if err == nil && msg.Type == "ACK" {
				SERVERHEALTH[peerAddr] = true
				respBody := []byte{}
				resp.Body.Read(respBody)

				fmt.Println("[debug] - ", peerAddr, " - ", resp.Status, " - ", string(respBody))

			}
		}
		MAPMUTEX.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func runCommand(command api.Command) {
	UpdateLogs(LOGSPATH, command)
	switch command.Resource {
	case "player":
		players := models.RetrievePlayers(PLAYERSPATH)
		var player models.Player
		json.Unmarshal(command.Value, &player)
		found := models.RetrievePlayerByName(player.Username, players)
		if found == nil {
			players = append(players, player)
		} else {
			models.UpdatePlayer(player.Password, player, players)
		}
		models.SavePlayers(PLAYERSPATH, players)
	case "match":
		matches := models.RetrieveMatches(MATCHESPATH)
		var match models.Match
		json.Unmarshal(command.Value, &match)
		found := models.RetrieveMatch(match.Id, matches)
		if found == nil {
			matches = append(matches, match)
		} else {
			models.UpdateMatch(match, matches)
		}
		models.SaveMatches(MATCHESPATH, matches)
	case "trade":
		trades := models.RetrieveTrades(TRADESPATH)
		var trade models.Trade
		json.Unmarshal(command.Value, &trade)
		found := models.RetrieveTrade(trade.Id, trades)
		if found == nil {
			trades = append(trades, trade)
		} else {
			models.UpdateTrade(trade, trades)
		}
		models.SaveTrades(TRADESPATH, trades)
	case "queue", "players":
		path := QUEUEPATH
		if command.Resource == "players" {
			path = PLAYERSPATH
		}
		file, err := os.Create(path)
		if err == nil {
			file.Write(command.Value)
			file.Close()
		}
	default:
		fmt.Println("[debug]: Unknown Command")

	}

}

func executeCommands() {
	for {
		if len(COMMANDQUEUE) > 0 {
			Request(COMMANDQUEUE.Front().TimeStamp)
			fmt.Println("[debug] executing a command")
			c := COMMANDQUEUE.Pop()
			switch c.Operation {
			case "signup":
				var cred communication.Credentials
				err := json.Unmarshal(c.Value, &cred)
				if err == nil {
					bytes := Signup(cred)
					if bytes != nil {
						c.Value = *bytes
						c.Operation = "create"
						propagate(*c)
					}
				}
			case "buy":
				var cred communication.Credentials
				err := json.Unmarshal(c.Value, &cred)
				if err == nil {
					bytes := Buy(cred)
					if bytes != nil {
						c.Value = *bytes
						c.Operation = "update"
						propagate(*c)
					}
				}
			case "enqueue":
				var cred communication.Credentials
				err := json.Unmarshal(c.Value, &cred)
				if err == nil {
					paired, bytes := Enqueue(cred)
					if bytes != nil {
						c.Value = *bytes
						if paired {
							var match models.Match
							json.Unmarshal(*bytes, &match)
							c.Operation = "create"
							c.Resource = "match"
							c.ResourceID = fmt.Sprint(match.Id)
							c.ID = fmt.Sprintf("%d%d", match.Id, c.TimeStamp)
							propagate(*c)
							queue := []int{}
							*bytes, _ = json.Marshal(queue)
							c.Value = *bytes
							propagate(*c)
						} else {
							c.Operation = "update"
							propagate(*c)
						}
					}
				}
			case "playCard":
				var msg communication.MatchMessage
				err := json.Unmarshal(c.Value, &msg)
				if err == nil {
					bytes := PlayCard(msg)
					if bytes != nil {
						c.Operation = "update"
						c.Value = *bytes
						propagate(*c)
					}
				}
			case "surrender":
				var msg communication.MatchMessage
				err := json.Unmarshal(c.Value, &msg)
				if err == nil {
					bytes := Surrender(msg)
					if bytes != nil {
						c.Operation = "update"
						c.Value = *bytes
						propagate(*c)
					}
				}
			case "createTrade":
				var tradeMsg communication.TradeMessage
				err := json.Unmarshal(c.Value, &tradeMsg)
				if err == nil {
					bytes := CreateTrade(tradeMsg)
					if bytes != nil {
						var trade models.Trade
						json.Unmarshal(*bytes, &trade)
						c.ID = fmt.Sprintf("%d%d", trade.Id, c.TimeStamp)
						c.ResourceID = fmt.Sprint(tradeMsg.TradeID)
						c.Operation = "create"
						c.Value = *bytes
						propagate(*c)
					}
				}
			// adicionar o commandID e o resource ID

			case "acceptTrade":
				var tradeMsg communication.TradeMessage
				err := json.Unmarshal(c.Value, &tradeMsg)
				if err == nil {
					bytes := AcceptTrade(tradeMsg)
					if bytes != nil {
						var trade models.Trade
						json.Unmarshal(*bytes, &trade)
						c.Operation = "update"
						c.Value = *bytes
						propagate(*c)
						c.Resource = "Players"
						c.ResourceID = "players"
						c.ID = fmt.Sprintf("%s%d", c.ResourceID, c.TimeStamp)
						players := models.RetrievePlayers(PLAYERSPATH)
						*bytes, _ = json.Marshal(players)
						c.Value = *bytes
						propagate(*c)
					}
				}
			case "denyTrade":
				var tradeMsg communication.TradeMessage
				err := json.Unmarshal(c.Value, &tradeMsg)
				if err == nil {
					bytes := DenyTrade(tradeMsg)
					if bytes != nil {
						c.Operation = "update"
						c.Value = *bytes
						propagate(*c)
					}
				}

			case "suggestTrade":
				var tradeMsg communication.TradeMessage
				err := json.Unmarshal(c.Value, &tradeMsg)
				if err == nil {
					bytes := SuggestTrade(tradeMsg)
					if bytes != nil {
						c.Operation = "update"
						c.Value = *bytes
						propagate(*c)
					}
				}
			default:
				fmt.Println("Uknown Command")
			}
			// propagando informação
			ONCRITICALREGION = false
		}
	}
}

func propagate(command api.Command) {
	UpdateLogs(LOGSPATH, command)
	fmt.Println("propadando")
	MAPMUTEX.Lock()
	for peer, alive := range SERVERHEALTH {
		if alive {
			com, _ := json.Marshal(command)
			_, err := http.Post(peer+"/api/update", "application/json", bytes.NewBuffer(com))
			fmt.Println("err; ", err)
		}
	}
	MAPMUTEX.Unlock()
}

func subscribeChannels(broker net.Conn) bool {
	topics := []string{"login", "signup", "buy",
		"createTrade", "acceptTrade", "tradableCards", "denyTrade",
		"suggestTrade", "enqueue", "playCard", "getMatchData", "surrender", "getCards"}

	for _, topic := range topics {
		message := communication.Message{
			Cmd: communication.SUBSCRIBE,
			Tpc: topic,
		}
		err := communication.SendMessage(broker, message)
		if err != nil {
			return false
		}
	}

	return true
}

func topicHandler(broker net.Conn) {
	for {
		var message communication.Message
		err := communication.ReceiveMessage(broker, &message)
		now := time.Now().UnixMilli()
		if err != nil {

			log.Fatal(err)
		}
		switch message.Cmd {
		case "getCards":
			DATAMUTEX.Lock()
			var cred communication.Credentials
			err := communication.UnmarshalMessage(message.Msg, &cred)
			if err == nil {
				cards := []models.Card{}
				players := models.RetrievePlayers(PLAYERSPATH)
				allCards := models.RetrieveCards(CARDSPATH)
				player := models.RetrievePlayerByName(cred.Username, players)
				for _, c := range player.Cards {
					card := models.RetrieveCard(c, allCards)
					if card != nil {
						cards = append(cards, *card)
					}
				}
				bytes, _ := json.Marshal(cards)
				msg := communication.Message{
					Cmd: communication.PUBLISH,
					Tpc: cred.Username,
					Msg: bytes,
				}
				communication.SendMessage(broker, msg)

			}
			DATAMUTEX.Unlock()
		case "login":
			DATAMUTEX.Lock()
			var cred communication.Credentials
			err := communication.UnmarshalMessage(message.Msg, &cred)
			if err == nil {
				players := models.RetrievePlayers(PLAYERSPATH)
				player := models.RetrievePlayerByName(cred.Username, players)
				ok := player != nil
				bytes, _ := json.Marshal(ok)
				msg := communication.Message{
					Cmd: communication.PUBLISH,
					Tpc: cred.Username,
					Msg: bytes,
				}
				communication.SendMessage(broker, msg)
			}
			DATAMUTEX.Unlock()

		case "tradableCards":
			DATAMUTEX.Lock()
			var cred communication.Credentials
			err := communication.UnmarshalMessage(message.Msg, &cred)
			if err == nil {
				trades := models.RetrieveTrades(TRADESPATH)
				openTrades := models.RetrieveOpenTrades(trades)
				bytes, _ := json.Marshal(openTrades)
				msg := communication.Message{
					Cmd: communication.PUBLISH,
					Tpc: cred.Username,
					Msg: bytes,
				}
				communication.SendMessage(broker, msg)
			}
			DATAMUTEX.Unlock()

		case "getMatchData":
			DATAMUTEX.Lock()
			var msg communication.MatchMessage
			err := communication.UnmarshalMessage(message.Msg, &msg)
			if err == nil {
				matches := models.RetrieveMatches(MATCHESPATH)
				match := models.RetrieveMatch(msg.MatchID, matches)
				if match != nil {
					bytes, _ := json.Marshal(*match)
					msg := communication.Message{
						Cmd: communication.PUBLISH,
						Tpc: msg.Credentials.Username,
						Msg: bytes,
					}
					communication.SendMessage(broker, msg)
				}
			}
			DATAMUTEX.Unlock()

		case "signup", "buy", "enqueue":
			QUEUEMUTEX.Lock()
			var cred communication.Credentials
			err := communication.UnmarshalMessage(message.Msg, &cred)
			if err == nil {
				command := api.Command{
					ID:         cred.Username + fmt.Sprint(now),
					ResourceID: cred.Username,
					NodeID:     SERVERNAME,
					TimeStamp:  now,
					Operation:  message.Cmd,
					Resource:   "player",
					Value:      message.Msg,
				}
				if message.Cmd == "enqueue" {
					command.Resource = "queue"
				}
				COMMANDQUEUE.Push(command)
			}
			QUEUEMUTEX.Unlock()
		case "createTrade", "acceptTrade", "denyTrade", "suggestTrade":
			QUEUEMUTEX.Lock()
			var trademMessage communication.TradeMessage
			err := communication.UnmarshalMessage(message.Msg, &trademMessage)
			if err == nil {
				command := api.Command{
					NodeID:    SERVERNAME,
					TimeStamp: now,
					Operation: message.Cmd,
					Value:     message.Msg,
					Resource:  "trade",
				}
				if message.Cmd != "createTrade" {
					command.ResourceID = fmt.Sprint(trademMessage.TradeID)
					command.ID = fmt.Sprintf("%d%d", trademMessage.TradeID, now)
				}

				COMMANDQUEUE.Push(command)
			}
			QUEUEMUTEX.Unlock()

		case "playCard", "surrender":
			QUEUEMUTEX.Lock()
			var matchMessage communication.MatchMessage
			err := communication.UnmarshalMessage(message.Msg, &matchMessage)
			if err == nil {
				command := api.Command{
					ID:         fmt.Sprintf("%d%d", matchMessage.MatchID, now),
					ResourceID: fmt.Sprint(matchMessage.MatchID),
					TimeStamp:  now,
					NodeID:     SERVERNAME,
					Operation:  message.Cmd,
					Value:      message.Msg,
					Resource:   "match",
				}
				COMMANDQUEUE.Push(command)
			}
			QUEUEMUTEX.Unlock()
		}
	}
}

func main() {
	SERVERNAME = utils.GetSelfAddres()
	fmt.Println(SERVERNAME)
	NETIP = utils.GetNetworkAddress(SERVERNAME)

	// quando um server inicia, ele procura por todos os servidores de 0 a 10 e adiciona no SERVERHEALTH
	fmt.Println("Server is starting")
	for i := range 255 {
		var peername = fmt.Sprintf("%s.%d", NETIP, i)
		fmt.Println(peername)
		if peername != SERVERNAME {
			go checkPeerHealth("http://" + peername + DEFAULTPORT)
		}
	}

	time.Sleep(time.Second * 5)
	MAPMUTEX.Lock()
	for peer, alive := range SERVERHEALTH {
		if alive {
			fmt.Println("Syncing logs = ", peer)
			var logs, err = os.ReadFile(LOGSPATH)
			if err != nil {
				logs, _ = json.Marshal([]api.Command{}) // inicia um vetor vazio caso não consiga abrir o arquivo de logs
			}
			resp, err := http.Post(peer+"/api/sync", "application/json", bytes.NewBuffer(logs))
			fmt.Println("error -> ", err)

			if err == nil {
				fmt.Println("resp status -> ", resp.Status)
				logs := []api.Command{}
				json.NewDecoder(resp.Body).Decode(&logs)
				fmt.Println("logs: ", logs)
				for _, log := range logs {
					fmt.Println("running command: ", log)
					runCommand(log)
				}
			}
			break // so precisa sincronizar com um servidor
		}
	}
	MAPMUTEX.Unlock()
	go executeCommands()
	broker, err := net.Dial("tcp", SERVERNAME+communication.BROKERPORT)
	if err != nil {
		fmt.Println("[error] - impossible to connect to broker")
		return
	}
	ok := subscribeChannels(broker)
	if !ok {
		fmt.Println("[error] - couldn'd subscribe in the topics")
		return
	}
	go topicHandler(broker)
	handleRequests()
}
