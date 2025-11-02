package models

import (
	"encoding/json"
	"os"
)

func CreatePlayer(username, password string, players []Player) []Player {
	id := 0
	for _, p := range players { // procura pelo maior ID
		if p.Id >= id {
			id = p.Id + 1
		}
	}
	player := Player{
		Id:       id,
		Username: username,
		Password: password,
		Cards:    []int{},
		Coins:    10, // padrão, da pra comprar 2 boosters
	}
	players = append(players, player)
	return players
}

func RetrievePlayers(filepath string) []Player {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return []Player{}
	}
	var data []Player
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return []Player{}
	}
	return data
}

func RetrievePlayer(id int, password string, players []Player) *Player {
	for _, p := range players {
		if p.Id == id {
			if password == p.Password {
				return &p
			}
			break
		}
	}
	return nil
}

func UpdatePlayer(id int, password string, newdata Player, players []Player) []Player {
	index := 0
	for i, p := range players {
		if p.Id == id {
			if password == p.Password {
				index = i
				break
			} else {
				return players
			}
		}
	}
	players = append(players[:index], players[index+1:]...)
	players = append(players, newdata)
	return players
}

func DeletePlayer(id int, password string, players []Player) []Player {
	index := -1
	for i, player := range players {
		if player.Id == id {
			if password == player.Password {
				index = i
			}
			break
		}
	}
	if index >= 0 {
		players = append(players[:index], players[index+1:]...)
	}
	return players
}
