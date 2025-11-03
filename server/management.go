package main

import (
	"encoding/json"
	"math/rand"
	"wizard-duel-distributed/api"
	"wizard-duel-distributed/communication"
	"wizard-duel-distributed/models"
)

func RunCommand(command api.Command) *communication.Message {
	switch command.Operation {
	case "signin":
		var cred communication.Credentials
		err := json.Unmarshal(command.Value, &cred)
		ok := err == nil
		if ok {
			ok = signin(cred)
			msg := "ok"
			if !ok {
				msg = "err"
			}
			return &communication.Message{
				Cmd: communication.PUBLISH,
				Tpc: cred.Username,
				Msg: []byte(msg),
			}
		}
		return nil

	case "buy":
		var cred communication.Credentials
		err := json.Unmarshal(command.Value, &cred)
		ok := err == nil
		if ok {
			booster := buy(cred)
			if booster != nil {
				bytes, _ := json.Marshal(*booster)
				return &communication.Message{
					Cmd: communication.PUBLISH,
					Tpc: cred.Username,
					Msg: bytes,
				}
			}
		}
		return nil
	case "":
	}
	return nil
}

func signin(credentials communication.Credentials) bool {
	players := models.RetrievePlayers(PLAYERSPATH)
	player := models.CreatePlayer(credentials.Username, credentials.Username, &players)
	ok := models.SavePlayers(PLAYERSPATH, players)
	return player != nil && ok
}

func buy(credentials communication.Credentials) *[]models.Card {
	players := models.RetrievePlayers(PLAYERSPATH)
	player := models.RetrievePlayerByName(credentials.Username, credentials.Password, players)
	if player == nil || player.Coins < 5 {
		return nil
	}
	booster := []models.Card{}
	cards := models.RetrieveCards(CARDSPATH)
	total := 0

	for _, card := range cards {
		total += card.Quantity
	}
	// é como se as cartas estivessem em sequencia
	for range 5 { // escolhendo 5 cards
		rnd := rand.Intn(total)
		for _, card := range cards {
			if rnd < card.Quantity {
				booster = append(booster, card)

				card.Quantity--
				total--
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

				if card.Quantity < minqnt { // previne que o estoque fique vazio
					total += minqnt - card.Quantity
					card.Quantity = minqnt
				}
				break
			} else {
				rnd -= card.Quantity
			}
		}
	}
	for _, card := range booster {
		player.Cards = append(player.Cards, card.Id)
	}
	models.SavePlayers(PLAYERSPATH, players)
	return &booster
}
