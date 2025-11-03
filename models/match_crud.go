package models

import (
	"encoding/json"
	"os"
)

func CreateMatch(playerA, playerB string, matches *[]Match) Match {
	id := 0
	for _, m := range *matches {
		if m.Id >= id {
			id = m.Id + 1
		}
	}
	match := Match{Id: id, Players: make(map[string]int)}
	match.Players[playerA] = 10
	match.Players[playerB] = 10

	*matches = append(*matches, match)
	return match
}

func RetrieveMatches(filepath string) []Match {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return []Match{}
	}
	var data []Match
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return []Match{}
	}
	return data
}

func RetrieveMatch(id int, matches []Match) *Match {
	for _, m := range matches {
		if m.Id == id {
			return &m
		}
	}
	return nil
}

// retorna a partida não finalizada
func RetrieveMatchByPlayer(player string, matches []Match) *Match {
	for _, m := range matches {
		_, ok := m.Players[player]
		if !m.Over && ok {
			return &m
		}
	}
	return nil
}

func UpdateMatch(newdata Match, matches []Match) []Match {
	index := -1
	for i, m := range matches {
		if m.Id == newdata.Id {
			index = i
		}
	}
	if index >= 0 {
		matches = append(matches[:index], matches[index+1:]...)
		matches = append(matches, newdata)
	}
	return matches
}

func SaveMatches(filepath string, matches[]Match) bool {
	file, err := os.Create(filepath)
	if err != nil {
		return false
	}
	defer file.Close()
	bytes, err := json.Marshal(matches)
	if err != nil {
		bytes, _ = json.Marshal([]Match{})
	}
	_, err = file.Write(bytes)
	return err == nil
}

// Não preciso de um DeleteMatch é importante manter o historico de partidas
