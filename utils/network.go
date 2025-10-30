package utils

import "net"

func GetSelfAddres() string{
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if (err != nil){
		return "0.0.0.0"
	}
	defer conn.Close()
	localAddress := conn.LocalAddr().(*net.UDPAddr)
	return localAddress.IP.String()
}

func GetNetworkAddress() string{
	// procurando pelo terceiro ponto
	return Input("digite o ip da rede: ")
}
