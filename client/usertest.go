package main

import (
	"encoding/json"
	"net"
	"wizard-duel-distributed/communication"
	"wizard-duel-distributed/utils"
)

var HOSTNAME = utils.GetSelfAddres()

func main() {
	conn, err := net.Dial("tcp", HOSTNAME+communication.BROKERPORT)
	if err != nil {
		panic(err)
	}

	username := utils.Input("Digite seu username: ")
	password := utils.Input("Digite seu password: ")
	credentials := communication.Credentials{
	 Username: username,
		Password: password,
	}
	bytes, _ := json.Marshal(credentials)
	msg := communication.Message {
		Cmd: communication.PUBLISH,
		Tpc: "signup",
		Msg: bytes,
	}
	communication.SendMessage(conn, msg)
}
