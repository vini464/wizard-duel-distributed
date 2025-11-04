package main

import (
	"fmt"
	"os"
	"time"
)

var PLAYERSPATH = "database/players.json"
var MATCHESPATH = "database/matches.json"
var TRADESPATH = "database/trades.json"
var QUEUEPATH = "database/queue.json" // é um vetor de inteiro

func main() {
	for {
		players, err := os.ReadFile(PLAYERSPATH)
		if err == nil {
			fmt.Println("Players: ", string(players))
		}
		matches, err := os.ReadFile(MATCHESPATH)
		if err == nil {
			fmt.Println("Matches: ", string(matches))
		}
		trades, err := os.ReadFile(TRADESPATH)
		if err == nil {
			fmt.Println("Trades: ", string(trades))
		}
		queue, err := os.ReadFile(QUEUEPATH)
		if err == nil {
			fmt.Println("Queue: ", string(queue))
		}

		time.Sleep(time.Second * 2)
	}

}
