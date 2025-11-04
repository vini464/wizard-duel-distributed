package utils

import (
	"net"
)

func GetSelfAddres() string{
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if (err != nil){
		return "0.0.0.0"
	}
	defer conn.Close()
	localAddress := conn.LocalAddr().(*net.UDPAddr)
	return localAddress.IP.String()
}

func GetNetworkAddress(ip_address string) string{
	// procurando pelo terceiro ponto
	times := 0
	id := 0
	for i, c := range ip_address {
		if c == '.' {
			times ++ 
			if times == 3 {
				id  = i
				break
			}
		}
	}
	return ip_address[:id]
}
