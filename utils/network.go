package utils

import "net"

func GetSelfAddres() string{
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if (err != nil){
		return "0.0.0.0"
	}
	defer conn.Close()
	localAddress := conn.LocalAddr().(*net.UDPAddr)
	return string(localAddress.IP)
}

func GetNetworkAddress() string{
	localIp := GetNetworkAddress()
	// procurando pelo terceiro ponto
	p := 0
	pointId := len(localIp)
	for i, c := range localIp {
		if p < 3 && c == '.' {
			p++
			pointId = i;
		}
		if p == 3 {
			break
		}
	}
	return localIp[:pointId+1] // retorna o endereço sem o id do pc
}
