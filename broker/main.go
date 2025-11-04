package main

import (
	"fmt"
	"net"

	"wizard-duel-distributed/communication"
	"wizard-duel-distributed/utils"
)

var HOSTNAME = utils.GetSelfAddres()

func main() {
	listener, err := net.Listen(communication.SERVERTYPE, HOSTNAME+ communication.BROKERPORT)
	if err != nil {
		panic(err)
	}
	fmt.Println("[debug] - broker online")

	broker := NewBroker()

	for {
		conn, err := listener.Accept()
		if err == nil {
			go handle_connection(conn, broker)
		}
	}

}

// Essa função lida com a comunicação dos clientes
// [4 bytes - tamanho da mesagem][X bytes - messagem ]
// Client -> Broker
// Subscribe:
// [CMD: Subscribe; TOPIC: topic_name]
// Publish:
// [CMD: Publish; TOPIC: topic_name; MESSAGE: message_body]
//
// Broker -> Client
// [TYPE: status; VALUE: OK/ERROR]
// [TYPE: msg; VALUE: msg_body]
func handle_connection(conn net.Conn, broker *Broker) {
	fmt.Println("[debug] - New client connected: ", conn.RemoteAddr())
	defer conn.Close()
	communication_chan := make(chan communication.Message)

	// receiving a client message
	go func() {
		var msg communication.Message
		for {
			err := communication.ReceiveMessage(conn, &msg)
			if err != nil {
				return
			}
			communication_chan <- msg
		}
	}()

	for {
		select {
		case msg := <-communication_chan:
			switch msg.Cmd {
			case communication.SUBSCRIBE:
				fmt.Println("[debug] - subs")
				broker.Subscribe(msg.Tpc, conn)
			case communication.UNSUB:
				fmt.Println("[debug] - unsub")
				broker.Unsubscribe(msg.Tpc, conn)
			case communication.PUBLISH:
				fmt.Println("[debug] - Publish")
				broker.Publish(msg.Tpc, msg.Msg)
			}
		case <-broker.Quit:
			fmt.Println("[DEBUG] - Broker is stopped")
			return
		}
	}
}
