package utils

import "net"

func PortInUse(port string) bool {
	listener, er := net.Listen("tcp", port)
	if er != nil {
		return true
	}
	defer listener.Close()

	return false
}
