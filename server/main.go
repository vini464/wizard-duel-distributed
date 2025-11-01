package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
	"wizard-duel-distributed/api"
	"wizard-duel-distributed/utils"
)

var SERVERNAME string
var NETIP string

const SERVERPREFIX = 6                                   // quantidade de letras no prefixo do servername
var SERVERHEALTH map[string]bool = make(map[string]bool) // SERVERNAME: isAlive
var DEFAULTPORT = ":8080"
var LOGSPATH = "logs/logs.json"
var COMMANDQUEUE = make(utils.PriorityQueue, 0)
var MAPMUTEX sync.Mutex
var QUEUEMUTEX sync.Mutex

func checkPeerHealth(peerAddr string) {
	for {
		resp, err := http.Get(peerAddr + "/api/checkhealth")

		MAPMUTEX.Lock()
		if err != nil || resp.Status != "200 OK" {
			SERVERHEALTH[peerAddr] = false
			if peerAddr == "172.16.201.7" {
				fmt.Println("[debug] - Unable to connect with peer: ", peerAddr)
				fmt.Println("[debug] - err ", err)
				if err == nil {
					fmt.Println("[debug] - resp ", resp.Status)
				}
			}
		} else {
			SERVERHEALTH[peerAddr] = true
			respBody := []byte{}
			resp.Body.Read(respBody)
			fmt.Println("[debug] - ", peerAddr, " - ", resp.Status, " - ", string(respBody))
		}
		MAPMUTEX.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func executeCommands() {
	for {
		if len(COMMANDQUEUE) > 0 {
			Request(COMMANDQUEUE.Front().TimeStamp)
			fmt.Println("[debug] executing a command")
			c := COMMANDQUEUE.Pop()
			// propagando informação
			propagate(*c)
			ONCRITICALREGION = false
		}
	}
}

func propagate(command api.Command) {
	MAPMUTEX.Lock()
	for peer, alive := range SERVERHEALTH {
		if alive {
			com, _ := json.Marshal(command)
			http.Post(peer+"api/update", "application/json", bytes.NewBuffer(com))
		}
	}
}

func Instantiate() {
	SERVERNAME = utils.GetSelfAddres()
	fmt.Println(SERVERNAME)
	NETIP = utils.GetNetworkAddress()

	// quando um server inicia, ele procura por todos os servidores de 0 a 10 e adiciona no SERVERHEALTH
	fmt.Println("Server is starting")
	for i := range 255 {
		var peername = fmt.Sprintf("%s.%d", NETIP, i)
		fmt.Println(peername)
		if peername != SERVERNAME {
			go checkPeerHealth("http://" + peername + DEFAULTPORT)
		}
	}
	for peer, alive := range SERVERHEALTH {
		if alive {
			var logs, err = os.ReadFile(LOGSPATH)
			if err != nil {
				logs, _ = json.Marshal([]api.Command{}) // inicia um vetor vazio caso não consiga abrir o arquivo de logs
			}
			http.Post(peer+DEFAULTPORT+"/api/sync", "application/json", bytes.NewBuffer(logs))
		}
	}
	go executeCommands()
	handleRequests()
}
