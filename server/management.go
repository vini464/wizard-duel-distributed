package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"slices"
	"wizard-duel-distributed/communication"
	"wizard-duel-distributed/models"
)

func Signup(credentials communication.Credentials) *[]byte {
	players := models.RetrievePlayers(PLAYERSPATH)
	player := models.CreatePlayer(credentials.Username, credentials.Password, &players)
	if player == nil {
		return nil
	}
	ok := models.SavePlayers(PLAYERSPATH, players)
	if !ok {
		return nil
	}
	bytes, _ := json.Marshal(*player)
	return &bytes
}

func Buy(credentials communication.Credentials) *[]byte {
	players := models.RetrievePlayers(PLAYERSPATH)
	player := models.RetrievePlayerByName(credentials.Username, players)
	if player == nil {
		return nil
	}
	booster := []models.Card{}
	cards := models.RetrieveCards(CARDSPATH)

	// é como se as cartas estivessem em sequencia
	for range 5 { // escolhendo 5 cards
		rnd := rand.Intn(75)
		for _, card := range cards {
			minqnt := 0
			switch card.Rarity {
			case "common":
				minqnt = 40
			case "uncommon":
				minqnt = 20
			case "rare":
				minqnt = 10
			case "legend":
				minqnt = 5
			}
			if rnd < minqnt {
				booster = append(booster, card)
				break
			} else {
				rnd -= minqnt
			}
		}
	}
	for _, card := range booster {
		player.Cards = append(player.Cards, card.Id)
	}
	models.UpdatePlayer(player.Password, *player, players)
	models.SavePlayers(PLAYERSPATH, players)
	bytes, _ := json.Marshal(player)
	return &bytes
}

func CreateTrade(msg communication.TradeMessage) *[]byte {
	players := models.RetrievePlayers(PLAYERSPATH)
	player := models.RetrievePlayerByName(msg.Credentials.Username, players)

	if player == nil {
		return nil
	}

	trades := models.RetrieveTrades(TRADESPATH)
	trade := models.CreateTrade(msg.Credentials.Username, msg.CardID, &trades)
	models.SaveTrades(TRADESPATH, trades)
	bytes, _ := json.Marshal(trade)
	return &bytes
}

func AcceptTrade(msg communication.TradeMessage) *[]byte {
	players := models.RetrievePlayers(PLAYERSPATH)
	player := models.RetrievePlayerByName(msg.Credentials.Username, players)

	if player == nil {
		return nil
	}
	trades := models.RetrieveTrades(TRADESPATH)
	trade := models.RetrieveTrade(msg.TradeID, trades)
	if trade == nil || trade.Accepted {
		return nil
	}
	p2 := models.RetrievePlayerByName(trade.PlayerB, players)
	if p2 == nil {
		return nil
	}
	idA := -1
	idB := -1
	for i, c := range player.Cards {
		if c == trade.CardA {
			idA = i
			break
		}
	}
	for i, c := range p2.Cards {
		if c == trade.CardB {
			idB = i
			break
		}
	}
	if idB < 0 || idA < 0 {
		return nil
	}
	player.Cards = append(player.Cards[:idA], player.Cards[idA+1:]...)
	player.Cards = append(player.Cards, trade.CardB)

	p2.Cards = append(p2.Cards[:idB], p2.Cards[idB+1:]...)
	p2.Cards = append(p2.Cards, trade.CardA)

	trade.Accepted = true
	trades = models.UpdateTrade(*trade, trades)
	models.UpdatePlayer(player.Password, *player, players)
	models.UpdatePlayer(p2.Password, *p2, players)
	models.SaveTrades(TRADESPATH, trades)
	models.SavePlayers(PLAYERSPATH, players)
	bytes, _ := json.Marshal(*trade)
	return &bytes
}

