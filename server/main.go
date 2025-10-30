package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"wizard-duel-distributed/api"
)

var HASTOKEN = false

const SERVERPREFIX = 6                                   // quantidade de letras no prefixo do servername
var SERVERNAME = os.Getenv("SERVERNAME")                 // env var SERVERNAME ex SERVER1 - basicamente o id do servidor
var SERVERHEALTH map[string]bool = make(map[string]bool) // SERVERNAME: isAlive
var DEFAULTPORT = ":8080"
var LOGSPATH = "logs/logs.json"

func getToken(timestamp int64) {
	var wg sync.WaitGroup

	for peer, _ := range SERVERHEALTH {
		wg.Add(1)
		go func() {
			buff, err := json.Marshal(api.Message{Type: "Request", Commands: []byte(strconv.FormatInt(timestamp, 10))})
			if err != nil {
				return
			}
			http.Post(peer+DEFAULTPORT+"/api/token", "application/json", bytes.NewBuffer(buff))
			defer wg.Done()
		}()
	}
	wg.Wait()
	HASTOKEN = true
}

func removeDuplicates(array []api.Command) []api.Command {
	seen := make(map[string]bool)
	unique := []api.Command{}
	for _, e := range array {
		if !seen[e.ID] {
			seen[e.ID] = true
			unique = append(unique, e)
		}
	}
	return unique
}

func getLatest(logs []api.Command) []api.Command {
	latests := make(map[string]api.Command) // resource: command
	for _, command := range logs {
		com, ok := latests[command.ResourceID]
		if !ok || com.TimeStamp < command.TimeStamp {
			latests[command.ResourceID] = command
		}
	}
	uniqueLogs := []api.Command{}
	for _, command := range logs {
		uniqueLogs = append(uniqueLogs, command)
	}
	return uniqueLogs
}

func getHealthCheck(w http.ResponseWriter, r *http.Request) {
	message := api.Message{Type: "ACK"} // indica que recebeu a mensagem e que ta vivo
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}

func syncLogs(w http.ResponseWriter, r *http.Request) {
	var bufferedBody []byte
	var i, err = r.Body.Read(bufferedBody)
	if err != nil || i <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.Message{Type: "SyncResponse"})
		return
	}
	var bodyMessage []api.Command
	err = json.Unmarshal(bufferedBody, &bodyMessage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.Message{Type: "SyncResponse"})
		return
	}
	selfLogsBytes, err := os.ReadFile(LOGSPATH)
	if err != nil {
		selfLogsBytes, _ = json.Marshal("[]") // inicia um vetor vazio caso não consiga abrir o arquivo de logs
	}
	var logs []api.Command
	err = json.Unmarshal(selfLogsBytes, &logs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.Message{Type: "SyncResponse"})
		return
	}
	logs = append(logs, bodyMessage...)
	logs = removeDuplicates(logs) // sem comandos repetidos
	logs = getLatest(logs)        // diff completa
	// TODO: Executar commandos do log  -> vai entrar na danger zone
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logs)
}

func checkPeerHealth(peerAddr string) {
	for {
		resp, err := http.Get(peerAddr + DEFAULTPORT + "/api/checkhealth")
		respBody, _ := io.ReadAll(resp.Body)
		if err != nil || resp.Status != "200 OK" {
			SERVERHEALTH[peerAddr] = false
			fmt.Println("[debug] - Unable to connect with peer: ", peerAddr)
		} else {
			SERVERHEALTH[peerAddr] = true
		}
		fmt.Println("[debug] - ", peerAddr, " - ", resp.Status, " - ", string(respBody))
		time.Sleep(1 * time.Second)
	}
}

func handleRequests() {
	http.Handle("GET /api/checkhealth", http.HandlerFunc(getHealthCheck))
	http.Handle("POST /api/sync", http.HandlerFunc(syncLogs))
}

func main() {

	// quando um server inicia, ele procura por todos os servidores de 0 a 10 e adiciona no SERVERHEALTH
	fmt.Println("Server is starting")
	for i := range 10 {
		var peername = fmt.Sprintf("%s-%d", SERVERNAME[:SERVERPREFIX], i)
		if peername != SERVERNAME {
			go checkPeerHealth(peername)
		}
	}
	for peer, alive := range SERVERHEALTH {
		if alive {
			var logs, err = os.ReadFile(LOGSPATH)
			if err != nil {
				logs, _ = json.Marshal("[]") // inicia um vetor vazio caso não consiga abrir o arquivo de logs
			}
			http.Post(peer+DEFAULTPORT+"/api/sync", "application/json", bytes.NewBuffer(logs))
		}
	}

}
