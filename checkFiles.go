package main

import (
	"fmt"
	"os"
)

var PLAYERSPATH = "database/players.json"
var MATCHESPATH = "database/matches.json"
var CARDSPATH = "database/cards.json"
var TRADESPATH = "database/trades.json"
var QUEUEPATH = "database/queue.json" // é um vetor de inteiro

func main(){
for {
		players, err := os.ReadFile(PLAYERSPATH)
		if err == nil {
			fmt.Println("Players: ", string(players))
		}
		matches, err := os.ReadFile(MATCHESPATH)
		if err == nil {
			fmt.Println("Matches: ", string(matches))
		}
		cards, err := os.ReadFile(CARDSPATH)
		if err == nil {
			fmt.Println("Cards: ", string(cards))
		}
		trades, err := os.ReadFile(CARDSPATH)
		if err == nil {
			fmt.Println("Trades: ", string(trades))
		}
		queue, err := os.ReadFile(QUEUEPATH)
		if err == nil {
			fmt.Println("Queue: ", string(queue))
		}
	}

}
