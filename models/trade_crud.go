package models

import (
	"encoding/json"
	"os"
)

func CreateTrade(userId, cardId int, trades []Trade) []Trade {
	id := 0
	for _, t := range trades {
		if t.Id >= id {
			id = t.Id + 1
		}
	}
	trade := Trade{
		Accepted: false,
		Id:       id,
		PlayerA:  userId,
		CardA:    cardId,
	}
	trades = append(trades, trade)
	return trades
}

func RetrieveTrades(filepath string) []Trade {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return []Trade{}
	}
	var data []Trade
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return []Trade{}
	}
	return data
}

func RetrieveTrade(id int, trades []Trade) *Trade {
	for _, t := range trades {
		if t.Id == id {
			return &t
		}
	}
	return nil
}

func UpdateTrade(newdata Trade, trades []Trade) []Trade {
	index := -1
	for i, t := range trades {
		if t.Id == newdata.Id {
			index = i;
			break
		} 	
	}
	if index >= 0 {
		trades = append(trades[:index], trades[index+1:]...)
		trades = append(trades, newdata)
	}
	return  trades
}


// não preciso de um delete trade, é bom manter o histórico de trocas
