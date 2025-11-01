package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"wizard-duel-distributed/api"
)

var ONCRITICALREGION bool = false

func Request(timestamp int64) {
	var wg sync.WaitGroup

	for peer := range SERVERHEALTH {
		wg.Add(1)
		go func() {
			buff, err := json.Marshal(api.Message{Type: "Request", Commands: []byte(strconv.FormatInt(timestamp, 10))})
			if err != nil {
				return
			}
			resp, err := http.Post("http://"+peer+DEFAULTPORT+"/api/token", "application/json", bytes.NewBuffer(buff))
			for err != nil || resp.Status != "200 OK" {
				resp, err = http.Post("http://"+peer+DEFAULTPORT+"/api/token", "application/json", bytes.NewBuffer(buff))
			}
			defer wg.Done()
		}()
	}
	wg.Wait()
	ONCRITICALREGION = true
}

func Reply(w http.ResponseWriter, r *http.Request) {
	var bufferedBody []byte
	var _, err = r.Body.Read(bufferedBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.Message{Type: "ERR"})
		return
	}
	var message api.Message
	var askingtimestamp int64
	err = json.NewDecoder(r.Body).Decode(&message)
	askingtimestamp, err = strconv.ParseInt(string(message.Commands), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.Message{Type: "ERR"})
		return
	}

	for ONCRITICALREGION || COMMANDQUEUE.Front() != nil && COMMANDQUEUE.Front().TimeStamp < askingtimestamp { // fica preso aqui até o outro servidor ter prioridade
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(api.Message{Type: "ACK"})
}
