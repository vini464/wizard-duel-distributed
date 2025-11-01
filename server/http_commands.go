package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"wizard-duel-distributed/api"
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
	fmt.Println("Executando codigo")
	w.WriteHeader(http.StatusOK)
}
