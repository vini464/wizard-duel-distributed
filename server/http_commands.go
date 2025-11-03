package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"wizard-duel-distributed/api"
	"wizard-duel-distributed/models"
	"wizard-duel-distributed/utils"
)

func handleRequests() {
	http.Handle("GET /api/checkhealth", http.HandlerFunc(getHealthCheck))
	http.Handle("POST /api/sync", http.HandlerFunc(syncLogs))
	http.Handle("POST /api/request", http.HandlerFunc(Reply))
	http.Handle("POST /api/update", http.HandlerFunc(update))
	log.Fatal(http.ListenAndServe(SERVERNAME+DEFAULTPORT, nil))
}

func getHealthCheck(w http.ResponseWriter, r *http.Request) {
	message := api.Message{Type: "ACK"} // indica que recebeu a mensagem e que ta vivo
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}

func syncLogs(w http.ResponseWriter, r *http.Request) {
	var bodyMessage []api.Command
	err := json.NewDecoder(r.Body).Decode(&bodyMessage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.Message{Type: "ERROR"})
		return
	}

	var logs []api.Command
	err = utils.DecodeFromFile(LOGSPATH, &logs, []api.Command{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.Message{Type: "ERROR"})
		return
	}
	logs = append(logs, bodyMessage...)
	logs = removeDuplicates(logs) // sem comandos repetidos
	logs = getLatest(logs)        // diff completa

	// executar comandos do log, sem pedir permissão pois é pra sincronizar os logs

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logs)
}

func update(w http.ResponseWriter, r *http.Request) {
	DATAMUTEX.Lock()
	defer DATAMUTEX.Unlock()
	fmt.Println("Executando codigo")
	var command api.Command
	err := json.NewDecoder(r.Body).Decode(&command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.Message{Type: "ERROR"})
		return
	}
	UpdateLogs(LOGSPATH, command)
	switch command.Resource {
	case "player":
		players := models.RetrievePlayers(PLAYERSPATH)
		var player models.Player
		json.Unmarshal(command.Value, &player)

		switch command.Operation {
		case "create":
			players = append(players, player)
		case "update":
			models.UpdatePlayer(player.Password, player, players)
		}
		models.SavePlayers(PLAYERSPATH, players)
	case "match":
		matches := models.RetrieveMatches(MATCHESPATH)
		var match models.Match
		json.Unmarshal(command.Value, &match)
		switch command.Operation {
		case "create":
			matches = append(matches, match)
		case "update":
			models.UpdateMatch(match, matches)
		}
		models.SaveMatches(MATCHESPATH, matches)
	case "trade":
		trades := models.RetrieveTrades(TRADESPATH)
		var trade models.Trade
		json.Unmarshal(command.Value, &trade)
		switch command.Operation {
		case "create":
			trades = append(trades, trade)
		case "update":
			models.UpdateTrade(trade, trades)
		}
		models.SaveTrades(TRADESPATH, trades)
	case "queue", "players":
		path := QUEUEPATH
		if command.Resource == "players" {
			path = PLAYERSPATH
		}
		file, err := os.Create(path)
		if err == nil {
			file.Write(command.Value)
			file.Close()
		}
	default:
		fmt.Println("[debug]: Unknown Command")

	}
	w.WriteHeader(http.StatusOK)
}
