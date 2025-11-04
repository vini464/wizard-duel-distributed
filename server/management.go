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
			} else {
				rnd -= card.Quantity
			}
		}
	}
	for _, card := range booster {
		player.Cards = append(player.Cards, card.Id)
	}
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
	if trade == nil {
		return nil
	}
	id := -1
	for _, c := range player.Cards {
		if c == trade.CardA {
			id = c
			break
		}
	}
	if id >= 0 {
		player.Cards = append(player.Cards[:id], player.Cards[id+1:]...)
		player.Cards = append(player.Cards, trade.CardB)
	}

	id = 0
	p2 := models.RetrievePlayerByName(trade.PlayerB, players)
	if id >= 0 {
		p2.Cards = append(p2.Cards[:id], p2.Cards[id+1:]...)
		p2.Cards = append(p2.Cards, trade.CardA)
	}

	trade.Accepted = true
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
	if trade == nil {
		return nil
	}
	*trade = models.Trade{
		Id:       trade.Id,
		Accepted: false,
		PlayerA:  trade.PlayerA,
		CardA:    trade.CardA,
	}
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
	if trade == nil {
		return nil
	}
	trade.CardB = msg.CardID
	trade.PlayerB = msg.Credentials.Username

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

	if player == nil || match == nil || match.Turn != player.Username {
		return nil
	}

	match.Over = true
	models.SaveMatches(MATCHESPATH, matches)
	bytes, _ := json.Marshal(*match)
	return &bytes
}

func Enqueue(msg communication.Credentials) (bool, *[]byte) {
	players := models.RetrievePlayers(PLAYERSPATH)
	player := models.CreatePlayer(msg.Username, msg.Username, &players)
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
	}

	matches := models.RetrieveMatches(MATCHESPATH)
	match := models.CreateMatch(msg.Username, queue[0], &matches)
	models.SaveMatches(MATCHESPATH, matches)
	bytes, _ := json.Marshal(match)

	file, err := os.Create(QUEUEPATH)
	if err == nil {
		queue = []string{}
		queueBytes, _ = json.Marshal(queue)
		file.Write(queueBytes)
	}
	return true, &bytes
}