func DenyTrade(msg communication.TradeMessage) *[]byte {
	players := models.RetrievePlayers(PLAYERSPATH)
	player := models.RetrievePlayerByName(msg.Credentials.Username, players)

	if player == nil {
		return nil
	}
	trades := models.RetrieveTrades(TRADESPATH)
	trade := models.RetrieveTrade(msg.TradeID, trades)
	if trade == nil || trade.Accepted {
		return nil
	}
	*trade = models.Trade{
		Id:       trade.Id,
		Accepted: false,
		PlayerA:  trade.PlayerA,
		CardA:    trade.CardA,
	}
	trades = models.UpdateTrade(*trade, trades)
	models.SaveTrades(TRADESPATH, trades)
	bytes, _ := json.Marshal(*trade)
	return &bytes
}

func SuggestTrade(msg communication.TradeMessage) *[]byte {
	players := models.RetrievePlayers(PLAYERSPATH)
	player := models.RetrievePlayerByName(msg.Credentials.Username, players)

	if player == nil {
		return nil
	}
	trades := models.RetrieveTrades(TRADESPATH)
	trade := models.RetrieveTrade(msg.TradeID, trades)
	if trade == nil || trade.Accepted {
		return nil
	}
	trade.CardB = msg.CardID
	trade.PlayerB = msg.Credentials.Username

	trades = models.UpdateTrade(*trade, trades)
	models.SaveTrades(TRADESPATH, trades)
	bytes, _ := json.Marshal(*trade)
	return &bytes
}

func PlayCard(msg communication.MatchMessage) *[]byte {
	players := models.RetrievePlayers(PLAYERSPATH)
	player := models.RetrievePlayerByName(msg.Credentials.Username, players)
	cards := models.RetrieveCards(CARDSPATH)
	card := models.RetrieveCard(msg.CardID, cards)
	matches := models.RetrieveMatches(MATCHESPATH)
	match := models.RetrieveMatch(msg.MatchID, matches)

	if match.Over || card == nil || player == nil || match == nil || match.Turn != player.Username {
		return nil
	}

	hasCard := slices.Contains(player.Cards, msg.CardID)
	if !hasCard {
		return nil
	}
	op := ""
	for n := range match.Players {
		if n != player.Username {
			op = n
		}
	}
	match.Players[op] -= card.Power
	if match.Players[op] > 0 {
		match.Turn = op
	} else {
		match.Over = true
	}
	models.SaveMatches(MATCHESPATH, matches)
	bytes, _ := json.Marshal(*match)
	return &bytes
}

func Surrender(msg communication.MatchMessage) *[]byte {
	players := models.RetrievePlayers(PLAYERSPATH)
	player := models.RetrievePlayerByName(msg.Credentials.Username, players)
	matches := models.RetrieveMatches(MATCHESPATH)
	match := models.RetrieveMatch(msg.MatchID, matches)

	if match.Over || player == nil || match == nil || match.Turn != player.Username {
		return nil
	}

	match.Over = true
	models.SaveMatches(MATCHESPATH, matches)
	bytes, _ := json.Marshal(*match)
	return &bytes
}

func Enqueue(msg communication.Credentials) (bool, *[]byte) {
	players := models.RetrievePlayers(PLAYERSPATH)
	player := models.RetrievePlayerByName(msg.Username, players)
	if player == nil {
		return false, nil
	}
	queueBytes, err := os.ReadFile(QUEUEPATH)
	queue := []string{}
	if err == nil {
		err = json.Unmarshal(queueBytes, &queue)
		if err != nil {
			queue = []string{}
		}
	}
	if len(queue) == 0 {
		queue = append(queue, msg.Username)
		queueBytes, _ = json.Marshal(queue)
		file, err := os.Create(QUEUEPATH)
		if err == nil {
			file.Write(queueBytes)
			file.Close()
		}
		return false, &queueBytes
	} else if queue[0] == player.Username {
		return false, nil
	}

	matches := models.RetrieveMatches(MATCHESPATH)
	match := models.CreateMatch(msg.Username, queue[0], &matches)
	match.Turn = player.Username
	models.SaveMatches(MATCHESPATH, matches)
	bytes, _ := json.Marshal(match)

	file, err := os.Create(QUEUEPATH)
	if err == nil {
		queue = []string{}
		queueBytes, _ = json.Marshal(queue)
		file.Write(queueBytes)
		file.Close()
	}
	return true, &bytes
}
