package models

import (
	"encoding/json"
	"os"
)

func CreateCard(power int, cardname, rarity string, cards *[]Card) Card {
	id := 0
	for _, c := range *cards {
		if c.Id >= id {
			id = c.Id + 1
		}
	}
	card := Card{
		Id:       id,
		Power:    power,
		Cardname: cardname,
		Rarity:   rarity,
	}
	*cards = append(*cards, card)

	return card
}

func RetrieveCards(filepath string) []Card {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return []Card{}
	}
	var cards []Card
	err = json.Unmarshal(bytes, &cards)
	if err != nil {
		return []Card{}
	}
	return cards
}

func RetrieveCard(id int, cards []Card) *Card {
	for _, card := range cards {
		if card.Id == id {
			return &card
		}
	}
	return nil
}

func SaveCards(filepath string, cards []Card) bool {
	file, err := os.Create(filepath)
	if err != nil {
		return false
	}
	defer file.Close()
	bytes, err := json.Marshal(cards)
	if err != nil {
		bytes, _ = json.Marshal([]Card{})
	}
	_, err = file.Write(bytes)
	return err == nil
}

// não preciso de update e delete para cartas
