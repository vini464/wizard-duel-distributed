package models

import (
	"encoding/json"
	"os"
)

func CreatePlayer(username, password string, players *[]Player) *Player {
	for _, p := range *players { // procura por algum username já registrado
		if p.Username == username {
			return nil
		}
	}
	player := Player{
		Username: username,
		Password: password,
		Cards:    []int{},
		Coins:    10, // padrão, da pra comprar 2 boosters
	}
	*players = append(*players, player)
	return &player
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

func RetrievePlayerByName(username string, players []Player) *Player {
	for _, p := range players {
		if p.Username == username {
			return &p
		}
	}
	return nil
}

func UpdatePlayer(password string, newdata Player, players []Player) []Player {
	index := -1
	for i, p := range players {
		if p.Username == newdata.Username {
			if password == p.Password {
				index = i
				break
			}
		}
	}
	if index >= 0 {
		players = append(players[:index], players[index+1:]...)
		players = append(players, newdata)
	}
	return players
}

func DeletePlayer(username, password string, players []Player) []Player {
	index := -1
	for i, player := range players {
		if player.Username == username {
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

func SavePlayers(filepath string, players []Player) bool {
	file, err := os.Create(filepath)
	if err != nil {
		return false
	}
	defer file.Close()
	bytes, err := json.Marshal(players)
	if err != nil {
		bytes, _ = json.Marshal([]Player{})
	}
	_, err = file.Write(bytes)
	return err == nil
}
